package utils

import (
	"encoding/base64"
	"fmt"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"os"
)

func BlurAndResizeImage(file *os.File) string {
	_, err := file.Seek(0, 0)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	imageToResize, _, err := image.Decode(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	resizedImage := resize.Resize(9, 0, imageToResize, resize.Lanczos3)
	tempFile, _ := os.CreateTemp("", "")
	defer func(name string) {
		err = os.Remove(name)
		if err != nil {
			fmt.Println(err)
		}
	}(tempFile.Name())
	err = jpeg.Encode(tempFile, resizedImage, nil)
	readFile, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(readFile)
}
