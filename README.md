# PocketHealth

Backend Coding Challenge


# Usage
To test, simply check out this project and start it locally using 'go run main.go' within the project folder.
Make sure you have go installed as well.

# API
This project gives a user access to 3 endpoints:

- Store DICOM images
To store images, call endpoint 'http://localhost:8080/dicom' with a POST request. The DICOM file that you intend to store should be attached to the request as a file using the Multipart form. On success, the request returns the id of the newly stored DICOM file. The request fails if the attached file is not a valid DICOM file (e.g. a file that can't be parsed using [the dicom package for go](https://pkg.go.dev/github.com/suyashkumar/dicom). The request also fails if a file with the same id has already been stored. A file's id is generated using the DICOM's StudyInstanceUID, SeriesInstanceUID and SOPInstanceUID. These three values in combination should be unique to a DICOM.
For an example request, see screenshots below

- Request PNG image for a stored DICOM
To request a PNG image for a stored DICOM, call endpoint 'http://localhost:8080/dicom/:studyId', with :studyId being the unique identifier that was returned when the DICOM image was stored. On success. the request returns the PNG file. The request fails if no DICOM file with the given id is stored.
For an example request, see screenshots below

- Request Attribute Header content for a stored DICOM and a tag
To request a certain attribute value for a stored DICOM, call endpoint 'http://localhost:8080/dicom/:studyId/:tag', with :studyId being the unique identifier that was returned when the DICOM image was stored and :tag being the name of a DICOM tag. A list of names can be found [here](https://pkg.go.dev/github.com/suyashkumar/dicom@v1.0.6/pkg/tag). On success, the request returns the value of the given tag for the stored DICOM image. The request fails if no DICOM file with the given id is stored. The request also fails if no valid tag name was provided. The request also files if the tag is not present in the DICOM file.
For an example request, see screenshots below

# Data storage
For the scope of this project, DICOM files are stored within the project folder, in a subfolder with path temp/dicom-management. To 'reset' the storage, simply delete that folder. It will automatically be re-created the next time you try to store a DICOM image.

# What else?
This is my first project using Go. Yay! You might have been able to tell by looking at the code, or by the fact that you just had to read through this README, instad of looking at a porper yaml-file. Despite that, I had a lot of fun learning more about go and am looking forward to further extending the tests and adding some API documentation later on. 
And if you're the one reviewing this project: Thanks for being here! Have a cookie 🍪
