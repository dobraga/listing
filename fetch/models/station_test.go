package models

import (
	"fmt"
	"testing"
)

func TestStation(t *testing.T) {
	// settings.LoadSettings()

	s := Station{URL: "https://pt.wikipedia.org/wiki/Esta%C3%A7%C3%A3o_Saens_Pe%C3%B1a_/_Tijuca"}
	err := s.Extract()
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("%+v", s)

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
