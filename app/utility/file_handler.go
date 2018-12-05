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

func UploadImage(file *multipart.FileHeader, folder string) (bool ,string, string) {

	arrayOfFileName := strings.Split(file.Filename, ".")

	mimeType := strings.ToLower(arrayOfFileName[len(arrayOfFileName)-1])

	imageValidTypes := map[string]string{"jpg": "jpg", "jpeg": "jpeg", "png": "png"}

	//Validation image type
	if _, ok := imageValidTypes[mimeType]; !ok {
		return false, "",  "image type error"
	}

	t := time.Now().Unix()
	year := strconv.FormatInt(int64(time.Now().Year()), 10)
	month := strconv.FormatInt(int64(time.Now().Month()), 10)
	image := year + "_" + month + "_" + strconv.FormatInt(t, 10) + "." + mimeType

	imagePath := "static/assets/images/" + folder + "/" + year + month + "/" + image

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		os.Mkdir("static/assets/images/"+folder+"/"+year+month, os.ModePerm)
	}

	return true, imagePath, ""
}

func UploadImageCustom(imageFile *multipart.FileHeader, folder string, width, height uint) (bool ,string, string) {

	arrayOfFileName := strings.Split(imageFile.Filename, ".")

	mimeType := strings.ToLower(arrayOfFileName[len(arrayOfFileName)-1])

	imageValidTypes := map[string]string{"jpg": "jpg", "jpeg": "jpeg"}

	//Validation image type
	if _, ok := imageValidTypes[mimeType]; !ok {
		return false, "",  "فرمت غیرمجاز میباشد"
	}

	t := time.Now().Unix()
	year := strconv.FormatInt(int64(time.Now().Year()), 10)
	month := strconv.FormatInt(int64(time.Now().Month()), 10)
	image := year + "_" + month + "_" + strconv.FormatInt(t, 10) + "." + mimeType

	imagePath := "static/assets/images/" + folder + "/" + year + month + "/" + image

	if _, err := os.Stat(imagePath); os.IsNotExist(err) {
		os.Mkdir("static/assets/images/"+folder+"/"+year+month, os.ModePerm)
	}

	file, err := imageFile.Open()
	if err != nil {
		return false, "",  err.Error()
	}
	defer file.Close()

	// decode jpeg into image.Image
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

	// write new image to file
	jpeg.Encode(out, m, nil)

	return true, imagePath, ""
}
