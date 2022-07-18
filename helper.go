package kmsg

import (
	"bytes"
	"text/template"

	"github.com/sebarcode/codekit"
)

func translate(source string, data codekit.M) (string, error) {
	w := bytes.NewBufferString("")
	tt, e := template.New("tmp").Parse(source)
	if e != nil {
		return source, e
	}

	e = tt.Execute(w, data)
	if e != nil {
		return source, e
	}

	return w.String(), nil
}
