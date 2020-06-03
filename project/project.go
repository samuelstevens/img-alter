package project

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/samuelstevens/goimglabeler/api"
	"github.com/samuelstevens/goimglabeler/caption"
	"github.com/samuelstevens/goimglabeler/webpage"
)

type Project struct {
	dir    string
	client *api.Client
}

func New(dir string, client *api.Client) *Project {
	err := caption.LoadCaptions()

	if err != nil {
		log.Printf("couldn't load captions: %s", err.Error())
	}

	project := Project{dir: dir, client: client}

	return &project
}

// UpdateImages walks all .html files in the project directory and
func (p *Project) UpdateImages() error {

	err := filepath.Walk(p.dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}

		if filepath.Ext(path) != ".html" {
			return nil
		}

		page, err := webpage.New(path, p.dir)

		if err != nil {
			return err
		}

		err = page.UpdateImgTags(p.client)

		if err != nil {
			log.Fatal(err.Error())
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
