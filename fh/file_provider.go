// Package fh - file handler
// You can use any storage type
package fh

import (
	"path/filepath"

	"github.com/myfantasy/mdp"
)

// FileProvider - file or directory create update and remove
type FileProvider interface {
	Exists(path string) (bool, error)
	// data, exists file, error
	Read(path string) (data []byte, e bool, err error)
	MkDirIfNotExists(path string) (err error)
	Write(path string, data []byte) error
	Remove(path string) error
	Rename(pathOld string, pathNew string) error
}

// FileProviderA - file or directory create update and remove and additional methods
type FileProviderA interface {
	FileProvider

	FileReplace(path string, data []byte) error
	FileLoad(path string) (data []byte, e bool, needResave bool, err error)
	FileLoadAndFix(path string) (data []byte, e bool, err error)
	// clone provider
	Clone() FileProviderA
}

// FileReplace - Create or replace file (remove .old && .new if exists -> create .new -> move current to .old -> move .new to current -> remove .old)
func FileReplace(fp FileProvider, path string, data []byte) error {
	pathOld := path + ".old"
	pathNew := path + ".new"

	e, err := fp.Exists(pathOld)
	if err != nil {
		return mdp.ErrorNew("FileReplace Check old file "+pathOld, err)
	}
	if e {
		err = fp.Remove(pathOld)
		if err != nil {
			return mdp.ErrorNew("FileReplace Remove old file (1) "+pathOld, err)
		}
	}

	e, err = fp.Exists(pathNew)
	if err != nil {
		return mdp.ErrorNew("FileReplace Check new file "+pathNew, err)
	}
	if e {
		err = fp.Remove(pathNew)
		if err != nil {
			return mdp.ErrorNew("FileReplace Remove new file "+pathNew, err)
		}
	}

	err = fp.Write(pathNew, data)
	if err != nil {
		return mdp.ErrorNew("FileReplace Write new file "+pathNew, err)
	}

	e, err = fp.Exists(path)
	if err != nil {
		return mdp.ErrorNew("FileReplace Check file "+path, err)
	}
	if e {
		err = fp.Rename(path, pathOld)
		if err != nil {
			return mdp.ErrorNew("FileReplace move "+path+" to new file "+pathOld, err)
		}
	}

	err = fp.Rename(pathNew, path)
	if err != nil {
		return mdp.ErrorNew("FileReplace move new file "+pathNew+" to "+path, err)
	}

	e, err = fp.Exists(pathOld)
	if err != nil {
		return mdp.ErrorNew("FileReplace Check old file "+pathOld, err)
	}
	if e {
		err = fp.Remove(pathOld)
		if err != nil {
			return mdp.ErrorNew("FileReplace Remove old file (2) "+pathOld, err)
		}
	}

	return nil
}

//FileLoad load file path -> .new -> .old
func FileLoad(fp FileProvider, path string) (data []byte, e bool, needResave bool, err error) {

	path = filepath.FromSlash(path)

	pathOld := path + ".old"
	pathNew := path + ".new"

	data, e, err = fp.Read(path)
	if err != nil {
		return data, e, needResave, mdp.ErrorNew("FileLoad Check file "+path, err)
	}
	if e {
		return data, e, needResave, nil
	}

	data, e, err = fp.Read(pathNew)
	if err != nil {
		return data, e, needResave, mdp.ErrorNew("FileLoad Check new file "+pathNew, err)
	}
	if e {
		needResave = true
		return data, e, needResave, nil
	}

	data, e, err = fp.Read(pathOld)
	if err != nil {
		return data, e, needResave, mdp.ErrorNew("FileLoad Check old file "+pathOld, err)
	}
	if e {
		needResave = true
		return data, e, needResave, nil
	}

	return data, false, false, nil
}

//FileLoadAndFix load file path -> .new -> .old and move file to path if path is not exists
func FileLoadAndFix(fp FileProvider, path string) (data []byte, e bool, err error) {

	path = filepath.FromSlash(path)

	pathOld := path + ".old"
	pathNew := path + ".new"

	data, e, err = fp.Read(path)
	if err != nil {
		return data, e, mdp.ErrorNew("FileLoad Check file "+path, err)
	}
	if e {
		return data, e, nil
	}

	data, e, err = fp.Read(pathNew)
	if err != nil {
		return data, e, mdp.ErrorNew("FileLoad Check new file "+pathNew, err)
	}
	if e {
		err = fp.Rename(pathNew, path)
		if err != nil {
			err = mdp.ErrorNew("FileLoadAndFix move new file "+pathNew+" to "+path, err)
		}
		return data, e, err
	}

	data, e, err = fp.Read(pathOld)
	if err != nil {
		return data, e, mdp.ErrorNew("FileLoad Check old file "+pathOld, err)
	}
	if e {
		err = fp.Rename(pathOld, path)
		if err != nil {
			err = mdp.ErrorNew("FileLoadAndFix move new file "+pathOld+" to "+path, err)
		}
		return data, e, err
	}

	return data, false, nil
}
