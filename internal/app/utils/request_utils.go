package utils

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"go.uber.org/zap"
)

func FileUploadReceiver(r *http.Request) (result []*UploadMeta, err error) {
	// GET FILE FROM JSON
	if err = r.ParseMultipartForm(5242880); err != nil {
		zap.S().Error(err)
		return
	}

	formData := r.MultipartForm
	files := formData.File["file[]"]
	result = []*UploadMeta{}

	// CREATE NEW FOLDER TO DOWNLOAD DOCUMENT
	// if _, err = os.Stat(storagePath); os.IsNotExist(err) {
	// 	if err = os.MkdirAll(storagePath, 0777); err != nil {
	// 		return
	// 	}
	// }
	// if _, err = os.Stat(GenerateFilePath()); os.IsNotExist(err) {
	// 	if err = os.MkdirAll(GenerateFilePath(), 0777); err != nil {
	// 		return
	// 	}
	// }
	for index := range files {
		file, errTemp := files[index].Open()
		defer file.Close()

		var buf bytes.Buffer

		if errTemp != nil {
			err = errTemp
			zap.S().Error(err)
			return
		}
		defer file.Close()

		finalName := files[index].Filename

		// count := 0
		// name := strings.Split(files[index].Filename, ".")

		// if name[len(name)-1] != "xlsx" && name[len(name)-1] != "xls" {
		// 	err = errors.New("File Format is not Supported. File extension must be .xlsx or .xls")
		// 	return
		// }

		// for _, err := os.Stat(GenerateFilePath() + finalName); !os.IsNotExist(err); _, err = os.Stat(GenerateFilePath() + finalName) {
		// 	name = strings.Split(files[index].Filename, ".")
		// 	name = append(name, strconv.Itoa(count))
		// 	name[len(name)-1], name[len(name)-2] = name[len(name)-2], name[len(name)-1]
		// 	finalName = strings.Join(name, ".")

		// 	count++
		// }

		io.Copy(&buf, file)
		if bytes.NewReader(buf.Bytes()).Size() > 10485760 { //10MB
			err = errors.New("File Size Must Be 10 MB Max")
			zap.S().Error(err)
			return
		}
		// if err = ioutil.WriteFile(GenerateFilePath()+finalName, buf.Bytes(), 0644); err != nil {
		// 	return
		// }
		result = append(result, &UploadMeta{
			Buffer: buf,
			Name:   finalName,
		})
		buf.Reset()
	}
	return
}

func ExcelFileReceiver(r *http.Request) (result []*excelize.File, err error) {
	result = []*excelize.File{}
	uploadFiles, err := FileUploadReceiver(r)
	if err != nil {
		zap.S().Error(err)
		return
	}

	// CREATE NEW FOLDER TO DOWNLOAD DOCUMENT
	if _, err = os.Stat(StoragePath); os.IsNotExist(err) {
		if err = os.MkdirAll(StoragePath, 0777); err != nil {
			zap.S().Error(err)
			return
		}
	}
	if _, err = os.Stat(GenerateFilePath()); os.IsNotExist(err) {
		if err = os.MkdirAll(GenerateFilePath(), 0777); err != nil {
			zap.S().Error(err)
			return
		}
	}

	for _, file := range uploadFiles {
		count := 0
		name := strings.Split(file.Name, ".")

		if name[len(name)-1] != "xlsx" && name[len(name)-1] != "xls" {
			err = errors.New("File Format is not Supported. File extension must be .xlsx or .xls")
			zap.S().Error(err)
			return
		}

		for _, err := os.Stat(GenerateFilePath() + file.Name); !os.IsNotExist(err); _, err = os.Stat(GenerateFilePath() + file.Name) {
			name = strings.Split(file.Name, ".")
			name = append(name, strconv.Itoa(count))
			name[len(name)-1], name[len(name)-2] = name[len(name)-2], name[len(name)-1]
			file.Name = strings.Join(name, ".")

			count++
		}

		if err = ioutil.WriteFile(GenerateFilePath()+file.Name, file.Buffer.Bytes(), 0644); err != nil {
			zap.S().Error(err)
			return
		}

		excelFile, errExcel := excelize.OpenFile(GenerateFilePath() + file.Name)
		if errExcel != nil {
			err = errExcel
			zap.S().Error(err)
			return
		}
		result = append(result, excelFile)
	}

	return
}

func RequestFormReceiver(r *http.Request, keys ...string) (result map[string]string) {
	result = map[string]string{}
	for _, key := range keys {
		result[key] = r.FormValue(key)
	}
	return
}
