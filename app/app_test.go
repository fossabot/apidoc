package app

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/spaceavocado/apidoc/extract"
	"github.com/spaceavocado/apidoc/token"
)

type errorResolver struct{}

func (r *errorResolver) Resolve(endpoints []extract.Block) error {
	return errors.New("simulated error")
}

type errorGenerator struct{}

func (g *errorGenerator) Generate(main []token.Token, endpoints [][]token.Token, file string) error {
	return errors.New("simulated error")
}

func TestStart(t *testing.T) {
	content := []string{
		`
		// @summary Refresh ID Token
		// @desc Use the refresh token
		// to receive a new ID token.
		// It must be in a valid format.
		`,
		"func main(){}",
		"func main(){}",
	}
	files := []string{
		"tmp1",
		"tmp2",
		"tmp/tmp.go",
	}
	err := os.MkdirAll("tmp", os.ModePerm)
	if err != nil {
		t.Errorf("Unexpected error %v", err)
		return
	}

	for i, f := range files {
		err := ioutil.WriteFile(f, []byte(content[i]), 0644)
		if err != nil {
			t.Errorf("Unexpected error %v", err)
			return
		}
	}
	defer func() {
		for _, f := range files {
			os.Remove(f)
		}
		os.Remove("tmp")
	}()

	b := &bytes.Buffer{}
	log.SetOutput(b)
	hook := test.NewGlobal()
	a := New(Configuration{
		MainFile: "",
		EndsRoot: "tmp",
		Output:   "tmp/output",
	})
	a.Start()
	if len(hook.Entries) != 1 {
		t.Errorf("Expected %d log entries, got %d", 1, len(hook.Entries))
	}

	hook.Reset()
	a = New(Configuration{
		MainFile: "tmp1",
		EndsRoot: "tmp/",
		Output:   "tmp/output",
	})
	a.refResolver = &errorResolver{}
	a.Start()
	if len(hook.Entries) != 1 {
		t.Errorf("Expected %d log entries, got %d", 1, len(hook.Entries))
	}

	hook.Reset()
	a = New(Configuration{
		MainFile: "tmp1",
		EndsRoot: "tmp/",
		Output:   "tmp/output",
	})
	a.tokenParser = &mockParser{
		returns: -1,
		tokens:  [][]token.Token{make([]token.Token, 0)},
		err:     []error{errors.New("simulated error")},
	}
	a.Start()
	if len(hook.Entries) != 1 {
		t.Errorf("Expected %d log entries, got %d", 1, len(hook.Entries))
	}

	hook.Reset()
	a = New(Configuration{
		MainFile: "tmp1",
		EndsRoot: "tmp/",
		Output:   "tmp/output",
	})
	a.generator = &errorGenerator{}
	a.Start()
	if len(hook.Entries) != 1 {
		t.Errorf("Expected %d log entries, got %d", 1, len(hook.Entries))
	}

}
