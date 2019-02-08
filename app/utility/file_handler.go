package utility

import (
	"github.com/nfnt/resize"
	"image/jpeg"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"time"
)

func InitUploadImage(file *multipart.FileHeader, folder string) (bool ,string, string) {

	arrayOfFileName := strings.Split(file.Filename, ".")

	mimeType := strings.ToLower(arrayOfFileName[len(arrayOfFileName)-1])

	imageValidTypes := map[string]string{"jpg": "jpg", "jpeg": "jpeg", "png": "png"}

	//Validation image type
	if _, ok := imageValidTypes[mimeType]; !ok {
		return false, "",  "فرمت تصویر غیرمجاز میباشد"
	}

	t := time.Now().Unix()
	year := strconv.FormatInt(int64(time.Now().Year()), 10)
	month := strconv.FormatInt(int64(time.Now().Month()), 10)
	image := year + "_" + month + "_" + strconv.FormatInt(t, 10) + "." + mimeType

	imagePath := "static/assets/images/" + folder + "/" + year + month + "/" + image

	imageDir := "static/assets/images/" + folder + "/" + year + month

	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		if err := os.Mkdir(imageDir, os.ModePerm); err != nil {
			return false, "",  err.Error()
		}
	}

	return true, imagePath, ""
}

func UploadImageWithResize(imageFile *multipart.FileHeader, folder string, width, height uint) (bool ,string, string) {

	arrayOfFileName := strings.Split(imageFile.Filename, ".")

	mimeType := strings.ToLower(arrayOfFileName[len(arrayOfFileName)-1])

	imageValidTypes := map[string]string{"jpg": "jpg", "jpeg": "jpeg"}

	//Validation image type
	if _, ok := imageValidTypes[mimeType]; !ok {
		return false, "",  "فرمت تصویر غیرمجاز میباشد"
	}

	t := time.Now().Unix()
	year := strconv.FormatInt(int64(time.Now().Year()), 10)
	month := strconv.FormatInt(int64(time.Now().Month()), 10)
	image := year + "_" + month + "_" + strconv.FormatInt(t, 10) + "." + mimeType

	imagePath := "static/assets/images/" + folder + "/" + year + month + "/" + image
	imageDir := "static/assets/images/" + folder + "/" + year + month

	if _, err := os.Stat(imageDir); os.IsNotExist(err) {
		if err := os.Mkdir(imageDir, os.ModePerm); err != nil {
			return false, "",  err.Error()
		}
	}

	file, err := imageFile.Open()
	if err != nil {
		return false, "",  err.Error()
	}
	defer file.Close()

	//Decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return false, "",  err.Error()
	}

	file.Close()

	m := resize.Resize(width, height, img, resize.Lanczos3)

	out, err := os.Create(imagePath)
	if err != nil {
		return false, "",  err.Error()
	}
	defer out.Close()

	//Write new image to file
	if err := jpeg.Encode(out, m, nil); err != nil {
		return false, "",  err.Error()
	}

	return true, imagePath, ""
}
