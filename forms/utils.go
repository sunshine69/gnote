package forms

import (
	"github.com/gotk3/gotk3/gtk"
)

//MessageBox - display a message
func MessageBox(msg string) {
	d := gtk.MessageDialogNew(nil,gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, msg)
	d.Run()
	d.Destroy()
}