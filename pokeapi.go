package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	pokecache "github.com/atomk/pokedexcli/internal"
)

// Fields names MUST be exported (capitalized) otherwise unmarshaling will silently ignore them!
// If all struct fields are unexported, you will end up with an empty LocationAreasResponse object.
// encoding/json can only populate exported struct fields.
type LocationAreasResponse struct {
	Count    int
	Next     *string // nil on last page
	Previous *string // nil on first page
	Results  []LocationArea
}
type LocationArea struct {
	Name string
	Url  string
}

func getLocationAreas(url *string, cache *pokecache.Cache) (LocationAreasResponse, error) {
	if cache == nil {
		panic("cache is nil")
	}

	var stringUrl string
	if url == nil {
		fmt.Println("using default endpoint")
		stringUrl = "https://pokeapi.co/api/v2/location-area"
	} else {
		stringUrl = *url
	}

	bytes, ok := cache.Get(stringUrl)
	if !ok {
		response, err := http.Get(stringUrl)
		if err != nil {
			return LocationAreasResponse{}, err
		}
		defer response.Body.Close()

		bytes, err = io.ReadAll(response.Body)
		if err != nil {
			return LocationAreasResponse{}, err
		}

		cache.Add(stringUrl, bytes)
	}

	// Unmarshal
	var locations LocationAreasResponse
	if err := json.Unmarshal(bytes, &locations); err != nil {
		return LocationAreasResponse{}, err
	}

	return locations, nil
}

func getLocationAreas_(limit, page int) (LocationAreasResponse, error) {
	offset := page * limit
	query := fmt.Sprintf("?limit=%d&offset=%d", limit, offset)
	res, err := http.Get("https://pokeapi.co/api/v2/location-area/" + query)
	if err != nil {
		return LocationAreasResponse{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationAreasResponse{}, err
	}

	// Unmarshal
	var locations LocationAreasResponse
	if err := json.Unmarshal(body, &locations); err != nil {
		return LocationAreasResponse{}, err
	}

	return locations, nil
}
