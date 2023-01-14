package engine

import (
	"errors"
	"fmt"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
)

func createImageReference(addr, namespace, image string) string {
	var reference string = ""
	if namespace != "" {
		reference = fmt.Sprintf("%s/%s/%s", addr, namespace, image)
	} else {
		reference = fmt.Sprintf("%s/%s", addr, image)
	}

	return reference
}

func (im *imageManager) validateLocalImageReference() (map[string]name.Reference, error) {
	var validateLocalImageReferenceOk = true
	var localImageReferences = make(map[string]name.Reference)

	im.images.Range(func(localImage, image interface{}) bool {
		imageReference := createImageReference(im.registryOption.registry, im.registryOption.namespace, image.(string))
		ref, err := name.ParseReference(imageReference, name.StrictValidation)
		if err != nil {
			validateLocalImageReferenceOk = false
			fmt.Printf("Error: image reference %s parse failed with error %v.\n", imageReference, err.Error())
		}

		localImageReferences[localImage.(string)] = ref
		return true
	})

	if !validateLocalImageReferenceOk {
		return nil, errors.New("image reference parse failed")
	}

	return localImageReferences, nil
}

func (im *imageManager) validateImages() (map[v1.Image]name.Reference, error) {
	localImageReferences, err := im.validateLocalImageReference()
	if err != nil {
		return nil, err
	}

	var imageReferences = make(map[v1.Image]name.Reference)
	var validateLocalImageOk = true
	for localImage, ref := range localImageReferences {
		image, err := tarball.ImageFromPath(localImage, nil)
		if err != nil {
			validateLocalImageOk = false
			fmt.Printf("Error: load image from path %s failed with error %v.\n", localImage, err.Error())
			continue
		}

		imageReferences[image] = ref
	}

	if !validateLocalImageOk {
		return nil, errors.New("load image failed")
	}

	return imageReferences, nil
}

func (im *imageManager) Onboard() error {
	imageReferences, err := im.validateImages()
	if err != nil {
		return err
	}

	options, err := im.createOptions()
	if err != nil {
		return err
	}

	jobs = make(chan job, numberOfJobs)
	results = make(chan result, len(imageReferences))
	startWorkers(jobs, results)

	for image, ref := range imageReferences {
		jobs <- job{
			ref:     ref,
			image:   image,
			options: options,
		}
	}

	isOnboardSuccessed := true
	numberOfReusult := 0
	for result := range results {
		numberOfReusult++
		if result.err != nil {
			isOnboardSuccessed = false
			fmt.Printf("Error: onboard image %s failed with error %v.\n", result.imageName, result.err.Error())
		} else {
			if result.imageExist {
				fmt.Printf("Info: image %v already exist.\n", result.imageName)
			} else {
				fmt.Printf("Info: onboard image (%d/%d)\n", numberOfReusult, len(imageReferences))
			}
		}

		if numberOfReusult == len(imageReferences) {
			close(results)
		}
	}

	close(jobs)

	if !isOnboardSuccessed {
		return errors.New("onboard image failed")
	}

	return nil
}
