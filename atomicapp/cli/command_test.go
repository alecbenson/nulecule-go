package cli

import (
	"flag"
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestGetVal(t *testing.T) {
	expected := "expectedValue"
	testFlagSet := flag.NewFlagSet("test", flag.PanicOnError)
	testFlagSet.String("testFlag", expected, "This is simply a test flag")
	result := getVal(testFlagSet, "testFlag").(string)
	if result != expected {
		logrus.Fatalf("Expected to get flag value '%s', got '%s' instead", expected, result)
	}
}
