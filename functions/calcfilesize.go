package functions

import (
	"math"
	"mime/multipart"
	"strconv"
)

const KB = 1024
var SIZES = []string{"B", "KB", "MB", "GB"}

type Size interface {
	Size() int64
}

func CalcFileSize(file *multipart.FileHeader) (*string, error) {
	fp, err := file.Open()

	if err != nil {
		return nil, err
	}

	size := fp.(Size).Size()

	selectv := math.Floor(math.Log(float64(size)) / math.Log(float64(KB)))

	result := strconv.Itoa(int(float64(size) / math.Pow(float64(KB), float64(selectv)))) + SIZES[int(selectv)]

	return &result, nil
}