package db

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type FileOption struct {
	DbPath       string
	registryPath string
	file         *os.File
	Registry
}

func (o *FileOption) checkDbExistFor(name ...string) error {
	switch len(name) {
	case 0:
		f, err := os.Stat(o.DbPath)
		if !(err == nil && f.IsDir()) {
			return fmt.Errorf("database %s not found", o.DbPath)
		}
	case 1:
		o.registryPath = filepath.Join(o.DbPath, name[0])
		f, err := os.Stat(o.registryPath)
		if !(err == nil && !f.IsDir()) {
			return fmt.Errorf("registry %s not found", name)
		}
	default:
		return fmt.Errorf("unexpected name %v found", name)
	}

	return nil
}

func (o *FileOption) initRegistry(name string) error {
	var err error

	registryPath := filepath.Join(o.DbPath, name)
	if f, err := os.Stat(registryPath); err == nil {
		if !f.IsDir() {
			return fmt.Errorf("registry %s exists", name)
		}
	}

	o.file, err = os.Create(registryPath)
	if err != nil {
		return fmt.Errorf("create registry %s failed", name)
	}

	return nil
}

func (o *FileOption) Add(option []byte) error {
	if err := json.Unmarshal(option, &o.Registry); err != nil {
		return err
	}

	if err := o.initRegistry(o.Name); err != nil {
		return err
	}

	data, err := json.MarshalIndent(o.Registry, "", "	")
	if err != nil {
		return err
	}

	if _, err := o.file.Write(data); err != nil {
		return err
	}

	return nil
}

func (o *FileOption) Delete(name string) error {
	if err := o.checkDbExistFor(name); err != nil {
		return err
	}

	if err := os.Remove(o.registryPath); err != nil {
		return fmt.Errorf("remove registry %s failed: %v", name, err)
	}

	return nil
}

func (o *FileOption) Update(option []byte) error {
	if err := json.Unmarshal(option, &o.Registry); err != nil {
		return err
	}

	if err := o.Delete(o.Name); err != nil {
		return err
	}

	if err := o.Add(option); err != nil {
		return err
	}

	return nil
}

func (o *FileOption) Get(name string) (*Registry, error) {
	if err := o.checkDbExistFor(name); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(o.registryPath)
	if err != nil {
		return nil, err
	}

	var dbregistry Registry
	if err := json.Unmarshal(data, &dbregistry); err != nil {
		return nil, err
	}

	return &dbregistry, nil
}

func (o *FileOption) Show() ([]Registry, error) {
	if err := o.checkDbExistFor(); err != nil {
		return nil, err
	}

	registries, err := os.ReadDir(o.DbPath)
	if err != nil {
		return nil, err
	}

	var outputs []Registry
	for _, registry := range registries {
		var dbregistry Registry
		if !registry.IsDir() {
			registryPath := filepath.Join(o.DbPath, registry.Name())
			data, err := os.ReadFile(registryPath)
			if err != nil {
				return nil, err
			}

			if err := json.Unmarshal(data, &dbregistry); err != nil {
				return nil, err
			}

			outputs = append(outputs, dbregistry)
		}
	}

	return outputs, nil
}

func NewFileOption(dbpath string) *FileOption {
	return &FileOption{DbPath: dbpath}
}
