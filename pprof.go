package openai

import (
	"context"
	"os"
	"runtime/pprof"
)

var tokenUsageProfile = pprof.NewProfile("openai.token.usage")

type PprofClient struct {
	*Client
}

// NewPprofClient creates new OpenAI API client that enables pprof.
func NewPprofClient(authToken string) *PprofClient {
	config := DefaultConfig(authToken)
	return &PprofClient{
		Client: NewClientWithConfig(config),
	}
}

// NewPprofClientWithConfig creates new OpenAI API client for specified config.
// It also enables pprof.
func NewPprofClientWithConfig(config ClientConfig) *PprofClient {
	return &PprofClient{
		Client: &Client{
			config:         config,
			requestBuilder: newRequestBuilder(),
		},
	}
}

func (c *PprofClient) WriteToFile(profileOutPath string) error {
	out, err := os.Create(profileOutPath)
	if err != nil {
		return err
	}
	if err := tokenUsageProfile.WriteTo(out, 0); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}

func (c *PprofClient) CreateCompletion(
	ctx context.Context,
	request CompletionRequest,
) (response CompletionResponse, err error) {
	response, err = c.Client.CreateCompletion(ctx, request)
	if err != nil {
		return
	}
	tokenUsageProfile.Add(response.Usage.TotalTokens, 2)
	return
}

func (c *PprofClient) CreateChatCompletion(
	ctx context.Context,
	request ChatCompletionRequest,
) (response ChatCompletionResponse, err error) {
	response, err = c.Client.CreateChatCompletion(ctx, request)
	if err != nil {
		return
	}
	tokenUsageProfile.Add(response.Usage.TotalTokens, 2)
	return
}
