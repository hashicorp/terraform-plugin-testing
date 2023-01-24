package plugintest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"regexp"

	tfjson "github.com/hashicorp/terraform-json"
)

type TerraformJSONDiagnostics []tfjson.Diagnostic

func (d TerraformJSONDiagnostics) Contains(r *regexp.Regexp, severity tfjson.DiagnosticSeverity) bool {
	for _, v := range d {
		if v.Severity != severity {
			continue
		}

		if !r.MatchString(v.Summary) && !r.MatchString(v.Detail) {
			continue
		}

		return true
	}

	return false
}

var _ io.Writer = &TerraformJSONBuffer{}
var _ io.Reader = &TerraformJSONBuffer{}

// TerraformJSONBuffer is used for storing and processing streaming
// JSON output generated when running terraform commands with the
// `-json` flag.
type TerraformJSONBuffer struct {
	buf         *bytes.Buffer
	diagnostics TerraformJSONDiagnostics
	jsonOutput  []string
	rawOutput   string
	parsed      bool
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

func (b *TerraformJSONBuffer) Parse() {
	if b.buf == nil {
		log.Fatal("call NewTerraformJSONBuffer to initialise buffer")
	}

	scanner := bufio.NewScanner(b.buf)

	for scanner.Scan() {
		txt := scanner.Text()

		b.rawOutput += "\n" + txt

		var outer struct {
			Diagnostic tfjson.Diagnostic
		}

		// This will only capture buffer entries that can be unmarshalled
		// as JSON. If there are entries in the buffer that are not JSON
		// they will be discarded.
		if json.Unmarshal([]byte(txt), &outer) == nil {
			b.jsonOutput = append(b.jsonOutput, txt)

			if outer.Diagnostic.Severity != "" {
				b.diagnostics = append(b.diagnostics, outer.Diagnostic)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	b.parsed = true
}

func (b *TerraformJSONBuffer) Diagnostics() TerraformJSONDiagnostics {
	if !b.parsed {
		b.Parse()
	}

	return b.diagnostics
}

func (b *TerraformJSONBuffer) JsonOutput() []string {
	if !b.parsed {
		b.Parse()
	}

	return b.jsonOutput
}

func (b *TerraformJSONBuffer) RawOutput() string {
	if !b.parsed {
		b.Parse()
	}

	return b.rawOutput
}
