package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
)

func musicHand(c *gin.Context) {
	c.HTML(http.StatusOK, "music.tmpl", gin.H{
		"Name": "Music Library", "Contents": getDirContent("", false),
	})
}

func artistHand(c *gin.Context) {
	c.HTML(http.StatusOK, "artist.tmpl", gin.H{
		"Name": c.Param("artist"), "Contents": getDirContent(c.Param("artist"), false),
		"Link": c.Request.URL.Path,
	})
}

func albumHand(c *gin.Context) {
	c.HTML(http.StatusOK, "album.tmpl", gin.H{
		"Name": c.Param("album"), "Contents": getDirContent(
			c.Param("artist")+"/"+c.Param("album"), true,
		),
		"Link": c.Request.URL.Path,
	})
}

func downloadZipHand(c *gin.Context) {
	album := c.Param("album")
	artist := c.Param("artist")
	folder := fmt.Sprintf("%s/%s/%s", MusicDir, artist, album)
	filename := folder + ".zip"

	// create zip if not exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := ZipFiles(filename,
			addFullPath(folder, getDirContent(artist+"/"+album, true)),
		)
		if err != nil {
			println("error", err.Error())
			return
		}
	}

	header := c.Writer.Header()
	header["Content-type"] = []string{"application/zip"}
	header["Content-Disposition"] = []string{`attachment; filename= "` + album + `.zip"`}

	file, err := os.Open(filename)
	if err != nil {
		c.String(http.StatusNotFound, "%v", err)
		return
	}
	defer file.Close()

	io.Copy(c.Writer, file)

}

func errorHand(c *gin.Context, msg string) {
	c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{
		"Message": msg,
	})
}

func messageHand(c *gin.Context, msg string) {
	c.HTML(http.StatusOK, "message.tmpl", gin.H{
		"Message": msg,
	})
}

func downloadHand(c *gin.Context) {

	url := c.Query("url")

	if url != "" {
		qType, qId := getValues(url)
		if qType == "" {
			errorHand(c, "Can't get URL")
			return
		}
		if qType == "artist" {
			errorHand(c, "Artist downloads disabled")
			return
		}
		err := download(qType, qId)
		if err != nil {
			errorHand(c, "Download failed")
			log.Printf("Command finished with error: %v", err)
			return
		}
		messageHand(c, "Downloaded")
	} else {
		c.HTML(http.StatusOK, "download.tmpl", gin.H{})
	}

}
