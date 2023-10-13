package models

import "fmt"

type Transformation struct {
	height int32
	width  int32
	url    string
}

func NewTransformation(
	height int32,
	width int32,
	url string,
) *Transformation {
	return &Transformation{height, width, url}
}

func (t Transformation) Identity() string {
	return fmt.Sprintf("%d/%d/%s", t.GetHeight(), t.GetWidth(), t.GetUrl())
}

func (t Transformation) GetUrl() string {
	return t.url
}

func (t Transformation) GetWidth() int32 {
	return t.width
}

func (t Transformation) GetHeight() int32 {
	return t.height
}
