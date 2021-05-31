package main

import (
	"embed"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"io/fs"
	"log"
	"net/http"
	"os"
)

//go:embed assets
var embeddedAssets embed.FS

type embedFileSystem struct {
	http.FileSystem
}

func (e embedFileSystem) Exists(_ string, path string) bool {
	_, err := e.Open(path)
	if err != nil {
		return false
	}
	return true
}

func EmbedFolder(fsEmbed embed.FS, targetPath string) static.ServeFileSystem {
	fileSystem, err := fs.Sub(fsEmbed, targetPath)
	if err != nil {
		panic(err)
	}
	return embedFileSystem{
		FileSystem: http.FS(fileSystem),
	}
}

func main() {
	r := gin.Default()

	r.Use(static.Serve("/", EmbedFolder(embeddedAssets, "assets")))

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run("127.0.0.1:" + os.Getenv("PORT")) // Only listen locally for nginx
	if err != nil {
		log.Fatal("failed to start gin server" + err.Error())
	}
}
