package radiooooo

// Playlist has all the properties of a playlist
type Playlist struct {
	ID int
}

// PlaylistResponse holds a Radiooooo playlist response from the API
type PlaylistResponse struct {
	Playlist  int      `json:"playlist"`
	Song      Song     `json:"song"`
	Countries []string `json:"countries"`
}
