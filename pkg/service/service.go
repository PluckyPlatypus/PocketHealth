package service

import (
	"bytes"
	"image/png"
	"strings"

	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"

	"pocket-health/pkg/repository"
)

type Service interface {
	StoreDicom(newDicom dicom.Dataset) (string, error)
	GetPngImageForDicom(key string) (*bytes.Buffer, error)
	GetAttributeHeaderContentForDicom(key string, tag tag.Tag) (string, error)
}

type DicomSerivce struct {
	dicomRepository repository.Repository
}

func NewService() Service {
	return &DicomSerivce{
		dicomRepository: repository.NewRepository("temp/dicom-management"),
	}
}

// body as parsed dicom image or just as body object, to be parsed and validated here
func (s *DicomSerivce) StoreDicom(newDicom dicom.Dataset) (string, error) {
	key, err := generateId(newDicom)
	if err != nil {
		return "", err
	}

	if err := s.dicomRepository.CreateFile(newDicom, key); err != nil {
		return "", err
	}

	return key, nil
}

func (s *DicomSerivce) GetPngImageForDicom(key string) (*bytes.Buffer, error) {
	buffer := new(bytes.Buffer)
	dicomFile, err := getDicomFile(key, s)
	if err != nil {
		return buffer, err // dicom file with given key does not exist
	}

	pixelDataElement, err := dicomFile.FindElementByTag(tag.PixelData)
	if err != nil {
		return buffer, err // could not read pixel data of dicom file
	}

	pixelDataInfo := dicom.MustGetPixelDataInfo(pixelDataElement.Value)
	for _, fr := range pixelDataInfo.Frames {
		img, err := fr.GetImage()
		if err != nil {
			return buffer, err // could not extract image from dicom file
		}
		if err := png.Encode(buffer, img); err != nil {
			return buffer, err // could not encode file to png image
		}
	}
	return buffer, err
}

func (s *DicomSerivce) GetAttributeHeaderContentForDicom(key string, tag tag.Tag) (string, error) {
	dicomFile, err := getDicomFile(key, s)
	if err != nil {
		return "", err // dicom file with given key does not exist
	}

	attr, err := dicomFile.FindElementByTag(tag)
	if err != nil {
		return "", err // no value for given tag in dicom image
	}
	return attr.Value.String(), err
}

func getDicomFile(key string, s *DicomSerivce) (dicom.Dataset, error) {
	dataset, err := s.dicomRepository.ReadFile(key)
	if err != nil {
		return dicom.Dataset{}, err // could not read dicom file, e.g. because file does not exist
	}
	return dataset, err
}

func generateId(dicomFile dicom.Dataset) (string, error) {
	studyInstance, errStudyInstance := dicomFile.FindElementByTag(tag.StudyInstanceUID)
	seriesInstance, errSeriesInstance := dicomFile.FindElementByTag(tag.SeriesInstanceUID)
	sopInstance, errSopInstance := dicomFile.FindElementByTag(tag.SOPInstanceUID)

	if errStudyInstance != nil || errSeriesInstance != nil || errSopInstance != nil {
		// dicom tags for unique id are not all present -> dicom is invalid!
		// return one of the present errors
		var err error
		if errStudyInstance != nil {
			err = errStudyInstance
		} else if errSeriesInstance != nil {
			err = errSeriesInstance
		} else {
			err = errSopInstance
		}
		return "", err
	}

	// dicom images are uniquely identified by the combination of StudyInstanceUID, SeriesInstanceUID and SOPInstanceUID
	// to create one unique identifier, merge the three values, then re-format the string
	guidArray := []string{studyInstance.Value.String(), seriesInstance.Value.String(), sopInstance.Value.String()}
	var dicomId string
	dicomId = strings.Join(guidArray, "-")
	dicomId = strings.ReplaceAll(dicomId, ".", "")
	dicomId = strings.ReplaceAll(dicomId, "[", "")
	dicomId = strings.ReplaceAll(dicomId, "]", "")
	return dicomId, nil
}
