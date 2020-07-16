package utils

import (
	"encoding/base64"
	"github.com/gabriel-vasile/mimetype"
	"github.com/gosimple/slug"
	c "github.com/otiai10/copy"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
	m "file-manager/models"
)

func MoveFile(sourceFilePath, moveToDir string) error {
	inputFile, err := os.Open(sourceFilePath)
	if err != nil {
		return err

	}

	outputFile, err := os.Create(filepath.Join(moveToDir + "/" + filepath.Base(inputFile.Name())))
	if err != nil {
		inputFile.Close()
		return err
	}
	defer outputFile.Close()
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return err
	}
	err = RemoveFile(sourceFilePath, "")
	if err != nil {
		return err
	}
	return nil
}





func CopyDir(src string, dest string) error {
	return c.Copy(src, dest)
}

func CopyFile(src string, dest string) error {
	in, err := os.Open(src); if err != nil {
		return err
	}

	defer in.Close()

	out, err := os.Create(dest); if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, in); if err != nil {
		return err
	}

	return out.Close()
}

func CreateDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0777)
	}
	return nil
}

func RemoveDir(dir string) error {
	err := RemoveAllFromDir(dir); if err != nil {
		return err
	}
	err = os.Remove(dir); if err != nil {
		return err
	}
	return nil
}


func RemoveAllFromDir(dir string) error {
	directory, err := ioutil.ReadDir(dir); if err != nil {
		return err
	}
	for _, d := range directory {
		err := os.RemoveAll(path.Join([]string{dir, d.Name()}...)); if err !=nil {
			return err
		}
	}
	return nil
}



func RemoveFile(fileName string, dir string) error {
	err := os.Remove(filepath.Join(dir, fileName)); if err != nil {
		return err
	}
	return nil
}

func CreateFileIfNotExists(fileName string, dir string, base64String string) (*m.FileInfo,  error) {

	if !strings.Contains(fileName, ".") {
		return nil, errors.New("Extension was not provided")
	}

	slugStr := getSlug(fileName)

	_, err := os.Stat(dir + "/"+ fileName)
	if !os.IsNotExist(err) {
		return nil, errors.New("file or directory exists")
	}

	filePath := filepath.Join(dir, filepath.Base(slugStr))
	file, err := os.Create(filePath); if err != nil {
		return nil, err
	}

	defer file.Close()
	dec, err := base64.StdEncoding.DecodeString(base64String)

	_, err = file.Write(dec); if err != nil {
		return nil, err
	}

	var fileInfo = &m.FileInfo{
		Name:         slugStr,
		OriginalName: fileName,
		Dir:          dir,
		Path:         filePath,
		Extension:    filepath.Ext(fileName),
		IsDir:		  false,
		ModTime:      time.Now().Add(6 * time.Hour),

	}
	return fileInfo, nil
}


func Listing(dir string) m.ListingInfo {
	list, _ := ioutil.ReadDir(dir)
	var listingInfo m.ListingInfo
	for _, file := range list {
		mime, _ := mimetype.DetectFile(filepath.Join(dir, file.Name()))
		var mimeType = mime.String()
		if file.IsDir() {
			mimeType = "inode/directory"
		}
		var fileInfo = m.FileInfo{
			Name:         file.Name(),
			OriginalName: "",
			Dir:          dir,
			Path:         filepath.Join(dir, file.Name()),
			Extension:    filepath.Ext(file.Name()),
			Size:         file.Size(),
			ModTime:      file.ModTime(),
			IsDir:        file.IsDir(),
			MimeType:	  mimeType,
		}
		listingInfo.Files = append(listingInfo.Files, fileInfo)
	}

	return listingInfo
}


func getSlug(fileName string) string {

	slugStr := slug.Make(strings.Split(fileName, ".")[0]) + "." + strings.Split(fileName, ".")[1]
	return slugStr
}
