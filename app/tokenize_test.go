package app

import (
	"errors"
	"testing"

	"github.com/spaceavocado/apidoc/extract"
	"github.com/spaceavocado/apidoc/token"
)

type mockParser struct {
	returns int
	tokens  [][]token.Token
	err     []error
}

func (p *mockParser) Parse(b extract.Block) ([]token.Token, error) {
	p.returns++
	return p.tokens[p.returns], p.err[p.returns]
}

func TestTokenize(t *testing.T) {
	// Valid
	a := New(Configuration{})
	_, err := a.Tokenize(ExtractResult{
		Main: extract.Block{
			Lines: []string{
				"title Sample API",
			},
		},
		Endpoints: []extract.Block{
			{
				Lines: []string{
					"summary Sample endpoint",
				},
			},
		},
	})
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	// Errors
	a.tokenParser = &mockParser{
		returns: -1,
		tokens:  [][]token.Token{make([]token.Token, 0)},
		err:     []error{errors.New("simulated error")},
	}
	_, err = a.Tokenize(ExtractResult{})
	if err == nil {
		t.Errorf("Expected error, got nil")
		return
	}

	a.tokenParser = &mockParser{
		returns: -1,
		tokens: [][]token.Token{
			make([]token.Token, 1),
			make([]token.Token, 0),
		},
		err: []error{
			nil,
			errors.New("simulated error"),
		},
	}
	_, err = a.Tokenize(ExtractResult{
		Endpoints: make([]extract.Block, 1),
	})
	if err != nil {
		t.Errorf("Expected error, got nil")
		return
	}
}
