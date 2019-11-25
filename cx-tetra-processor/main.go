package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/arizard/tetra-processor/factory"
	"github.com/arizard/tetra-processor/usecase"

	"github.com/golang/glog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.S3Event) {
	configBucket := os.Getenv("CONFIG_BUCKET")
	tetraConfigKey := os.Getenv("TETRA_CONFIG_KEY")
	// tesseractConfigKey := os.Getenv("TESSERACT_CONFIG_KEY")
	outputBucket := os.Getenv("OUTPUT_BUCKET")
	outputNamePattern := os.Getenv("OUTPUT_NAME_PATTERN")

	for _, record := range request.Records {
		if record.EventName != "ObjectCreated:Put" {
			continue
		}

		sess, err := session.NewSession(
			&aws.Config{
				Region: aws.String("us-west-2"),
			},
		)

		if err != nil {
			glog.Fatalf("error: could not create new session (%s)\n", err)
		}

		svc := s3.New(sess)

		tetraBucket := record.S3.Bucket.Name
		keySplit := strings.Split(record.S3.Object.Key, "/")
		companyKey, fileName := keySplit[0], keySplit[1]
		csvFileKey := record.S3.Object.Key

		csvOutputKey := fmt.Sprintf(outputNamePattern, companyKey, fileName)
		tetraConfigPath := fmt.Sprintf("%s/%s", companyKey, tetraConfigKey)

		stringCSVLoader := factory.BuildAWSCSVLoader(
			tetraBucket,
			csvFileKey,
			svc,
		)

		stringCSVSaver := factory.BuildAWSCSVSaver(
			csvOutputKey,
			outputBucket,
			svc,
		)

		tetraConfigLoader := factory.BuildAWSTetraConfigGetter(
			tetraConfigPath,
			configBucket,
			svc,
		)

		usecase.TransformCSV(tetraConfigLoader, stringCSVLoader, stringCSVSaver)

	}
}

func main() {
	flag.Parse()
	flag.Set("logtostderr", "1")
	glog.Infof("cx-tetra-processor init\n")
	lambda.Start(Handler)
}
