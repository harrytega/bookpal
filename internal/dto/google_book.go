package dto

type BookSearchResult struct {
	Books      []BookSummary `json:"items"`
	TotalItems int           `json:"totalItems"`
}

type BookSummary struct {
	GoogleID    string      `json:"id"`
	BookDetails BookDetails `json:"volumeInfo"`
}

type BookDetails struct {
	Title       string     `json:"title"`
	Authors     []string   `json:"authors"`
	Publisher   string     `json:"publisher"`
	Description string     `json:"description"`
	Genre       []string   `json:"categories"`
	Pages       int        `json:"pageCount"`
	ImageLinks  ImageLinks `json:"imageLinks"`
}

type ImageLinks struct {
	SmallThumbnail string `json:"smallThumbnail"`
	Thumbnail      string `json:"thumbnail"`
}
