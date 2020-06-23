package webpage

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/samuelstevens/gocaption/api"
	"github.com/samuelstevens/gocaption/caption"
	"github.com/samuelstevens/gocaption/util"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// WebPage represents an HTML file that will have its <img/>
// tags updated with an "alt" attribute
type WebPage struct {
	absolutePath string
	content      string
	Captions     []*caption.Caption
}

type LabelFunc func(imgPath string, prevDescription string) string

// New returns a new WebPage
func New(path string) (*WebPage, error) {
	if filepath.Ext(path) != ".html" {
		return nil, &FileTypeError{path}
	}

	path, err := filepath.Abs(path)

	if err != nil {
		return nil, err
	}

	return &WebPage{
		absolutePath: path,
		content:      "",
		Captions:     []*caption.Caption{},
	}, nil
}

func (wp *WebPage) Write() error {
	return ioutil.WriteFile(wp.absolutePath, []byte(wp.content), 0644)
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}

// LabelNode recursively searches through an html node
// and adds an attribute to any image nodes it finds
func LabelNode(n *html.Node, labelFunc LabelFunc) {

	switch n.Type {
	case html.DocumentNode, html.ElementNode:
		if n.DataAtom == atom.Img || n.DataAtom == atom.Image {
			var imgSrc, imgAlt string

			for _, a := range n.Attr {
				if a.Key == "src" {
					imgSrc = a.Val
				}

				if a.Key == "alt" {
					imgAlt = a.Val
				}
			}

			caption := labelFunc(imgSrc, imgAlt)
			newAlt := html.Attribute{Key: "alt", Val: caption}
			newAttr := []html.Attribute{newAlt}

			for _, a := range n.Attr {
				if a.Key != "alt" {
					newAttr = append(newAttr, a)
				}
			}

			n.Attr = newAttr
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			LabelNode(child, labelFunc)
		}
	}
}

// LabelImages takes an unescaped HTML string and returns a new string containing labeled images
func LabelImages(inputHTML string, labelFunc LabelFunc) (string, error) {
	doc, err := html.Parse(strings.NewReader(inputHTML))

	if err != nil {
		return "", err
	}

	LabelNode(doc, labelFunc)

	return renderNode(doc), nil
}

// LabelImages takes all the <img> in an .html document and adds
// an "alt" attribute if it is missing.
func (wp *WebPage) LabelImages(client *api.Client) error {

	file, err := os.Open(wp.absolutePath)

	if err != nil {
		return err
	}

	defer file.Close()

	rawDoc, err := ioutil.ReadAll(file)

	if err != nil {
		return err
	}

	updatedDoc, err := LabelImages(string(rawDoc), func(relativeImgPath string, prevDescription string) string {
		absImgPath, err := util.MakeAbsRelativeTo(wp.absolutePath, relativeImgPath)

		if err != nil {
			return ""
		}

		caption, err := caption.New(absImgPath, prevDescription, client)

		if err != nil {
			return ""
		}

		wp.Captions = append(wp.Captions, caption)

		return caption.Description
	})

	if err != nil {
		return err
	}

	wp.content = updatedDoc

	return nil
}
