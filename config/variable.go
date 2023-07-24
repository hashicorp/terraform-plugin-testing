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

type Variables map[string]Variable

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

type Variable interface {
	json.Marshaler
}

func BoolVariable(value bool) boolVariable {
	return boolVariable{
		value: value,
	}
}

type boolVariable struct {
	value bool
}

func (t boolVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func ListVariable(value ...Variable) listVariable {
	return listVariable{
		value: value,
	}
}

type listVariable struct {
	value []Variable
}

func (t listVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func MapVariable(value map[string]Variable) mapVariable {
	return mapVariable{
		value: value,
	}
}

type mapVariable struct {
	value map[string]Variable
}

func (t mapVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func ObjectVariable(value map[string]Variable) objectVariable {
	return objectVariable{
		value: value,
	}
}

type objectVariable struct {
	value map[string]Variable
}

func (t objectVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

type number interface {
	constraints.Float | constraints.Integer | *big.Float
}

func NumberVariable[T number](value T) numberVariable {
	return numberVariable{
		value: value,
	}
}

type numberVariable struct {
	value any
}

func (t numberVariable) MarshalJSON() ([]byte, error) {
	switch v := t.value.(type) {
	case *big.Float:
		return []byte(v.Text('g', -1)), nil
	}

	return json.Marshal(t.value)
}

func SetVariable(value ...Variable) setVariable {
	return setVariable{
		value: value,
	}
}

type setVariable struct {
	value []Variable
}

func (t setVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func StringVariable(value string) stringVariable {
	return stringVariable{
		value: value,
	}
}

type stringVariable struct {
	value string
}

func (t stringVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}

func TupleVariable(value ...Variable) tupleVariable {
	return tupleVariable{
		value: value,
	}
}

type tupleVariable struct {
	value []Variable
}

func (t tupleVariable) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.value)
}
