package resource

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"regexp"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"
)

var _ io.Writer = &Stdout{}
var _ io.Reader = &Stdout{}

type Stdout struct {
	buf *bytes.Buffer
}

func NewStdout() *Stdout {
	return &Stdout{
		buf: new(bytes.Buffer),
	}
}

func (s *Stdout) Write(p []byte) (n int, err error) {
	if s.buf == nil {
		log.Fatal("call NewStdout to initialise buffer")
	}

	return s.buf.Write(p)
}

func (s *Stdout) Read(p []byte) (n int, err error) {
	if s.buf == nil {
		log.Fatal("call NewStdout to initialise buffer")
	}

	return s.buf.Read(p)
}

func (s *Stdout) GetJSONOutput() []string {
	if s.buf == nil {
		log.Fatal("call NewStdout to initialise buffer")
	}

	scanner := bufio.NewScanner(s.buf)
	var jsonOutput []string

	for scanner.Scan() {
		var outer struct{}

		txt := scanner.Text()

		if json.Unmarshal([]byte(txt), &outer) == nil {
			jsonOutput = append(jsonOutput, txt)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return jsonOutput
}

func (s *Stdout) GetJSONOutputStr() string {
	return strings.Join(s.GetJSONOutput(), "\n")
}

func (s *Stdout) DiagnosticFound(r *regexp.Regexp, severity tfjson.DiagnosticSeverity) (bool, []string) {
	if s.buf == nil {
		log.Fatal("call NewStdout to initialise buffer")
	}

	scanner := bufio.NewScanner(s.buf)
	var jsonOutput []string

	for scanner.Scan() {
		var outer struct {
			Diagnostic tfjson.Diagnostic
		}

		txt := scanner.Text()

		if json.Unmarshal([]byte(txt), &outer) == nil {
			jsonOutput = append(jsonOutput, txt)

			if outer.Diagnostic.Severity == "" {
				continue
			}

			if !r.MatchString(outer.Diagnostic.Summary) && !r.MatchString(outer.Diagnostic.Detail) {
				continue
			}

			if outer.Diagnostic.Severity != severity {
				continue
			}

			return true, jsonOutput
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return false, jsonOutput
}
