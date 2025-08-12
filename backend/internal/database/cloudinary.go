package database

import (
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

var cldClient *cloudinary.Cloudinary

func ConnectCloudinary(cldURL string) (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromURL(cldURL)
	if err != nil {
		log.Fatalf("Cloudinary initialization error: %v", err)
		return nil, err
	}
	return cld, nil
}

func SetCloudinary(c *cloudinary.Cloudinary) {
	cldClient = c
}

func GetCloudinary() *cloudinary.Cloudinary {
	return cldClient
}
