package engine

import (
	"context"

	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

func (im *imageManager) Login() error {
	options, err := im.createOptions()
	if err != nil {
		return err
	}

	registry, err := name.NewRegistry(im.registryOption.registry)
	if err != nil {
		return err
	}

	if _, err = remote.Catalog(context.Background(), registry, options...); err != nil {
		return err
	}

	return nil
}
