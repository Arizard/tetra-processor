package resource

import (
	"bytes"
	"encoding/json"
	"github.com/arizard/tetra-processor/types"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"
)

func GetAWSTetraMetaConfig(
	configBucket string,
	tetraProcessorConfigPath string,
	svc *s3.S3,
) types.MetaConfig {

	tetraProcessorConfigOutput, err := svc.GetObject(
		&s3.GetObjectInput{
			Bucket: &configBucket,
			Key:    &tetraProcessorConfigPath,
		},
	)

	if err != nil {
		glog.Fatalf("error: could not get the file from s3 (%s)\n", err)
	}

	cfgBuf := new(bytes.Buffer)
	if _, err := cfgBuf.ReadFrom(tetraProcessorConfigOutput.Body); err != nil {
		glog.Fatalf(
			"error: failed to read bytes from s3 file body (%s)\n",
			err,
		)
	}

	var tetraMetaConfig types.MetaConfig

	if err = json.Unmarshal(cfgBuf.Bytes(), &tetraMetaConfig); err != nil {
		glog.Fatalf(
			"error: could not unmarshal json data (%s)\n",
			err,
		)
	}
	return tetraMetaConfig
}
