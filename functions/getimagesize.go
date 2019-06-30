package functions

import (
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"mime/multipart"
)

func GetImageSize(file *multipart.FileHeader) (int, int, error) {
	fp, err := file.Open()

	defer fp.Close()

	if err != nil {
		return 0, 0, err
	}

	iconfig, _, err := image.DecodeConfig(fp)

	if err != nil {
		return 0, 0, err
	}

	return iconfig.Width, iconfig.Height, nil
}