package fs

import (
	"io/fs"
	"os"
	"path/filepath"
	"syscall"
	"time"

	"github.com/spf13/afero"
)

type Matcher interface {
	Match(s string) bool
}

func NewMatcherFs(source afero.Fs, m Matcher) afero.Fs {
	return &MatcherFs{source: source, m: m}
}

type MatherFile struct {
	m Matcher
	f afero.File
}

func (mf *MatherFile) Close() error {
	return mf.f.Close()
}

func (mf *MatherFile) Read(p []byte) (n int, err error) {
	return mf.f.Read(p)
}

func (mf *MatherFile) ReadAt(p []byte, off int64) (n int, err error) {
	return mf.f.ReadAt(p, off)
}

func (mf *MatherFile) Seek(offset int64, whence int) (int64, error) {
	return mf.f.Seek(offset, whence)
}

func (mf *MatherFile) Write(p []byte) (n int, err error) {
	return mf.f.Write(p)
}

func (mf *MatherFile) WriteAt(p []byte, off int64) (n int, err error) {
	return mf.f.WriteAt(p, off)
}

func (mf *MatherFile) Name() string {
	return mf.f.Name()
}

func (mf *MatherFile) Readdir(count int) ([]os.FileInfo, error) {
	rfi, err := mf.f.Readdir(count)
	if err != nil {
		return nil, err
	}
	fi := make([]os.FileInfo, 0)
	for _, i := range rfi {
		if i.IsDir() || mf.m.Match(filepath.Join(mf.f.Name(), i.Name())) {
			fi = append(fi, i)
		}
	}
	return fi, nil
}

func (mf *MatherFile) Readdirnames(count int) ([]string, error) {
	fi, err := mf.Readdir(count)
	if err != nil {
		return nil, err
	}
	n := make([]string, 0)
	for _, s := range fi {
		n = append(n, s.Name())
	}
	return n, nil
}

func (mf *MatherFile) Stat() (os.FileInfo, error) {
	return mf.f.Stat()
}

func (mf *MatherFile) Sync() error {
	return mf.f.Sync()
}

func (mf *MatherFile) Truncate(size int64) error {
	return mf.f.Truncate(size)
}

func (mf *MatherFile) WriteString(s string) (ret int, err error) {
	return mf.f.WriteString(s)
}

type MatcherFs struct {
	m      Matcher
	source afero.Fs
}

func (mfs *MatcherFs) matchesName(name string) error {
	if mfs.m == nil {
		return nil
	}
	if mfs.m.Match(name) {
		return nil
	}
	return syscall.ENOENT
}

func (mfs *MatcherFs) dirOrMatches(name string) error {
	dir, err := afero.IsDir(mfs.source, name)
	if err != nil {
		return err
	}
	if dir {
		return nil
	}
	return mfs.matchesName(name)
}

// Chmod implements afero.Fs
func (mfs *MatcherFs) Chmod(name string, mode fs.FileMode) error {
	if err := mfs.dirOrMatches(name); err != nil {
		return err
	}
	return mfs.source.Chmod(name, mode)
}

// Chown implements afero.Fs
func (mfs *MatcherFs) Chown(name string, uid int, gid int) error {
	if err := mfs.dirOrMatches(name); err != nil {
		return err
	}
	return mfs.source.Chown(name, uid, gid)
}

// Chtimes implements afero.Fs
func (mfs *MatcherFs) Chtimes(name string, atime time.Time, mtime time.Time) error {
	if err := mfs.dirOrMatches(name); err != nil {
		return err
	}
	return mfs.source.Chtimes(name, atime, mtime)
}

// Create implements afero.Fs
func (mfs *MatcherFs) Create(name string) (afero.File, error) {
	if err := mfs.matchesName(name); err != nil {
		return nil, err
	}
	return mfs.source.Create(name)
}

// Mkdir implements afero.Fs
func (mfs *MatcherFs) Mkdir(name string, perm fs.FileMode) error {
	return mfs.source.Mkdir(name, perm)
}

// MkdirAll implements afero.Fs
func (mfs *MatcherFs) MkdirAll(path string, perm fs.FileMode) error {
	return mfs.source.MkdirAll(path, perm)
}

// Name implements afero.Fs
func (*MatcherFs) Name() string {
	return "MatcherFs"
}

// Open implements afero.Fs
func (mfs *MatcherFs) Open(name string) (afero.File, error) {
	dir, err := afero.IsDir(mfs.source, name)
	if err != nil {
		return nil, err
	}
	if !dir {
		if err := mfs.matchesName(name); err != nil {
			return nil, err
		}
	}
	f, err := mfs.source.Open(name)
	if err != nil {
		return nil, err
	}
	return &MatherFile{f: f, m: mfs.m}, nil
}

// OpenFile implements afero.Fs
func (mfs *MatcherFs) OpenFile(name string, flag int, perm fs.FileMode) (afero.File, error) {
	if err := mfs.dirOrMatches(name); err != nil {
		return nil, err
	}
	return mfs.source.OpenFile(name, flag, perm)
}

// Remove implements afero.Fs
func (mfs *MatcherFs) Remove(name string) error {
	if err := mfs.matchesName(name); err != nil {
		return err
	}
	return mfs.source.Remove(name)
}

// RemoveAll implements afero.Fs
func (mfs *MatcherFs) RemoveAll(path string) error {
	if err := mfs.dirOrMatches(path); err != nil {
		return err
	}
	return mfs.source.RemoveAll(path)
}

// Rename implements afero.Fs
func (mfs *MatcherFs) Rename(oldname string, newname string) error {
	dir, err := afero.IsDir(mfs.source, oldname)
	if err != nil {
		return err
	}
	if dir {
		return nil
	}
	if err := mfs.matchesName(oldname); err != nil {
		return err
	}
	if err := mfs.matchesName(newname); err != nil {
		return err
	}
	return mfs.source.Rename(oldname, newname)
}

// Stat implements afero.Fs
func (mfs *MatcherFs) Stat(name string) (fs.FileInfo, error) {
	if err := mfs.dirOrMatches(name); err != nil {
		return nil, err
	}
	return mfs.source.Stat(name)
}
