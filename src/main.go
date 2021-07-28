package main

import (
    "fmt"
    "io"
    "net/http"
    "os"
    "./ris"
    "encoding/json"
    "github.com/go-http-utils/logger"
)

const (
    imagesFolder = "./images"
    videoServiceUrl = "http://video-service:3003/videos"
)

func CreateDirIfNotExist(dir string) {
      if _, err := os.Stat(dir); os.IsNotExist(err) {
              err = os.MkdirAll(dir, 0755)
              if err != nil {
                      panic(err)
              }
      }
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func risHandler(w http.ResponseWriter, r *http.Request) {
    enableCors(&w)
    if (*r).Method == "OPTIONS" {
		return
    }
    
    videoId:= r.URL.Query()["videoId"][0]
    offset := r.URL.Query()["offset"][0]
    getVideoFrame(videoId, offset)
    result:=ris.ImgFromFile(imagesFolder+"/videoframe.jpg")
    encodedResult, err := json.Marshal(result)
    if  err!=nil { 
        return  
    }
    w.Header().Set("Content-Type", "application/json")
    w.Write(encodedResult)
}

func getVideoFrame(videoId string, offset string) {
    req, _ := http.NewRequest("GET", videoServiceUrl + "/frame", nil)

	q := req.URL.Query()
	q.Add("videoId", videoId)
	q.Add("offset", offset)
    req.URL.RawQuery = q.Encode()
    
    client := &http.Client{}

    resp, err := client.Do(req)
    if err != nil {
        fmt.Println("Error getting the frame.")
        return
    }
    defer resp.Body.Close()

    // Create the file
    out, err := os.Create(imagesFolder+"/videoframe.jpg")
    if err != nil {
        return
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)

    if err != nil {
		fmt.Println("Error while downloading")
		return
	}
}


func setupRoutes() {
    mux := http.NewServeMux()
    mux.HandleFunc("/ris", risHandler)
    http.ListenAndServe(":3006", logger.Handler(mux, os.Stdout, logger.DevLoggerType))
}

func main() {
    CreateDirIfNotExist(imagesFolder)
    fmt.Println("ris-service is listening on port 3006.")
    setupRoutes()
}