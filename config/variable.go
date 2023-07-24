package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"golang.org/x/exp/constraints"
)

const autoTFVarsJson = "generated.auto.tfvars.json"

// Variable interface is an alias to json.Marshaler.
type Variable interface {
	json.Marshaler
}

// Variables is a type holding a key-value map of variable names
// to types implementing Variable interface.
type Variables map[string]Variable

// Write iterates over each element in v and assembles a JSON
// file which is named autoTFVarsJson and written to dest.
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

	outFilename := filepath.Join(dest, autoTFVarsJson)

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
func (t boolVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
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
func (t listVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
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
func (t mapVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
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
func (t objectVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
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
	constraints.Float | constraints.Integer | *big.Float
}

type numberVariable struct {
	value any
}

// MarshalJSON returns the JSON encoding of numberVariable.
// If the value of numberVariable is *bigFloat then the
// representation of the value is the smallest number of
// digits required to uniquely identify the value using the
// precision of the *bigFloat that was supplied when
// numberVariable was instantiated.
func (t numberVariable) MarshalJSON() ([]byte, error) {
	switch v := t.value.(type) {
	case *big.Float:
		return []byte(v.Text('g', -1)), nil
	}

	return json.Marshal(t.value)
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
func (t setVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
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
func (t stringVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
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
func (t tupleVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

// TupleVariable instantiates an instance of tupleVariable,
// which implements Variable.
func TupleVariable(value ...Variable) tupleVariable {
	return tupleVariable{
		value: value,
	}
}
