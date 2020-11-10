package files

import (
	"fmt"
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

	fmt.Println(name)
	file, err := os.Open(name)

	if err != nil {
		log.Println(err)
	}
	defer file.Close()
	// decode jpeg into image.Image
	img, _, err := image.Decode(file)
	if err != nil {
		log.Println(err)
	}

	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio

	m := resize.Resize(90, 0, img, resize.Lanczos3)

	name = strings.Replace(name, ".jpg", "-min.jpg", 1)
	name = strings.Replace(name, "Пользователи", "Пользователи-min", 1)
	fmt.Println(name)
	out, err := os.Create(name)
	if err != nil {
		log.Println(err)
	}
	defer out.Close()

	// write new image to file
	jpeg.Encode(out, m, nil)

}
