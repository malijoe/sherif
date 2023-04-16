package sherif

import (
	"encoding/json"
	"strings"
	"testing"
)

type tester struct {
	Field1 string
	Field2 int
	Field3 []string
}

func (t *tester) Unmarshal(unmarshal func(any) error) (err error) {
	var obj struct {
		Field1 string   `json:"field1" yaml:"field1"`
		Field2 int      `json:"field2" yaml:"field2"`
		Field3 []string `json:"field3" yaml:"field3"`
	}

	if err = unmarshal(&obj); err != nil {
		return
	}

	*t = tester{
		Field1: obj.Field1,
		Field2: obj.Field2,
		Field3: obj.Field3,
	}
	return
}

func (t tester) Marshal() any {
	return struct {
		Field1 string   `json:"field1" yaml:"field1"`
		Field2 int      `json:"field2" yaml:"field2"`
		Field3 []string `json:"field3" yaml:"field3"`
	}{
		Field1: t.Field1,
		Field2: t.Field2,
		Field3: t.Field3,
	}
}

func TestSherif(t *testing.T) {
	jsonStr := `{"field1":"value1","field2":2,"field3":["value","three"]}`
	var test tester
	err := Unmarshal(JSONDecoderFunc([]byte(jsonStr)), &test)
	if err != nil {
		t.Errorf("failed to unmarshal json: %s", err)
		return
	}
	if test.Field1 != "value1" {
		t.Errorf("test.Field1 -> got '%s' want '%s'", test.Field1, "value1")
	}
	if test.Field2 != 2 {
		t.Errorf("test.Field2 -> got '%d' want '%d'", test.Field2, 2)
	}
	if test.Field3[0] != "value" || test.Field3[1] != "three" {
		t.Errorf("test.Field3 -> got [%s] want [%s]", strings.Join(test.Field3, ","), strings.Join([]string{"value", "three"}, ","))
	}
	var data []byte
	data, err = json.Marshal(Marshal(test))
	if err != nil {
		t.Errorf("failed to marshal json: %s", err)
		return
	}
	if string(data) != jsonStr {
		t.Errorf("marshal json -> got '%s'; want '%s'", string(data), jsonStr)
	}
}
