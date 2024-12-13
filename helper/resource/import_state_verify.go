package resource

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	tfjson "github.com/hashicorp/terraform-json"
	"github.com/hashicorp/terraform-plugin-testing/internal/logging"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/mitchellh/go-testing-interface"
	"github.com/zclconf/go-cty/cty"
)

// importStateVerify compares the oldState with the new state and verifies that all attributes are equal
func importStateVerify(ctx context.Context, t testing.T, step TestStep, oldState *tfjson.State, newState *tfjson.State, providerSchemas *tfjson.ProviderSchemas) error {
	logging.HelperResourceTrace(ctx, "Using TestStep ImportStateVerify")

	if newState == nil {
		return errors.New("new import state is nil")
	}

	if newState.Values == nil {
		return errors.New("new import state does not contain any state values")
	}

	if newState.Values.RootModule == nil {
		return errors.New("new import state does not contain a root module")
	}

	if oldState == nil {
		return errors.New("old state is nil")
	}

	if oldState.Values == nil {
		return errors.New("old state does not contain any state values")
	}

	if oldState.Values.RootModule == nil {
		return errors.New("old state does not contain a root module")
	}

	// TODO: logic for data source removal shouldn't be relevant?

	identifierAttribute := step.ImportStateVerifyIdentifierAttribute
	if identifierAttribute == "" {
		identifierAttribute = "id"
	}

	// TODO: the old logic has some weird logic:
	//   - Remove data sources from shim? Shouldn't be a problem, since we are only looping through resources
	for _, newResource := range newState.Values.RootModule.Resources {

		r1Identifier, err := tfjsonpath.Traverse(newResource.AttributeValues, tfjsonpath.New(identifierAttribute)) // TODO: convert from flatmap path to tfjson path lol
		if err != nil {
			return fmt.Errorf("New resource missing identifier attribute %q, ensure attribute value is properly set or use ImportStateVerifyIdentifierAttribute to choose different attribute", identifierAttribute)
		}

		// Find the existing resource
		var oldResource *tfjson.StateResource
		for _, r2 := range oldState.Values.RootModule.Resources {
			if r2.Type != newResource.Type || r2.ProviderName != newResource.ProviderName {
				continue
			}

			r2Identifier, err := tfjsonpath.Traverse(r2.AttributeValues, tfjsonpath.New(identifierAttribute)) // TODO: convert from flatmap path to tfjson path lol
			if err != nil {
				return fmt.Errorf("Old resource missing identifier attribute %q, ensure attribute value is properly set or use ImportStateVerifyIdentifierAttribute to choose different attribute", identifierAttribute)
			}

			// TODO: probably need to write a generic reflection compare for tfjson values (see check value implementations)
			if reflect.DeepEqual(r2Identifier, r1Identifier) {
				oldResource = r2
				break
			}
		}

		if oldResource == nil {
			return fmt.Errorf("Failed state verification, resource with ID %s not found", r1Identifier)
		}

		if providerSchemas == nil || providerSchemas.Schemas[oldResource.ProviderName] == nil || providerSchemas.Schemas[oldResource.ProviderName].ResourceSchemas == nil {
			return errors.New("Failed to retrieve provider schema")
		}

		schema, ok := providerSchemas.Schemas[oldResource.ProviderName].ResourceSchemas[oldResource.Type]
		if !ok || schema.Block == nil {
			return errors.New("Failed to retrieve resource schema")
		}

		// TODO: the old logic has some weird logic:
		//   - Skip empty containers? Don't think we need to do that
		//   - Skip attributes via the "step.ImportStateVerifyIgnore" field, which is a list of flatmapped prefixes, need to be converted to tfjsonpath
		//   - At the root, skip all objects in "timeouts" (legacy behavior, since they are not always present)
		var errs error
		for name, attr := range schema.Block.Attributes {
			attrPath := tfjsonpath.New(name)
			oldVal, _ := tfjsonpath.Traverse(oldResource.AttributeValues, attrPath)
			newVal, _ := tfjsonpath.Traverse(newResource.AttributeValues, attrPath)

			errs = errors.Join(errs, compareWithAttribute(attrPath, attr, oldVal, newVal))
		}

		for name, nestedBlock := range schema.Block.NestedBlocks {
			blockPath := tfjsonpath.New(name)
			// TODO: add function for traversing nested blocks
			fmt.Println(blockPath, nestedBlock)
		}
	}

	return nil
}

// This is called once we find an attribute, loops through nested attributes, with the goal of eventually finding an attribute type to end the recursion
func compareWithAttribute(p tfjsonpath.Path, attr *tfjson.SchemaAttribute, oldVal any, newVal any) error {
	var errs error

	// No more nested attributes to traverse
	if attr.AttributeNestedType == nil {
		// TODO: error handling of course!
		errs = errors.Join(errs, compareValue(p, attr.AttributeType, oldVal, newVal))
	}

	// TODO: add support for nested attributes

	return errs
}

func compareValue(p tfjsonpath.Path, typ cty.Type, oldVal any, newVal any) error {
	var errs error
	if typ.IsPrimitiveType() {
		// Compare the two primitive values with reflection
		if !reflect.DeepEqual(oldVal, newVal) {
			errs = errors.Join(errs, fmt.Errorf("found a diff between: %s and %s, at path: %s", oldVal, newVal, p))
		}
	} else if typ.IsListType() {
		// TODO: error handling of course!
		oldList, _ := oldVal.([]any)
		newList, _ := newVal.([]any)

		if len(oldList) != len(newList) {
			errs = errors.Join(errs, fmt.Errorf("found a diff between: %s and %s, at path: %s", oldVal, newVal, p))
		}

		for i, oldV := range oldList {
			newV := newList[i]
			elementPath := p.AtSliceIndex(i)

			errs = errors.Join(errs, compareValue(elementPath, typ.ElementType(), oldV, newV))
		}
	} else if typ.IsMapType() || typ.IsObjectType() {
		// TODO: error handling of course!
		oldMap, _ := oldVal.(map[string]any)
		newMap, _ := newVal.(map[string]any)

		if len(oldMap) != len(newMap) {
			errs = errors.Join(errs, fmt.Errorf("found a diff between: %s and %s, at path: %s", oldVal, newVal, p))
		}

		for key, oldV := range oldMap {
			// TODO: check for mismatched key
			newV := newMap[key]
			elementPath := p.AtMapKey(key)

			var nextTyp cty.Type
			if typ.IsMapType() {
				nextTyp = typ.ElementType()
			} else {
				nextTyp = typ.AttributeType(key)
			}

			errs = errors.Join(errs, compareValue(elementPath, nextTyp, oldV, newV))
		}
	} else if typ.IsSetType() {
		// TODO: error handling of course!
		oldList, _ := oldVal.([]any)
		newList, _ := newVal.([]any)

		if len(oldList) != len(newList) {
			errs = errors.Join(errs, fmt.Errorf("found a diff between: %s and %s, at path: %s", oldVal, newVal, p))
		}

		for i, oldV := range oldList {
			// TODO: we don't have a AtValue for sets :)
			elementPath := p.AtSliceIndex(i)
			foundMatch := false
			for i, newV := range newList {
				err := compareValue(elementPath, typ.ElementType(), oldV, newV)
				if err == nil {
					foundMatch = true
					// Remove the matching element from the slice of candidates
					newList = append(newList[:i], newList[i+1:]...)
					break
				}
			}

			// No matching element in the new set
			if !foundMatch {
				errs = errors.Join(errs, fmt.Errorf("old set value %s not found in new set, at path: %s", oldVal, p))
			}
		}
	} else {
		panic("ayo, what is this type?")
	}

	return errs
}
