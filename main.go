package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"code.google.com/p/go-uuid/uuid"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kavkaz/pavo/attachment"
	"github.com/kavkaz/pavo/upload"
)

func main() {
	flag.Parse()

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Use(static.Serve(*storage))

	r.POST("/files", CreateAttachment)

	log.Printf("Storage place in: %s", *storage)
	log.Printf("Start server on %s", *host)
	r.Run(*host)

}

func CreateAttachment(c *gin.Context) {
	converts, err := GetConvertParams(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  fmt.Sprintf("Query params: %s", err),
		})
		return
	}
	converts["original"] = ""

	pavo, _ := c.Request.Cookie("pavo")
	if pavo == nil {

		pavo = &http.Cookie{
			Name:    "pavo",
			Value:   uuid.New(),
			Expires: time.Now().Add(10 * 356 * 24 * time.Hour),
			Path:    "/",
		}
		c.Request.AddCookie(pavo)
		http.SetCookie(c.Writer, pavo)
	}

	files, err := upload.Process(c.Request, *storage)
	if err == upload.Incomplete {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"file":   gin.H{"size": files[0].Size},
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  fmt.Sprintf("Upload error: %q", err.Error()),
		})
		return
	}

	data := make([]map[string]interface{}, 0)
	for _, ofile := range files {
		attachment, err := attachment.Create(*storage, ofile, converts)
		if err != nil {
			data = append(data, map[string]interface{}{
				"name":  ofile.Filename,
				"size":  ofile.Size,
				"error": err.Error(),
			})
			continue
		}
		data = append(data, attachment.ToJson())
	}

	c.JSON(http.StatusCreated, gin.H{"status": "ok", "files": data})

}

// Get parameters for convert from Request query string
func GetConvertParams(req *http.Request) (map[string]string, error) {
	raw_converts := req.URL.Query().Get("converts")

	if raw_converts == "" {
		raw_converts = "{}"
	}

	convert := make(map[string]string)

	err := json.Unmarshal([]byte(raw_converts), &convert)
	if err != nil {
		return nil, err
	}

	return convert, nil
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set(
			"Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, PATCH, DELETE")
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Content-Length, Accept-Encoding, Content-Range, Content-Disposition, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.Abort(200)
			return
		}
		// c.Next()
	}
}
