package webpage

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/samuelstevens/goimglabeler/api"
	"github.com/samuelstevens/goimglabeler/img"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type WebPage struct {
	absolutePath string
	content      *strings.Builder
	rootDir      string
	file         *os.File
}

func New(absolutePath string, rootDir string) (*WebPage, error) {
	if filepath.Ext(absolutePath) != ".html" {
		return nil, &FileTypeError{absolutePath}
	}

	return &WebPage{
		absolutePath: absolutePath,
		content:      &strings.Builder{},
		rootDir:      rootDir,
		file:         nil,
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

func (wp *WebPage) saveContent() error {
	if err := wp.close(); err != nil {
		return err
	}

	return ioutil.WriteFile(wp.absolutePath, []byte(wp.content.String()), 0644)
}

// UpdateImgTags updates all <img/> tags in a webpage to have an alt="" attribute,
// if they don't already.
// Doesn't update <img></img> tags.
func (wp *WebPage) UpdateImgTags(client *api.Client) error {
	if err := wp.open(); err != nil {
		return err
	}

	defer wp.close()

	tokenizer := html.NewTokenizer(wp.file)

	for {
		switch tokenType := tokenizer.Next(); tokenType {

		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return wp.saveContent()
			}
			return tokenizer.Err()

		case html.TextToken, html.StartTagToken, html.EndTagToken, html.DoctypeToken:
			t := tokenizer.Token()
			fmt.Fprintf(wp.content, t.String())

		case html.SelfClosingTagToken:
			t := tokenizer.Token()

			if t.DataAtom == atom.Img || t.DataAtom == atom.Image {
				imgTag := img.New(&t, wp.rootDir, client)
				fmt.Fprintf(wp.content, imgTag.String())
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
