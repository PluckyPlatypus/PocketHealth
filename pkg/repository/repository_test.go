package repository

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/suyashkumar/dicom"
)

func TestCreateFileSuccess(t *testing.T) {
	// CreateFile SHOULD succeed given a dicom file and a key
	tempDir, err := ioutil.TempDir(".", "temp")
	if err != nil {
		t.Errorf("Error while creating temp dir: %v", err)
	}
	testRepository := NewRepository(tempDir)

	dataset, err := dicom.ParseFile("../testdata/IM000001", nil)
	if err != nil {
		t.Errorf("Error while accessing test dicom: %v", err)
	}

	got := testRepository.CreateFile(dataset, "uniqueTestkey")
	if got != nil {
		t.Errorf("Got %q, Expected function to run without errors", got)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })
}

func TestCreateFileFailure(t *testing.T) {
	// CreateFile SHOULD succeed fail if file with name already exists
	notUniqueFilename := "notUniqueFilename"

	tempDir, err := ioutil.TempDir(".", "temp")
	if err != nil {
		t.Errorf("Error while creating temp dir: %v", err)
	}
	testRepository := NewRepository(tempDir)
	// create file
	os.Create(fmt.Sprintf("%s/%s", tempDir, notUniqueFilename))

	dataset, err := dicom.ParseFile("../testdata/IM000001", nil)
	if err != nil {
		t.Errorf("Error while accessing test dicom: %v", err)
	}

	got := testRepository.CreateFile(dataset, notUniqueFilename)
	if got == nil {
		t.Errorf("Got no error, Expected function to return error")
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })
}

func TestReadFileSuccess(t *testing.T) {
	// ReadFile SHOULD succeed in reading a file GIVEN the filename of an existing file

	filename := "existingFile"

	tempDir, err := ioutil.TempDir(".", "temp")
	if err != nil {
		t.Errorf("Error while creating temp dir: %v", err)
	}
	testRepository := NewRepository(tempDir)

	dataset, err := dicom.ParseFile("../testdata/IM000001", nil)
	if err != nil {
		t.Errorf("Error while accessing test dicom: %v", err)
	}

	// create file for reading
	if err := testRepository.CreateFile(dataset, filename); err != nil {
		t.Errorf("Error while creating test file for reading: %v", err)
	}

	got, err := testRepository.ReadFile(filename)
	if err != nil {
		t.Errorf("Got %q, Expected function to run without errors", err)
	}
	if got.String() != dataset.String() {
		t.Error("Read file does not match expected file")
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })
}

func TestReadFileFailure(t *testing.T) {
	// ReadFile SHOULD fail in reading a file GIVEN the name of a file that doesnt exist
	filename := "nonExistingFile"

	tempDir, err := ioutil.TempDir(".", "temp")
	if err != nil {
		t.Errorf("Error while creating temp dir: %v", err)
	}

	testRepository := NewRepository(tempDir)

	got, err := testRepository.ReadFile(filename)
	if err == nil {
		t.Errorf("Got %q, Expected function to throw error instead", got)
	}
	t.Cleanup(func() { os.RemoveAll(tempDir) })
}
