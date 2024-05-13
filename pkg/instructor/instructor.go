package instructor

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	anthropic "github.com/liushuangls/go-anthropic/v2"
	openai "github.com/sashabaranov/go-openai"
)

type Instructor[T any] struct {
	Client     Client[T]
	Provider   Provider
	Mode       Mode
	MaxRetries int

	Schema *Schema[T]
	Type   reflect.Type
}

func FromOpenAI[T any](client *openai.Client, opts ...Options) (*Instructor[T], error) {

	options := mergeOptions(opts...)

	schema, err := NewSchema[T]()
	if err != nil {
		return nil, err
	}

	cli, err := NewOpenAIClient(client, schema, *options.Mode)
	if err != nil {
		return nil, err
	}

	t := reflect.TypeOf(new(T))

	i := &Instructor[T]{
		Client:     cli,
		Provider:   OpenAI,
		Mode:       *options.Mode,
		MaxRetries: *options.MaxRetries,
		Schema:     schema,
		Type:       t,
	}
	return i, nil
}

func FromAnthropic[T any](client *anthropic.Client, opts ...Options) (*Instructor[T], error) {

	options := mergeOptions(opts...)

	schema, err := NewSchema[T]()
	if err != nil {
		return nil, err
	}

	cli, err := NewAnthropicClient(client, schema, *options.Mode)
	if err != nil {
		return nil, err
	}

	t := reflect.TypeOf(new(T))

	i := &Instructor[T]{
		Client:     cli,
		Provider:   OpenAI,
		Mode:       *options.Mode,
		MaxRetries: *options.MaxRetries,
		Schema:     schema,
		Type:       t,
	}
	return i, nil
}

func (i *Instructor[T]) CreateChatCompletion(ctx context.Context, request Request) (*T, error) {

	for attempt := 0; attempt < i.MaxRetries; attempt++ {

		text, err := i.Client.CreateChatCompletion(ctx, request)
		if err != nil {
			// no retry on non-marshalling/validation errors
			println(text)
			println(err.Error())
			return nil, err
		}

		t, err := i.processResponse(text)
		if err != nil {
			// TODO:
			// add more sophisticated retry logic (send back json and parse error for model to fix).
			//
			// Currently, its just recalling with no new information
			// or attempt to fix the error with the last generated JSON
			println(text)
			println(err.Error())
			continue
		}

		return t, nil
	}

	return nil, errors.New("hit max retry attempts")
}

func (i *Instructor[T]) processResponse(response string) (*T, error) {

	t := new(T)

	err := json.Unmarshal([]byte(response), t)
	if err != nil {
		return nil, err
	}

	// TODO: if direct unmarshal fails: check common erors like wrapping struct with key name of struct, instead of just the value

	return t, nil
}