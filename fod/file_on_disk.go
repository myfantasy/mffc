// Package fod is FileProvider over disk
package fod

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/myfantasy/mdp"
	"github.com/myfantasy/mffc/fh"
)

// FileOnDisk file on disk profider
type FileOnDisk struct {
	FolderPerm os.FileMode
	FilePerm   os.FileMode
}

// FileOnDiskSimpleCreate creates simple FileOnDisk with perms 0760 & 0660
func FileOnDiskSimpleCreate() *FileOnDisk {
	return &FileOnDisk{
		FolderPerm: 0760,
		FilePerm:   0660,
	}
}

// FileReplace replace file on disk like fh.FileReplace
func (fod *FileOnDisk) FileReplace(path string, data []byte) error {
	return fh.FileReplace(fod, path, data)
}

// FileLoad load file from disk like fh.FileLoad
func (fod *FileOnDisk) FileLoad(path string) (data []byte, e bool, needResave bool, err error) {
	return fh.FileLoad(fod, path)
}

// FileLoadAndFix load file from disk like fh.FileLoadAndFix
func (fod *FileOnDisk) FileLoadAndFix(path string) (data []byte, e bool, err error) {
	return fh.FileLoadAndFix(fod, path)
}

// Clone file provider
func (fod *FileOnDisk) Clone() fh.FileProviderA {
	return &FileOnDisk{
		FolderPerm: fod.FolderPerm,
		FilePerm:   fod.FilePerm,
	}
}

// Exists file on disk
func (fod *FileOnDisk) Exists(path string) (bool, error) {
	path = filepath.FromSlash(path)

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Read from disk data, flag that exists file, or error
func (fod *FileOnDisk) Read(path string) (data []byte, e bool, err error) {

	path = filepath.FromSlash(path)

	data, err = ioutil.ReadFile(path)

	if err == nil {
		return data, true, nil
	} else if os.IsNotExist(err) {
		return data, false, nil
	}
	return data, false, err
}

// MkDirIfNotExists make directory
func (fod *FileOnDisk) MkDirIfNotExists(path string) (err error) {

	path = filepath.FromSlash(path)

	ok, err := fod.Exists(path)
	if err != nil {
		return mdp.ErrorNew("MkDirIfNotExists Check directory "+path, err)
	}
	if !ok {
		err = os.MkdirAll(path, fod.FolderPerm)
		if err != nil {
			return mdp.ErrorNew("MkDirIfNotExists Mkdir file "+path, err)
		}
	}

	return nil
}

// Write file on disk
func (fod *FileOnDisk) Write(path string, data []byte) error {
	path = filepath.FromSlash(path)

	return ioutil.WriteFile(path, data, fod.FilePerm)
}

// Remove file from disk
func (fod *FileOnDisk) Remove(path string) error {
	path = filepath.FromSlash(path)

	return os.RemoveAll(path)
}

// Rename files
func (fod *FileOnDisk) Rename(pathOld string, pathNew string) error {
	pathOld = filepath.FromSlash(pathOld)
	pathNew = filepath.FromSlash(pathNew)

	return os.Rename(pathOld, pathNew)
}
