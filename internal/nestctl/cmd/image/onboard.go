package image

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/hxia043/nest/internal/engine"
	"github.com/hxia043/nest/internal/nestctl/db"
	"github.com/hxia043/nest/pkg/util/filesys"
	"github.com/hxia043/nest/types/exception"
)

type ImageOnboardOption struct {
	Name      string
	ImagePath string
	Namespace string
	imagesOption
}

type imagesOption struct {
	wg               sync.WaitGroup
	lock             sync.Mutex
	imagesMap        sync.Map
	invalidImageName []string
}

func isNumber(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func (o *ImageOnboardOption) validateOnboardOption() error {
	var err error
	if strings.HasPrefix(o.ImagePath, ".") {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}

		o.ImagePath = strings.Replace(o.ImagePath, ".", wd, -1)
	} else {
		o.ImagePath, err = os.Getwd()
		if err != nil {
			return err
		}
	}

	if o.Namespace == "" {
		return fmt.Errorf("empty namespace")
	}

	return nil
}

func (o *ImageOnboardOption) parseImages() ([]string, error) {
	var images []string

	if filesys.IsDirectory(o.ImagePath) {
		files, _ := os.ReadDir(o.ImagePath)
		for _, file := range files {
			if file.IsDir() || !strings.HasSuffix(file.Name(), ".tar") {
				continue
			}
			images = append(images, file.Name())
		}
	} else {
		return nil, fmt.Errorf("%s is not a local directory of images", o.ImagePath)
	}

	return images, nil
}

func (o *ImageOnboardOption) parseImagesMapFrom(images []string) (sync.Map, error) {
	for index := range images {
		o.wg.Add(1)
		go func(i int) {
			defer o.wg.Done()
			imageFileNameList := strings.Split(images[i], "-")
			if len(imageFileNameList) < 2 {
				o.lock.Lock()
				o.invalidImageName = append(o.invalidImageName, images[i])
				o.lock.Unlock()
				return
			}

			namePool, tagPool, tagBegin := make([]string, 0), make([]string, 0), false
			for j := range imageFileNameList {
				if tagBegin {
					tagPool = append(tagPool, imageFileNameList[j])
					continue
				}

				if isNumber(string(imageFileNameList[j][0])) {
					tagBegin = true
					tagPool = append(tagPool, imageFileNameList[j])
					continue
				}

				namePool = append(namePool, imageFileNameList[j])
			}

			tag := fmt.Sprintf("%s:%s", strings.Join(namePool, "-"), strings.Join(tagPool, "-"))
			o.imagesMap.Store(filepath.Join(o.ImagePath, images[i]), strings.Replace(tag, ".tar", "", -1))
		}(index)
	}

	o.wg.Wait()
	if len(o.imagesOption.invalidImageName) > 0 {
		return o.imagesMap, exception.NewInvalidImageNameError(o.imagesOption.invalidImageName...)
	}

	return o.imagesMap, nil
}

func (o *ImageOnboardOption) Onboard(args []string) error {
	var err error
	if err = o.validateOnboardOption(); err != nil {
		return err
	}

	o.Name, err = getRegistryName(args)
	if err != nil {
		return err
	}

	dbOption, err := db.InitDB()
	if err != nil {
		return err
	}

	registry, err := dbOption.Get(o.Name)
	if err != nil {
		return fmt.Errorf("no registry %s found", o.Name)
	}

	images, err := o.parseImages()
	if err != nil {
		return err
	}

	imagesMap, err := o.parseImagesMapFrom(images)
	if err != nil {
		return err
	}

	options := NewEngineOption(o.Namespace, registry)
	im := engine.NewImageManager(options, imagesMap)
	if err := im.Login(); err != nil {
		return err
	}

	return im.Onboard()
}
