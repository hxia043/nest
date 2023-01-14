package engine

import "sync"

type Options map[string]string

// images is a map which to mapping the local image path to image reference
// for example key="/home/docker-pause-0.1.0.tar", value="docker-pause:0.1.0"
type imageManager struct {
	images sync.Map
	registryOption
}

type registryOption struct {
	registry  string
	username  string
	password  string
	caCert    string
	namespace string
}

func NewImageManager(options Options, images sync.Map) *imageManager {
	o := registryOption{}
	for option, value := range options {
		switch option {
		case "username":
			o.username = value
		case "password":
			o.password = value
		case "ca-cert":
			o.caCert = value
		case "namespace":
			o.namespace = value
		case "registry":
			o.registry = value
		}
	}

	return &imageManager{images: images, registryOption: o}
}
