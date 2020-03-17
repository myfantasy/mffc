package fod

import (
	"testing"
)

func TestFileDeflateOnDisk(t *testing.T) {
	fod, err := FileDeflateOnDiskSimpleCreate()

	if err != nil {
		t.Fatal("Create", err)
	}

	err = fod.MkDirIfNotExists("test")
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

	err = fod.Append("test/t3.txt", []byte(`aaassdd`))
	if err != nil {
		t.Fatal("Append1", err)
	}
	err = fod.Append("test/t3.txt", []byte("\r\nddaaassdd"))
	if err != nil {
		t.Fatal("Append2", err)
	}
	data, e, err = fod.Read("test/t3.txt")
	if err != nil {
		t.Fatal("Read after append", err)
	}
	if !e {
		t.Fatal("Read after append Not exists")
	}

	if string(data) != "aaassdd\r\nddaaassdd" {
		t.Fatal("Read after append Not correct data")
	}

	err = fod.Remove("test/t3.txt")
	if err != nil {
		t.Fatal("Remove File", err)
	}

	err = fod.Remove("test")
	if err != nil {
		t.Fatal("Remove Folder", err)
	}

}
