package files

import (
	"encoding/base64"
	"errors"
	"github.com/gosimple/slug"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileInfo struct {
	Name 		 *string
	OriginalName *string
	Dir 		 *string
	Path		 *string
	Extension 	 *string
	Size     	 *int64
	ModTime  	 *time.Time
	Mode     	 *os.FileMode
}


func RewriteIfExists(fileName string, dir string, base64String string) (*FileInfo, error) {
	slugStr, err := getSlug(fileName); if err != nil {
		return nil, err
	}
	isExists, err := fileExists(dir, fileName); if !isExists {
		return nil, errors.New("file does not exists")
	}

	filePath := filepath.Join(dir, filepath.Base(*slugStr))
	fileExtension, _ := getExtension(fileName)

	err = os.Remove(filePath); if err != nil {
		return nil, err
	}

	file, err := os.Create(filePath); if err != nil {
		return nil, err
	}
	defer file.Close()

	dec, err := base64.StdEncoding.DecodeString(base64String)

	file.Write(dec)

	var fileInfo = &FileInfo{
		Name:         slugStr,
		OriginalName: &fileName,
		Dir:          &dir,
		Path:         &filePath,
		Extension:    fileExtension,
	}

	return fileInfo, nil
}

func CreateIfNotExists(fileName string, dir string, base64String string) (*FileInfo,  error) {
	slugStr, err := getSlug(fileName); if err != nil {
		return nil, err
	}


	//isExists, _ := fileExists(dir, fileName); if isExists {
	//
	//	return nil, errors.New("File exists")
	//}


	info, err := os.Stat(dir + "/" + fileName)

	if err != nil {
		return nil, err
	}




	filePath := filepath.Join(dir, filepath.Base(*slugStr))
	fileExtension, _ := getExtension(fileName)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	dec, err := base64.StdEncoding.DecodeString(base64String)

	_, err = file.Write(dec)
	if err != nil {
		return nil, err
	}

	var fileInfo = &FileInfo{
		Name:         slugStr,
		OriginalName: &fileName,
		Dir:          &dir,
		Path:         &filePath,
		Extension:    fileExtension,
		IsDir:
	}
	return fileInfo, nil
}


func fileExists(path string, fileName string) (bool, error) {
	info, err := os.Stat(path + "/" + fileName)
	if os.IsNotExist(err) {
		return  false, err
	}
	return true, nil
}


func getSlug(fileName string) (*string, error) {
	if !strings.Contains(fileName, ".") {
		return nil, errors.New("extension was not provided")
	}
	slugStr := slug.Make(strings.Split(fileName, ".")[0]) + "." + strings.Split(fileName, ".")[1]
	return &slugStr, nil
}

func getExtension(fileName string) (*string, error) {
	if !strings.Contains(fileName, ".") {
		return nil, errors.New("extension was not provided")
	}

	extension := strings.Split(fileName, ".")[1]
	return &extension, nil
}




