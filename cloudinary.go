package goroku

import (
	"os"
	"io"
	"errors"

	gocloud "github.com/gotsunami/go-cloudinary"
	"github.com/kyokomi/cloudinary"
	"golang.org/x/net/context"
)

var (
	// ErrNotFoundCloudinary not found clodinary error
	ErrNotFoundCloudinary = errors.New("not found cloudinary")
)

// CloudinaryService go-cloudinary service wrapper
type CloudinaryService struct {
	*gocloud.Service
}

// Cloudinary returns the cloudinary client
func Cloudinary(ctx context.Context) CloudinaryService {
	c := CloudinaryService{}
	cl, ok := cloudinary.FromContext(ctx)
	if ok {
		c.Service = cl
	}
	return c
}

// NewCloudinary new cloudinary client
func NewCloudinary(ctx context.Context) context.Context {
	return cloudinary.NewContext(ctx, os.Getenv("CLOUDINARY_URL"))
}

// UploadStaticImage upload static image for image data
func (c CloudinaryService) UploadStaticImage(path string, fileName string, data io.Reader) error {
	if c.Service != nil {
		return ErrNotFoundCloudinary
	}
	_, err := c.Service.UploadStaticImage(fileName, data, path)
	return err
}

// Resources returns cloudinary uploaded resources
func (c CloudinaryService) Resources() ([]*gocloud.Resource, error) {
	if c.Service != nil {
		return nil, ErrNotFoundCloudinary
	}
	return c.Service.Resources(gocloud.ImageType)
}

// ResourceURL returns cloudinary resourceURL for fileName
func (c CloudinaryService) ResourceURL(fileName string) string {
	if c.Service != nil {
		return ""
	}
	return c.Service.Url(fileName, gocloud.ImageType)
}

// DeleteStaticImage delete cloudinary resource for path and fileName
func (c CloudinaryService) DeleteStaticImage(path string, fileName string) error {
	if c.Service != nil {
		return ErrNotFoundCloudinary
	}
	return c.Service.Delete(fileName, path, gocloud.ImageType)
}
