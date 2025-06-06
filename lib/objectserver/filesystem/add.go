package filesystem

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/Cloud-Foundations/Dominator/lib/backoffdelay"
	"github.com/Cloud-Foundations/Dominator/lib/fsutil"
	"github.com/Cloud-Foundations/Dominator/lib/hash"
	"github.com/Cloud-Foundations/Dominator/lib/objectcache"
	"github.com/Cloud-Foundations/Dominator/lib/objectserver"
)

const (
	buflen = 65536
)

// This must be called with the lock held. The object must not already exist.
func (objSrv *ObjectServer) add(object *objectType) {
	objSrv.objects[object.hash] = object
	objSrv.addUnreferenced(object)
	objSrv.lastMutationTime = time.Now()
	objSrv.totalBytes += object.size
}

func (objSrv *ObjectServer) addObject(reader io.Reader, length uint64,
	expectedHash *hash.Hash) (hash.Hash, bool, error) {
	hashVal, data, err := objectcache.ReadObject(reader, length, expectedHash)
	if err != nil {
		return hashVal, false, err
	}
	filename := path.Join(objSrv.BaseDirectory,
		objectcache.HashToFilename(hashVal))
	// Check for existing object and collision.
	if isNew, err := objSrv.addOrCompare(hashVal, data, filename); err != nil {
		return hashVal, false, err
	} else {
		object := &objectType{hash: hashVal, size: uint64(len(data))}
		objSrv.rwLock.Lock()
		if _, ok := objSrv.objects[object.hash]; !ok {
			objSrv.add(object)
		}
		objSrv.rwLock.Unlock()
		if objSrv.addCallback != nil {
			objSrv.addCallback(hashVal, uint64(len(data)), isNew)
		}
		return hashVal, isNew, nil
	}
}

// addOrCompare returns the following:
//
//	a boolean which is true if the object is new
//	an error or nil if no error.
func (objSrv *ObjectServer) addOrCompare(hashVal hash.Hash, data []byte,
	filename string) (bool, error) {
	sleeper := backoffdelay.NewExponential(time.Duration(len(data)),
		time.Second, 1)
	gc := objSrv.gc
	var firstRetryTime time.Time
	var loggedRetry bool
	for {
		isNew, err := objSrv.addOrCompareOnce(hashVal, data, filename, &gc)
		if err == nil {
			return isNew, nil
		}
		if !os.IsExist(err) {
			return false, err
		}
		if !loggedRetry {
			if firstRetryTime.IsZero() {
				firstRetryTime = time.Now()
			} else if time.Since(firstRetryTime) > time.Second {
				objSrv.Logger.Printf("retrying: %s\n", err)
				loggedRetry = true
			}
		}
		sleeper.Sleep()
	}
}

func (objSrv *ObjectServer) addOrCompareOnce(hashVal hash.Hash, data []byte,
	filename string, gc *objectserver.GarbageCollector) (bool, error) {
	fi, err := os.Lstat(filename)
	if err == nil {
		if !fi.Mode().IsRegular() {
			return false, errors.New("existing non-file: " + filename)
		}
		if err := collisionCheck(data, filename, fi.Size()); err != nil {
			return false, errors.New("collision detected: " + err.Error())
		}
		// No collision and no error: it's the same object. Go home early.
		return false, nil
	}
	if *gc != nil { // Have external garbage collector: trigger it inline.
		objSrv.garbageCollector(nil)
		*gc = nil
	}
	err = os.MkdirAll(path.Dir(filename), fsutil.PrivateDirPerms)
	if err != nil {
		return false, err
	}
	err = fsutil.CopyToFileExclusive(filename, fsutil.PrivateFilePerms,
		bytes.NewReader(data), uint64(len(data)))
	if err != nil {
		return false, err
	}
	return true, nil
}

func collisionCheck(data []byte, filename string, size int64) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	if int64(len(data)) != size {
		return fmt.Errorf("length mismatch. Data=%d, existing object=%d",
			len(data), size)
	}
	reader := bufio.NewReader(file)
	buffer := make([]byte, 0, buflen)
	for len(data) > 0 {
		numToRead := len(data)
		if numToRead > cap(buffer) {
			numToRead = cap(buffer)
		}
		buf := buffer[:numToRead]
		nread, err := reader.Read(buf)
		if err != nil {
			return err
		}
		if bytes.Compare(data[:nread], buf[:nread]) != 0 {
			return errors.New("content mismatch")
		}
		data = data[nread:]
	}
	return nil
}
