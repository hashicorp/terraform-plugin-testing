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

// QueryConfig reads the query configuration file at path.
func (c *Collector) QueryConfig(path string) {
	if c == nil {
		return
	}
	data, err := os.ReadFile(path)
	if err != nil {
		c.Warning(WarningConfigUnavailable)
		return
	}
	c.Parse(filepath.Base(path), data)
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

// configBlock describes a configuration block type that declares a feature.
type configBlock struct {
	kind Kind
	// addressed blocks have two labels, the type and the name, and an
	// address of addressPrefix + type + "." + name. Other blocks have one
	// label, the name, and no address.
	addressed     bool
	addressPrefix string
	// nestedConfig blocks keep their schema attributes in a nested "config"
	// block.
	nestedConfig bool
}

// configBlocks maps configuration block types to the features they declare.
var configBlocks = map[string]configBlock{
	"provider":    {kind: KindProvider},
	"state_store": {kind: KindStateStore, nestedConfig: true},
	"resource":    {kind: KindResource, addressed: true},
	"data":        {kind: KindDataSource, addressed: true, addressPrefix: "data."},
	"ephemeral":   {kind: KindEphemeralResource, addressed: true, addressPrefix: "ephemeral."},
	"list":        {kind: KindListResource, addressed: true, addressPrefix: "list.", nestedConfig: true},
	"action":      {kind: KindAction, addressed: true, addressPrefix: "action.", nestedConfig: true},
}

// blockFeature returns the feature that a block with the given labels
// declares, or nil if the labels are incomplete.
func (c *Collector) blockFeature(block configBlock, labels ...string) *featureEvidence {
	if !block.addressed {
		if len(labels) < 1 {
			return nil
		}
		return c.feature(block.kind, labels[0], "")
	}
	if len(labels) < 2 {
		return nil
	}
	return c.feature(block.kind, labels[0], block.addressPrefix+labels[0]+"."+labels[1])
}

func (c *Collector) blocks(blocks []*hclsyntax.Block) {
	for _, b := range blocks {
		block, ok := configBlocks[b.Type]
		// Providers in a query file configure the query, not a test step.
		if ok && !(block.kind == KindProvider && strings.HasSuffix(b.DefRange().Filename, ".tfquery.hcl")) {
			if f := c.blockFeature(block, b.Labels...); f != nil {
				if block.nestedConfig {
					for _, nested := range b.Body.Blocks {
						if nested.Type == "config" {
							walkInputs(f, nested.Body, "")
						}
					}
				} else {
					walkInputs(f, b.Body, "")
				}
			}
		}
		c.functions(b.Body)
		c.blocks(b.Body.Blocks)
	}
}

func walkInputs(f *featureEvidence, body *hclsyntax.Body, prefix string) {
	for name := range body.Attributes {
		f.inputs.add(prefix + name)
	}
	for _, b := range body.Blocks {
		path := prefix + b.Type
		f.inputs.add(path)
		walkInputs(f, b.Body, path+".")
	}
}

func (c *Collector) functions(body *hclsyntax.Body) {
	for _, attr := range body.Attributes {
		diags := hclsyntax.VisitAll(attr.Expr, func(node hclsyntax.Node) hcl.Diagnostics {
			if call, ok := node.(*hclsyntax.FunctionCallExpr); ok && strings.HasPrefix(call.Name, "provider::") {
				parts := strings.Split(call.Name, "::")
				if len(parts) == 3 {
					c.feature(KindFunction, parts[2], "").provider = parts[1]
				}
			}
			return nil
		})
		if diags.HasErrors() {
			c.warnings.add(WarningConfigUnavailable)
		}
	}
}

func (c *Collector) jsonBlocks(root map[string]any) {
	for blockType, value := range root {
		block, ok := configBlocks[blockType]
		if !ok {
			continue
		}
		for name, raw := range asMap(value) {
			if !block.addressed {
				c.jsonFeature(block, raw, name)
				continue
			}
			for label, cfg := range asMap(raw) {
				c.jsonFeature(block, cfg, name, label)
			}
		}
	}
}

// jsonFeature records the feature that a Terraform JSON block with the given
// labels and body declares.
func (c *Collector) jsonFeature(block configBlock, body any, labels ...string) {
	f := c.blockFeature(block, labels...)
	attributes := asMap(body)
	if block.nestedConfig {
		attributes = asMap(attributes["config"])
	}
	for name := range attributes {
		f.inputs.add(name)
	}
}

func asMap(value any) map[string]any { m, _ := value.(map[string]any); return m }
