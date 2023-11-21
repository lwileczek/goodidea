package goodidea

import (
	"fmt"
	"os"
)

type FileManager interface {
	StoreFile(b []byte, ext string) error
}

type localStorage struct {
	dirName string
}

type objectStorage struct {
	//The region the bucket is within
	region    string
	accessKey string
	secretKey string
	//The URI for the Object Storage, can be AWS or otherwise
	endpoint string
	//The name of the bucket
	bucket string
}

// StoreFile - store a file locally with the provided file extension
func (ls *localStorage) StoreFile(b []byte, ext string) error {
	if ls.dirName == "" {
		ls.dirName = os.TempDir()
	}
	f, err := os.CreateTemp(ls.dirName, fmt.Sprintf("idea-*.%s", ext))
	if err != nil {
		Logr.Error("Error Creating Local TMP file", "err", err.Error())
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		Logr.Error("Error writing bytes to a temp file", "err", err.Error())
		return err
	}

	return nil
}

func NewFileManager() FileManager {
	fm := localStorage{
		dirName: "tmp",
	}
	return &fm
}
