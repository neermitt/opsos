package fs

import (
	"io/fs"

	"github.com/spf13/afero"
)

func AllFiles(afs afero.Fs) ([]string, error) {
	match := make([]string, 0)
	err := afero.Walk(afs, "", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			match = append(match, path)
		}

		return nil
	})
	return match, err
}
