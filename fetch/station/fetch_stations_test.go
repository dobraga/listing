package station

import (
	"fmt"
	"testing"
)

func TestFetchStations(t *testing.T) {
	transports_systems := []TransportSystem{
		{URL: "https://pt.wikipedia.org/wiki/Metr%C3%B4_do_Rio_de_Janeiro", Uf: "RJ"},
	}

	stations, err := FetchStations(transports_systems)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", stations[0])
}
