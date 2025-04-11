package googlebooks

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"test-project/internal/config"
	"test-project/internal/dto"

	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
)

type Service struct {
	apiKey  string
	baseURL string
}

func NewService(config config.GoogleBooks) *Service {
	return &Service{
		apiKey:  config.APIKey,
		baseURL: config.BaseURL,
	}
}

func (s *Service) SearchBooks(ctx context.Context, query string, pageSize, page int) (*dto.BookSearchResult, int64, error) {
	logger := log.Ctx(ctx).With().
		Str("query", query).
		Int("pageSize", pageSize).
		Int("page", page).
		Logger()

	baseURL := s.baseURL

	if s.apiKey == "" {
		logger.Warn().Msg("Google Books API Key not configured.")
		return nil, 0, fmt.Errorf("Google Books API Key not configured")
	}

	startIndex := (page - 1) * pageSize

	params := url.Values{}
	params.Add("q", query)
	params.Add("key", s.apiKey)
	params.Add("maxResults", fmt.Sprintf("%d", pageSize))
	params.Add("startIndex", fmt.Sprintf("%d", startIndex))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to create HTTP request")
		return nil, 0, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).Msg("HTTP request failed")
		return nil, 0, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Error().Int("statusCode", resp.StatusCode).Msg("Google Books API returned non-OK status")
		return nil, 0, fmt.Errorf("Google Books API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Error().Err(err).Msg("Failed to read response body")
		return nil, 0, fmt.Errorf("Failed to read response body: %w", err)
	}

	var apiResponse dto.BookSearchResult
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Error().Err(err).Msg("Failed to parse JSON response")
		return nil, 0, fmt.Errorf("Failed to parse JSON response: %w", err)
	}

	totalItems := int64(apiResponse.TotalItems)

	logger.Info().
		Int64("totalItems", totalItems).
		Int("resultsReturned", len(apiResponse.Books)).
		Msg("successful search")

	return &apiResponse, totalItems, nil
}

func (s *Service) GetBookByID(ctx context.Context, bookID string) (*dto.BookSummary, error) {
	logger := log.Ctx(ctx).With().Str("book_id", bookID).Logger()

	baseURL := s.baseURL

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

	logger.Info().Msg("book succesfully fetched")

	return &apiResponse, nil
}
