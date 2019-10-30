package forms

import (
	"strconv"
	"time"
	"github.com/gotk3/gotk3/gtk"
)

//MessageBox - display a message
func MessageBox(msg string) {
	d := gtk.MessageDialogNew(nil,gtk.DIALOG_DESTROY_WITH_PARENT, gtk.MESSAGE_ERROR, gtk.BUTTONS_CLOSE, msg)
	d.Run()
	d.Destroy()
}

func ticks2Time(ticks string) time.Time {
    i, err := strconv.ParseInt(ticks, 10, 64)
    if err != nil {
        panic(err)
    }
    tm := time.Unix(i, 0)
    return tm
}