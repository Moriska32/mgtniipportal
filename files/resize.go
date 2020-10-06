package files

import (
	"image"
	"image/jpeg"
	"io"
	"log"
	"os"
	"strings"

	"github.com/nfnt/resize"
)

//Resize pictures

func guessImageFormat(r io.Reader) (format string, err error) {
	_, format, err = image.DecodeConfig(r)
	return
}

func guessImageMimeTypes(r io.Reader) bool {
	format, _ := guessImageFormat(r)
	if format == "" {
		return true
	}
	return false
}

//Resize resize image
func Resize(name string) {

	file, err := os.Open(name)

	if guessImageMimeTypes(file) {
		return
	}

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}

	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio

	m := resize.Resize(90, 0, img, resize.Lanczos3)

	name = strings.Replace(name, ".jpg", "-min.jpg", 0)
	name = strings.Replace(name, "Пользователи", "Пользователи-min.jpg", 0)

	out, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

}

func main() {

	Resize("D:/progi/Microsoft VS Code/Programs/mgtniipportal/public/photos/ss/43.jpg")

}
