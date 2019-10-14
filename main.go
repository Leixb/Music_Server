package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
)

type MusicEntry struct {
	Name     string
	Contents []string
	Link     string
}

var MusicDir string

func main() {

	portPtr := flag.Int("port", 8090, "port to listen to")
	musicDirPtr := flag.String("MusicDir", "/data/Music", "Music Folder")

	flag.Parse()

	MusicDir = *musicDirPtr

	r := gin.Default()

	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"guest": "password",
	}))

	r.LoadHTMLGlob("templates/*")

	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	r.StaticFile("/style.css", "./resources/main.css")

	authorized.GET("/", func(c *gin.Context) {
		c.Request.URL.Path = "/music"
		r.HandleContext(c)
	})

	authorized.Static("/a/music", "/data/Music")

	rMusic := authorized.Group("/music")
	{
		rMusic.GET("/", musicHand)
		rMusic.GET("/:artist/", artistHand)
		rMusic.GET("/:artist/:album/", albumHand)
	}

	authorized.GET("/download", downloadHand)

	authorized.GET("/d/music/:artist/:album", downloadZipHand)

	r.Run(fmt.Sprintf(":%d", *portPtr))
}
