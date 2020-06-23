package main

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	files "path/filepath"

	"github.com/samuelstevens/gocaption/api"
	"github.com/samuelstevens/gocaption/caption"
	"github.com/samuelstevens/gocaption/cli"
	"github.com/samuelstevens/gocaption/webpage"
)

type fileType int

const (
	image fileType = iota
	html
	unknown
)

func getFileType(filepath string) fileType {
	switch files.Ext(filepath) {
	case ".jpg", ".png", ".jpeg":
		return image
	case ".html", ".htm":
		return html
	default:
		return unknown
	}
}

func displayCaption(path string, description string, opts *cli.Options) {
	if !opts.Silent {
		fmt.Printf("%s\t%s\n", filepath.Base(path), description)
	}
}

func displayError(path string, err error) {
	log.Printf("Can't caption %s; %s.\n", filepath.Base(path), err.Error())
}

func captionHTML(filepath string, opts *cli.Options, client *api.Client) {
	page, err := webpage.New(filepath)

	if err != nil {
		displayError(filepath, err)
		return
	}

	err = page.LabelImages(client)

	if err != nil {
		displayError(filepath, err)
		return
	}

	if opts.Write {
		err = page.Write()
		if err != nil {
			fmt.Printf("Couldn't update file: %s.\n", err.Error())
		}
	}

	for _, caption := range page.Captions {
		displayCaption(caption.FilePath, caption.Description, opts)
	}
}

func main() {
	opts := cli.Cli()

	if len(opts.Files) == 0 {
		fmt.Println("Please supply file(s) or directory.")
		return
	}

	caption.InitializeCache(opts.CacheFile)

	client, err := api.New(opts.APIKey, opts.Endpoint, opts.Threshold, opts.Loud)

	if err != nil {
		if errors.Is(err, api.ErrorAuth) {
			// provide some help in the form of missing config

			fmt.Println("You need to specify some keys for MS Azure. You can also specify a config file with --config.")
		}
		log.Fatal(err.Error())
	}

	for _, filepath := range opts.Files {
		switch getFileType(filepath) {
		case image:
			caption, err := caption.New(filepath, "", client)

			if err != nil {
				displayError(filepath, err)
				continue
			}

			displayCaption(filepath, caption.Description, opts)

		case html:
			captionHTML(filepath, opts, client)

		case unknown:

		default:
			log.Fatalf("Unreachable code.\n")
		}
	}

}
