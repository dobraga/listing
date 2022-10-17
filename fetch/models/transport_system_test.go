package models

import (
	"fmt"
	"testing"
)

func TestTransportSystem(t *testing.T) {
	s := TransportSystem{URL: "https://pt.wikipedia.org/wiki/Metr%C3%B4_do_Rio_de_Janeiro"}
	lines, err := s.Extract()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v\n", lines)

	// pages, err := ListStations()

	// if err == nil {
	// 	t.Error(err)
	// }
	// wg := new(sync.WaitGroup)

	// for _, page := range pages {
	// 	_, err := FetchStationRetry(page)
	// 	if err == nil {
	// 		t.Error(err)
	// 	}

	// 	// wg.Add(1)

	// 	// go func(p PageStation, w *sync.WaitGroup) {
	// 	// 	defer w.Done()

	// 	// 	_, err := FetchStationRetry(p)
	// 	// 	if err == nil {
	// 	// 		t.Error(err)
	// 	// 	}
	// 	// }(page, wg)
	// }

	// wg.Wait()

}
