package factory

import (
	"bytes"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"
)

func BuildAWSCSVSaver(
	csvOutputKey string,
	outputBucket string,
	svc *s3.S3,
) func(s string) {
	stringSaver := func(s string) {
		body := bytes.NewReader([]byte(s))
		// csvOutputKey := fmt.Sprintf(outputNamePattern, companyKey, fileName)
		_, err := svc.PutObject(
			&s3.PutObjectInput{
				Bucket: &outputBucket,
				Key:    &csvOutputKey,
				Body:   body,
			},
		)

		if err != nil {
			glog.Fatalf("error: could not write the file to s3 (%s)\n", err)
		}
	}

	return stringSaver
}
