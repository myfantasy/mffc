// Package fod is FileProvider over disk
package fod

import (
	"compress/flate"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/myfantasy/mdp"
	"github.com/myfantasy/mffc/compress"
	"github.com/myfantasy/mffc/fh"
)

// FileDeflateOnDisk file on disk profider
type FileDeflateOnDisk struct {
	FolderPerm os.FileMode
	FilePerm   os.FileMode
	Compressor *compress.DeflateCompressor
	FileTail   string
}

// FileDeflateOnDiskSimpleCreate creates simple FileDeflateOnDisk with perms 0760 & 0660
// and Best compression
// and FileTail .gzips
func FileDeflateOnDiskSimpleCreate() (*FileDeflateOnDisk, error) {
	compressor, err := compress.DeflateCompressorCreate(flate.BestCompression)

	return &FileDeflateOnDisk{
		FolderPerm: 0760,
		FilePerm:   0660,
		Compressor: compressor,
		FileTail:   ".gzips",
	}, err
}

// FileDeflateOnDiskLevelCreate creates simple FileDeflateOnDisk with perms 0760 & 0660 and compressor
func FileDeflateOnDiskLevelCreate(level int) (*FileDeflateOnDisk, error) {
	compressor, err := compress.DeflateCompressorCreate(level)

	return &FileDeflateOnDisk{
		FolderPerm: 0760,
		FilePerm:   0660,
		Compressor: compressor,
	}, err
}

// FileReplace replace file on disk like fh.FileReplace
func (fod *FileDeflateOnDisk) FileReplace(path string, data []byte) error {
	return fh.FileReplace(fod, path, data)
}

// FileLoad load file from disk like fh.FileLoad
func (fod *FileDeflateOnDisk) FileLoad(path string) (data []byte, e bool, needResave bool, err error) {
	return fh.FileLoad(fod, path)
}

// FileLoadAndFix load file from disk like fh.FileLoadAndFix
func (fod *FileDeflateOnDisk) FileLoadAndFix(path string) (data []byte, e bool, err error) {
	return fh.FileLoadAndFix(fod, path)
}

// Clone file provider
func (fod *FileDeflateOnDisk) Clone() fh.FileProviderA {

	return &FileDeflateOnDisk{
		FolderPerm: fod.FolderPerm,
		FilePerm:   fod.FilePerm,
		Compressor: fod.Compressor.Clone(),
		FileTail:   fod.FileTail,
	}
}

// Exists file on disk
func (fod *FileDeflateOnDisk) Exists(path string) (bool, error) {
	path = filepath.FromSlash(path + fod.FileTail)

	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// Read from disk data, flag that exists file, or error
func (fod *FileDeflateOnDisk) Read(path string) (data []byte, e bool, err error) {

	path = filepath.FromSlash(path + fod.FileTail)

	data, err = ioutil.ReadFile(path)

	if err == nil {

		data, err = fod.Compressor.Restore(data)

		return data, true, err
	} else if os.IsNotExist(err) {
		return data, false, nil
	}
	return data, false, err
}

// MkDirIfNotExists make directory
func (fod *FileDeflateOnDisk) MkDirIfNotExists(path string) (err error) {

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
func (fod *FileDeflateOnDisk) Write(path string, data []byte) error {
	path = filepath.FromSlash(path + fod.FileTail)

	d, err := fod.Compressor.Compress(data)
	if err != nil {
		return mdp.ErrorNew("Write fail to compress "+path, err)
	}

	return ioutil.WriteFile(path, d, fod.FilePerm)
}

// Remove file from disk
func (fod *FileDeflateOnDisk) Remove(path string) error {
	path = filepath.FromSlash(path + fod.FileTail)

	return os.RemoveAll(path)
}

// Rename files
func (fod *FileDeflateOnDisk) Rename(pathOld string, pathNew string) error {
	pathOld = filepath.FromSlash(pathOld + fod.FileTail)
	pathNew = filepath.FromSlash(pathNew + fod.FileTail)

	return os.Rename(pathOld, pathNew)
}
