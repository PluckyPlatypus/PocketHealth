package repository

import "github.com/suyashkumar/dicom"

type MockRepository struct {
	databaseName string
}

func (m *MockRepository) CreateFile(newDicom dicom.Dataset, key string) error {
	return nil
}

func (m *MockRepository) ReadFile(key string) (dicom.Dataset, error) {
	return dicom.ParseFile(key, nil)
}
