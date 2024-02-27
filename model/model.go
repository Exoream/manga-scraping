package model

type Manga struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Author  string `json:"author"`
	Rating  string `json:"rating"`
	Image   string `json:"image"`
	Updated string `json:"updated"`
	View    string `json:"view"`
}

type MangaDetails struct {
	Name            string    `json:"name"`
	AlternativeName string    `json:"alternative_name"`
	Author          string    `json:"author"`
	Status          string    `json:"status"`
	Updated         string    `json:"updated"`
	View            string    `json:"view"`
	Rating          string    `json:"rating"`
	Description     string    `json:"description"`
	Genre           []string  `json:"genre"`
	Chapters        []Chapter `json:"chapters"`
}

type Chapter struct {
	Title    string `json:"title"`
	Number   string `json:"number"`
	URL      string `json:"url"`
	Uploaded string `json:"uploaded"`
}
