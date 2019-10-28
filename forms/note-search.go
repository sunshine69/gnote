package forms

import (
	"github.com/gotk3/gotk3/gtk"
)

//NoteSearch - GUI related
type NoteSearch struct {
	w *gtk.Window
	builder *gtk.Builder
	np *NotePad
}

//NewNoteSearch - Create new  NotePad
func NewNoteSearch(np *NotePad) *NoteSearch {
	ns := &NoteSearch{np: np}
	builder, err := gtk.BuilderNewFromFile("glade/note-search.glade")
	np.builder = builder
	if err != nil {
		panic(err)
	}
	_obj, err := builder.GetObject("note_search")
	if err != nil {
		panic(err)
	}
	ns.w = _obj.(*gtk.Window)

	ns.w.ShowAll()
	return ns
}