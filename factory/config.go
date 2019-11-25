package factory

import (
	"bytes"

	"github.com/arizard/tetra"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"
)

// NewAWSTetraConfigGetter creates a new config getter function which depends
// on AWS.
func BuildAWSTetraConfigGetter(
	tetraConfigPath string,
	configBucket string,
	svc *s3.S3,
) func() tetra.Config {

	tetraConfigLoader := func() tetra.Config {
		tetraConfigOutput, err := svc.GetObject(
			&s3.GetObjectInput{
				Bucket: &configBucket,
				Key:    &tetraConfigPath,
			},
		)

		if err != nil {
			glog.Fatalf("error: could not get the file from s3 (%s)\n", err)
		}

		cfgBuf := new(bytes.Buffer)
		cfgBuf.ReadFrom(tetraConfigOutput.Body)

		var tetraConfig tetra.Config

		err = tetraConfig.LoadFromJSON(cfgBuf.Bytes())

		if err != nil {
			glog.Fatalf("error: could not unmarshal file (%s)\n", err)
		}

		glog.Infof("%v", tetraConfig)

		return tetraConfig
	}

	return tetraConfigLoader
}

// NewTestTetraConfigGetter creates a new config getter function which can be
// used for tests.
func BuildTestTetraConfigGetter(cfg tetra.Config) func() tetra.Config {
	getter := func() tetra.Config {
		return cfg
	}
	return getter
}
