package usecase

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/rs/xid"
)

// FileUploadUC ...
type FileUploadUC struct {
	*ContractUC
}

// CreateFolder ...
func (uc FileUploadUC) CreateFolder(name string) (err error) {
	path := uc.ContractUC.EnvConfig["FILE_STATIC_FILE"] + name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
	}

	return err
}

// Upload ...
func (uc FileUploadUC) Upload(fileName, Type string, file []byte) (fullFileName string, err error) {
	folder := "/" + Type
	err = uc.CreateFolder(folder)
	if err != nil {
		return "", errors.New("Invalid File Type")
	}

	uploadPath := uc.ContractUC.EnvConfig["FILE_STATIC_FILE"]
	filetype := http.DetectContentType(file)

	if filetype != "image/jpeg" && filetype != "image/jpg" && filetype != "image/gif" && filetype != "image/png" && filetype != "application/pdf" {
		return "", errors.New("invalid_file_type")
	}

	id := xid.New().String()
	fullFileName = folder + id + "_" + strings.ReplaceAll(fileName, " ", "-")
	uploadURL := uploadPath + fullFileName

	err = ioutil.WriteFile(uploadURL, file, 0644)
	fmt.Println(err)

	return fullFileName, err
}
