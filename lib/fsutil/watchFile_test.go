package fsutil

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"
	"time"

	"github.com/Cloud-Foundations/Dominator/lib/log/testlogger"
)

var errorTimeout = errors.New("timeout")

func TestWatchFileDir(t *testing.T) {
	dirname, err := ioutil.TempDir("", "WatchFileTests")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirname)
	testWatchFile(t, dirname)
}

func TestWatchFileCwd(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)
	dirname, err := ioutil.TempDir("", "WatchFileTests")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dirname)
	if err := os.Chdir(dirname); err != nil {
		t.Fatal(err)
	}
	testWatchFile(t, "")
}

func testWatchFile(t *testing.T, dirname string) {
	logger := testlogger.New(t)
	pathNotExist := path.Join(dirname, "never-exists")
	ch := WatchFile(pathNotExist, logger)
	rc, err := watchTimeout(ch, time.Millisecond*50)
	if err != errorTimeout {
		rc.Close()
		t.Fatal("Expected timeout error for non-existent file")
	}
	pathExists := path.Join(dirname, "exists")
	file, err := os.Create(pathExists)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	ch = WatchFile(pathExists, logger)
	rc, err = watchTimeout(ch, time.Millisecond*50)
	if err != nil {
		t.Fatal(err)
	} else {
		rc.Close()
	}
	pathExistsLater := path.Join(dirname, "exists-later")
	errorChannel := make(chan error, 1)
	go func() {
		time.Sleep(time.Millisecond * 50)
		file, err := os.Create(pathExistsLater)
		if file != nil {
			file.Close()
		}
		errorChannel <- err
	}()
	ch = WatchFile(pathExistsLater, logger)
	_, err = watchTimeout(ch, time.Millisecond*10)
	if err != errorTimeout {
		t.Fatal("Expected timeout error for non-existent file")
	}
	rc, err = watchTimeout(ch, time.Millisecond*90)
	if err != nil {
		t.Fatal(err)
	} else {
		rc.Close()
	}
	pathWillBeRenamed := path.Join(dirname, "will-be-renamed")
	file, err = os.Create(pathWillBeRenamed)
	if err != nil {
		t.Fatal(err)
	}
	file.Close()
	rc, err = watchTimeout(ch, time.Millisecond*10)
	if err != errorTimeout {
		rc.Close()
		t.Fatal("Expected timeout error for unchanged file")
	}
	if err := os.Rename(pathWillBeRenamed, pathExistsLater); err != nil {
		t.Fatal(err)
	}
	rc, err = watchTimeout(ch, time.Millisecond*50)
	if err != nil {
		t.Fatal(err)
	} else {
		rc.Close()
	}
	if err := <-errorChannel; err != nil {
		t.Fatal(err)
	}
}

func watchTimeout(channel <-chan io.ReadCloser, timeout time.Duration) (
	io.ReadCloser, error) {
	select {
	case readCloser := <-channel:
		return readCloser, nil
	case <-time.After(timeout):
		return nil, errorTimeout
	}
}
