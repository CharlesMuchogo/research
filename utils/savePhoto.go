package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

func SavePhoto(c *gin.Context, file *multipart.FileHeader, userID string) (string, error) {
	assetsDir := os.Getenv("PHOTO_DIRECTORY")
	log.Printf("assetsDir: %s", assetsDir)
	if _, err := os.Stat(assetsDir); os.IsNotExist(err) {
		log.Printf("Directory %s does not exist. Creating...", assetsDir)
		if err := os.Mkdir(assetsDir, 0755); err != nil {
			return "", err
		}
	}

	fileExt := filepath.Ext(file.Filename)
	currentTimestamp := time.Now().UnixNano() / int64(time.Millisecond)

	filename := fmt.Sprintf("%d_%s%s", currentTimestamp, userID, fileExt)

	dst := filepath.Join(assetsDir, filename)
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return "", err
	}
	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	photoPath := fmt.Sprintf("/images/%s", filename)
	return photoPath, nil
}
