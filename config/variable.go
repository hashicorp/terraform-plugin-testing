// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"golang.org/x/exp/constraints"
)

const AutoTFVarsJson = "generated.auto.tfvars.json"

// Variable interface is an alias to json.Marshaler.
type Variable interface {
	json.Marshaler
}

// Variables is a type holding a key-value map of variable names
// to types implementing Variable interface.
type Variables map[string]Variable

// Write iterates over each element in v and assembles a JSON
// file which is named AutoTFVarsJson and written to dest.
func (v Variables) Write(dest string) error {
	buf := bytes.NewBuffer(nil)

	buf.Write([]byte(`{`))

	for k, val := range v {
		j, err := val.MarshalJSON()

		if err != nil {
			return err
		}

		buf.Write([]byte(fmt.Sprintf("%q: ", k)))
		buf.Write(j)
		buf.Write([]byte(","))
	}

	b := bytes.TrimRight(buf.Bytes(), ",")

	buf = bytes.NewBuffer(b)

	buf.Write([]byte(`}`))

	outFilename := filepath.Join(dest, AutoTFVarsJson)

	err := os.WriteFile(outFilename, buf.Bytes(), 0700)

	if err != nil {
		return err
	}

	return nil
}

var _ Variable = boolVariable{}

type boolVariable struct {
	value bool
}

// MarshalJSON returns the JSON encoding of boolVariable.
func (v boolVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

// BoolVariable instantiates an instance of boolVariable,
// which implements Variable.
func BoolVariable(value bool) boolVariable {
	return boolVariable{
		value: value,
	}
}

var _ Variable = listVariable{}

type listVariable struct {
	value []Variable
}

// MarshalJSON returns the JSON encoding of listVariable.
func (v listVariable) MarshalJSON() ([]byte, error) {
	if !typesEq(v.value) {
		return nil, errors.New("lists must contain the same type")
	}

	return json.Marshal(v.value)
}

// ListVariable instantiates an instance of listVariable,
// which implements Variable.
func ListVariable(value ...Variable) listVariable {
	return listVariable{
		value: value,
	}
}

var _ Variable = mapVariable{}

type mapVariable struct {
	value map[string]Variable
}

// MarshalJSON returns the JSON encoding of mapVariable.
func (v mapVariable) MarshalJSON() ([]byte, error) {
	var variables []Variable

	for _, variable := range v.value {
		variables = append(variables, variable)
	}

	if !typesEq(variables) {
		return nil, errors.New("maps must contain the same type")
	}

	return json.Marshal(v.value)
}

// MapVariable instantiates an instance of mapVariable,
// which implements Variable.
func MapVariable(value map[string]Variable) mapVariable {
	return mapVariable{
		value: value,
	}
}

var _ Variable = objectVariable{}

type objectVariable struct {
	value map[string]Variable
}

// MarshalJSON returns the JSON encoding of objectVariable.
func (v objectVariable) MarshalJSON() ([]byte, error) {
	b, err := json.Marshal(v.value)

	if err != nil {
		innerErr := err

		for errors.Unwrap(innerErr) != nil {
			innerErr = errors.Unwrap(err)
		}

		return nil, innerErr
	}

	return b, nil
}

// ObjectVariable instantiates an instance of objectVariable,
// which implements Variable.
func ObjectVariable(value map[string]Variable) objectVariable {
	return objectVariable{
		value: value,
	}
}

var _ Variable = numberVariable{}

type number interface {
	constraints.Float | constraints.Integer | string
}

type numberVariable struct {
	value any
}

// MarshalJSON returns the JSON encoding of numberVariable.
// NumberVariable allows initialising a number with any floating
// point or integer type. NumberVariable can be initialised
// with a string for values that do not fit into a floating point
// or integer type.
// TODO: Impose restrictions on what can be held in numberVariable
// to match restrictions imposed by Terraform.
func (v numberVariable) MarshalJSON() ([]byte, error) {
	switch v := v.value.(type) {
	case string:
		return []byte(v), nil
	}

	return json.Marshal(v.value)
}

// NumberVariable instantiates an instance of numberVariable,
// which implements Variable.
func NumberVariable[T number](value T) numberVariable {
	return numberVariable{
		value: value,
	}
}

var _ Variable = setVariable{}

type setVariable struct {
	value []Variable
}

// MarshalJSON returns the JSON encoding of setVariable.
func (v setVariable) MarshalJSON() ([]byte, error) {
	for kx, x := range v.value {
		for ky, y := range v.value {
			if kx == ky {
				continue
			}

			if _, ok := x.(setVariable); !ok {
				continue
			}

			if _, ok := y.(setVariable); !ok {
				continue
			}

			if reflect.DeepEqual(x, y) {
				return nil, errors.New("sets must contain unique elements")
			}
		}

	}

	if !typesEq(v.value) {
		return nil, errors.New("sets must contain the same type")
	}

	return json.Marshal(v.value)
}

// SetVariable instantiates an instance of setVariable,
// which implements Variable.
func SetVariable(value ...Variable) setVariable {
	return setVariable{
		value: value,
	}
}

var _ Variable = stringVariable{}

type stringVariable struct {
	value string
}

// MarshalJSON returns the JSON encoding of stringVariable.
func (v stringVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

// StringVariable instantiates an instance of stringVariable,
// which implements Variable.
func StringVariable(value string) stringVariable {
	return stringVariable{
		value: value,
	}
}

type tupleVariable struct {
	value []Variable
}

// MarshalJSON returns the JSON encoding of tupleVariable.
func (v tupleVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

// TupleVariable instantiates an instance of tupleVariable,
// which implements Variable.
func TupleVariable(value ...Variable) tupleVariable {
	return tupleVariable{
		value: value,
	}
}

func typesEq(variables []Variable) bool {
	var t reflect.Type

	for _, variable := range variables {
		switch x := variable.(type) {
		case listVariable:
			if !typesEq(x.value) {
				return false
			}
		case mapVariable:
			var vars []Variable

			for _, v := range x.value {
				vars = append(vars, v)
			}

			if !typesEq(vars) {
				return false
			}
		case setVariable:
			if !typesEq(x.value) {
				return false
			}
		}

		typeOfVariable := reflect.TypeOf(variable)

		if t == nil {
			t = typeOfVariable
			continue
		}

		if t != typeOfVariable {
			return false
		}
	}

	return true
}
