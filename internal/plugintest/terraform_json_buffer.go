package plugintest

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

var _ io.Writer = &TerraformJSONBuffer{}
var _ io.Reader = &TerraformJSONBuffer{}

type TerraformJSONBuffer struct {
	buf *bytes.Buffer
}

func NewTerraformJSONBuffer() *TerraformJSONBuffer {
	return &TerraformJSONBuffer{
		buf: new(bytes.Buffer),
	}
}

func (b *TerraformJSONBuffer) Write(p []byte) (n int, err error) {
	if b.buf == nil {
		log.Fatal("call NewTerraformJSONBuffer to initialise buffer")
	}

	return b.buf.Write(p)
}

func (b *TerraformJSONBuffer) Read(p []byte) (n int, err error) {
	if b.buf == nil {
		log.Fatal("call NewTerraformJSONBuffer to initialise buffer")
	}

	return b.buf.Read(p)
}

func (b *TerraformJSONBuffer) GetJSONOutput() []string {
	if b.buf == nil {
		log.Fatal("call NewTerraformJSONBuffer to initialise buffer")
	}

	scanner := bufio.NewScanner(b.buf)
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

func (b *TerraformJSONBuffer) GetJSONOutputStr() string {
	return strings.Join(b.GetJSONOutput(), "\n")
}

func (b *TerraformJSONBuffer) DiagnosticFound(r *regexp.Regexp, severity tfjson.DiagnosticSeverity) (bool, []string) {
	if b.buf == nil {
		log.Fatal("call NewTerraformJSONBuffer to initialise buffer")
	}

	scanner := bufio.NewScanner(b.buf)
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
