package forms

import (
	"os"
	"github.com/gotk3/gotk3/gtk"
	"fmt"
	"testing"
)

func TestInputDialog(t *testing.T) {
	gtk.Init(&os.Args)
	o := InputDialog(map[string]interface{} {"title": "Test title"} )
	fmt.Printf("%s\n", o)
	gtk.Main()
}