package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/arizard/tetra-processor/resource"

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
func Handler(_ context.Context, request events.S3Event) {
	configBucket := os.Getenv("CONFIG_BUCKET")
	tesseractConfigKey := os.Getenv("TESSERACT_CONFIG_KEY")
	tetraConfigKey := os.Getenv("TETRA_CONFIG_KEY")
	tetraProcessorConfigKey := os.Getenv("TETRA_PROCESSOR_CONFIG_KEY")
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

		tetraProcessorConfigPath := fmt.Sprintf(
			"%s/%s",
			companyKey,
			tetraProcessorConfigKey,
		)

		metaCfg := resource.GetAWSTetraMetaConfig(
			configBucket,
			tetraProcessorConfigPath,
			svc,
		)

		destFileName := fmt.Sprintf(metaCfg.DestFileNamePattern, fileName)
		csvOutputKey := fmt.Sprintf(outputNamePattern, companyKey, destFileName)
		tetraConfigPath := fmt.Sprintf("%s/%s", companyKey, tetraConfigKey)
		tesseractConfigPath := fmt.Sprintf(
			"%s/%s",
			companyKey,
			tesseractConfigKey,
		)

		tesseractConfig := resource.GetAWSTesseractConfig(
			configBucket,
			tesseractConfigPath,
			svc,
		)

		decryptFunc := factory.BuildNullDecrypt()
		glog.Infof("config: %+v\n", tesseractConfig)
		if tesseractConfig.Inbound.Action == "decrypt" {
			if tesseractConfig.Inbound.Encryption == "" {
				glog.Fatal("error: 'encryption' not provided in config")
			}
			if tesseractConfig.Inbound.PrivateKeyFileName == "" {
				glog.Fatal(
					"error: 'privateKeyFileName' not provided in config",
				)
			}
			privateKeyPath := fmt.Sprintf(
				"%s/%s",
				companyKey,
				tesseractConfig.Inbound.PrivateKeyFileName,
			)
			privateKeyOutput, err := svc.GetObject(
				&s3.GetObjectInput{
					Bucket: &configBucket,
					Key:    &privateKeyPath,
				},
			)
			if err != nil {
				glog.Fatalf("error: could not retrieve private key from s3 (%s)", err)
			}
			decryptFunc = factory.BuildPGPDecrypt(
				privateKeyOutput.Body,
				[]byte(tesseractConfig.Inbound.Passphrase),
			)
			csvOutputKey = fmt.Sprintf(
				outputNamePattern,
				companyKey,
				fmt.Sprintf(
					"decrypted_%s",
					strings.TrimRight(destFileName, ".gpg"),
				),
			)

		}

		stringCSVLoader := factory.BuildAWSCSVLoader(
			tetraBucket,
			csvFileKey,
			svc,
			decryptFunc,
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
	if err := flag.Set("logtostderr", "1"); err != nil {
		glog.Fatalf(
			"error: could not force flag '-logtostderr' (%s)",
			err,
		)
	}
	lambda.Start(Handler)
}
