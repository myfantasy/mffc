package fod

import (
	"testing"
)

func TestFileOnDisk(t *testing.T) {
	fod := FileOnDiskSimpleCreate()
	err := fod.MkDirIfNotExists("test")
	if err != nil {
		t.Fatal("MkDirIfNotExists", err)
	}

	err = fod.FileReplace("test/test.json", []byte(`{"aaa":"bbb"}`))
	if err != nil {
		t.Fatal("FileReplace1", err)
	}

	err = fod.FileReplace("test/test.json", []byte(`{"aaa":"ccc"}`))
	if err != nil {
		t.Fatal("FileReplace2", err)
	}

	data, e, needResave, err := fod.FileLoad("test/test.json")
	if err != nil {
		t.Fatal("FileLoad", err)
	}
	if needResave {
		t.Fatal("FileLoad -> needResave")
	}
	if !e {
		t.Fatal("FileLoad -> e")
	}
	if string(data) != `{"aaa":"ccc"}` {
		t.Fatal("FileLoad -> data")
	}

	err = fod.FileReplace("test/test2.json", []byte(`{"aaa":"ccc"}`))

	err = fod.Remove("test/test2.json")
	if err != nil {
		t.Fatal("Remove File", err)
	}

	err = fod.Remove("test")
	if err != nil {
		t.Fatal("Remove Folder", err)
	}

}
