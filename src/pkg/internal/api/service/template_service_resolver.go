package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/antoniofrisenda/template-service/src/pkg/internal/assets/model"
)

type IResolver interface {
	ExtractVariables(ctx context.Context, text string, contentType string) ([]string, error)
	ParseFileContent(ctx context.Context, input []byte) (string, error)
}

type Resolver struct{}

func NewResolver() IResolver {
	return &Resolver{}
}

func (r *Resolver) ExtractVariables(ctx context.Context, text string, contentType string) ([]string, error) {
	_ = ctx

	if text == "" {
		return []string{}, nil
	}

	normalizedType := strings.ToUpper(strings.TrimSpace(contentType))
	switch normalizedType {
	case string(model.HTML), string(model.PLAIN_TEXT):
		matches := model.VarRegex.FindAllString(text, -1)
		if len(matches) == 0 {
			return []string{}, nil
		}

		result := make([]string, 0, len(matches))
		seen := make(map[string]struct{}, len(matches))
		for _, match := range matches {
			if _, ok := seen[match]; ok {
				continue
			}
			seen[match] = struct{}{}
			result = append(result, match)
		}

		return result, nil
	case string(model.PDF):
		panic("unimplemented")
	default:
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}
}

func (r *Resolver) ParseFileContent(ctx context.Context, input []byte) (string, error) {
	_ = ctx
	return string(input), nil
}
