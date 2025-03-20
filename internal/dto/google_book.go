package dto

type BookSearchResult struct {
	Books []BookSummary `json:"items"`
}

type BookSummary struct {
	GoogleID    string      `json:"id"`
	BookDetails BookDetails `json:"volumeInfo"`
}

type BookDetails struct {
	Title       string   `json:"title"`
	Authors     []string `json:"authors"`
	Publisher   string   `json:"publisher"`
	Description string   `json:"description"`
	Genre       []string `json:"categories"`
	Pages       int      `json:"pageCount"`
}
