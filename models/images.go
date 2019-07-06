package models

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

//Image is the image struct -- it will NOT go to the DB!
type Image struct {
	GalleryID uint
	Filename  string
}

//Path will return the Path to the image... you would have never guessed!
func (i *Image) Path() string {
	temp := url.URL{
		Path: "/" + i.OsPath(),
	}
	return temp.String()
}

//OsPath will return the Path to the image... but without the /
func (i *Image) OsPath() string {
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.Filename)
}

//ImageService is the interfacewe nned for our images
type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

//imageService is the available API
type imageService struct{}

//NewImageService will init the gallery service and will make it available
func NewImageService() ImageService {
	return &imageService{}
}

func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	defer r.Close()
	path, err := is.mkImagePath(galleryID)
	if err != nil {
		return err
	}

	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (is *imageService) Delete(i *Image) error {
	return os.Remove(i.OsPath())
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (is *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755) // hackerage...
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}

func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	galleryPath := is.imagePath(galleryID)
	images, err := filepath.Glob(galleryPath + "*")
	if err != nil {
		return nil, err
	}
	ret := make([]Image, len(images))
	for i := range images {
		images[i] = "\\" + images[i]
		images[i] = strings.ReplaceAll(images[i], "\\", "/")
		images[i] = strings.ReplaceAll(images[i], fmt.Sprintf("/images/galleries/%v/", galleryID), "")
		ret[i] = Image{
			Filename:  images[i],
			GalleryID: galleryID,
		}
	}
	return ret, nil
}
