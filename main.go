package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/kavkaz/pavo/attachment"
	"github.com/kavkaz/pavo/upload"
)

func main() {
	root_storage := "./"
	host := "localhost:9073"

	r := gin.Default()
	r.Use(CORSMiddleware())
	r.Use(static.Serve(root_storage))

	r.POST("/files", CreateAttachment)

	log.Printf("Storage place in: %s", root_storage)
	log.Printf("Start server on %s", host)
	r.Run(host)

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
	converts["thumbnail"] = "120x90"

	files, err := upload.Process(c.Request, "./")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "error",
			"error":  fmt.Sprintf("Upload error: %q", err.Error()),
		})
	}

	data := make([]map[string]interface{}, 0)
	for _, ofile := range files {
		attachment, err := attachment.Create("./example", ofile, converts)
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
			"Content-Type, Content-Length, Accept-Encoding, Content-Ragne, Content-Disposition, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.Abort(200)
			return
		}
		// c.Next()
	}
}
