package webpage

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/samuelstevens/gocaption/api"
	"github.com/samuelstevens/gocaption/caption"
	"github.com/samuelstevens/gocaption/img"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// WebPage represents an HTML file that will have its <img/>
// tags updated with an "alt" attribute
type WebPage struct {
	absolutePath string
	content      *strings.Builder
	file         *os.File
	Captions     []*caption.Caption
}

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
		content:      &strings.Builder{},
		file:         nil,
		Captions:     []*caption.Caption{},
	}, nil
}

func (wp *WebPage) open() error {
	if wp.file != nil {
		return errors.New("already an open file")
	}

	file, err := os.Open(wp.absolutePath)

	if err != nil {
		return err
	}

	wp.file = file

	return nil
}

func (wp *WebPage) close() error {
	if wp.file == nil {
		return nil
	}

	err := wp.file.Close()

	if err != nil {
		return err
	}

	wp.file = nil

	return nil
}

func (wp *WebPage) Write() error {
	if err := wp.close(); err != nil {
		return err
	}

	return ioutil.WriteFile(wp.absolutePath, []byte(wp.content.String()), 0644)
}

// LabelImages takes all the <img> in an .html document and adds
// an "alt" attribute if it is missing.
func (wp *WebPage) LabelImages(client *api.Client) error {
	if err := wp.open(); err != nil {
		return err
	}

	defer wp.close()

	tokenizer := html.NewTokenizer(wp.file)

	for {
		switch tokenType := tokenizer.Next(); tokenType {

		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return nil
			}
			return tokenizer.Err()

		case html.TextToken, html.StartTagToken, html.EndTagToken, html.DoctypeToken:
			t := tokenizer.Token()
			fmt.Fprintf(wp.content, t.String())

		case html.SelfClosingTagToken: // might be <img/> tag
			t := tokenizer.Token()

			if t.DataAtom == atom.Img || t.DataAtom == atom.Image { // not sure what the difference is
				imgTag, err := img.New(&t, wp.absolutePath, client)

				if err != nil {
					return err
				}

				fmt.Fprintf(wp.content, imgTag.String())

				wp.Captions = append(wp.Captions, imgTag.Caption)
			} else {
				fmt.Fprintf(wp.content, t.String())
			}
		case html.CommentToken:
			// ignore comments
		default:
			return fmt.Errorf("unknown type: %s", tokenizer.Token().String())
		}
	}
}
