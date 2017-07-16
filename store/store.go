package store

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"net/http"

	"github.com/dewey/go-radiooooo/radiooooo"
)

// Archive contains everything about our storage backend
type Archive struct {
	Path string
}

// NewArchive creates a new Archive backend
func NewArchive(path string) *Archive {
	return &Archive{Path: path}
}

// WriteCountry will write a directory for our country if it doesn't exist already
func (a *Archive) WriteCountry(country string) error {
	path := path.Join(a.Path, country)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0700)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// WriteYear will write a directory for the scraped year if it doesn't exist already. A year is always
// nested in a country
func (a *Archive) WriteYear(country string, year int) error {
	if country == "" {
		return errors.New("a year always needs to be nested in a country")
	}
	if year < 1 {
		return errors.New("a year always needs to be positive")
	}
	path := path.Join(a.Path, country, strconv.Itoa(year))
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0700)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

// WriteSong will write a directory for the scraped song if it doesn't exist.
func (a *Archive) WriteSong(country string, year int, s radiooooo.Song) error {
	if country == "" {
		return errors.New("a year always needs to be nested in a country")
	}
	if year < 1 {
		return errors.New("a year always needs to be positive")
	}
	if s.UUID == "" {
		return errors.New("song is has empty UUID")
	}
	path := path.Join(a.Path, country, strconv.Itoa(year), s.UUID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0700)
		if err != nil {
			return err
		}
		if err := writeSongFile(path, s); err != nil {
			return err
		}
		return nil
	}
	return nil
}

// writeSongDetails downloads the assets and stores the meta data along side the audio files
func writeSongFile(parentPath string, s radiooooo.Song) error {
	path := path.Join(parentPath, s.UUID) + ".mp3"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		out, err := os.Create(path)
		if err != nil {
			return err
		}
		defer out.Close()
		resp, err := http.Get(s.Mp3)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		if _, err = io.Copy(out, resp.Body); err != nil {
			return err
		}
	}

	if err := writeSongMetadata(parentPath, s); err != nil {
		return err
	}
	return nil
}

// writeSongMetadata downloads the raw JSON response and stores it alongside the audio files
func writeSongMetadata(parentPath string, s radiooooo.Song) error {
	path := path.Join(parentPath, s.UUID) + ".json"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		md, err := json.Marshal(s)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(path, md, 0644); err != nil {
			return err
		}
	}
	return nil
}
