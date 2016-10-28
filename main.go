package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/kpawlik/geojson"
)

type Coordinate struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func (c Coordinate) ToGeojsonCoordinate() (coord geojson.Coordinate, err error) {
	lat, err := strconv.ParseFloat(c.Latitude, 64)
	if err != nil {
		return
	}

	lon, err := strconv.ParseFloat(c.Longitude, 64)
	if err != nil {
		return
	}

	return geojson.Coordinate{
		geojson.CoordType(lon), // IN LON, LAT ORDER, Yes...
		geojson.CoordType(lat),
	}, nil
}

func main() {
	coordFilePtr := flag.String("f", "coordinates.json", "A file of coordinates")
	flag.Parse()
	byt, err := ioutil.ReadFile(*coordFilePtr)
	if err != nil {
		fmt.Println("Failed to read file:", err)
		return
	}

	coords := strings.Split(string(byt), "\n")
	fc := geojson.NewFeatureCollection([]*geojson.Feature{})

	for line, coordRaw := range coords {
		var coord Coordinate
		err = json.Unmarshal([]byte(coordRaw), &coord)
		if err != nil {
			fmt.Printf("Failed to unmarshal coordinate on file line %v: %v\n", line, err)
			return
		}

		c, err := coord.ToGeojsonCoordinate()
		if err != nil {
			fmt.Println("Could not read float from coordinate:", err)
			return
		}
		p := geojson.NewPoint(c)
		feature := geojson.NewFeature(p, nil, nil)
		fc.AddFeatures(feature)
	}

	geo, err := geojson.Marshal(fc)
	if err != nil {
		fmt.Println("Failed to marshal geojson:", err)
		return
	}

	fmt.Println(geo)
}
