package api

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/services/cognitiveservices/v2.0/computervision"
	"github.com/Azure/go-autorest/autorest"
)

// Client object to interact with Azure Computer Vision services cleanly
type Client struct {
	visionClient  computervision.BaseClient
	visionContext context.Context
}

// New Client object
func New() *Client {
	computerVisionKey := os.Getenv("COMPUTER_VISION_SUBSCRIPTION_KEY")
	if computerVisionKey == "" {
		log.Fatal("Need a computer vision API key!")
	}

	endpointURL := os.Getenv("COMPUTER_VISION_ENDPOINT")
	if endpointURL == "" {
		log.Fatal("Need a computer vision endpoint!")
	}

	a := Client{visionContext: context.Background(), visionClient: computervision.New(endpointURL)}

	a.visionClient.Authorizer = autorest.NewCognitiveServicesAuthorizer(computerVisionKey)

	return &a
}

// Describe a local image with the highest confidence guess.
func (c *Client) Describe(localImagePath string) (*computervision.ImageCaption, error) {
	fmt.Printf("Trying to describe %s\n", localImagePath)

	var localImage io.ReadCloser
	localImage, err := os.Open(localImagePath)

	if err != nil {
		return nil, err
	}

	maxNumberDescriptionCandidates := new(int32)
	*maxNumberDescriptionCandidates = 1

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
		return nil, &NoLabelError{localImagePath}
	}

	imageCaption := (*localImageDescription.Captions)[0]

	if *imageCaption.Confidence < 0.7 { // @TODO
		return &imageCaption, &ConfidenceError{*imageCaption.Confidence}
	}

	return &imageCaption, nil
}
