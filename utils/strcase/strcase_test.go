package strcase

import "testing"

func TestToDelimited(t *testing.T) {
	for _, delimiter := range []byte{'_', '-', '.', '/', 0, 0xff} {
		want := "json" + string([]byte{delimiter}) + "data" + string([]byte{delimiter}) + "2" + string([]byte{delimiter}) + "value"
		if got := ToDelimited("JSONData2.value", delimiter); got != want {
			t.Errorf("delimiter %x: got %q, want %q", delimiter, got, want)
		}

		if got := ToDelimited(" \t", delimiter); got != "" {
			t.Errorf("empty input: got %q", got)
		}
	}
}
