// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package resource

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	tfjson "github.com/hashicorp/terraform-json"

	"github.com/hashicorp/terraform-plugin-testing/internal/plugintest"
)

const stdoutFile = "stdout.txt"

func unmarshalJSON(data []byte, v interface{}) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	return dec.Decode(v)
}

func stdoutWriter(wd *plugintest.WorkingDir) (io.WriteCloser, error) {
	return os.Create(filepath.Join(wd.GetBaseDir(), stdoutFile))
}

func stdoutReader(wd *plugintest.WorkingDir) (io.ReadCloser, error) {
	return os.Open(filepath.Join(wd.GetBaseDir(), stdoutFile))
}

func getJSONOutput(wd *plugintest.WorkingDir) []string {
	reader, err := stdoutReader(wd)
	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()

	scanner := bufio.NewScanner(reader)
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

func getJSONOutputStr(wd *plugintest.WorkingDir) string {
	return strings.Join(getJSONOutput(wd), "\n")
}

func diagnosticFound(wd *plugintest.WorkingDir, r *regexp.Regexp, severity tfjson.DiagnosticSeverity) (bool, []string) {
	reader, err := stdoutReader(wd)
	if err != nil {
		log.Fatal(err)
	}

	defer reader.Close()

	scanner := bufio.NewScanner(reader)
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
