package file

import (
	"os"
	"path"
)

type FileOption struct {
	DbPath string
}

func (o *FileOption) createDBPath() (bool, error) {
	if f, err := os.Stat(o.DbPath); err == nil {
		if f.IsDir() {
			return true, nil
		}
	}

	if err := os.MkdirAll(o.DbPath, os.ModePerm); err != nil {
		return false, err
	}

	return true, nil
}

func (o *FileOption) InitDB() error {
	isCreated, err := o.createDBPath()
	if !isCreated {
		return err
	}

	return nil
}

func NewFileOption(workspace, name string) *FileOption {
	return &FileOption{
		DbPath: path.Join(workspace, name),
	}
}
