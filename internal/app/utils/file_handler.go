package utils

import (
	"bytes"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)

type UploadMeta struct {
	Name   string
	Buffer bytes.Buffer
}

var (
	// Default Variable
	// srv retro.Service

	// Static Variable
	StoragePath = "./internal/storage/"
)

func GenerateFilePath(folderNameDetail ...string) (result string) {
	time := time.Now()
	result = StoragePath + time.Format("02-01-2006") + "/"
	if len(folderNameDetail) > 0 {
		for _, value := range folderNameDetail {
			result += value + "/"
		}
	}
	return
}

func FileSaver(file *UploadMeta, pathDetail ...string) (fileName, path string, err error) {
	path = GenerateFilePath(pathDetail...)
	count := 0
	if _, err = os.Stat(StoragePath); os.IsNotExist(err) {
		if err = os.MkdirAll(StoragePath, 0777); err != nil {
			zap.S().Error(err)
			return
		}
	}
	if _, err = os.Stat(path); os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0777); err != nil {
			zap.S().Error(err)
			return
		}
	}
	name := strings.Split(file.Name, ".")
	for _, err := os.Stat(path + file.Name); !os.IsNotExist(err); _, err = os.Stat(path + file.Name) {
		name = strings.Split(file.Name, ".")
		name = append(name, strconv.Itoa(count))
		name[len(name)-1], name[len(name)-2] = name[len(name)-2], name[len(name)-1]
		file.Name = strings.Join(name, ".")

		count++
	}

	if err = ioutil.WriteFile(path+file.Name, file.Buffer.Bytes(), 0644); err != nil {
		zap.S().Error(err)
		return
	}
	fileName = file.Name

	return
}
