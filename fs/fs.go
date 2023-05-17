package fs

import (
	iofs "io/fs"
)

type RenameFS interface {
	iofs.FS
	Rename(oldPath, newPath string) error
}

type RemoveFS interface {
	iofs.FS
	Remove(name string) error
}

type WriteFileFS interface {
	WriteFile(name string, data []byte, perm iofs.FileMode) error
}

type ExistsFS interface {
	Exists(path string) (bool, error)
}

type MakeDirFS interface {
	Mkdir(path string, perm iofs.FileMode) error
	MkdirAll(path string, perm iofs.FileMode) error
}

type CreateFS interface {
	Create(path string) (iofs.File, error)
}

type FS interface {
	iofs.FS
	CreateFS
	RenameFS
	RemoveFS
	WriteFileFS
	ExistsFS
	iofs.GlobFS
	iofs.ReadFileFS
	iofs.ReadFileFS
	iofs.StatFS
	iofs.SubFS
	iofs.ReadDirFS
	MakeDirFS
}
