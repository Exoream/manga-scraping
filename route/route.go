package route

import (
	"manga/dto"
	"manga/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	mangaService := service.MangaService{}

	router.GET("/mangas", func(c *gin.Context) {
		mangas, err := mangaService.ScrapeMangaData()
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch manga data"})
			return
		}
		var mangaResponses []dto.MangaResponse
		for _, manga := range mangas {
			mangaResponses = append(mangaResponses, dto.MangaResponse{
				ID:     manga.ID,
				Name:   manga.Name,
				Author: manga.Author,
				Rating: manga.Rating,
				Image:  manga.Image,
			})
		}
		c.JSON(200, mangaResponses)
	})

	router.GET("/mangas/:id", func(c *gin.Context) {
		id := c.Param("id")

		// Call the GetMangaByID function
		mangaDetails, err := mangaService.GetMangaByID(id)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch manga details"})
			return
		}

		c.JSON(200, mangaDetails)
	})

	router.GET("/search", func(c *gin.Context) {
		// Retrieve the query parameter from the request
		searchQuery := c.Query("query")

		// Call the SearchManga function
		searchResults, err := mangaService.SearchManga(searchQuery)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to perform manga search"})
			return
		}

		c.JSON(200, searchResults)
	})

	router.GET("/mangas/latest/:page", func(c *gin.Context) {
		page := c.Param("page")
		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid page number"})
			return
		}

		mangas, err := mangaService.ScrapeMangaWithPagination(pageNumber)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch manga data"})
			return
		}
		c.JSON(200, mangas)
	})

	router.GET("/mangas/hot/:page", func(c *gin.Context) {
		page := c.Param("page")
		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid page number"})
			return
		}

		mangas, err := mangaService.ScrapeMangaByTopViewWithPagination(pageNumber)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch manga data"})
			return
		}
		c.JSON(200, mangas)
	})

	router.GET("/mangas/newest/:page", func(c *gin.Context) {
		page := c.Param("page")
		pageNumber, err := strconv.Atoi(page)
		if err != nil {
			c.JSON(400, gin.H{"error": "Invalid page number"})
			return
		}

		mangas, err := mangaService.ScrapeMangaByNewestWithPagination(pageNumber)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to fetch manga data"})
			return
		}
		c.JSON(200, mangas)
	})

	return router
}
