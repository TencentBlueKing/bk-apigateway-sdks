package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Pet example
type Pet struct {
	ID       int `json:"id"`
	Category struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"category"`
	Name      string   `json:"name"`
	PhotoUrls []string `json:"photoUrls"`
	Status    string   `json:"status"`
}

// GetStringByInt example
//
//	@Summary		Add a new pet to the store
//	@Description	get string by ID
//	@ID				get-string-by-int
//	@Accept			json
//	@Produce		json
//	@Param			some_id	path		int		true	"Some ID"
//	@Param			some_id	body		Pet		true	"Some ID"
//	@Success		200		{string}	string	"ok"
//	@Router			/testapi/get-string-by-int/{some_id} [get]
func GetStringByInt(c *gin.Context) {
	pet := Pet{
		ID: 1,
		Category: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			ID:   1,
			Name: "dog",
		},
		Name:      "dog",
		PhotoUrls: []string{"123"},
		Status:    "available",
	}
	c.JSON(http.StatusOK, gin.H{
		"data": pet,
	})
}
