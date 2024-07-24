package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/upload", uploadFile)
	r.Run(":8000")
}

func uploadFile(c *gin.Context) {
	f, err := c.FormFile("file")
	if err != nil {
		fmt.Println("Error while uploading file")
		return
	}
	filepath := filepath.Join("uploads", f.Filename)
	if err := c.SaveUploadedFile(f, filepath); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save file",
		})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"Status": "Saved",
	})
}
