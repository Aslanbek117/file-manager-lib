package utils

import (
	"encoding/base64"
	e "github.com/Aslanbek117/file-manager-lib/errors"
	m "github.com/Aslanbek117/file-manager-lib/models"
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

//
//oldFileName  and newFileName - full path to file with old/new names
func RenameFile(oldFileName  string, newFileName string)  error {
	return os.Rename(oldFileName, newFileName)
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
		return nil, e.New("Extension was not provided")
	}

	slugStr := getSlug(fileName)

	_, err := os.Stat(dir + "/"+ fileName)
	if !os.IsNotExist(err) {
		return nil, e.New("file or directory exists")
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


func Listing(dir string) *m.ListingInfo {
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

	return &listingInfo
}


func getSlug(fileName string) string {

	slugStr := slug.Make(strings.Split(fileName, ".")[0]) + "." + strings.Split(fileName, ".")[1]
	return slugStr
}




func fileInfoFromInterface(v os.FileInfo, directory string) *m.FileInfo {
	mime, _ := mimetype.DetectFile(v.Name())
	var mimeType = mime.String()

	if v.IsDir() {
		mimeType = "inode/directory"
	}
	return &m.FileInfo{v.Name(), v.Name(), directory, filepath.Join(directory, v.Name()), filepath.Ext(v.Name()), v.Size(), v.ModTime(), v.IsDir(), mimeType}
}

// Node represents a node in a directory tree.
type Node struct {
	Title string   		 `json:"title"`
	FullPath string 	 `json:"path"`
	Key  string    		 `json:"key"`
	Info     *m.FileInfo `json:"info"`
	Children []*Node  	 `json:"children"`
	Parent   *Node    	 `json:"-"`
}

// Create directory hierarchy.
func GetFileTree(root string, onlyDirectories bool) (result *Node, err error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		var key string
		if !info.IsDir() {
			key = ""
		} else {
			key = path
		}
		parents[path] = &Node{
			Title: filepath.Base(path),
			Key:   key,
			FullPath: path,
			Info:     fileInfoFromInterface(info, filepath.Base(path)),
			Children: make([]*Node, 0),
		}
		return nil
	}
	if err = filepath.Walk(absRoot, walkFunc); err != nil {
		return
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			result = node
		} else {
			node.Parent = parent
			if onlyDirectories {
				if node.Info.IsDir {
					parent.Children = append(parent.Children, node)
				}
			} else {
				parent.Children = append(parent.Children, node)
			}

		}
	}
	return
}


