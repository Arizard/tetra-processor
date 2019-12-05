package resource

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"

	"github.com/arizard/tetra-processor/types"
)

func GetAWSTesseractConfig(
	configBucket string,
	tesseractConfigPath string,
	svc *s3.S3,
) types.TesseractConfig {

	tesseractConfigOutput, err := svc.GetObject(
		&s3.GetObjectInput{
			Bucket: &configBucket,
			Key:    &tesseractConfigPath,
		},
	)

	if err != nil {
		glog.Fatalf("error: could not get the file from s3 (%s, %s, %s)\n", configBucket, tesseractConfigPath, err)
	}

	cfgBuf := new(bytes.Buffer)
	if _, err := cfgBuf.ReadFrom(tesseractConfigOutput.Body); err != nil {
		glog.Fatalf(
			"error: failed to read bytes from s3 file body (%s)\n",
			err,
		)
	}

	var tesseractConfig types.TesseractConfig

	if err = json.Unmarshal(cfgBuf.Bytes(), &tesseractConfig); err != nil {
		glog.Fatalf(
			"error: could not unmarshal json data (%s)\n",
			err,
		)
	}
	return tesseractConfig
}
