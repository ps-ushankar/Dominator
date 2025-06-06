package client

import (
	"crypto/sha512"
	"fmt"
	"io"

	"github.com/Cloud-Foundations/Dominator/lib/errors"
	"github.com/Cloud-Foundations/Dominator/lib/hash"
	"github.com/Cloud-Foundations/Dominator/lib/queue"
	"github.com/Cloud-Foundations/Dominator/lib/srpc"
	"github.com/Cloud-Foundations/Dominator/proto/objectserver"
)

func newObjectAdderQueue(client srpc.ClientI) (*ObjectAdderQueue, error) {
	var objQ ObjectAdderQueue
	var err error
	objQ.conn, err = client.Call("ObjectServer.AddObjects")
	if err != nil {
		return nil, err
	}
	getResponseSendChan, getResponseReceiveChan := queue.NewEventQueue()
	errorChan := make(chan error, 1024)
	objQ.getResponseChan = getResponseSendChan
	objQ.errorChan = errorChan
	objQ.sendSemaphore = make(chan struct{}, 1)
	go readResponses(objQ.conn, getResponseReceiveChan, errorChan)
	return &objQ, nil
}

func (objQ *ObjectAdderQueue) add(reader io.Reader, length uint64) (
	hash.Hash, error) {
	var hashVal hash.Hash
	data := make([]byte, length)
	nRead, err := io.ReadFull(reader, data)
	if err != nil {
		return hashVal, err
	}
	if uint64(nRead) != length {
		return hashVal, fmt.Errorf(
			"failed to read file data, wanted: %d, got: %d bytes",
			length, nRead)
	}
	hasher := sha512.New()
	if _, err := hasher.Write(data); err != nil {
		return hashVal, err
	}
	copy(hashVal[:], hasher.Sum(nil))
	err = objQ.addData(data, hashVal)
	return hashVal, err
}

func (objQ *ObjectAdderQueue) addData(data []byte, hashVal hash.Hash) error {
	if err := objQ.consumeErrors(false); err != nil {
		return err
	}
	// Send in a goroutine to increase concurrency. A small win.
	objQ.sendSemaphore <- struct{}{}
	go func() {
		defer func() {
			<-objQ.sendSemaphore
		}()
		var request objectserver.AddObjectRequest
		request.Length = uint64(len(data))
		request.ExpectedHash = &hashVal
		objQ.conn.Encode(request)
		objQ.conn.Write(data)
		objQ.getResponseChan <- struct{}{}
	}()
	return nil
}

func (objQ *ObjectAdderQueue) close() error {
	// Wait for any sends in progress to complete.
	objQ.sendSemaphore <- struct{}{}
	var request objectserver.AddObjectRequest
	err := objQ.conn.Encode(request)
	err = updateError(err, objQ.conn.Flush())
	close(objQ.getResponseChan)
	err = updateError(err, objQ.consumeErrors(true))
	return updateError(err, objQ.conn.Close())
}

func updateError(oldError, newError error) error {
	if oldError == nil {
		return newError
	}
	return oldError
}

func (objQ *ObjectAdderQueue) consumeErrors(untilClose bool) error {
	if untilClose {
		for err := range objQ.errorChan {
			if err != nil {
				return err
			}
		}
		return nil
	}
	for {
		select {
		case err := <-objQ.errorChan:
			if err != nil {
				return err
			}
		default:
			return nil
		}
	}
}

func readResponses(conn *srpc.Conn, getResponseChan <-chan struct{},
	errorChan chan<- error) {
	defer close(errorChan)
	for range getResponseChan {
		var reply objectserver.AddObjectResponse
		if err := conn.Decode(&reply); err != nil {
			errorChan <- err
			break // Transport error: abandon ship.
		}
		errorChan <- errors.New(reply.ErrorString)
	}
}
