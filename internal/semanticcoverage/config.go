// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package semanticcoverage

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

// ObserveConfigDir records the features and configured attributes declared by
// the configuration files in dir, the test working directory.
//
// Only HCL configuration is read. Terraform JSON configuration is rare in
// provider tests, so it isn't parsed; a .tf.json file only records
// WarningConfigUnavailable, so that the report shows the missing evidence.
func (s *Step) ObserveConfigDir(dir string) {
	if s == nil {
		return
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		s.warn(WarningConfigUnavailable)
		return
	}
	for _, entry := range entries {
		name := entry.Name()
		switch {
		case entry.IsDir():
		case strings.HasSuffix(name, ".tf.json"):
			s.warn(WarningConfigUnavailable)
		case strings.HasSuffix(name, ".tf"), strings.HasSuffix(name, ".tfquery.hcl"):
			s.observeConfigFile(filepath.Join(dir, name))
		}
	}
}

// ObserveQueryConfig records the features and configured attributes declared
// by the query configuration file at path.
func (s *Step) ObserveQueryConfig(path string) {
	if s == nil {
		return
	}
	s.observeConfigFile(path)
}

// observeConfigFile records the features and configured attributes declared by
// the HCL configuration file at path.
func (s *Step) observeConfigFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		s.warn(WarningConfigUnavailable)
		return
	}
	file, diags := hclparse.NewParser().ParseHCL(data, filepath.Base(path))
	if diags.HasErrors() {
		s.warn(WarningConfigUnavailable)
		return
	}
	body, ok := file.Body.(*hclsyntax.Body)
	if !ok { // ParseHCL always returns an hclsyntax body.
		s.warn(WarningConfigUnavailable)
		return
	}
	// Providers in a query file configure the query, not a test step.
	s.observeBlocks(body.Blocks, strings.HasSuffix(path, ".tfquery.hcl"))
	s.observeFunctionCalls(body)
}

// configBlock describes a configuration block type that declares a feature.
type configBlock struct {
	kind Kind
	// nestedConfig blocks keep their schema attributes in a nested "config"
	// block.
	nestedConfig bool
}

// configBlocks maps configuration block types to the features they declare.
// The first label of each is the feature's type name.
var configBlocks = map[string]configBlock{
	"provider":    {kind: KindProvider},
	"state_store": {kind: KindStateStore},
	"resource":    {kind: KindResource},
	"data":        {kind: KindDataSource},
	"ephemeral":   {kind: KindEphemeralResource},
	"list":        {kind: KindListResource, nestedConfig: true},
	"action":      {kind: KindAction, nestedConfig: true},
}

func (s *Step) observeBlocks(blocks []*hclsyntax.Block, queryFile bool) {
	for _, b := range blocks {
		block, ok := configBlocks[b.Type]
		if ok && len(b.Labels) > 0 && !(queryFile && block.kind == KindProvider) {
			f := s.feature(block.kind, b.Labels[0], "")
			for _, body := range schemaBodies(b, block) {
				recordConfigured(f, body)
			}
			// The body of a feature follows the provider's schema, so its
			// nested blocks are schema blocks, not declarations. A state
			// store is the exception: it holds the provider block it uses.
			if block.kind != KindStateStore {
				continue
			}
		}
		// Other blocks, such as terraform and check, can declare features.
		s.observeBlocks(b.Body.Blocks, queryFile)
	}
}

// schemaBodies returns the bodies of b that hold the feature's schema
// attributes.
func schemaBodies(b *hclsyntax.Block, block configBlock) []*hclsyntax.Body {
	if !block.nestedConfig {
		return []*hclsyntax.Body{b.Body}
	}
	var bodies []*hclsyntax.Body
	for _, nested := range b.Body.Blocks {
		if nested.Type == "config" {
			bodies = append(bodies, nested.Body)
		}
	}
	return bodies
}

// recordConfigured records the top-level attributes and blocks set in body.
func recordConfigured(f *Feature, body *hclsyntax.Body) {
	for name := range body.Attributes {
		add(&f.Configured, name)
	}
	for _, b := range body.Blocks {
		name := b.Type
		if name == "dynamic" && len(b.Labels) == 1 {
			name = b.Labels[0]
		}
		add(&f.Configured, name)
	}
}

// observeFunctionCalls records the provider-defined functions called in body.
func (s *Step) observeFunctionCalls(body *hclsyntax.Body) {
	// The callback returns no diagnostics, so neither does VisitAll.
	_ = hclsyntax.VisitAll(body, func(node hclsyntax.Node) hcl.Diagnostics {
		if call, ok := node.(*hclsyntax.FunctionCallExpr); ok {
			if parts := strings.Split(call.Name, "::"); len(parts) == 3 && parts[0] == "provider" {
				s.feature(KindFunction, parts[2], parts[1])
			}
		}
		return nil
	})
}
