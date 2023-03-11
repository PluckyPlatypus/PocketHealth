package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"

	"pocket-health/pkg/service"
)

type Controller interface {
	PostDicom(c *gin.Context)
	GetDicomPng(c *gin.Context)
	GetDicomAttributeHeader(c *gin.Context)
}

type DicomController struct {
	dicomService service.Service
}

func NewController() Controller {
	return &DicomController{
		dicomService: service.NewService(),
	}
}

// TODO: swagger or other rest annotations?
func (controller *DicomController) PostDicom(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "No file in request", "error": err.Error()})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "File in request could not be opened", "error": err.Error()})
		return
	}

	dataset, err := dicom.Parse(openedFile, file.Size, nil)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "File is not a valid DICOM file", "error": err.Error()})
		return
	}

	newId, err := controller.dicomService.StoreDicom(dataset)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Saving DICOM file failed", "error": err.Error()})
	} else {
		c.IndentedJSON(http.StatusOK, gin.H{"id": newId})
	}
}

// GET png image for DICOM image from local storage
func (controller *DicomController) GetDicomPng(c *gin.Context) {
	id := c.Param("studyId")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No DICOM id provided"})
		return
	}

	buffer, err := controller.dicomService.GetPngImageForDicom(id)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not parse PNG image for given dicom", "error": err.Error()})
		return
	}

	c.Writer.Header().Set("Content-Type", "image/png")
	if _, err := c.Writer.Write(buffer.Bytes()); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not send PNG image for given DICOM", "error": err.Error()})
		return
	}
	c.IndentedJSON(http.StatusNotImplemented, "GetDicomPng has not been implemented yet")
}

// GET attribute header content for DICOM image from local storage
func (controller *DicomController) GetDicomAttributeHeader(c *gin.Context) {
	id := c.Param("studyId")
	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "No DICOM id provided"})
		return
	}

	requestedTag := c.Param("tag")
	tagInfo, err := tag.FindByName(requestedTag)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Provided invalid tag parameter", "error": err.Error()})
		return
	}

	attrValue, err := controller.dicomService.GetAttributeHeaderContentForDicom(id, tagInfo.Tag)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Could not retrieve attribute header for given file and tag", "error": err.Error()})
		return
	}

	c.IndentedJSON(http.StatusAccepted, gin.H{"attribute": attrValue})
}
