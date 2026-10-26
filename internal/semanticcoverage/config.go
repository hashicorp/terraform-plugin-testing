// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// Config reads the files actually copied into the test working directory.
func (c *Collector) Config(dir string) {
	if c == nil {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		c.Warning(WarningConfigUnavailable)
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if !strings.HasSuffix(name, ".tf") && !strings.HasSuffix(name, ".tf.json") && !strings.HasSuffix(name, ".tfquery.hcl") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			c.Warning(WarningConfigUnavailable)
			continue
		}
		c.Parse(name, data)
	}
}

func (c *Collector) QueryConfig(dir string) {
	if c == nil {
		return
	}
	data, err := os.ReadFile(filepath.Join(dir, "terraform_plugin_test.tfquery.hcl"))
	if err != nil {
		c.Warning(WarningConfigUnavailable)
		return
	}
	c.Parse("terraform_plugin_test.tfquery.hcl", data)
}

func (c *Collector) Parse(name string, data []byte) {
	if c == nil {
		return
	}
	if strings.HasSuffix(name, ".json") {
		var object map[string]any
		if json.Unmarshal(data, &object) != nil {
			c.Warning(WarningConfigUnavailable)
			return
		}
		if c.active {
			c.jsonBlocks(object)
		}
		return
	}
	parser := hclparse.NewParser()
	file, diags := parser.ParseHCL(data, name)
	if diags.HasErrors() {
		c.Warning(WarningConfigUnavailable)
		return
	}
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok {
		c.Warning(WarningConfigUnavailable)
		return
	}
	if c.active {
		c.blocks(body.Blocks)
		c.functions(body)
	}
}

func (c *Collector) blocks(blocks []*hclsyntax.Block) {
	for _, b := range blocks {
		kind, name, address := "", "", ""
		switch b.Type {
		case "provider":
			if len(b.Labels) > 0 && !strings.HasSuffix(b.DefRange().Filename, ".tfquery.hcl") {
				kind, name = "provider", b.Labels[0]
			}
		case "resource", "data", "ephemeral", "action", "list":
			if len(b.Labels) > 1 {
				kind, name, address = b.Type, b.Labels[0], b.Labels[0]+"."+b.Labels[1]
				switch kind {
				case "data":
					kind, address = "data_source", "data."+address
				case "list":
					kind, address = "list_resource", "list."+address
				case "ephemeral":
					kind, address = "ephemeral_resource", "ephemeral."+address
				case "action":
					address = "action." + address
				}
			}
		case "state_store":
			if len(b.Labels) > 0 {
				kind, name = "state_store", b.Labels[0]
			}
		}
		if kind != "" {
			f := c.feature(kind, name, address)
			switch kind {
			case "action", "list_resource", "state_store":
				// Schema attributes live in a nested config block.
				for _, nested := range b.Body.Blocks {
					if nested.Type == "config" {
						walkInputs(f, nested.Body, "")
					}
				}
			default:
				walkInputs(f, b.Body, "")
			}
		}
		c.functions(b.Body)
		c.blocks(b.Body.Blocks)
	}
}

func walkInputs(f *Feature, body *hclsyntax.Body, prefix string) {
	for name := range body.Attributes {
		path := prefix + name
		add(&f.Inputs, path)
	}
	for _, b := range body.Blocks {
		path := prefix + b.Type
		add(&f.Inputs, path)
		walkInputs(f, b.Body, path+".")
	}
}

func (c *Collector) functions(body *hclsyntax.Body) {
	for _, attr := range body.Attributes {
		diags := hclsyntax.VisitAll(attr.Expr, func(node hclsyntax.Node) hcl.Diagnostics {
			if call, ok := node.(*hclsyntax.FunctionCallExpr); ok && strings.HasPrefix(call.Name, "provider::") {
				parts := strings.Split(call.Name, "::")
				if len(parts) == 3 {
					c.feature("function", parts[2], "").Provider = parts[1]
				}
			}
			return nil
		})
		if diags.HasErrors() {
			add(&c.warnings, WarningConfigUnavailable)
		}
	}
}

func (c *Collector) jsonBlocks(root map[string]any) {
	for rawKind, value := range root {
		kind := rawKind
		if kind == "data" {
			kind = "data_source"
		}
		if kind == "list" {
			kind = "list_resource"
		}
		if kind == "ephemeral" {
			kind = "ephemeral_resource"
		}
		if kind == "resource" || kind == "data_source" || kind == "ephemeral_resource" || kind == "action" || kind == "list_resource" || kind == "provider" || kind == "state_store" {
			for name, raw := range asMap(value) {
				if kind == "provider" {
					f := c.feature(kind, name, "")
					jsonInputs(f, asMap(raw))
					continue
				}
				if kind == "state_store" {
					jsonInputs(c.feature(kind, name, ""), asMap(asMap(raw)["config"]))
					continue
				}
				for label, cfg := range asMap(raw) {
					address := name + "." + label
					if kind == "data_source" {
						address = "data." + address
					}
					if kind == "ephemeral_resource" {
						address = "ephemeral." + address
					}
					if kind == "action" || kind == "list_resource" {
						cfg = asMap(cfg)["config"]
					}
					if kind == "action" {
						address = "action." + address
					}
					if kind == "list_resource" {
						address = "list." + name + "." + label
					}
					jsonInputs(c.feature(kind, name, address), asMap(cfg))
				}
			}
		}
	}
}

func asMap(value any) map[string]any { m, _ := value.(map[string]any); return m }

func jsonInputs(f *Feature, object map[string]any) {
	for name := range object {
		add(&f.Inputs, name)
	}
}
