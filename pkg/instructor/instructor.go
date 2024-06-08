package instructor

import (
	"context"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type Instructor interface {
	Provider() Provider
	Mode() Mode
	MaxRetries() int
	WithValidator() bool

	// Chat / Messages

	chat(
		ctx context.Context,
		request interface{},
		schema *Schema,
	) (string, interface{}, error)

	// Streaming Chat / Messages

	chatStream(
		ctx context.Context,
		request interface{},
		schema *Schema,
	) (<-chan string, error)
}
