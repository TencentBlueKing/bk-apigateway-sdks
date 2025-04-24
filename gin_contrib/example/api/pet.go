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

// GetPetByID example
//
//	@Summary		get a  pet
//	@Description	get pet by ID
//	@ID				get_pet_by_id
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"pet id"
//	@Success		200	{string}	Pet
//	@Router			/testapi/pets/{id}/ [get]
func GetPetByID(c *gin.Context) {
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
