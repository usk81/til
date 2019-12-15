package main

import (
	"os"
	"sort"

	"github.com/davecgh/go-spew/spew"
	"github.com/gocarina/gocsv"
)

type book struct {
	Title string `csv:"Title"`
	Count int    `csv:"Count"`
	Brand string `csv:"Brand"`
	Maker string `csv:"Maker"`
}

type rack struct {
	Total int
	Items []book
}

var (
	groupMax = 22
	rackMax  = 44
)

func main() {
	books := []*book{}
	if err := loadCSVFile("books.csv", &books); err != nil {
		panic(err)
	}

	sort.Slice(books, func(i, j int) bool {
		return books[i].Count > books[j].Count
	})

	racks := []rack{}
	for _, b := range books {
		items := []book{}
		if b.Count > groupMax {
			items = append(items, *b)
			racks = append(racks, rack{
				Total: b.Count,
				Items: items,
			})
		} else {
			var isSet bool
			for i, r := range racks {
				if r.Total >= rackMax {
					continue
				}
				if r.Total >= groupMax && r.Total+b.Count > rackMax {
					continue
				}
				r.Total += b.Count
				r.Items = append(r.Items, *b)
				racks[i] = r
				isSet = true
				break
			}
			if !isSet {
				items = append(items, *b)
				racks = append(racks, rack{
					Total: b.Count,
					Items: items,
				})
			}
		}
	}
	spew.Dump(racks)
}

func loadCSVFile(fp string, v interface{}) (err error) {
	f, err := os.OpenFile(fp, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return
	}
	defer f.Close()
	return gocsv.UnmarshalFile(f, v)
}
