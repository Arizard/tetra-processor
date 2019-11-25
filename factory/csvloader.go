package factory

import (
	"bytes"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"
)

func BuildAWSCSVLoader(
	tetraBucket string,
	csvFileKey string,
	svc *s3.S3,
) func() string {
	stringCSVLoader := func() string {

		csvFileOutput, err := svc.GetObject(
			&s3.GetObjectInput{
				Bucket: &tetraBucket,
				Key:    &csvFileKey,
			},
		)

		if err != nil {
			glog.Fatalf("error: could not get the file from s3 (%s)\n", err)
		}

		csvBuf := new(bytes.Buffer)
		csvBuf.ReadFrom(csvFileOutput.Body)
		csvString := csvBuf.String()

		return csvString
	}

	return stringCSVLoader
}
