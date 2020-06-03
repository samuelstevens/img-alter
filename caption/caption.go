package caption

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/samuelstevens/goimglabeler/api"
	"github.com/samuelstevens/goimglabeler/util"
)

const (
	captionFileDir  = "."
	captionFileName = "captions.json"
)

type Caption struct {
	Hash        string
	Filename    string
	Description string
	Confidence  float64
}

type Captions struct {
	lookup   map[string]*Caption
	filepath string
}

var captionLookup = Captions{filepath: filepath.Join(captionFileDir, captionFileName)}

func New(imgPath string, prevDescription string, client *api.Client) (*Caption, error) {
	defaultCaption := Caption{Description: prevDescription}

	hash, err := util.HashFile(imgPath)

	if err != nil {
		return &defaultCaption, err
	}

	caption, ok := captionLookup.get(hash)

	if ok {
		return caption, nil
	}

	description := prevDescription
	confidence := 1.0

	if description == "" {
		confidence = 0.0

		imageCaption, err := client.Describe(imgPath)

		if err != nil {
			if _, ok := err.(*api.ConfidenceError); ok {
				description = fmt.Sprintf("Possibly inaccurate: %s", *imageCaption.Text)
			} else {
				return &defaultCaption, err // only if not a confidence error
			}

		} else {
			description = *imageCaption.Text
		}

		confidence = *imageCaption.Confidence
	}

	c := Caption{
		Hash:        hash,
		Filename:    filepath.Base(imgPath),
		Description: description,
		Confidence:  confidence,
	}

	return &c, captionLookup.set(&c)
}

func (c *Captions) save() error {
	jsonRep, err := json.Marshal(c)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(c.filepath, jsonRep, 0644)

}

func LoadCaptions() error {
	jsonRep, err := ioutil.ReadFile(captionLookup.filepath)

	if err != nil {
		return err
	}

	if len(jsonRep) == 0 {
		return nil
	}

	err = json.Unmarshal(jsonRep, &captionLookup.lookup)

	if err != nil {
		return err
	}

	return nil
}

func (c *Captions) set(caption *Caption) error {
	if caption.Description == "" {
		return nil
	}

	c.lookup[caption.Hash] = caption

	return c.save()
}

func (c *Captions) get(hash string) (*Caption, bool) {
	caption, ok := c.lookup[hash]

	return caption, ok
}
