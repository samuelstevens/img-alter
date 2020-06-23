package api

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/computervision"
	"github.com/Azure/go-autorest/autorest"
)

// Client object to interact with Azure Computer Vision services cleanly
type Client struct {
	visionClient  computervision.BaseClient
	visionContext context.Context
	threshold     float64
	loud          bool
}

var ErrorAuth = errors.New("no key or endpoint")

// New Client object
func New(key string, endpoint string, threshold float64, loud bool) (*Client, error) {
	if key == "" || endpoint == "" {
		return nil, ErrorAuth
	}

	computerVisionKey := key

	endpointURL := endpoint

	client := Client{
		visionContext: context.Background(),
		visionClient:  computervision.New(endpointURL),
		threshold:     threshold,
		loud:          loud,
	}

	client.visionClient.Authorizer = autorest.NewCognitiveServicesAuthorizer(computerVisionKey)

	return &client, nil
}

// Describe a local image with the highest confidence guess.
func (c *Client) Describe(localImagePath string) (*computervision.ImageCaption, error) {

	if c.loud {
		fmt.Printf("Trying to describe %s\n", localImagePath)
	}
	var localImage io.ReadCloser
	localImage, err := os.Open(localImagePath)

	if err != nil {
		return nil, err
	}

	maxNumberDescriptionCandidates := new(int32)
	*maxNumberDescriptionCandidates = 1

	// @TODO: check file size

	localImageDescription, err := c.visionClient.DescribeImageInStream(
		c.visionContext,
		localImage,
		maxNumberDescriptionCandidates,
		"", // language
	)

	if err != nil {
		return nil, err
	}

	if len(*localImageDescription.Captions) == 0 {
		return nil, ErrorNoLabel
	}

	imageCaption := (*localImageDescription.Captions)[0]

	if *imageCaption.Confidence < c.threshold {
		return &imageCaption, &ConfidenceError{*imageCaption.Confidence}
	}

	return &imageCaption, nil
}
