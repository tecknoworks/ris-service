package ris

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	"github.com/PuerkitoBio/goquery"
)

const (
	baseurl = "https://www.google.com"
)

// requestParams : Parameters for fetchURL
type requestParams struct {
	Method      string
	URL         string
	Contenttype string
	Data        io.Reader
	Client      *http.Client
}

type RisResult struct {
	Title		string
	Subtitle 	string
	Description	string
	ImageUrl 	string
	Links		[]map[string]string
}

// Imgdata : Image URL
type Imgdata struct {
	OU      string `json:"ou"`
	WebPage bool
}

// fetchURL : Fetch method
func (r *requestParams) fetchURL() *http.Response {
	req, err := http.NewRequest(
		r.Method,
		r.URL,
		r.Data,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		os.Exit(1)
	}
	if len(r.Contenttype) > 0 {
		req.Header.Set("Content-Type", r.Contenttype)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 Firefox/26.0")
	res, _ := r.Client.Do(req)
	return res
}

func ImgFromFile(file string) RisResult {
	// var url string
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	fs, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		os.Exit(1)
	}
	defer fs.Close()
	data, err := w.CreateFormFile("encoded_image", file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		os.Exit(1)
	}
	if _, err = io.Copy(data, fs); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v. ", err)
		os.Exit(1)
	}
	w.Close()
	r := &requestParams{
		Method: "POST",
		URL:    baseurl + "/searchbyimage/upload",
		Data:   &b,
		Client: &http.Client{
			Timeout:       time.Duration(10) * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error { return errors.New("Redirect") },
		},
		Contenttype: w.FormDataContentType(),
	}
	var res *http.Response
	for {
		res = r.fetchURL()
		if res.StatusCode == 200 {
			break
		}
		reurl, _ := res.Location()
		r.URL = reurl.String()
		r.Method = "GET"
		r.Data = nil
		r.Contenttype = ""
	}
	defer res.Body.Close()
	doc, _ := goquery.NewDocumentFromResponse(res)

	//web scrapping the result
	var result RisResult
	doc.Find(".SPZz6b").Children().Each(func(i int, s *goquery.Selection){
		if i==0{
			result.Title=s.Text()
		}else if i==1 {
			result.Subtitle=s.Text()
		}
	})
	doc.Find(".i4J0ge").Each(func(i int, s *goquery.Selection){
		
		s.Find(".bNg8Rb").Each(func(i int,s *goquery.Selection){
			s.Parent().Children().Each(func(i int, s *goquery.Selection){
				result.Description+=s.Text()+"<br>"
			})
			result.Description+="<br>"
			
		})
	})
	result.ImageUrl, _= doc.Find("#dimg_1").Attr("src")
	result.Links = getWebPages(doc)

	return result
}

func getWebPages(doc *goquery.Document) []map[string]string {
	var ar []map[string]string
	doc.Find(".r").Each(func(i int, s *goquery.Selection) {

		linkMap:=make(map[string]string)
		s.Find("a").Each(func(i int, s *goquery.Selection) {
			if i==0{
				url, _ := s.Attr("href")
				linkMap["url"]=url
			}
		})

		s.Find(".LC20lb").Each(func(i int, s *goquery.Selection){
			s.Find(".ellip").Each(func(i int, s *goquery.Selection){
				linkMap["title"]=s.Text()
			})
		})
		ar = append(ar, linkMap)
	})
	return ar
}