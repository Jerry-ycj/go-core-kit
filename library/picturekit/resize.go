package picturekit

import (
	"github.com/mizuki1412/go-core-kit/class/exception"
	"github.com/mizuki1412/go-core-kit/library/filekit"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
)

// Resize 压缩图片,保持长宽不变,命名为原来基础上+.jpg
func Resize(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		panic(exception.New(err.Error()))
	}
	defer file.Close()
	bin := filekit.ReadBytes(filepath)
	contentType := http.DetectContentType(bin)
	var img image.Image
	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(file)
	case "image/png":
		img, err = png.Decode(file)
	default:
		panic(exception.New("不支持的类型"))
	}
	if err != nil {
		panic(exception.New(err.Error()))
	}
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	m := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	out, err := os.Create(filepath + ".jpeg")
	if err != nil {
		panic(exception.New(err.Error()))
	}
	err = jpeg.Encode(out, m, nil)
	if err != nil {
		panic(exception.New(err.Error()))
	}
}
