package googlebooks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"test-project/internal/dto"

	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	err := godotenv.Load(".env.local")
	if err != nil {
		fmt.Println("No .env found")
	}
}

type Service struct {
	apiKey string
}

func NewService() *Service {
	apiKey := os.Getenv("GOOGLE_BOOKS_API_KEY")
	return &Service{
		apiKey: apiKey,
	}
}

func (s *Service) SearchBooks(ctx context.Context, query string, maxResults int) (*dto.BookSearchResult, error) {
	logger := log.Ctx(ctx).With().
		Str("query", query).
		Int("maxresults", maxResults).
		Logger()

	baseURL := "https://www.googleapis.com/books/v1/volumes"

	if s.apiKey == "" {
		logger.Warn().Msg("Google Books API Key not configured.")
		return nil, fmt.Errorf("Google Books API Key not configured")
	}

	params := url.Values{}
	params.Add("q", query)
	params.Add("key", s.apiKey)
	params.Add("maxResults", fmt.Sprintf("%d", maxResults))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create HTTP request")
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("HTTP request failed")
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("statusCode", resp.StatusCode).Msg("Google Books API returned non-OK status")
		return nil, fmt.Errorf("Google Books API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read response body")
		return nil, fmt.Errorf("Failed to read respone body: %w", err)
	}

	var apiResponse dto.BookSearchResult
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Error().Err(err).Msg("Failed to parse JSON response")
		return nil, fmt.Errorf("Failed to parse JSON response: %w", err)
	}

	return &apiResponse, nil
}

func (s *Service) GetBookByID(ctx context.Context, bookID string) (*dto.BookSummary, error) {
	logger := log.Ctx(ctx).With().Str("book_id", bookID).Logger()

	baseURL := "https://www.googleapis.com/books/v1/volumes"

	if s.apiKey == "" {
		logger.Warn().Msg("Google Books API Key not configured")
		return nil, fmt.Errorf("Google Books API Key not configured")
	}

	requestURL := baseURL + "/" + bookID + "?key=" + s.apiKey

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create HTTP request")
		return nil, fmt.Errorf("Failed to create HTTP request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("HTTP request failed")
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("statusCode", resp.StatusCode).Msg("Google Books API returned non-OK status")
		return nil, fmt.Errorf("Google Books API returned non-OK status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read response body")
		return nil, fmt.Errorf("Failed to read response body: %w", err)
	}

	var apiResponse dto.BookSummary
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Error().Err(err).Msg("Failed to parse JSON responese")
		return nil, fmt.Errorf("Failed to parse JSON response: %w", err)
	}

	return &apiResponse, nil
}
