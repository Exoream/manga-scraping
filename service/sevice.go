package service

import (
	"fmt"
	"log"
	"manga/dto"
	"manga/model"
	"regexp"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
)

type MangaService struct{}

func (s *MangaService) ScrapeMangaData() ([]model.Manga, error) {
	c := colly.NewCollector()

	var mangas []model.Manga

	c.OnHTML(".content-homepage-item", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		idParts := strings.Split(link, "-")
		id := idParts[len(idParts)-1]

		name := e.ChildText("h3")
		author := e.ChildText(".item-author")
		rating := e.ChildText(".item-rate")
		image := e.ChildAttr("img", "src")

		manga := model.Manga{
			ID:     id,
			Name:   name,
			Author: author,
			Rating: rating,
			Image:  image,
		}
		mangas = append(mangas, manga)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit("https://manganato.com/")
	if err != nil {
		return nil, err
	}

	return mangas, nil
}

func (s *MangaService) GetMangaByID(id string) (*model.MangaDetails, error) {
	c := colly.NewCollector()

	var mangaDetails model.MangaDetails

	// URL for chapmanganato.to
	chapManganatoURL := "https://chapmanganato.to/manga-" + id
	// URL for manganato.com
	manganatoURL := "https://manganato.com/manga-" + id

	c.OnHTML(".story-info-right", func(e *colly.HTMLElement) {
		mangaDetails.Name = e.ChildText("h1")
	})

	c.OnHTML(".variations-tableInfo", func(e *colly.HTMLElement) {
		mangaDetails.AlternativeName = e.ChildText("tr:nth-child(1) .table-value")
		mangaDetails.Author = e.ChildText("tr:nth-child(2) .table-value")
		mangaDetails.Status = e.ChildText("tr:nth-child(3) .table-value")

		e.ForEach("tr:nth-child(4) .table-value a", func(_ int, elem *colly.HTMLElement) {
			genre := elem.Text
			mangaDetails.Genre = append(mangaDetails.Genre, genre)
		})
	})

	c.OnHTML(".story-info-right-extent", func(e *colly.HTMLElement) {
		mangaDetails.Updated = e.ChildText("p:nth-child(1) .stre-value")
		mangaDetails.View = e.ChildText("p:nth-child(2) .stre-value")
		mangaDetails.Rating = s.GetMangaRating(e.ChildText("em#rate_row_cmd"))
	})

	c.OnHTML(".panel-story-info-description", func(e *colly.HTMLElement) {
		s.GetMangaDescription(e.Text, &mangaDetails)
	})

	c.OnHTML(".row-content-chapter li.a-h", func(h *colly.HTMLElement) {
		chapter := model.Chapter{
			Title:    h.ChildText("a.chapter-name"),
			Number:   h.ChildText("a.chapter-name"),
			URL:      h.ChildAttr("a.chapter-name", "href"),
			Uploaded: h.ChildText("span.chapter-time"),
		}
		mangaDetails.Chapters = append(mangaDetails.Chapters, chapter)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	errChapManganato := c.Visit(chapManganatoURL)
	if errChapManganato != nil {
		errManganato := c.Visit(manganatoURL)
		if errManganato != nil {
			return nil, errManganato
		}
	}

	return &mangaDetails, nil
}

func (s *MangaService) SearchManga(query string) ([]model.Manga, error) {
	c := colly.NewCollector()

	var mangas []model.Manga

	// Assuming the URL structure for manga search is like: https://manganato.com/search/story/{search}
	url := "https://manganato.com/search/story/" + query

	c.OnHTML(".search-story-item", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		idParts := strings.Split(link, "-")
		id := idParts[len(idParts)-1]

		name := e.ChildText("h3 a")
		img := e.ChildAttr("img", "src")

		combinedInfo := e.ChildText(".item-time")
		re := regexp.MustCompile(`Updated : (.+?)View : (.+)`)
		matches := re.FindStringSubmatch(combinedInfo)

		var updated, view string
		if len(matches) == 3 {
			updated = strings.TrimSpace(matches[1])
			view = strings.TrimSpace(matches[2])
		}
		author := e.ChildText(".item-author")
		rating := e.ChildText(".item-rate")

		manga := model.Manga{
			ID:      id,
			Name:    name,
			Author:  author,
			Image:   img,
			Updated: updated,
			View:    view,
			Rating:  rating,
		}
		mangas = append(mangas, manga)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(url)
	if err != nil {
		return nil, err
	}

	return mangas, nil
}

func (s *MangaService) ScrapeMangaWithPagination(pageNumber int) (dto.PaginationResponse, error) {
	c := colly.NewCollector()

	var mangas []model.Manga
	baseURL := fmt.Sprintf("https://manganato.com/genre-all/%d", pageNumber)
	c.OnHTML(".content-genres-item", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		idParts := strings.Split(link, "-")
		id := idParts[len(idParts)-1]

		name := e.ChildText("h3")
		author := e.ChildText(".genres-item-author")
		rating := e.ChildText(".genres-item-rate")
		image := e.ChildAttr("img", "src")
		updated := e.ChildText(".genres-item-time")
		view := e.ChildText(".genres-item-view")

		manga := model.Manga{
			ID:      id,
			Name:    name,
			Author:  author,
			Rating:  rating,
			Image:   image,
			Updated: updated,
			View:    view,
		}
		mangas = append(mangas, manga)
	})

	var lastPages int
	c.OnHTML(".pagination li", func(e *colly.HTMLElement) {
		pageNum, err := strconv.Atoi(e.Text)
		if err == nil && pageNum > lastPages {
			lastPages = pageNum
		}
	})

	// Check if it is the last page link
	c.OnHTML(".page-blue.page-last", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		lastPages, _ = extractLastPageNumber(href)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(baseURL)
	if err != nil {
		return dto.PaginationResponse{}, err
	}

	paginationResponse := dto.PaginationResponse{
		CurrentPage: pageNumber,
		LastPage:    lastPages,
		Mangas:      mangas,
	}

	return paginationResponse, nil
}

func (s *MangaService) ScrapeMangaByTopViewWithPagination(pageNumber int) ([]model.Manga, error) {
	c := colly.NewCollector()

	var mangas []model.Manga
	baseURL := fmt.Sprintf("https://manganato.com/genre-all/%d?type=topview", pageNumber)
	c.OnHTML(".content-genres-item", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		idParts := strings.Split(link, "-")
		id := idParts[len(idParts)-1]

		name := e.ChildText("h3")
		author := e.ChildText(".genres-item-author")
		rating := e.ChildText(".genres-item-rate")
		image := e.ChildAttr("img", "src")
		updated := e.ChildText(".genres-item-time")
		view := e.ChildText(".genres-item-view")

		manga := model.Manga{
			ID:      id,
			Name:    name,
			Author:  author,
			Rating:  rating,
			Image:   image,
			Updated: updated,
			View:    view,
		}
		mangas = append(mangas, manga)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(baseURL)
	if err != nil {
		return nil, err
	}

	return mangas, nil
}

func (s *MangaService) ScrapeMangaByNewestWithPagination(pageNumber int) ([]model.Manga, error) {
	c := colly.NewCollector()

	var mangas []model.Manga
	baseURL := fmt.Sprintf("https://manganato.com/genre-all/%d?type=newest", pageNumber)
	c.OnHTML(".content-genres-item", func(e *colly.HTMLElement) {
		link := e.ChildAttr("a", "href")
		idParts := strings.Split(link, "-")
		id := idParts[len(idParts)-1]

		name := e.ChildText("h3")
		author := e.ChildText(".genres-item-author")
		rating := e.ChildText(".genres-item-rate")
		image := e.ChildAttr("img", "src")
		updated := e.ChildText(".genres-item-time")
		view := e.ChildText(".genres-item-view")

		manga := model.Manga{
			ID:      id,
			Name:    name,
			Author:  author,
			Rating:  rating,
			Image:   image,
			Updated: updated,
			View:    view,
		}
		mangas = append(mangas, manga)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	err := c.Visit(baseURL)
	if err != nil {
		return nil, err
	}

	return mangas, nil
}

func extractLastPageNumber(link string) (int, error) {
	// Mencari indeks awal angka setelah tanda '/'
	startIndex := strings.LastIndex(link, "/") + 1

	// Mengambil substring yang berisi angka halaman
	lastPageStr := link[startIndex:]

	// Konversi string ke integer
	lastPageNumber, err := strconv.Atoi(lastPageStr)
	if err != nil {
		return 0, err
	}

	return lastPageNumber, nil
}

func (s *MangaService) GetMangaRating(rating string) string {
	// Split the string into words
	words := strings.Fields(rating)

	// Extract the relevant part of the rating and join them
	tmp := words[3:]
	return strings.Join(tmp, " ")
}

func (s *MangaService) GetMangaDescription(desc string, mangaDetails *model.MangaDetails) {
	pref := "Description :\n"
	desc = strings.Trim(desc, "\n")
	mangaDetails.Description = strings.TrimPrefix(desc, pref)
}
