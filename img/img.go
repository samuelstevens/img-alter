package img

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/samuelstevens/goimglabeler/api"
	"github.com/samuelstevens/goimglabeler/caption"

	"golang.org/x/net/html"
)

// ImgTag represents an HTML image with attributes
type ImgTag struct {
	relativePath string
	absPath      string
	caption      *caption.Caption
	attributes   []html.Attribute
}

func New(tok *html.Token, rootDir string, api *api.Client) *ImgTag {
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

	absPath := filepath.Join(rootDir, relativePath)

	// get caption
	caption, err := caption.New(absPath, prevDescription, api)

	if err != nil {
		log.Printf(err.Error())
	}

	// make tag
	img := ImgTag{
		relativePath: relativePath,
		absPath:      absPath,
		caption:      caption,
		attributes:   attributes,
	}

	return &img
}

func (img *ImgTag) Filename() string {
	return filepath.Base(img.relativePath)
}

func (img *ImgTag) String() string {
	tag := fmt.Sprintf("<img src=\"%s\" alt=\"%s\"", img.relativePath, img.caption.Description)

	for _, attr := range img.attributes {
		tag = tag + fmt.Sprintf(" %s=\"%s\"", attr.Key, attr.Val)
	}

	tag = tag + " />"

	return tag
}
