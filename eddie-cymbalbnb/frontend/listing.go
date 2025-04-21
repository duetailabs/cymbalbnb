// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"slices"

	"github.com/minherz/bnb-demo/frontend/utils"
)

// following code is synchronized with catalog/src/main/java/com/example/bnb/catalog/Listing.java
type Listing struct {
	Id              string            `json:"id"`
	Name            string            `json:"name"`
	Location        string            `json:"location"`
	Description     string            `json:"description"`
	Categories      []ListingCategory `json:"categories"`
	Price           float32           `json:"price"`
	Images          []Image           `json:"images"`
	FrontPictureURI string            `json:"frontPictureUri"`
	VideoURI        string            `json:"videoUri"`
	featured        bool
}

type Image struct {
	Label string `json:"label"`
	URI   string `json:"uri"`
}

type ListingCategory string

const (
	UnspecifiedListing ListingCategory = "UNSPECIFIED"
	HouseListing       ListingCategory = "HOUSE"
	ApartmentListing   ListingCategory = "APARTMENT"
	SharedRoomsListing ListingCategory = "SHARED_ROOMS"
	CabinListing       ListingCategory = "CABIN"
)

var DebugListings = []Listing{
	{
		Id:              "listing1",
		Name:            "Chic lakefront with hot tub, Sauna, Firepits, Dock",
		Location:        "Sherman, CT",
		Description:     "Unique opportunity to live in a mug. Feel yourself in Alice's Wonderland. Spend wonderful time living in a real mug. Not for fat people.",
		Price:           100.0,
		FrontPictureURI: "https://storage.googleapis.com/bnb-demo-static/listing/listing1/photos/large_lake.jpg",
		Categories:      []ListingCategory{ApartmentListing},
		Images: []Image{
			{URI: "https://storage.googleapis.com/bnb-demo-static/listing/listing1/photos/lake_exterior_stone.jpg", Label: "Exterior"},
			{URI: "https://storage.googleapis.com/bnb-demo-static/listing/listing1/photos/lake_primary_bedroom.jpg", Label: "Master BedRoom"},
			{URI: "https://storage.googleapis.com/bnb-demo-static/listing/listing1/photos/lake_outdoor_sauna.jpg", Label: "Sauna"},
		},
		VideoURI: "https://storage.googleapis.com/bnb-demo-static/listing/listing1/videos/tour.mp4",
	},
	{
		Id:              "listing2",
		Name:            "Stunning city retreat, in the heart of everything, trains, restaurants, quiet apt",
		Location:        "Boston, MA",
		Description:     "If you are looking to hide from sun, this is the house for you. This summer house is behind the watch, protected by glass, metal and complex mechanism of cogs and springs. Do not forget about screws. All of it is there to protect you against the sun in this wonderful summer house.",
		Price:           300.0,
		FrontPictureURI: "https://storage.googleapis.com/bnb-demo-static/listing/listing2/photos/large_nyc.jpg",
		Categories:      []ListingCategory{CabinListing},
		Images: []Image{
			{URI: "https://storage.googleapis.com/bnb-demo-static/listing/listing2/photos/nyc_primary_bedroom.jpg", Label: "Primary Bedroom"},
			{URI: "https://storage.googleapis.com/bnb-demo-static/listing/listing2/photos/nyc_ensuite_bathroom.jpg", Label: "Bathroom"},
			{URI: "https://storage.googleapis.com/bnb-demo-static/listing/listing2/photos/nyc_luxury_building.jpg", Label: "Building View"},
		},
		VideoURI: "https://storage.googleapis.com/bnb-demo-static/listing/listing2/videos/tour_tbd.mp4",
	},
	{
		Id:              "listing3",
		Name:            "Spring-ready hot tub, fire pit, s’mores and views",
		Location:        "Luray, VA",
		Description:     "The beach is there. Enjoy sand, sun and sea with our sunglasses. Yes, it is not the villa that you might expected. But... Sunglasses!!!",
		Price:           200.0,
		FrontPictureURI: "https://storage.googleapis.com/bnb-demo-static/listing/listing3/photos/large_virginia.jpg",
		Categories:      []ListingCategory{HouseListing},
		Images: []Image{
			{URI: "https://storage.googleapis.com/bnb-demo-static/listing/listing3/photos/virginia_living_room.jpg", Label: "Living Room"},
			{URI: "https://storage.googleapis.com/bnb-demo-static/listing/listing3/photos/virginia_backyard.jpg", Label: "Backyard"},
			{URI: "https://storage.googleapis.com/bnb-demo-static/listing/listing3/photos/virginia_primary_bedroom.jpg", Label: "BedRoom"},
		},
		VideoURI: "https://storage.googleapis.com/bnb-demo-static/listing/listing3/videos/tour_tbd.mp4",
	},
}

func (cl ListingCategory) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(cl))
}

func (cl *ListingCategory) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	if s == string(UnspecifiedListing) {
		*cl = UnspecifiedListing
	} else if s == string(HouseListing) {
		*cl = HouseListing
	} else if s == string(ApartmentListing) {
		*cl = ApartmentListing
	} else if s == string(SharedRoomsListing) {
		*cl = SharedRoomsListing
	} else if s == string(CabinListing) {
		*cl = CabinListing
	} else {
		return fmt.Errorf("unknown listing category: %s", s)
	}

	return nil
}

var yesValues = []string{"1", "y", "Y", "yes", "Yes", "YES"}

func localDebuggingEnabled() bool {
	s := os.Getenv("LOCAL_DEBUGGING")
	return slices.Contains(yesValues, s)
}
func listings(ctx context.Context, serviceURI string) ([]Listing, error) {
	if localDebuggingEnabled() {
		return DebugListings, nil
	}
	uri, err := url.JoinPath(serviceURI, "/listing")
	if err != nil {
		return []Listing{}, err
	}
	fmt.Printf("*** TEMP DEBUG: %s", uri)
	data, err := utils.RestCall(ctx, uri, http.MethodGet, []byte{})
	if err != nil {
		return []Listing{}, err
	}
	var listings []Listing
	if err := json.Unmarshal(data, &listings); err != nil {
		return []Listing{}, err
	}
	return listings, nil
}

func listing(ctx context.Context, serviceURI string, id string) (Listing, error) {
	if localDebuggingEnabled() {
		for _, l := range DebugListings {
			if l.Id == id {
				log.InfoContextf(ctx, "Listing %s found in local debugging mode", id)
				return l, nil
			}
		}
		return Listing{}, fmt.Errorf("listing '%s' is not found", id)
	}
	uri, err := url.JoinPath(serviceURI, "/listing/"+id)
	if err != nil {
		return Listing{}, err
	}
	data, err := utils.RestCall(ctx, uri, http.MethodGet, []byte{})
	if err != nil {
		return Listing{}, err
	}
	var l Listing
	if err := json.Unmarshal(data, &l); err != nil {
		return Listing{}, err
	}
	log.InfoContextf(ctx, "Listing %s found in service: %s", id, serviceURI)

	return l, nil
}
