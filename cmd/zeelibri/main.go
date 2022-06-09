package main

import (
	"fmt"
	// "log"
	"os"
	"strings"

	z "zeelibri"
)

func main() {
	query := strings.Join(os.Args[1:], " ")
	if query == "" {
		panic("Search for something\n")
	}
	fmt.Printf("Searching for %s...\n\n", query)
	books, err := z.Search(query)
	if err != nil {
		panic(err)
	}
	b := books[0]
	fmt.Printf("Title : %s\n", b.Title)
	fmt.Printf("Author: %s\n", b.Author)
	fmt.Printf("Size: %v %s\n", b.Size.Bytes, b.Size.Unit)

	err = b.Download("download")
	if err != nil {
		panic(err)
	}
}
