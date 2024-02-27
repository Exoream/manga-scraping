package dto

import "manga/model"

type MangaResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Author string `json:"author"`
	Rating string `json:"rating"`
	Image  string `json:"image"`
}

type PaginationResponse struct {
	CurrentPage int           `json:"current_page"`
	LastPage    int           `json:"last_pages"`
	Mangas      []model.Manga `json:"mangas"`
}
