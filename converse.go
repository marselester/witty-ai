package witty

import (
	"net/url"
)

type Context map[string]interface{}

type BotNextStep struct {
	Type       string
	Msg        string
	Action     string
	Entities   map[string]interface{}
	Confidence float64
}

// chatService handles communication with converse API resource.
type chatService struct {
	client *Client
}

// Converse gets your bot's next step.
func (s *chatService) Converse(sessID, msg string, ctx Context) (*BotNextStep, error) {
	params := &url.Values{}
	params.Set("session_id", sessID)
	if msg != "" {
		params.Set("q", msg)
	}

	req, err := s.client.NewRequest("POST", "converse", params, ctx)
	if err != nil {
		return nil, err
	}

	v := new(BotNextStep)
	_, err = s.client.Do(req, v)
	if err != nil {
		return nil, err
	}

	return v, nil
}
