package main

import (
	"errors"
	"os"

	"github.com/gocarina/gocsv"
)

// station_cd,station_g_cd,station_name,station_name_k,station_name_r,line_cd,pref_cd,post,add,lon,lat,open_ymd,close_ymd,e_status,e_sort
type Station struct {
	StationCD    int     `csv:"station_cd"`
	StationGCD   int     `csv:"station_g_cd"`
	StationName  string  `csv:"station_name"`
	StationNameK string  `csv:"station_name_k"`
	StationNameR string  `csv:"station_name_r"`
	LineCD       string  `csv:"line_cd"`
	PrefCD       string  `csv:"pref_cd"`
	Post         string  `csv:"post"`
	Add          string  `csv:"add"`
	Lng          float64 `csv:"lon"`
	Lat          float64 `csv:"lat"`
	OpenYMD      string  `csv:"open_ymd"`
	CloseYMD     string  `csv:"close_ymd"`
	EStatus      string  `csv:"e_status"`
	ESort        string  `csv:"e_sort"`
}

func getStationsFromCSV(path string) (result []*Station, err error) {
	f, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()

	if err = gocsv.UnmarshalFile(f, &result); err != nil {
		return
	}
	return
}

func nearestStation(l Location, stations []*Station) (result Location, dist float64, err error) {
	if len(stations) == 0 {
		err = errors.New("staions are empty")
		return
	}

	for i, st := range stations {
		r := Location{
			Lat: st.Lat,
			Lng: st.Lng,
		}
		d := distance(r, l)
		if i == 0 || dist > d {
			dist = d
			result = r
		}
	}
	return
}
