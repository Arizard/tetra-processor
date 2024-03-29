package factory

import (
	"bytes"
	"io"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/golang/glog"
)

func BuildAWSCSVLoader(
	tetraBucket string,
	csvFileKey string,
	svc *s3.S3,
	decrypt func(io.Reader) (io.Reader, error),
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
		csvDecrypted, err := decrypt(csvFileOutput.Body)

		if err != nil {
			glog.Fatalf(
				"error: could not decrypt the file (%s)\n",
				err,
			)
		}
		_, err = csvBuf.ReadFrom(csvDecrypted)
		if err != nil {
			glog.Fatalf("error: buffer could not read decrypted (%s)", err)
		}
		csvString := csvBuf.String()

		return csvString
	}

	return stringCSVLoader
}
