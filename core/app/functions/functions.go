package functions

import (
	"bytes"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"strings"

	draw "golang.org/x/image/draw"
)

func DownloadByHTTP(URL string, headers http.Header) ([]byte, error) {
	// Why not "net/url" - "10.0.0.1:800": first path segment in URL cannot contain colon
	normalizedURL := URL
	if strings.HasPrefix(normalizedURL, "https://") {
		normalizedURL = normalizedURL[len("https://"):len(URL)]
	}
	if !strings.HasPrefix(normalizedURL, "http://") {
		normalizedURL = strings.Join([]string{"http://", normalizedURL}, "")
	}
	client := http.Client{}
	request, err := http.NewRequest("GET", normalizedURL, nil)
	if err != nil {
		return nil, err
	}
	request.Header = headers
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return bodyBytes, nil
}

func TransformByNearestNeighbor(height int, width int, data []byte) ([]byte, error) {
	if width == 0 && height == 0 {
		return data, nil
	}
	src, err := jpeg.Decode(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	srcBounds := src.Bounds()
	if height == 0 {
		height = srcBounds.Dy()
	}
	if width == 0 {
		width = srcBounds.Dx()
	}
	dst := image.NewRGBA(image.Rect(0, 0, height, width))
	draw.NearestNeighbor.Scale(dst, dst.Rect, src, srcBounds, draw.Over, nil)
	buff := bytes.NewBuffer([]byte{})
	err = jpeg.Encode(buff, dst, &jpeg.Options{Quality: 100})
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}
