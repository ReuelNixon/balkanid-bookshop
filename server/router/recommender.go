package router

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Book struct {
	Title string `json:"title"`
}

func BookRecommendationEngine(bookName string, idlist [][]int, books []Book) []string {
	var bookListName []string
	// Find book index
	var bookIndex int
	for i, book := range books {
		if book.Title == bookName {
			bookIndex = i
			break
		}
	}	

	// Get recommended book titles
	for _, newID := range idlist[bookIndex] {
		bookListName = append(bookListName, books[newID].Title)
	}

	return bookListName
}

func GetRecommendationsHelper(bookName string) []string {
	currentDir, err := os.Getwd()
    if err != nil {
        fmt.Println("Error getting current directory:", err)
		os.Exit(1)
    }

	// Load idlist from JSON
	idlistFile, err := os.ReadFile(filepath.Join(currentDir, "/rec_data/idlist.json"))
	if err != nil {
		fmt.Println("Error reading idlist.json:", err)
		os.Exit(1)
	}

	var idlist [][]int
	err = json.Unmarshal(idlistFile, &idlist)
	if err != nil {
		fmt.Println("Error unmarshaling idlist:", err)
		os.Exit(1)
	}

	// Load new_data from JSON
	dataFile, err := os.ReadFile(filepath.Join(currentDir, "/rec_data/data.json"))
	if err != nil {
		fmt.Println("Error reading data.json:", err)
		os.Exit(1)
	}

	var books []Book
	err = json.Unmarshal(dataFile, &books)
	if err != nil {
		fmt.Println("Error unmarshaling data:", err)
		os.Exit(1)
	}

	bookListName := BookRecommendationEngine(bookName, idlist, books)

	return bookListName
}
