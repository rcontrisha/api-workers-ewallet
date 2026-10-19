package pkg

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
	"uuid"
)

type DokuClient struct {
	ClientID   string
	SecretKey  string
	BaseURL    string
	HTTPClient *http.Client
}

func NewDokuClient() *DokuClient {
	// baseURL := os.Getenv("DOKU_BASE_URL")
	// if baseURL == "" {
	// 	baseURL = "https://api-sandbox.doku.com"
	// }
	baseURL := "https://api-sandbox.doku.com"

	return &DokuClient{
		ClientID:  os.Getenv("DOKU_CLIENT_ID"),
		SecretKey: os.Getenv("DOKU_SECRET_KEY"),
		BaseURL:   baseURL,
		HTTPClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (d *DokuClient) GenerateSignature(requestID, timestamp, requestTarget string, bodyBytes []byte) string {
	bodyHash := sha256.Sum256(bodyBytes)
	digest := base64.StdEncoding.EncodeToString(bodyHash[:])

	componentStr := fmt.Sprintf("Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
		d.ClientID, requestID, timestamp, requestTarget, digest)

	h := hmac.New(sha256.New, []byte(d.SecretKey))
	h.Write([]byte(componentStr))
	rawSignature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return "HMACSHA256=" + rawSignature
}

type DokuCheckoutPayload struct {
	Order struct {
		Amount        int    `json:"amount"`
		InvoiceNumber string `json:"invoice_number"`
		Currency      string `json:"currency"`
		CallbackURL   string `json:"callback_url,omitempty"`
		AutoRedirect  bool   `json:"auto_redirect"`
		LineItems     []struct {
			ID       string `json:"id"`
			Name     string `json:"name"`
			Price    int64  `json:"price"`
			Quantity int    `json:"quantity"`
		} `json:"line_items,omitempty"`
	} `json:"order"`
	Payment struct {
		PaymentDueDate     int      `json:"payment_due_date,omitempty"`
		PaymentMethodTypes []string `json:"payment_method_types,omitempty"`
	} `json:"payment"`
	Customer struct {
		ID    string `json:"id,omitempty"`
		Name  string `json:"name,omitempty"`
		Email string `json:"email,omitempty"`
		Phone string `json:"phone,omitempty"`
	} `json:"customer"`
}

type DokuCheckoutResponse struct {
	Message  []string `json:"message"`
	Response struct {
		Order struct {
			InvoiceNumber string `json:"invoice_number"`
			Amount        string `json:"amount"`
		} `json:"order"`
		Payment struct {
			URL            string `json:"url"`
			PaymentDueDate int    `json:"payment_due_date"`
		} `json:"payment"`
	} `json:"response"`
}

func (d *DokuClient) CreateCheckout(ctx context.Context, payload DokuCheckoutPayload) (*DokuCheckoutResponse, error) {
	requestTarget := "/checkout/v1/payment"
	fullURL := d.BaseURL + requestTarget

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal doku payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fullURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	requestID := uuid.New().String()
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	signature := d.GenerateSignature(requestID, timestamp, requestTarget, bodyBytes)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Client-Id", d.ClientID)
	req.Header.Set("Request-Id", requestID)
	req.Header.Set("Request-Timestamp", timestamp)
	req.Header.Set("Signature", signature)

	resp, err := d.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to doku: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read doku response: %w", err)
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("doku API error (status %d): %s", resp.StatusCode, string(respBody))
	}

	var dokuResp DokuCheckoutResponse
	if err := json.Unmarshal(respBody, &dokuResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal doku response: %w", err)
	}

	return &dokuResp, nil
}
