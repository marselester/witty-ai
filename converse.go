package witty

import (
	"fmt"
	"log"
	"net/url"
)

type Context map[string]interface{}

type Entities map[string][]map[string]interface{}

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

func (s *chatService) RunActions(sessID, msg string, ctx Context, maxSteps int) (Context, error) {
	if maxSteps <= 0 {
		return ctx, ErrMaxSteps
	}

	step, err := s.Converse(sessID, msg, ctx)
	if err != nil {
		return ctx, err
	}

	switch step.Type {
	case "stop":
		log.Print("Executing stop")
		return ctx, nil

	case "msg":
		log.Printf("Executing say %q", step.Msg)
		s.client.SayAct(sessID, ctx, step.Msg)

	case "merge":
		log.Print("Executing merge")
		ctx = s.client.MergeAct(sessID, ctx, step.Entities)

	case "action":
		if action, ok := s.client.Actions[step.Action]; ok {
			log.Printf("Executing action %q", step.Action)
			ctx = action(sessID, ctx)
		} else {
			log.Printf("Executing action %q: not found", step.Action)
		}

	case "error":
		log.Print("Executing error")
		return ctx, ErrWitStep

	default:
		log.Printf("Unknown type %q", step.Type)
		return ctx, ErrUnkownStep
	}

	return s.RunActions(sessID, "", ctx, maxSteps-1)
}

func DefaultSayAct(sessID string, ctx Context, msg string) {
	fmt.Printf("💬  %v\n", msg)
}

func DefaultMergeAct(sessID string, ctx Context, entities Entities) Context {
	return ctx
}
