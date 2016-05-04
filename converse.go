package witty

import (
	"fmt"
	"log"
	"net/url"
)

type Context map[string]interface{}

type Entities map[string][]interface{}

type BotNextStep struct {
	Type       string
	Msg        string
	Action     string
	Entities   Entities
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

func (s *chatService) RunActions(sessID, msg string, ctx Context) Context {
	step, err := s.Converse(sessID, msg, ctx)
	if err != nil {
		log.Fatal(err)
	}

	switch step.Type {
	case "stop":
		return ctx

	case "msg":
		log.Printf("Executing say %q", step.Msg)
		s.client.SayAct(sessID, ctx, step.Msg)

	case "merge":
		log.Print("Executing merge")
		ctx = s.client.MergeAct(sessID, ctx, step.Entities)

	case "action":
		log.Printf("Executing action %q", step.Action)
		if action, ok := s.client.Actions[step.Action]; ok {
			ctx = action(sessID, ctx)
		}

	case "error":
		log.Print("Executing error")
		s.client.ErrorAct(sessID, ctx, "oops")

	default:
		log.Print("Unknown type")
	}

	return s.RunActions(sessID, "", ctx)
}

func DefaultSayAct(sessID string, ctx Context, msg string) {
	fmt.Printf("> %v\n", msg)
}

func DefaultMergeAct(sessID string, ctx Context, entities Entities) Context {
	return ctx
}

func DefaultErrorAct(sessID string, ctx Context, msg string) {
	fmt.Printf("Bot error: %v", msg)
}
