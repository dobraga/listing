package station

import (
	"fmt"
	"testing"
)

func TestLine(t *testing.T) {
	l := Line{URL: "https://pt.wikipedia.org/wiki/Linha_1_do_Metr%C3%B4_do_Rio_de_Janeiro"}

	stations, err := l.Extract()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v", stations)
}
