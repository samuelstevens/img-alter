package img

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/samuelstevens/gocaption/api"
	"github.com/samuelstevens/gocaption/caption"
	"github.com/samuelstevens/gocaption/util"

	"golang.org/x/net/html"
)

// ImgTag represents an HTML image with attributes
type ImgTag struct {
	relativePath string
	absPath      string
	Caption      *caption.Caption
	attributes   []html.Attribute
}

func New(tok *html.Token, webpagePath string, api *api.Client) (*ImgTag, error) {
	// parse token
	prevDescription := ""
	relativePath := ""
	attributes := []html.Attribute{}

	for _, attr := range tok.Attr {
		switch attr.Key {
		case "alt":
			prevDescription = attr.Val
		case "src":
			relativePath = attr.Val
		default:
			attributes = append(attributes, attr)
		}
	}

	absPath, err := util.MakeAbsRelativeTo(webpagePath, relativePath)

	if err != nil {
		log.Fatalf("%s: %s\n", webpagePath, err.Error())
		return nil, err
	}

	// get caption
	caption, err := caption.New(absPath, prevDescription, api)

	if err != nil {
		log.Printf(err.Error())
	}

	// make tag
	img := ImgTag{
		absPath:      absPath,
		relativePath: relativePath,
		Caption:      caption,
		attributes:   attributes,
	}

	return &img, nil
}

func (img *ImgTag) Filename() string {
	return filepath.Base(img.relativePath)
}

func (img *ImgTag) String() string {
	tag := fmt.Sprintf("<img src=\"%s\" alt=\"%s\"", img.relativePath, img.Caption.Description)

	for _, attr := range img.attributes {
		tag = tag + fmt.Sprintf(" %s=\"%s\"", attr.Key, attr.Val)
	}

	tag = tag + " />"

	return tag
}
