package utils

import (
	"strings"
	"testing"
)

func TestFullPath(t *testing.T) {
	p, err := FullPath("reltest.txt")
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	if !strings.HasPrefix(p, "/") {
		t.Error("expected an absolute path, got ", p)
	}
	if !strings.HasSuffix(p, "reltest.txt") {
		t.Error("expected path to end with filename, got ", p)
	}

	p, err = FullPath("~/hometest.txt")
	if err != nil {
		t.Error("Expected no error, got", err)
	}

	if !strings.HasPrefix(p, "/home/") && !strings.HasPrefix(p, "/Users/") {
		t.Error("expected an absolute path, got ", p)
	}
	if !strings.HasSuffix(p, "hometest.txt") {
		t.Error("expected path to end with filename, got ", p)
	}

}
