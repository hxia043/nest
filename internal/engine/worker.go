package engine

import (
	"runtime"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
)

const (
	numberOfJobs int = 100
)

var jobs chan job
var results chan result

type result struct {
	err        error
	imageName  string
	imageExist bool
}

type job struct {
	ref     name.Reference
	image   v1.Image
	options []remote.Option
}

func verifyImageExist(localImage v1.Image, remoteImage v1.Image) bool {
	if localImage == nil || remoteImage == nil {
		return false
	}

	localCfgHash, err := localImage.ConfigName()
	if err != nil {
		return false
	}

	remoteCfgHash, err := remoteImage.ConfigName()
	if err != nil {
		return false
	}

	if localCfgHash.Algorithm != remoteCfgHash.Algorithm || localCfgHash.Hex != remoteCfgHash.Hex {
		return false
	}

	return true
}

func worker(id int, jobs <-chan job, results chan<- result) {
	for job := range jobs {
		remoteImg, _ := remote.Image(job.ref, job.options...)
		if verifyImageExist(job.image, remoteImg) {
			results <- result{nil, job.ref.Name(), true}
		} else {
			err := remote.Write(job.ref, job.image, job.options...)
			if err != nil {
				results <- result{err, job.ref.Name(), false}
			}

			results <- result{nil, job.ref.Name(), false}
		}
	}
}

func startWorkers(jobs chan job, results chan result) {
	numberOfWorkers := runtime.GOMAXPROCS(runtime.NumCPU())
	for i := 0; i < numberOfWorkers; i++ {
		go worker(i, jobs, results)
	}
}
