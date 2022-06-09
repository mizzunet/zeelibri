package zeelibri

import (
	"errors"
	// "fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gocolly/colly/v2"
	"github.com/headzoo/surf"
)

type Book struct {
	Title  string
	Author string
	URL    string
	Size   struct {
		Bytes int
		Unit  string
	}
	Format string
}

const Filters = "?extensions[]=epub"

const Domain = "https://1lib.in/s/"

// const Domain = "https://u1lib.org/s/"

func Search(Query string) ([]Book, error) {
	var book Book
	var Results []Book
	URL := Domain + strings.ReplaceAll(Query, " ", "+") + Filters

	c := colly.NewCollector(colly.CacheDir("cache"))
	// Grab details from search results
	c.OnHTML("#searchResultBox > .resItemBox", func(e *colly.HTMLElement) {
		book.Title = strings.TrimSpace(e.ChildText("h3[itemprop='name']"))
		book.URL = e.Request.AbsoluteURL(e.ChildAttr("h3[itemprop='name'] > a", "href"))
		book.Author = strings.TrimSpace(e.ChildText("a[itemprop='author']"))
		FormatAndSize := strings.Split(e.ChildText("div.bookProperty.property__file > div.property_value"), " ")
		book.Format = strings.ReplaceAll(FormatAndSize[0], ",", "")
		book.Size.Bytes, _ = strconv.Atoi(FormatAndSize[1])
		book.Size.Unit = FormatAndSize[2]

		Results = append(Results, book)
	})
	c.Visit(URL)

	// error on no results
	if len(Results) == 0 {
		return Results, errors.New("Coudln't detect any books " + URL)
	}

	return Results, nil
}

func (b Book) Download(Path string) error {
	s := surf.NewBrowser()
	Path = Path + "/"
	FullPath := Path + strings.ReplaceAll(b.Title, " ", "_") + "." + b.Format
	s.Open(b.URL)

	// if b.Size.Unit == "MB" {
	// 	if b.Size.Bytes > 8 {
	// 		errors.New("Can't download as the size is higher")
	// 	}
	// }

	// Detect download button
	_, bool := s.Find("a.addDownloadedBook").Attr("href")
	if bool == false {
		return errors.New("Failed to get downoad link form " + b.URL)
	}

	// Click on dowlnoad
	s.Click("a.addDownloadedBook")

	// Check if limit reached
	if s.Find(".download-limits-error__header").Text() == "Daily limit reached" {
		return errors.New("Daily limit reached")
	}

	// Folder making
	if _, err := os.Stat(Path); os.IsNotExist(err) {
		os.Mkdir(Path, 0755)
	}

	// Check if file exists
	if _, err := os.Stat(FullPath); !os.IsNotExist(err) {
		return errors.New("File already been downloaded at " + FullPath)
	}

	// Create file
	File, err := os.Create(FullPath)
	if err != nil {
		return errors.New("Failed to create file " + err.Error())
	}

	// Download the book
	_, err = s.Download(File)
	if err != nil {
		return errors.New("Failed to download: " + err.Error())
	}

	defer File.Close()

	return nil
}
