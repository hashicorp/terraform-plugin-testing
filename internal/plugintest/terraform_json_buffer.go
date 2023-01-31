package plugintest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (d TerraformJSONDiagnostics) Errors() TerraformJSONDiagnostics {
	var tfJSONDiagnostics TerraformJSONDiagnostics

	for _, v := range d {
		if v.Severity == tfjson.DiagnosticSeverityError {
			tfJSONDiagnostics = append(tfJSONDiagnostics, v)
		}
	}

	return tfJSONDiagnostics
}

func (d TerraformJSONDiagnostics) Warnings() TerraformJSONDiagnostics {
	var tfJSONDiagnostics TerraformJSONDiagnostics

	for _, v := range d {
		if v.Severity == tfjson.DiagnosticSeverityWarning {
			tfJSONDiagnostics = append(tfJSONDiagnostics, v)
		}
	}

	return tfJSONDiagnostics
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
		return 0, fmt.Errorf("cannot write to uninitialized buffer, use NewTerraformJSONBuffer")
	}

	return b.buf.Write(p)
}

func (b *TerraformJSONBuffer) Read(p []byte) (n int, err error) {
	if b.buf == nil {
		return 0, fmt.Errorf("cannot write to uninitialized buffer, use NewTerraformJSONBuffer")
	}

	return b.buf.Read(p)
}

func read(r *bufio.Reader) ([]byte, error) {
	var (
		isPrefix = true
		err      error
		line, ln []byte
	)

	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}

	return ln, err
}

func (b *TerraformJSONBuffer) Parse() error {
	if b.buf == nil {
		return fmt.Errorf("cannot write to uninitialized buffer, use NewTerraformJSONBuffer")
	}

	reader := bufio.NewReader(b.buf)

	for {
		line, err := read(reader)
		if err != nil {
			if err == io.EOF {
				break
			}

			return fmt.Errorf("cannot read line: %s", err)
		}

		txt := string(line)

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

	b.parsed = true

	return nil
}

func (b *TerraformJSONBuffer) Diagnostics() (TerraformJSONDiagnostics, error) {
	if !b.parsed {
		err := b.Parse()
		if err != nil {
			return nil, err
		}
	}

	return b.diagnostics, nil
}

func (b *TerraformJSONBuffer) JsonOutput() ([]string, error) {
	if !b.parsed {
		err := b.Parse()
		if err != nil {
			return nil, err
		}
	}

	return b.jsonOutput, nil
}

func (b *TerraformJSONBuffer) RawOutput() (string, error) {
	if !b.parsed {
		err := b.Parse()
		if err != nil {
			return "", err
		}
	}

	return b.rawOutput, nil
}
