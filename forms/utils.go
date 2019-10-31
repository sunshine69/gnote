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

//ChunkString - 
func ChunkString(s string, chunkSize int) []string {
	var chunks []string
	runes := []rune(s)

	if len(runes) == 0 {
		return []string{s}
	}

	for i := 0; i < len(runes); i += chunkSize {
		nn := i + chunkSize
		if nn > len(runes) {
			nn = len(runes)
		}
		chunks = append(chunks, string(runes[i:nn]))
	}
	return chunks
}