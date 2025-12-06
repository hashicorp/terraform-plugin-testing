// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package terraform

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/hashicorp/go-cty/cty"

	"github.com/hashicorp/terraform-plugin-testing/internal/configs/configschema"
	"github.com/hashicorp/terraform-plugin-testing/internal/configs/hcl2shim"
)

func TestResourceConfigGet(t *testing.T) {
	t.Parallel()

	fooStringSchema := &configschema.Block{
		Attributes: map[string]*configschema.Attribute{
			"foo": {Type: cty.String, Optional: true},
		},
	}
	fooListSchema := &configschema.Block{
		Attributes: map[string]*configschema.Attribute{
			"foo": {Type: cty.List(cty.Number), Optional: true},
		},
	}

	cases := []struct {
		Config cty.Value
		Schema *configschema.Block
		Key    string
		Value  any
	}{
		{
			Config: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
			}),
			Schema: fooStringSchema,
			Key:    "foo",
			Value:  "bar",
		},

		{
			Config: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.UnknownVal(cty.String),
			}),
			Schema: fooStringSchema,
			Key:    "foo",
			Value:  hcl2shim.UnknownVariableValue,
		},

		{
			Config: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.ListVal([]cty.Value{
					cty.NumberIntVal(1),
					cty.NumberIntVal(2),
					cty.NumberIntVal(5),
				}),
			}),
			Schema: fooListSchema,
			Key:    "foo.0",
			Value:  1,
		},

		{
			Config: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.ListVal([]cty.Value{
					cty.NumberIntVal(1),
					cty.NumberIntVal(2),
					cty.NumberIntVal(5),
				}),
			}),
			Schema: fooListSchema,
			Key:    "foo.5",
			Value:  nil,
		},

		{
			Config: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.ListVal([]cty.Value{
					cty.NumberIntVal(1),
					cty.NumberIntVal(2),
					cty.NumberIntVal(5),
				}),
			}),
			Schema: fooListSchema,
			Key:    "foo.-1",
			Value:  nil,
		},

		// get from map
		{
			Config: cty.ObjectVal(map[string]cty.Value{
				"mapname": cty.ListVal([]cty.Value{
					cty.MapVal(map[string]cty.Value{
						"key": cty.NumberIntVal(1),
					}),
				}),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"mapname": {Type: cty.List(cty.Map(cty.Number)), Optional: true},
				},
			},
			Key:   "mapname.0.key",
			Value: 1,
		},

		// get from map with dot in key
		{
			Config: cty.ObjectVal(map[string]cty.Value{
				"mapname": cty.ListVal([]cty.Value{
					cty.MapVal(map[string]cty.Value{
						"key.name": cty.NumberIntVal(1),
					}),
				}),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"mapname": {Type: cty.List(cty.Map(cty.Number)), Optional: true},
				},
			},
			Key:   "mapname.0.key.name",
			Value: 1,
		},

		// get from map with overlapping key names
		{
			Config: cty.ObjectVal(map[string]cty.Value{
				"mapname": cty.ListVal([]cty.Value{
					cty.MapVal(map[string]cty.Value{
						"key.name":   cty.NumberIntVal(1),
						"key.name.2": cty.NumberIntVal(2),
					}),
				}),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"mapname": {Type: cty.List(cty.Map(cty.Number)), Optional: true},
				},
			},
			Key:   "mapname.0.key.name.2",
			Value: 2,
		},
		{
			Config: cty.ObjectVal(map[string]cty.Value{
				"mapname": cty.ListVal([]cty.Value{
					cty.MapVal(map[string]cty.Value{
						"key.name":     cty.NumberIntVal(1),
						"key.name.foo": cty.NumberIntVal(2),
					}),
				}),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"mapname": {Type: cty.List(cty.Map(cty.Number)), Optional: true},
				},
			},
			Key:   "mapname.0.key.name",
			Value: 1,
		},
		{
			Config: cty.ObjectVal(map[string]cty.Value{
				"mapname": cty.ListVal([]cty.Value{
					cty.MapVal(map[string]cty.Value{
						"listkey": cty.ListVal([]cty.Value{
							cty.MapVal(map[string]cty.Value{
								"key": cty.NumberIntVal(3),
							}),
						}),
					}),
				}),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"mapname": {Type: cty.List(cty.Map(cty.List(cty.Map(cty.Number)))), Optional: true},
				},
			},
			Key:   "mapname.0.listkey.0.key",
			Value: 3,
		},
	}

	for i, tc := range cases {
		rc := NewResourceConfigShimmed(tc.Config, tc.Schema)

		// Test getting a key
		t.Run(fmt.Sprintf("get-%d", i), func(t *testing.T) {
			t.Parallel()

			v, ok := rc.Get(tc.Key)
			if ok && v == nil {
				t.Fatal("(nil, true) returned from Get")
			}

			if !reflect.DeepEqual(v, tc.Value) {
				t.Fatalf("%d bad: %#v", i, v)
			}
		})

		// Test copying and equality
		t.Run(fmt.Sprintf("copy-and-equal-%d", i), func(t *testing.T) {
			t.Parallel()

			copiedConfig := rc.DeepCopy()
			if !reflect.DeepEqual(copiedConfig, rc) {
				t.Fatalf("bad:\n\n%#v\n\n%#v", copiedConfig, rc)
			}

			if !copiedConfig.Equal(rc) {
				t.Fatalf("copiedConfig != rc:\n\n%#v\n\n%#v", copiedConfig, rc)
			}
			if !rc.Equal(copiedConfig) {
				t.Fatalf("rc != copiedConfig:\n\n%#v\n\n%#v", copiedConfig, rc)
			}
		})
	}
}

func TestResourceConfigDeepCopy_nil(t *testing.T) {
	t.Parallel()

	var nilRc *ResourceConfig
	actual := nilRc.DeepCopy()
	if actual != nil {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestResourceConfigDeepCopy_nilComputed(t *testing.T) {
	t.Parallel()

	rc := &ResourceConfig{}
	actual := rc.DeepCopy()
	if actual.ComputedKeys != nil {
		t.Fatalf("bad: %#v", actual)
	}
}

func TestResourceConfigEqual_nil(t *testing.T) {
	t.Parallel()

	var nilRc *ResourceConfig
	notNil := NewResourceConfigShimmed(cty.EmptyObjectVal, &configschema.Block{})

	if nilRc.Equal(notNil) {
		t.Fatal("should not be equal")
	}

	if notNil.Equal(nilRc) {
		t.Fatal("should not be equal")
	}
}

func TestResourceConfigEqual_computedKeyOrder(t *testing.T) {
	t.Parallel()

	v := cty.ObjectVal(map[string]cty.Value{
		"foo": cty.UnknownVal(cty.String),
	})
	schema := &configschema.Block{
		Attributes: map[string]*configschema.Attribute{
			"foo": {Type: cty.String, Optional: true},
		},
	}
	rc := NewResourceConfigShimmed(v, schema)
	rc2 := NewResourceConfigShimmed(v, schema)

	// Set the computed keys manually to force ordering to differ
	rc.ComputedKeys = []string{"foo", "bar"}
	rc2.ComputedKeys = []string{"bar", "foo"}

	if !rc.Equal(rc2) {
		t.Fatal("should be equal")
	}
}

func TestNewResourceConfigShimmed(t *testing.T) {
	t.Parallel()

	for _, tc := range []struct {
		Name     string
		Val      cty.Value
		Schema   *configschema.Block
		Expected *ResourceConfig
	}{
		{
			Name: "empty object",
			Val:  cty.NullVal(cty.EmptyObject),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"foo": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			Expected: &ResourceConfig{
				Raw:    map[string]any{},
				Config: map[string]any{},
			},
		},
		{
			Name: "basic",
			Val: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.StringVal("bar"),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"foo": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			Expected: &ResourceConfig{
				Raw: map[string]any{
					"foo": "bar",
				},
				Config: map[string]any{
					"foo": "bar",
				},
			},
		},
		{
			Name: "null string",
			Val: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.NullVal(cty.String),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"foo": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			Expected: &ResourceConfig{
				Raw:    map[string]any{},
				Config: map[string]any{},
			},
		},
		{
			Name: "unknown string",
			Val: cty.ObjectVal(map[string]cty.Value{
				"foo": cty.UnknownVal(cty.String),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"foo": {
						Type:     cty.String,
						Optional: true,
					},
				},
			},
			Expected: &ResourceConfig{
				ComputedKeys: []string{"foo"},
				Raw: map[string]any{
					"foo": hcl2shim.UnknownVariableValue,
				},
				Config: map[string]any{
					"foo": hcl2shim.UnknownVariableValue,
				},
			},
		},
		{
			Name: "unknown collections",
			Val: cty.ObjectVal(map[string]cty.Value{
				"bar": cty.UnknownVal(cty.Map(cty.String)),
				"baz": cty.UnknownVal(cty.List(cty.String)),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"bar": {
						Type:     cty.Map(cty.String),
						Required: true,
					},
					"baz": {
						Type:     cty.List(cty.String),
						Optional: true,
					},
				},
			},
			Expected: &ResourceConfig{
				ComputedKeys: []string{"bar", "baz"},
				Raw: map[string]any{
					"bar": hcl2shim.UnknownVariableValue,
					"baz": hcl2shim.UnknownVariableValue,
				},
				Config: map[string]any{
					"bar": hcl2shim.UnknownVariableValue,
					"baz": hcl2shim.UnknownVariableValue,
				},
			},
		},
		{
			Name: "null collections",
			Val: cty.ObjectVal(map[string]cty.Value{
				"bar": cty.NullVal(cty.Map(cty.String)),
				"baz": cty.NullVal(cty.List(cty.String)),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"bar": {
						Type:     cty.Map(cty.String),
						Required: true,
					},
					"baz": {
						Type:     cty.List(cty.String),
						Optional: true,
					},
				},
			},
			Expected: &ResourceConfig{
				Raw:    map[string]any{},
				Config: map[string]any{},
			},
		},
		{
			Name: "unknown blocks",
			Val: cty.ObjectVal(map[string]cty.Value{
				"bar": cty.UnknownVal(cty.Map(cty.String)),
				"baz": cty.UnknownVal(cty.List(cty.String)),
			}),
			Schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"bar": {
						Block:   configschema.Block{},
						Nesting: configschema.NestingList,
					},
					"baz": {
						Block:   configschema.Block{},
						Nesting: configschema.NestingSet,
					},
				},
			},
			Expected: &ResourceConfig{
				ComputedKeys: []string{"bar", "baz"},
				Raw: map[string]any{
					"bar": hcl2shim.UnknownVariableValue,
					"baz": hcl2shim.UnknownVariableValue,
				},
				Config: map[string]any{
					"bar": hcl2shim.UnknownVariableValue,
					"baz": hcl2shim.UnknownVariableValue,
				},
			},
		},
		{
			Name: "unknown in nested blocks",
			Val: cty.ObjectVal(map[string]cty.Value{
				"bar": cty.ListVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"baz": cty.ListVal([]cty.Value{
							cty.ObjectVal(map[string]cty.Value{
								"list": cty.UnknownVal(cty.List(cty.String)),
							}),
						}),
					}),
				}),
			}),
			Schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"bar": {
						Block: configschema.Block{
							BlockTypes: map[string]*configschema.NestedBlock{
								"baz": {
									Block: configschema.Block{
										Attributes: map[string]*configschema.Attribute{
											"list": {Type: cty.List(cty.String),
												Optional: true,
											},
										},
									},
									Nesting: configschema.NestingList,
								},
							},
						},
						Nesting: configschema.NestingList,
					},
				},
			},
			Expected: &ResourceConfig{
				ComputedKeys: []string{"bar.0.baz.0.list"},
				Raw: map[string]any{
					"bar": []any{map[string]any{
						"baz": []any{map[string]any{
							"list": "74D93920-ED26-11E3-AC10-0800200C9A66",
						}},
					}},
				},
				Config: map[string]any{
					"bar": []any{map[string]any{
						"baz": []any{map[string]any{
							"list": "74D93920-ED26-11E3-AC10-0800200C9A66",
						}},
					}},
				},
			},
		},
		{
			Name: "unknown in set",
			Val: cty.ObjectVal(map[string]cty.Value{
				"bar": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"val": cty.UnknownVal(cty.String),
					}),
				}),
			}),
			Schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"bar": {
						Block: configschema.Block{
							Attributes: map[string]*configschema.Attribute{
								"val": {
									Type:     cty.String,
									Optional: true,
								},
							},
						},
						Nesting: configschema.NestingSet,
					},
				},
			},
			Expected: &ResourceConfig{
				ComputedKeys: []string{"bar.0.val"},
				Raw: map[string]any{
					"bar": []any{map[string]any{
						"val": "74D93920-ED26-11E3-AC10-0800200C9A66",
					}},
				},
				Config: map[string]any{
					"bar": []any{map[string]any{
						"val": "74D93920-ED26-11E3-AC10-0800200C9A66",
					}},
				},
			},
		},
		{
			Name: "unknown in attribute sets",
			Val: cty.ObjectVal(map[string]cty.Value{
				"bar": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"val": cty.UnknownVal(cty.String),
					}),
				}),
				"baz": cty.SetVal([]cty.Value{
					cty.ObjectVal(map[string]cty.Value{
						"obj": cty.UnknownVal(cty.Object(map[string]cty.Type{
							"attr": cty.List(cty.String),
						})),
					}),
					cty.ObjectVal(map[string]cty.Value{
						"obj": cty.ObjectVal(map[string]cty.Value{
							"attr": cty.UnknownVal(cty.List(cty.String)),
						}),
					}),
				}),
			}),
			Schema: &configschema.Block{
				Attributes: map[string]*configschema.Attribute{
					"bar": {
						Type: cty.Set(cty.Object(map[string]cty.Type{
							"val": cty.String,
						})),
					},
					"baz": {
						Type: cty.Set(cty.Object(map[string]cty.Type{
							"obj": cty.Object(map[string]cty.Type{
								"attr": cty.List(cty.String),
							}),
						})),
					},
				},
			},
			Expected: &ResourceConfig{
				ComputedKeys: []string{"bar.0.val", "baz.0.obj.attr", "baz.1.obj"},
				Raw: map[string]any{
					"bar": []any{map[string]any{
						"val": "74D93920-ED26-11E3-AC10-0800200C9A66",
					}},
					"baz": []any{
						map[string]any{
							"obj": map[string]any{
								"attr": "74D93920-ED26-11E3-AC10-0800200C9A66",
							},
						},
						map[string]any{
							"obj": "74D93920-ED26-11E3-AC10-0800200C9A66",
						},
					},
				},
				Config: map[string]any{
					"bar": []any{map[string]any{
						"val": "74D93920-ED26-11E3-AC10-0800200C9A66",
					}},
					"baz": []any{
						map[string]any{
							"obj": map[string]any{
								"attr": "74D93920-ED26-11E3-AC10-0800200C9A66",
							},
						},
						map[string]any{
							"obj": "74D93920-ED26-11E3-AC10-0800200C9A66",
						},
					},
				},
			},
		},
		{
			Name: "null blocks",
			Val: cty.ObjectVal(map[string]cty.Value{
				"bar": cty.NullVal(cty.Map(cty.String)),
				"baz": cty.NullVal(cty.List(cty.String)),
			}),
			Schema: &configschema.Block{
				BlockTypes: map[string]*configschema.NestedBlock{
					"bar": {
						Block:   configschema.Block{},
						Nesting: configschema.NestingMap,
					},
					"baz": {
						Block:   configschema.Block{},
						Nesting: configschema.NestingSingle,
					},
				},
			},
			Expected: &ResourceConfig{
				Raw:    map[string]any{},
				Config: map[string]any{},
			},
		},
	} {

		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			cfg := NewResourceConfigShimmed(tc.Val, tc.Schema)
			if !tc.Expected.Equal(cfg) {
				t.Fatalf("expected:\n%#v\ngot:\n%#v", tc.Expected, cfg)
			}
		})
	}
}
