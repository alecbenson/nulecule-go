package parser

import "testing"

func TestIsJson(t *testing.T) {
	testJSON := []byte("{'test': 'true'}")
	testJSON2 := []byte("    {'test': 'true'}")
	testNotJSON := []byte("notJson")

	if !isJSON(testJSON) {
		t.Fatalf("isJson returned %t for test json with data '%s', expected %t", false, string(testJSON), true)
	}
	if !isJSON(testJSON2) {
		t.Fatalf("isJson returned %t for test json with data '%s', expected %t", false, string(testJSON2), true)
	}
	if isJSON(testNotJSON) {
		t.Fatalf("isJson returned %t for test json with data '%s', expected %t", true, string(testNotJSON), false)
	}
}
