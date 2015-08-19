package nulecule

import (
	"testing"
)

func TestCheckConstraints(t *testing.T) {
	var (
		constraint Constraint
		testParam  Param
		valid      bool
		err        error
	)
	constraint = Constraint{AllowedPattern: "([A-Z-a-z])+", Description: "This is a test constraint"}
	testParam = Param{Name: "Test", Description: "Test parameter", Default: "T3st v4lu3", Constraints: []Constraint{constraint}}
	valid, err = checkConstraints(&testParam)
	if err != nil {
		t.Fatalf("CheckConstraints returned an error: %s", err)
	}
	if valid {
		t.Fatalf("Paramater with pattern %s and value %s should have failed the"+
			"constraint check but passed", constraint.AllowedPattern, testParam.Default)
	}

	testParam.Default = "TestValue"
	valid, err = checkConstraints(&testParam)
	if err != nil {
		t.Fatalf("CheckConstraints returned an error: %s", err)
	}
	if !valid {
		t.Fatalf("Paramater with pattern %s and value %s should have passed the"+
			"constraint check but failed", constraint.AllowedPattern, testParam.Default)
	}
}

func TestMakeTemplateReplacements(t *testing.T) {
	testBase := New("", "", true)
	testParam1 := Param{
		Name:    "parameter",
		Default: "test",
	}

	testParam2 := Param{
		Name:    "result",
		Default: "pass",
	}

	c := Component{Params: []Param{testParam1, testParam2}}
	template := []byte("this is a $parameter that should $result")
	data, err := testBase.makeTemplateReplacements(template, &c, false)
	if err != nil {
		t.Fatalf("Error making template replacements: %s", err)
	}

	expectedData := []byte("this is a test that should pass")
	if string(data) != string(expectedData) {
		t.Fatalf("Got back '%s' when applying template, expected '%s'", data, expectedData)
	}

}
