package mindkit

import (
	"context"
	"fmt"
)

// AI provides AI-powered features for Gitk repositories
type AI struct {
	client *Client
}

// NewAI creates a new AI service instance
func NewAI(client *Client) *AI {
	return &AI{
		client: client,
	}
}

// GenerateCommitMessage generates an AI-powered commit message
func (a *AI) GenerateCommitMessage(ctx context.Context, diff string) (string, error) {
	endpoint := fmt.Sprintf("%s/commit/generate", a.client.baseURL)
	
	req := struct {
		Diff string `json:"diff"`
	}{
		Diff: diff,
	}
	
	var result struct {
		Message string `json:"message"`
	}
	
	if err := a.client.post(ctx, endpoint, req, &result); err != nil {
		return "", err
	}
	
	return result.Message, nil
}

// AnalyzeCode performs AI analysis on code
func (a *AI) AnalyzeCode(ctx context.Context, code string) (*Analysis, error) {
	analysis := &Analysis{
		Suggestions: make([]Suggestion, 0),
		Score: 0.0,
	}
	
	// Call MindKit API to analyze code
	result, err := a.client.AnalyzeRepository(ctx, code)
	if err != nil {
		return nil, err
	}

	analysis.Suggestions = result.Suggestions
	analysis.Score = result.Score

	return analysis, nil
}

// ReviewCode performs an AI code review
func (a *AI) ReviewCode(ctx context.Context, diff string) (*Review, error) {
	review := &Review{
		Comments: make([]Comment, 0),
		Summary: "",
	}
	
	// Call MindKit API to review code
	result, err := a.client.ReviewCode(ctx, diff)
	if err != nil {
		return nil, err
	}

	review.Comments = result.Comments
	review.Summary = result.Summary

	return review, nil
}

// GenerateDocumentation generates AI-powered documentation
func (a *AI) GenerateDocumentation(ctx context.Context, path string) (*Documentation, error) {
	// Call MindKit API to generate documentation
	doc, err := a.client.GenerateDocumentation(ctx, path)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

// Types for AI analysis results
type Analysis struct {
	Suggestions []Suggestion `json:"suggestions"`
	Score       float64     `json:"score"`
}

type Review struct {
	Comments []Comment `json:"comments"`
	Summary  string   `json:"summary"`
}
