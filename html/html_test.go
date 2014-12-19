package html

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"testing"
)

func TestExec(t *testing.T) {

	os.MkdirAll("tpl/test/test2", 0755)
	data := "<html><body>Hello {{name}}</body></html>"
	ioutil.WriteFile("tpl/test/test2/test.html", []byte(data), 0755)
	ioutil.WriteFile("tpl/test/test2/test2.html", []byte(data), 0755)

	html := New(&Options{
		Directory:  "tpl",
		Extensions: []string{".html"},
	})

	w := httptest.NewRecorder()

	html.ExecW("test/test2/test2.html", Context{"name": "test"}, w)

	assert.Equal(t, "<html><body>Hello test</body></html>", w.Body.String())

	os.RemoveAll("tpl")

}
