package goodidea

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
)

type FileManager interface {
	StoreFile(b []byte, ext string) (string, error)
}

//TODO: https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/endpoints/
type objectStorage struct {
	//The region the bucket is within
	region    string
	accessKey string
	secretKey string
	//The URI for the Object Storage, can be AWS or otherwise
	endpoint string
	//The name of the bucket
	bucket string
	// cdn is an optional property which provides the base URL of the CDN used to serve the images
	// after it has been uploaded
	cdn string
	// The configuration setup to upload files to object storage
	cfg *aws.Config
}

func (obs *objectStorage) getAWSConfig() error {
	var c aws.Config
	var err error
	if obs.accessKey == "" {
		c, err = config.LoadDefaultConfig(context.TODO())
	} else {
		c, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(obs.region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(obs.accessKey, obs.secretKey, "")),
		)
	}
	if err != nil {
		return &ControllerError{Msg: "unable to load default config for object storage", Func: "getAWSConfig", Reason: err.Error()}
	}

	obs.cfg = &c
	return nil
}

type resolverV2 struct {
	URL *url.URL
}

func (r *resolverV2) ResolveEndpoint(ctx context.Context, params s3.EndpointParameters) (
	smithyendpoints.Endpoint, error,
) {
	if r.URL != nil {
		u := *r.URL
		u.Path += "/" + *params.Bucket
		return smithyendpoints.Endpoint{URI: u}, nil
	}

	// fallback to default
	return s3.NewDefaultEndpointResolverV2().ResolveEndpoint(ctx, params)
}

func (obs *objectStorage) getS3Client() (*s3.Client, error) {
	if obs.endpoint == "" {
		client := s3.NewFromConfig(*obs.cfg)
		return client, nil
	}

	endpointURL, err := url.Parse(obs.endpoint)
	if err != nil {
		return nil, &ControllerError{Msg: "unable to parse provided endpoint", Func: "getS3Client", Reason: err.Error()}
	}

	return s3.NewFromConfig(*obs.cfg, func(o *s3.Options) {
		o.EndpointResolverV2 = &resolverV2{URL: endpointURL}
	}), nil
}

//TODO: Create a file name with xid or something or prefix it with the task number idk
//createKeyName create a unique filename for uploading to object storage
func createKeyName(fe string) string {
	return fe
}

func (obs *objectStorage) StoreFile(b []byte, ext string) (string, error) {
	objectKey := createKeyName(ext)

	reader := bytes.NewReader(b)
	params := &s3.PutObjectInput{
		Bucket: &obs.bucket,
		Key:    &objectKey,
		Body:   reader,
	}

	_, err := client.PutObject(context.TODO(), params)
	if err != nil {
		fmt.Printf("Failed to upload object %s/%s, %s\n", obs.bucket, objectKey, err.Error())
		return "", err
	}
	return "tmp", nil
}

type localStorage struct {
	dirName string
}

// StoreFile - store a file locally with the provided file extension
func (ls *localStorage) StoreFile(b []byte, ext string) (string, error) {
	if ls.dirName == "" {
		ls.dirName = os.TempDir()
	}
	f, err := os.CreateTemp(ls.dirName, fmt.Sprintf("idea-*.%s", ext))
	if err != nil {
		Logr.Error("Error Creating Local TMP file", "err", err.Error())
		return "", err
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		Logr.Error("Error writing bytes to a temp file", "err", err.Error())
		return "", err
	}

	s := fmt.Sprintf("/%s", f.Name())
	return s, nil
}

func NewFileManager() FileManager {
	fm := localStorage{
		dirName: "static/img",
	}
	return &fm
}
