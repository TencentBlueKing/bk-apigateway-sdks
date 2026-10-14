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
//	@Summary		get a pet
//	@Description	get pet by ID
//	@ID				get_pet_by_id
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"pet id"
//	@Success		200	{string}	Pet
//	@Router			/pets/{id}/ [get]
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
