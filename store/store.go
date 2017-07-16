package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/go-kit/kit/log"

	"net/http"

	"github.com/dewey/go-radiooooo/radiooooo"
)

// Archive contains everything about our storage backend
type Archive struct {
	Path string
	Log  log.Logger
}

// NewArchive creates a new Archive backend
func NewArchive(log log.Logger, path string) *Archive {
	return &Archive{
		Path: path,
		Log:  log,
	}
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
		return errors.New("song has empty UUID")
	}
	path := path.Join(a.Path, country, strconv.Itoa(year), s.UUID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0700)
		if err != nil {
			return err
		}
		if err := a.writeSongFile(path, s); err != nil {
			return err
		}
		return nil
	}
	return nil
}

// writeSongDetails downloads the assets and stores the meta data along side the audio files
func (a *Archive) writeSongFile(parentPath string, s radiooooo.Song) error {
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
		a.Log.Log("msg", "successfully archived song", "album", s.Album, "artist", s.Artists, "year", s.Year)
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

// ArchiveInfo contains everything in our archive
type ArchiveInfo struct {
	Countries      []Country
	CountriesTotal int
}

// Decade contains the information about a decade
type Decade struct {
	Name      string
	MP3Total  int
	JSONTotal int
}

func (d *Decade) String() string {
	return fmt.Sprintf(" - %s: [JSON: %d, MP3: %d]", d.Name, d.JSONTotal, d.MP3Total)
}

// Country contains all the information of a country
type Country struct {
	Name         string
	Decades      []Decade
	DecadesTotal int
}

func (c *Country) String() string {
	return fmt.Sprintf("%s", c.Name)

}

// GetArchiveInfo returns information about our local archive
func (a *Archive) GetArchiveInfo() (ArchiveInfo, error) {
	var ai ArchiveInfo
	countries, err := ioutil.ReadDir(a.Path)
	if err != nil {
		return ai, err
	}

	var cc []Country
	for _, c := range countries {
		// TODO: refactor to ignore every dot file
		if c.Name() == ".DS_Store" {
			continue
		}
		decades, err := ioutil.ReadDir(path.Join(a.Path, c.Name()))
		if err != nil {
			return ai, err
		}
		var dc []Decade
		for _, d := range decades {
			if d.Name() == ".DS_Store" {
				continue
			}
			trackDirs, err := ioutil.ReadDir(path.Join(a.Path, c.Name(), d.Name()))
			if err != nil {
				return ai, err
			}
			var mp3Counter, jsonCounter int
			for _, t := range trackDirs {
				if t.Name() == ".DS_Store" {
					continue
				}
				files, err := ioutil.ReadDir(path.Join(a.Path, c.Name(), d.Name(), t.Name()))
				if err != nil {
					return ai, err
				}
				for _, f := range files {
					if f.Name() == ".DS_Store" {
						continue
					}
					if path.Ext(f.Name()) == ".mp3" {
						mp3Counter++
					}
					if path.Ext(f.Name()) == ".json" {
						jsonCounter++
					}
				}

			}
			dc = append(dc, Decade{
				Name:      d.Name(),
				MP3Total:  mp3Counter,
				JSONTotal: jsonCounter,
			})
		}
		cc = append(cc, Country{Name: c.Name(), Decades: dc, DecadesTotal: len(dc)})
	}
	return ArchiveInfo{Countries: cc, CountriesTotal: len(cc)}, nil
}
