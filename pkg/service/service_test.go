package service

import (
	"pocket-health/pkg/repository"
	"testing"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

func NewTestService() Service {
	return &DicomSerivce{
		dicomRepository: &repository.MockRepository{},
	}
}

func TestStoreDicomSuccess(t *testing.T) {
	// StoreDicom SHOULD succeed for a valid dicom
	testService := NewTestService()

	dataset, err := dicom.ParseFile("../testdata/IM000001", nil)
	if err != nil {
		t.Errorf("Error while accessing test dicom: %v", err)
	}

	resultString, resultErr := testService.StoreDicom(dataset)
	if resultString == "" || resultErr != nil {
		t.Errorf("Got %q, Expected function to run without errors", resultErr)
	}
}

func TestStoreDicomFailure(t *testing.T) {
	// StoreDicom SHOULD fail if dicom is not valid, e.g. if mandatory attributes are missing
	testService := NewTestService()

	resultString, resultErr := testService.StoreDicom(dicom.Dataset{})
	if resultString != "" || resultErr == nil {
		t.Errorf("Got %q, Expected function to fail", resultString)
	}
}

func TestGetPngImageForDicomSuccess(t *testing.T) {
	// GetPngImageForDicom SHOULD succeed for a valid dicom that exists in storage
	testService := NewTestService()

	_, err := testService.GetPngImageForDicom("../testdata/IM000001")
	if err != nil {
		t.Errorf("Got %q, Expected function to run without errors", err)
	}
}

func TestGetPngImageForDicomFailure(t *testing.T) {
	// GetPngImageForDicom SHOULD fail if dicom file does not exist
	testService := NewTestService()

	_, err := testService.GetPngImageForDicom("invalid path")
	if err == nil {
		t.Errorf("Expected function to return error due to unavailable dicom file")
	}
}

// TODO: add test that tries to extract value from empty dicom

func TestGetAttributeHeaderContentForDicomSuccess(t *testing.T) {
	// GetAttributeHeaderContentForDicom SHOULD succeed for a valid dicom that exists in storage, and valid tag
	testService := NewTestService()

	tags := []tag.Tag{
		tag.StudyInstanceUID,
		tag.SeriesInstanceUID,
		tag.SOPInstanceUID,
		tag.SOPClassUID,
		tag.PatientID,
	}

	for _, test := range tags {
		if _, err := testService.GetAttributeHeaderContentForDicom("../testdata/IM000001", test); err != nil {
			t.Errorf("Got %q, Expected function to run without errors", err)
		}
	}
}

func TestGetAttributeHeaderContentForDicom(t *testing.T) {
	// GetAttributeHeaderContentForDicom SHOULD fail if dicom file does not exist
	testService := NewTestService()

	_, err := testService.GetAttributeHeaderContentForDicom("invalid path", tag.ACR_NEMA_2C_AdaptiveMapFormat)
	if err == nil {
		t.Errorf("Expected function to return error due to unavailable dicom file")
	}
}

// TODO: add test that tries to extract value from empty dicom
