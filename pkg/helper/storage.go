package helper

import (
	"fmt"
	"github.com/google/uuid"
	"time"
)

func GeneratePathAndFilenameStorage(field string, ext string) (string, string) {
	filename := fmt.Sprintf("%s-%s.%s", field, uuid.New().String(), ext)
	path := fmt.Sprintf("%s/%s", time.Now().Format("2006-01-02"), filename)

	return filename, path
}
