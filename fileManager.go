package goodidea

import (
	"bytes"
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithyendpoints "github.com/aws/smithy-go/endpoints"
	"github.com/rs/xid"
)

type FileManager interface {
	StoreFile(b []byte, ext string) (string, error)
}

// TODO: https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/endpoints/
type objectStorage struct {
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

func createKeyName(fe string) string {
	x := xid.New()
	s := fmt.Sprintf("goodidea/task-img-%s.%s", x.String(), fe)
	return s
}

func (obs *objectStorage) constructlURL(key string) string {
	if obs.cdn != "" {
		return fmt.Sprintf("%s/%s", obs.cdn, key)
	}
	u := "s3.amazonaws.com"
	if obs.endpoint != "" {
		u = obs.endpoint
	}
	return fmt.Sprintf("https://%s.%s/%s", obs.bucket, u, key)
}

func (obs *objectStorage) StoreFile(b []byte, ext string) (string, error) {
	if obs.cfg == nil {
		c, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			return "", &ControllerError{Msg: "unable to load default config for object storage", Func: "ObjectStoreFile", Reason: err.Error()}
		}
		obs.cfg = &c
	}
	objectKey := createKeyName(ext)

	reader := bytes.NewReader(b)
	params := &s3.PutObjectInput{
		Bucket: &obs.bucket,
		Key:    &objectKey,
		Body:   reader,
	}

	client, err := obs.getS3Client()
	if err != nil {
		return "", err
	}
	_, err = client.PutObject(context.TODO(), params)
	if err != nil {
		fmt.Printf("Failed to upload object %s/%s, %s\n", obs.bucket, objectKey, err.Error())
		return "", err
	}

	imgURL := obs.constructlURL(objectKey)
	return imgURL, nil
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
	var fm FileManager
	if os.Getenv("AWS_BUCKET") != "" {
		fm = &objectStorage{
			bucket:   os.Getenv("AWS_BUCKET"),
			cdn:      os.Getenv("AWS_IMAGE_CDN"),
			endpoint: os.Getenv("AWS_ENDPOINT"),
		}
	} else {
		//TODO: Check if directory exists and if not try to create it
		d := "static/img"
		if os.Getenv("LOCAL_DIR") != "" {
			d = os.Getenv("LOCAL_DIR")
		}
		fm = &localStorage{
			dirName: d,
		}
	}
	return fm
}
