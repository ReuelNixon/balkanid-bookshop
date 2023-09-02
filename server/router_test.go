package main

import (
	"bookshop/router"
	"os"
	"testing"
)

func TestCalculatePagination(t *testing.T) {
	testCases := []struct {
		page           string
		pageSize       string
		totalCount     int64
		expectedOffset int
		expectedLimit  int
	}{
		{"1", "10", 50, 0, 10},
		{"2", "10", 50, 10, 10},
		{"3", "10", 50, 20, 10},
		{"1", "5", 50, 0, 5},
		{"2", "5", 50, 5, 5},
		{"3", "5", 50, 10, 5},
		{"0", "10", 50, 0, 10},
		{"1", "0", 50, 0, 10},
		{"1", "10", 0, 0, 10},
	}

	for _, testCase := range testCases {
		offset, limit := router.CalculatePagination(testCase.page, testCase.pageSize, testCase.totalCount)
		if offset != testCase.expectedOffset || limit != testCase.expectedLimit {
			t.Errorf("calculatePagination(%s, %s, %d) = (%d, %d); expected (%d, %d)",
				testCase.page, testCase.pageSize, testCase.totalCount,
				offset, limit, testCase.expectedOffset, testCase.expectedLimit)
		}
	}
}

func TestParsePageNumber(t *testing.T) {
	testCases := []struct {
		page            string
		expectedPageNum int
	}{
		{"1", 1},
		{"10", 10},
		{"-1", 1},
		{"abc", 1},
		{"0", 1},
	}

	for _, testCase := range testCases {
		pageNum := router.ParsePageNumber(testCase.page)
		if pageNum != testCase.expectedPageNum {
			t.Errorf("parsePageNumber(%s) = %d; expected %d",
				testCase.page, pageNum, testCase.expectedPageNum)
		}
	}
}

func TestParsePageSize(t *testing.T) {
	testCases := []struct {
		pageSize            string
		expectedPageSizeNum int
	}{
		{"1", 1},
		{"10", 10},
		{"-1", 10},
		{"abc", 10},
		{"0", 10},
	}

	for _, testCase := range testCases {
		pageSizeNum := router.ParsePageSize(testCase.pageSize)
		if pageSizeNum != testCase.expectedPageSizeNum {
			t.Errorf("parsePageSize(%s) = %d; expected %d",
				testCase.pageSize, pageSizeNum, testCase.expectedPageSizeNum)
		}
	}
}

func TestBookRecommendationEngine(t *testing.T) {
	books := []router.Book{
		{Title: "Book1"},
		{Title: "Book2"},
		{Title: "Book3"},
	}

	idlist := [][]int{
		{1, 2},
		{0, 2},
		{0, 1},
	}

	testCases := []struct {
		bookName                string
		expectedRecommendations []string
	}{
		{"Book1", []string{"Book2", "Book3"}},
		{"Book2", []string{"Book1", "Book3"}},
		{"Book3", []string{"Book1", "Book2"}},
	}

	for _, testCase := range testCases {
		recommendations := router.BookRecommendationEngine(testCase.bookName, idlist, books)
		if !stringSlicesEqual(recommendations, testCase.expectedRecommendations) {
			t.Errorf("bookRecommendationEngine(%s) = %v; expected %v",
				testCase.bookName, recommendations, testCase.expectedRecommendations)
		}
	}
}

func TestGetRecommendations(t *testing.T) {
    testCases := []struct {
        bookName                string
        expectedRecommendations []string
    }{
        {"Harry Potter and the Order of the Phoenix (Harry Potter  #5)", []string{
			"Harry Potter and the Order of the Phoenix (Harry Potter  #5)",
            "Harry Potter and the Half-Blood Prince (Harry Potter  #6)",
            "The Fellowship of the Ring (The Lord of the Rings  #1)",
            "Harry Potter and the Chamber of Secrets (Harry Potter  #2)",
            "Harry Potter and the Prisoner of Azkaban (Harry Potter  #3)",
            "The Hobbit  or There and Back Again",
            "The Lightning Thief (Percy Jackson and the Olympians  #1)",
            "The Book Thief",
            "The Giver (The Giver  #1)",
            "Little Women",
        }},
        {"The Hobbit", []string{
            "The Hobbit",
            "The Voyage of the Jerle Shannara Trilogy (Voyage of the Jerle Shannara  #1-3)",
            "Tsubasa: RESERVoir CHRoNiCLE  Vol. 8",
            "The Complete Plays",
            "Ariel: The Restored Edition",
            "The Forbidden (Vampire Huntress  #5)",
            "Agatha Christie: An Autobiography",
            "What Work Is",
            "How Europe Underdeveloped Africa",
            "Exclusion & Embrace: A Theological Exploration of Identity  Otherness  and Reconciliation",
        }},
    }

    for _, testCase := range testCases {
        recommendations := router.GetRecommendationsHelper(testCase.bookName)
        if !stringSlicesEqual(recommendations, testCase.expectedRecommendations) {
            t.Errorf("getRecommendations(%s) = %v; expected %v",
                testCase.bookName, recommendations, testCase.expectedRecommendations)
        }
    }
}

func stringSlicesEqual(slice1, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for i := range slice1 {
		if slice1[i] != slice2[i] {
			return false
		}
	}
	return true
}

func createTempFile(t *testing.T, filename string, content []byte) string {
	tempFile, err := os.CreateTemp("", filename)
	if err != nil {
		t.Fatal("Error creating temporary file:", err)
	}

	_, err = tempFile.Write(content)
	if err != nil {
		t.Fatal("Error writing to temporary file:", err)
	}

	return tempFile.Name()
}


/*
[Harry Potter and the Half-Blood Prince (Harry Potter  #6) The Fellowship of the Ring (The Lord of the Rings  #1) Harry Potter and the Chamber of Secrets (Harry Potter  #2) Harry Potter and the Prisoner of Azkaban (Harry Potter  #3) The Hobbit  or There and Back Again The Lightning Thief (Percy Jackson and the Olympians  #1) The Book Thief The Giver (The Giver  #1) Little Women]

[Harry Potter and the Half-Blood Prince (Harry Potter  #6) The Fellowship of the Ring (The Lord of the Rings  #1) Harry Potter and the Chamber of Secrets (Harry Potter  #2) Harry Potter and the Prisoner of Azkaban (Harry Potter  #3) The Hobbit  or There and Back Again The Lightning Thief (Percy Jackson and the Olympians  #1) The Book Thief The Giver (The Giver  #1) Little Women]
*/