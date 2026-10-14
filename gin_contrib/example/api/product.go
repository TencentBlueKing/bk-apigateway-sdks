/**
 * TencentBlueKing is pleased to support the open source community by
 * making 蓝鲸智云-蓝鲸 PaaS 平台(BlueKing-PaaS) available.
 * Copyright (C) 2025 Tencent. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 */

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
//	@Router		/update-product/{product_id} [post]
func UpdateProduct(c *gin.Context) {
	var pUpdates ProductUpdates
	if err := c.ShouldBindJSON(&pUpdates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{})
}
