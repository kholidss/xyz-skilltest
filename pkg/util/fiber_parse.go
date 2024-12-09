package util

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kholidss/xyz-skilltest/internal/presentation"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

func FiberParseFile(c *fiber.Ctx, key string) (*presentation.File, error) {
	formFile, err := c.FormFile(key)
	if err != nil {
		return nil, nil
	}

	// Open the uploaded file
	file, err := formFile.Open()
	if err != nil {
		return nil, errors.Wrap(err, "failed opening form file")
	}
	defer file.Close()

	// Read the file's content into a byte slice
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, errors.Wrap(err, "failed reading content file")
	}

	return &presentation.File{
		Name:     formFile.Filename,
		Mimetype: http.DetectContentType(fileBytes),
		Size:     int(formFile.Size),
		File:     fileBytes,
	}, nil

}
