// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

func main() {
	r := gin.Default()

	r.POST("/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errs := err.(validator.ValidationErrors)
			c.JSON(http.StatusBadRequest, gin.H{"error": errs.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
	})

	r.Run()
}
