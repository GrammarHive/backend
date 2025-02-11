// core/uploaderService/validator.go

package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"regexp"
)

var AllowedMimeTypes = map[string]bool{
	"text/plain": true,
	"application/json": true,
}

func (s *ProfileService) ValidateInput(name, username string) error {
	if err := s.ValidateName(name); err != nil {
		return err
	}
	if err := s.ValidateUsername(username); err != nil {
		return err
	}
	return nil
}

func (s *ProfileService) ValidateFile(file multipart.File) error {
	fileTypeBuffer := make([]byte, 512)
	if _, err := file.Read(fileTypeBuffer); err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// This resets the file after reading
	file.Seek(0, io.SeekStart)

	// Retrieve the MIME type
	mimeType := http.DetectContentType(fileTypeBuffer)
	if !AllowedMimeTypes[mimeType] {
		return fmt.Errorf("file type %s is not allowed", mimeType)
	}
	return nil
}

func (s *ProfileService) ValidateName(name string) error {
	if len(name) == 0 {
		return fmt.Errorf("name cannot be empty")
	}
	if len(name) > 100 {
		return fmt.Errorf("name cannot exceed 100 characters")
	}
	return nil
}

func (s *ProfileService) ValidateUsername(username string) error {
	if len(username) == 0 {
		return fmt.Errorf("username cannot be empty")
	}
	if len(username) > 50 {
		return fmt.Errorf("username cannot exceed 50 characters")
	}
	// Regex to allow alphanumeric characters and underscores
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !re.MatchString(username) {
		return fmt.Errorf("username can only contain alphanumeric characters and underscores")
	}
	return nil
}
