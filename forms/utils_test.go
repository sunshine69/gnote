package forms

import (
	"fmt"
	"testing"
)

func TestUtils(t *testing.T) {
	o := ticks2Time("1572337507")
	fmt.Printf("%v", o)
}