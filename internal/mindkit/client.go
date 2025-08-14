package mindkit

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Client represents a MindKit API client
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Config holds the configuration for the MindKit client
type Config struct {
	BaseURL string
	APIKey  string
}

// NewClient creates a new MindKit API client
func NewClient(config Config) *Client {
	return &Client{
		baseURL:    config.BaseURL,
		apiKey:     config.APIKey,
		httpClient: &http.Client{},
	}
}

// AnalyzeRepository performs AI analysis on the entire repository
func (c *Client) AnalyzeRepository(ctx context.Context, repoPath string) (*AnalysisResult, error) {
	endpoint := fmt.Sprintf("%s/analyze", c.baseURL)
	
	req := struct {
		Path string `json:"path"`
	}{
		Path: repoPath,
	}
	
	var result AnalysisResult
	if err := c.post(ctx, endpoint, req, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// ReviewCode performs AI code review
func (c *Client) ReviewCode(ctx context.Context, diff string) (*CodeReview, error) {
	endpoint := fmt.Sprintf("%s/review", c.baseURL)
	
	req := struct {
		Diff string `json:"diff"`
	}{
		Diff: diff,
	}
	
	var result CodeReview
	if err := c.post(ctx, endpoint, req, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// GenerateDocumentation generates documentation using AI
func (c *Client) GenerateDocumentation(ctx context.Context, repoPath string) (*Documentation, error) {
	endpoint := fmt.Sprintf("%s/docs/generate", c.baseURL)
	
	req := struct {
		Path string `json:"path"`
	}{
		Path: repoPath,
	}
	
	var result Documentation
	if err := c.post(ctx, endpoint, req, &result); err != nil {
		return nil, err
	}
	
	return &result, nil
}

// Helper method to make HTTP POST requests
func (c *Client) post(ctx context.Context, endpoint string, body interface{}, result interface{}) error {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}

// Types for API responses
type AnalysisResult struct {
	Suggestions []Suggestion `json:"suggestions"`
	Score       float64     `json:"score"`
}

type CodeReview struct {
	Comments []Comment `json:"comments"`
	Summary  string   `json:"summary"`
}

type Documentation struct {
	Content     string   `json:"content"`
	References  []string `json:"references"`
	Generated   bool     `json:"generated"`
}

type Suggestion struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Message  string `json:"message"`
	Severity string `json:"severity"`
}

type Comment struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Message string `json:"message"`
	Type    string `json:"type"`
}
