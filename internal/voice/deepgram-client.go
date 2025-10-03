package voice

import (
	"context"
	"errors"
	"mime/multipart"
	"os"

	api "github.com/deepgram/deepgram-go-sdk/v3/pkg/api/listen/v1/rest"
	interfaces "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/interfaces"
	client "github.com/deepgram/deepgram-go-sdk/v3/pkg/client/listen"
)

type Transcriber interface {
	Transcribe(ctx context.Context, file multipart.File) (string, error)
}

type deepgramClient struct {
	client *api.Client
}

func NewdeepgramClient() (Transcriber, error) {
	apiKey := os.Getenv("DEEPGRAM_API_KEY")
	if apiKey == "" {
		return nil, errors.New("api_key is missing")
	}

	restClient := client.NewREST(apiKey, &interfaces.ClientOptions{})
	dgClient := api.New(restClient)

	return &deepgramClient{client: dgClient}, nil
}

func (c *deepgramClient) Transcribe(ctx context.Context, file multipart.File) (string, error) {
	defer file.Close()

	opts := &interfaces.PreRecordedTranscriptionOptions{
		Model:       "nova-2",
		SmartFormat: true,
	}

	resp, err := c.client.FromStream(ctx, file, opts)
	if err != nil {
		return "", err
	}
	if resp == nil || resp.Results == nil || len(resp.Results.Channels) == 0 ||
		len(resp.Results.Channels[0].Alternatives) == 0 {
		return "", errors.New("deepgram: empty transcript response")
	}

	return resp.Results.Channels[0].Alternatives[0].Transcript, nil
}
