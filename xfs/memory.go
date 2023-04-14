package xfs

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	fstest "testing/fstest"

	"github.com/patrickhuber/go-xplat/xfilepath"
)

type memory struct {
	fs            fstest.MapFS
	pathSeperator xfilepath.PathSeparator
}

func NewMemory(options ...MemoryOption) FS {
	m := &memory{
		fs:            fstest.MapFS{},
		pathSeperator: xfilepath.Default,
	}
	for _, op := range options {
		op(m)
	}
	return m
}

type MemoryOption = func(*memory)

func WithPathSeperator(sep xfilepath.PathSeparator) MemoryOption {
	return func(m *memory) {
		m.pathSeperator = sep
	}
}

// Open implements FS
func (m *memory) Open(name string) (fs.File, error) {
	return m.fs.Open(name)
}

// Rename implements FS
func (m *memory) Rename(oldPath string, newPath string) error {
	file, ok := m.fs[oldPath]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.fs, oldPath)
	m.fs[newPath] = file
	return nil
}

// Remove implements FS
func (m *memory) Remove(path string) error {
	_, ok := m.fs[path]
	if !ok {
		return os.ErrNotExist
	}
	delete(m.fs, path)
	return nil
}

// RemoveAll implements FS
func (m *memory) RemoveAll(path string) error {
	paths := []string{}
	for p := range m.fs {
		if strings.HasPrefix(p, path) {
			paths = append(paths, p)
		}
	}
	for _, p := range paths {
		delete(m.fs, p)
	}
	return nil
}

// Glob implements FS
func (m *memory) Glob(pattern string) ([]string, error) {
	return m.fs.Glob(pattern)
}

// ReadDir implements FS
func (m *memory) ReadDir(name string) ([]fs.DirEntry, error) {
	return m.fs.ReadDir(name)
}

// ReadFile implements FS
func (m *memory) ReadFile(name string) ([]byte, error) {
	return m.fs.ReadFile(name)
}

// WriteFile implements FS
func (m *memory) WriteFile(name string, data []byte, perm os.FileMode) error {
	file, ok := m.fs[name]
	if !ok {
		file = &fstest.MapFile{}
		m.fs[name] = file
	}
	file.Data = data
	file.Mode = perm
	return nil
}

// Exists implements FS
func (m *memory) Exists(path string) (bool, error) {
	_, ok := m.fs[path]
	return ok, nil
}

// Stat implements FS
func (m *memory) Stat(name string) (fs.FileInfo, error) {
	return m.fs.Stat(name)
}

// Sub implements FS
func (m *memory) Sub(dir string) (fs.FS, error) {
	return m.fs.Sub(dir)
}

// Mkdir implements MakeDirFS
func (m *memory) Mkdir(path string, perm fs.FileMode) error {
	fp, err := xfilepath.Parse(path)
	if err != nil {
		return err
	}
	accumulator := fp.Root()

	// check each ancestor path
	for i := 0; i < len(fp.Segments); i++ {
		currentPath := accumulator.String(m.pathSeperator)
		_, ok := m.fs[currentPath]
		if !ok {
			return errNotExist(currentPath)
		}
		seg := fp.Segments[i]
		fpseg, err := xfilepath.Parse(seg)
		if err != nil {
			return err
		}
		accumulator = accumulator.Join(fpseg)
	}

	// write the segment
	m.fs[path] = &fstest.MapFile{
		Mode: perm | fs.ModeDir,
	}

	return nil
}

// MkdirAll implements MakeDirFS
func (m *memory) MkdirAll(path string, perm fs.FileMode) error {
	fp, err := xfilepath.Parse(path)
	if err != nil {
		return err
	}
	accumulator := fp.Root()

	// create each ancestor path
	for i := 0; i < len(fp.Segments); i++ {
		currentPath := accumulator.String(m.pathSeperator)
		_, ok := m.fs[currentPath]
		if !ok {
			m.fs[currentPath] = &fstest.MapFile{
				Mode: perm | fs.ModeDir,
			}
		}
		seg := fp.Segments[i]
		fpseg, err := xfilepath.Parse(seg)
		if err != nil {
			return err
		}
		accumulator = accumulator.Join(fpseg)
	}

	// create the path
	m.fs[path] = &fstest.MapFile{
		Mode: perm | fs.ModeDir,
	}

	return nil
}

func errNotExist(path string) error {
	return fmt.Errorf("'%s' %w", path, fs.ErrNotExist)
}