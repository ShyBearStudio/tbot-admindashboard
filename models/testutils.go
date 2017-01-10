package models

import (
	"testing"
)

func SkipTestIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test in short mode")
	}
}
