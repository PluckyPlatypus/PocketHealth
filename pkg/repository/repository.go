package repository

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/suyashkumar/dicom"
)

type Repository interface {
	CreateFile(newDicom dicom.Dataset, key string) error
	ReadFile(key string) (dicom.Dataset, error)
}

type DicomRepository struct {
	databaseName string
}

func NewRepository(dbName string) Repository {
	return &DicomRepository{
		databaseName: dbName,
	}
}

func (r *DicomRepository) CreateFile(newDicom dicom.Dataset, key string) error {
	filePath := fmt.Sprintf("%s/%s", r.databaseName, key)
	if checkFileExists(filePath) {
		return errors.New("Dicom file already exists")
	}

	file, err := create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	return dicom.Write(file, newDicom)
}

func (r *DicomRepository) ReadFile(key string) (dicom.Dataset, error) {
	fileName := fmt.Sprintf("%s/%s", r.databaseName, key)
	if !checkFileExists(fileName) {
		return dicom.Dataset{}, errors.New("Dicom file does not exist")
	}
	return dicom.ParseFile(fmt.Sprintf("%s/%s", r.databaseName, key), nil)

}

func create(filename string) (*os.File, error) {
	if err := os.MkdirAll(filepath.Dir(filename), 0770); err != nil {
		return nil, err
	}
	return os.Create(filename)
}

func checkFileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}
