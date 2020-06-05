package caption

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/samuelstevens/gocaption/api"
	"github.com/samuelstevens/gocaption/util"
)

const (
	captionFileDir  = "."
	captionFileName = "captions.json"
)

// Caption is a caption and confidence for a file
type Caption struct {
	hash        string
	FilePath    string
	Description string
	Confidence  float64
}

type captions struct {
	lookup   map[string]*Caption
	filepath string
}

var captionCache *captions

// New returns a new caption for an image.
func New(imgPath string, prevDescription string, client *api.Client) (*Caption, error) {

	defaultCaption := Caption{Description: prevDescription}

	hash, err := util.HashFile(imgPath)

	if err != nil {
		return &defaultCaption, err
	}

	caption, ok := captionCache.get(hash)

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
		hash:        hash,
		FilePath:    filepath.Base(imgPath),
		Description: description,
		Confidence:  confidence,
	}

	return &c, captionCache.set(&c)
}

func (c *captions) save() error {
	jsonRep, err := json.MarshalIndent(c.lookup, "", "\t")

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(c.filepath, jsonRep, 0666)

	if err != nil {
		fmt.Println("Saving error")
		return err
	}

	return nil

}

func InitializeCache(cacheFilepath string) {
	if captionCache != nil {
		return
	}

	captionCache = &captions{filepath: cacheFilepath}

	captionCache.lookup = loadLookup(captionCache.filepath)
}

func loadLookup(filepath string) map[string]*Caption {
	defaultMap := map[string]*Caption{}

	var res map[string]*Caption

	jsonRep, err := ioutil.ReadFile(filepath)

	if err != nil {
		return defaultMap
	}

	if len(jsonRep) == 0 {
		return defaultMap
	}

	err = json.Unmarshal(jsonRep, &res)

	if err != nil {
		return defaultMap
	}

	return res
}

func (c *captions) set(caption *Caption) error {
	if caption.Description == "" {
		return nil
	}

	c.lookup[caption.hash] = caption

	return c.save()
}

func (c *captions) get(hash string) (*Caption, bool) {
	caption, ok := c.lookup[hash]

	return caption, ok
}
