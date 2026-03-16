package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DefaultBaseURL = "https://phntm.sh"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func New(baseURL string) *Client {
	if baseURL == "" {
		baseURL = DefaultBaseURL
	}
	return &Client{
		BaseURL:    baseURL,
		HTTPClient: &http.Client{},
	}
}

// InitUploadResponse is the response from POST /api/upload.
type InitUploadResponse struct {
	ID        string `json:"id"`
	UploadURL string `json:"upload_url"`
	Token     string `json:"token"`
}

// FileMetadata is the response from GET /api/file/[id].
type FileMetadata struct {
	ID        string `json:"id"`
	FileName  string `json:"file_name"`
	FileSize  int64  `json:"file_size"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}

// InitUpload requests a presigned upload URL from the server.
func (c *Client) InitUpload(fileName string, fileSize int64, expiryHours int) (*InitUploadResponse, error) {
	body, _ := json.Marshal(map[string]interface{}{
		"file_name":    fileName,
		"file_size":    fileSize,
		"expiry_hours": expiryHours,
	})

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/upload", "application/json", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upload init failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	var result InitUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// UploadToStorage uploads encrypted data directly to Supabase Storage via presigned URL.
func (c *Client) UploadToStorage(uploadURL string, token string, data []byte) error {
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("storage upload failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("storage upload failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// ConfirmUpload creates the DB record after a successful storage upload.
func (c *Client) ConfirmUpload(id, fileName string, fileSize int64, expiryHours int) error {
	body, _ := json.Marshal(map[string]interface{}{
		"id":           id,
		"file_name":    fileName,
		"file_size":    fileSize,
		"expiry_hours": expiryHours,
	})

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/upload/confirm", "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("confirm request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("confirm failed (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

// GetFileMetadata fetches file info from the server.
func (c *Client) GetFileMetadata(id string) (*FileMetadata, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/file/" + id)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 410 {
		return nil, fmt.Errorf("TRANSMISSION_EXPIRED: DATA PURGED")
	}
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("TRANSMISSION_NOT_FOUND")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get file info (HTTP %d)", resp.StatusCode)
	}

	var meta FileMetadata
	if err := json.NewDecoder(resp.Body).Decode(&meta); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &meta, nil
}

// DownloadFile downloads the encrypted blob from the server.
func (c *Client) DownloadFile(id string) ([]byte, error) {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/download/" + id)
	if err != nil {
		return nil, fmt.Errorf("download request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 410 {
		return nil, fmt.Errorf("TRANSMISSION_EXPIRED: DATA PURGED")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("download failed (HTTP %d)", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	return data, nil
}
