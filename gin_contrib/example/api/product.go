package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ProductUpdates struct {
	Type        sql.NullString `json:"type"`
	Description sql.NullString `json:"description"`
	Stock       sql.NullInt64  `json:"stock"`
}

// UpdateProduct example
//
//	@Summary	Update product attributes
//	@ID			update_product
//	@Accept		json
//	@Param		product_id	path	int				true	"Product ID"
//	@Param		productInfo	body	ProductUpdates	true	" "
//	@Router		/testapi/update-product/{product_id} [post]
func UpdateProduct(c *gin.Context) {
	var pUpdates ProductUpdates
	if err := c.ShouldBindJSON(&pUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}
