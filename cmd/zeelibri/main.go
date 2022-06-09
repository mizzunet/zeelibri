package main

import (
	"flag"
	"fmt"
	// "log"
	"os"
	"strings"

	z "zeelibri"
)

func main() {
	var b z.Book
	Path := flag.String("p", "./", "Download path")
	IsFirst := flag.Bool("f", false, "Download first match")
	Items := flag.Int("n", 10, "Number of items to display")

	// Search
	Query := strings.Join(os.Args[1:], " ")
	if Query == "" {
		panic("Search for something\n")
	}
	books, err := z.Search(Query)
	if err != nil {
		panic(err)
	}

	// Pre-download
	var number int
	if *IsFirst == true {
		number = 0
	} else {
		for i := 0; i < *Items; i++ {
			fmt.Printf("%v. %s by %s\n", i+1, books[i].Title, books[i].Author)
			fmt.Printf("   %v %s - %s\n", books[i].Size.Bytes, books[i].Size.Unit, books[i].Format)
		}
		fmt.Scanf("\n > %s", &number)
	}

	// Download
	b = books[number]
	fmt.Printf("Downloading...\n")
	fmt.Printf("Title : %s\n", b.Title)
	fmt.Printf("Author: %s\n", b.Author)
	fmt.Printf("Size: %v %s %s\n", b.Size.Bytes, b.Size.Unit, b.Format)
	err = b.Download(*Path)
	if err != nil {
		panic(err)
	}
}
