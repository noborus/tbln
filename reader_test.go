package tbln

import (
	"bytes"
	"io"
	"testing"
)

func TestRead(t *testing.T) {
	data := `
| 1 | Bob | 19
| 2 | Alice | 14
  `
	b := bytes.NewBufferString(data)
	reader := NewReader(b)
	for {
		r, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				t.Errorf("expected Read got %v", err)
			}
			break
		}
		if len(r) != 3 {
			t.Errorf("expected 3 columns got %v", r)
		}
	}
}
