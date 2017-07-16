package radiooooo

// SongResponse holds a Radiooooo song response from the API
type SongResponse struct {
	Playlist int `json:"playlist"`
	Song     struct {
		ID          int         `json:"id"`
		Mp3         string      `json:"mp3"`
		Ogg         string      `json:"ogg"`
		Picture     interface{} `json:"picture"`
		Title       string      `json:"title"`
		Artists     string      `json:"artists"`
		ReleaseDate string      `json:"releaseDate"`
		Author      interface{} `json:"author"`
		Producer    interface{} `json:"producer"`
		Publisher   interface{} `json:"publisher"`
		Country     string      `json:"country"`
		Year        string      `json:"year"`
		Details     interface{} `json:"details"`
		Contributor struct {
			ID          int         `json:"id"`
			Fullname    string      `json:"fullname"`
			Surname     string      `json:"surname"`
			Avatar      string      `json:"avatar"`
			Country     string      `json:"country"`
			Ranking     float64     `json:"ranking"`
			Description string      `json:"description"`
			Email       string      `json:"email"`
			Role        string      `json:"role"`
			Birthdate   interface{} `json:"birthdate"`
			Medals      []struct {
				ID      int         `json:"id"`
				Name    string      `json:"name"`
				Picture string      `json:"picture"`
				Type    interface{} `json:"type"`
			} `json:"medals"`
			Taxidto              interface{} `json:"taxidto"`
			NbLiked              int         `json:"nbLiked"`
			LastConnexionNbLiked int         `json:"lastConnexionNbLiked"`
			DateLastArticleBlog  interface{} `json:"dateLastArticleBlog"`
		} `json:"contributor"`
		Likes      int         `json:"likes"`
		Length     int         `json:"length"`
		ItunesLink string      `json:"itunesLink"`
		UUID       string      `json:"uuid"`
		IslandID   int         `json:"islandId"`
		Album      interface{} `json:"album"`
	} `json:"song"`
	Countries []string `json:"countries"`
}

// Song contains a single Radiooooo song
type Song struct {
	ID          int         `json:"id"`
	Mp3         string      `json:"mp3"`
	Ogg         string      `json:"ogg"`
	Picture     interface{} `json:"picture"`
	Title       string      `json:"title"`
	Artists     string      `json:"artists"`
	ReleaseDate string      `json:"releaseDate"`
	Author      interface{} `json:"author"`
	Producer    interface{} `json:"producer"`
	Publisher   interface{} `json:"publisher"`
	Country     string      `json:"country"`
	Year        string      `json:"year"`
	Details     interface{} `json:"details"`
	Contributor struct {
		ID          int         `json:"id"`
		Fullname    string      `json:"fullname"`
		Surname     string      `json:"surname"`
		Avatar      string      `json:"avatar"`
		Country     string      `json:"country"`
		Ranking     float64     `json:"ranking"`
		Description string      `json:"description"`
		Email       string      `json:"email"`
		Role        string      `json:"role"`
		Birthdate   interface{} `json:"birthdate"`
		Medals      []struct {
			ID      int         `json:"id"`
			Name    string      `json:"name"`
			Picture string      `json:"picture"`
			Type    interface{} `json:"type"`
		} `json:"medals"`
		Taxidto              interface{} `json:"taxidto"`
		NbLiked              int         `json:"nbLiked"`
		LastConnexionNbLiked int         `json:"lastConnexionNbLiked"`
		DateLastArticleBlog  interface{} `json:"dateLastArticleBlog"`
	} `json:"contributor"`
	Likes      int         `json:"likes"`
	Length     int         `json:"length"`
	ItunesLink string      `json:"itunesLink"`
	UUID       string      `json:"uuid"`
	IslandID   int         `json:"islandId"`
	Album      interface{} `json:"album"`
}
