package rest

import (
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
)

const (
	contentTypeKey = "Content-Type"
	contentAudio   = "audio"
)

var allowedExtensions = map[string]struct{}{
	"mp3":  {},
	"wav":  {},
	"mpeg": {},
}

func parseMultipartForm(r *http.Request, key string) (multipart.File, *multipart.FileHeader, error) {
	// parse input, type multipart/form-data
	// 50 MB
	maxMemory := int64(50 << 20)

	err := r.ParseMultipartForm(maxMemory)
	if err != nil {
		log.Printf("Error while ParseMultipartForm: %v", err)
		return nil, nil, err
	}

	// retrieve file from posted form-data
	file, header, err := r.FormFile(key)
	if err != nil {
		return nil, nil, fmt.Errorf("Error retrieving file from form-data: %v\n", err)
	}

	return file, header, nil
}

func getFileExtension(contentType string) (string, error) {
	slashIdx := strings.Index(contentType, "/")
	if slashIdx == -1 {
		return "", fmt.Errorf("bad content type")
	}

	fileType := contentType[:slashIdx]
	extension := contentType[slashIdx+1:]

	if fileType != contentAudio {
		return "", fmt.Errorf("file type %s is not audio", fileType)
	}

	_, isAllowed := allowedExtensions[extension]
	if isAllowed {
		return extension, nil
	} else {
		return "", fmt.Errorf("extension %s is not allowed", extension)
	}
}
