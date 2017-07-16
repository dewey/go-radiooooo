package scrape

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dewey/go-radiooooo/radiooooo"
	"github.com/dewey/go-radiooooo/store"
)

// API defines our scraping API
type API struct {
	Endpoint string
	Client   *http.Client
	Storage  *store.Archive
}

// Country contains a country as defined by the Radiooooo API
type Country struct {
	Name string
	ISO3 string
}

// QueryPayload is used when we request something from the API
type QueryPayload struct {
	Decade  string   `json:"decade"`
	Country string   `json:"country"`
	Moods   []string `json:"moods"`
}

// CountriesInDecade contains all countries in a given decade
type CountriesInDecade map[string][]string

// getAllCountriesByDecade returns a list of all countries that are available in one decade
func (a *API) getAllCountriesByDecade(decade int, moods []string) (CountriesInDecade, error) {
	// TODO: Refactor with URL package
	resp, err := a.Client.Get(fmt.Sprintf("%s/api/playlist/countriesByTempos/%d?moods=%s", a.Endpoint, decade, strings.Join(moods, ",")))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var countries CountriesInDecade
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return nil, err
	}
	return countries, nil
}

// createNewPlaylist has to be called at the beginning so we get a playlist ID. Future requests will be done with that ID
func (a *API) createNewPlaylist(decade int, country string, moods []string) (radiooooo.Playlist, error) {
	body := new(bytes.Buffer)
	pl := QueryPayload{
		Decade:  strconv.Itoa(decade),
		Country: country,
		Moods:   moods,
	}
	if err := json.NewEncoder(body).Encode(pl); err != nil {
		return radiooooo.Playlist{}, err
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/playlist/next", a.Endpoint), body)
	if err != nil {
		return radiooooo.Playlist{}, err
	}
	// This request doesn't work without this content type
	req.Header.Set("Content-Type", "application/json")
	resp, err := a.Client.Do(req)
	if err != nil {
		return radiooooo.Playlist{}, err
	}
	defer resp.Body.Close()
	var playlist radiooooo.PlaylistResponse
	if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
		return radiooooo.Playlist{}, err
	}
	return radiooooo.Playlist{ID: playlist.Playlist}, nil
}

func (a *API) getAllTracks(playlist *radiooooo.Playlist, decade int, country string, moods []string) ([]radiooooo.Song, error) {
	// here we are collection all the songs we found
	var songs []radiooooo.Song

	// as the results seem to be randomly returned we just have to deduplicate and retry a bunch of times
	seen := make(map[string]int)

	var retry int
	for {
		body := new(bytes.Buffer)
		pl := QueryPayload{
			Decade:  strconv.Itoa(decade),
			Country: country,
			Moods:   moods,
		}
		if err := json.NewEncoder(body).Encode(pl); err != nil {
			return nil, err
		}
		req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/playlist/%d/next", a.Endpoint, playlist.ID), body)
		if err != nil {
			return nil, err
		}
		// This request doesn't work without this content type
		req.Header.Set("Content-Type", "application/json")
		resp, err := a.Client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		var sr radiooooo.SongResponse

		if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
			return nil, err
		}
		// If we haven't seen it already we add the key to our seen-list and the song to our song-list
		if _, ok := seen[sr.Song.UUID]; !ok {
			seen[sr.Song.UUID] = 0
			songs = append(songs, sr.Song)
		}
		// TODO: Each song should have a counter, if we have seen everything twice we assume we scraped everything and continue to the next search query, right now it's just a retry for 3 times
		retry++
		if retry > 3 {
			return songs, nil
		}
	}
}

// Start will start a full scrape
func (a *API) Start() (bool, error) {
	// We have to create a playlist for future requests, we just use a random country from our list for that request
	pl, err := a.createNewPlaylist(1900, "UZB", []string{"SLOW", "FAST", "WEIRD"})
	if err != nil {
		return false, err
	}
	fmt.Printf("%#v\n", pl)

	for _, d := range []int{1990} {
		countries, err := a.getAllCountriesByDecade(d, []string{"SLOW", "FAST", "WEIRD"})
		if err != nil {
			return false, err
		}
		fmt.Printf("%#v\n", countries)
		for _, c := range countries[strconv.Itoa(d)] {
			var err error
			err = a.Storage.WriteCountry(c)
			if err != nil {
				return false, errors.New("error creating new country subfolder in archive")
			}
			err = a.Storage.WriteYear(c, d)
			if err != nil {
				return false, errors.New("error creating new year folder for country")
			}
			sl, err := a.getAllTracks(&pl, d, c, []string{"SLOW", "FAST", "WEIRD"})
			if err != nil {
				return false, err
			}
			for _, i := range sl {
				err := a.Storage.WriteSong(c, d, i)
				if err != nil {
					return false, err
				}
				fmt.Println(i.ID)
			}
		}
	}

	return true, nil
}
