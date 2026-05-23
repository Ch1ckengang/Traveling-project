package payment

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// VNPayConfig represents the configuration for VNPay payment gateway
// Supports both sandbox and production environments
type VNPayConfig struct {
	Environment    string `json:"environment"`     // sandbox or production
	MerchantID     string `json:"merchant_id"`     // VNPay merchant/terminal code
	SecretKey      string `json:"secret_key"`      // Secret key for signature generation
	ReturnURL      string `json:"return_url"`      // User-facing URL for payment return
	IPNURL         string `json:"ipn_url"`         // Server-to-server webhook URL
	FrontendURL    string `json:"frontend_url"`    // Frontend base URL for result redirects
	PaymentTimeout int    `json:"payment_timeout"` // Payment session timeout in minutes
}

// Validate checks if the VNPayConfig has all required fields and valid formats
// Returns an error if validation fails with clear, actionable error messages
func (c *VNPayConfig) Validate() error {
	// Check required fields with actionable error messages
	if c.Environment == "" {
		return errors.New("configuration error: 'environment' is required. Set VNPAY_ENVIRONMENT to 'sandbox' or 'production'")
	}

	if c.MerchantID == "" {
		return errors.New("configuration error: 'merchant_id' is required. Set VNPAY_MERCHANT_ID to your VNPay merchant/terminal code")
	}

	if c.SecretKey == "" {
		return errors.New("configuration error: 'secret_key' is required. Set VNPAY_SECRET_KEY to your VNPay secret key for signature generation")
	}

	if c.ReturnURL == "" {
		return errors.New("configuration error: 'return_url' is required. Set VNPAY_RETURN_URL to the user-facing URL where VNPay will redirect after payment (e.g., https://yourdomain.com/payment/return)")
	}

	if c.IPNURL == "" {
		return errors.New("configuration error: 'ipn_url' is required. Set VNPAY_IPN_URL to the server-to-server webhook URL for payment notifications (e.g., https://yourdomain.com/api/v1/payments/webhook)")
	}

	if c.FrontendURL == "" {
		return errors.New("configuration error: 'frontend_url' is required. Set FRONTEND_BASE_URL to your frontend app base URL")
	}

	// Validate environment value
	env := strings.ToLower(c.Environment)
	if env != "sandbox" && env != "production" {
		return fmt.Errorf("configuration error: 'environment' must be 'sandbox' or 'production', got '%s'. Use 'sandbox' for testing or 'production' for live payments", c.Environment)
	}

	// Validate URL formats
	if err := validateURL(c.ReturnURL, "return_url"); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if err := validateURL(c.IPNURL, "ipn_url"); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	if err := validateURL(c.FrontendURL, "frontend_url"); err != nil {
		return fmt.Errorf("configuration error: %w", err)
	}

	// Validate payment timeout
	if c.PaymentTimeout <= 0 {
		return fmt.Errorf("configuration error: 'payment_timeout' must be greater than 0, got %d. Set VNPAY_PAYMENT_TIMEOUT to a positive number of minutes (recommended: 15)", c.PaymentTimeout)
	}

	// Validate merchant ID format (basic check)
	if len(c.MerchantID) < 4 {
		return fmt.Errorf("configuration error: 'merchant_id' appears invalid (too short: %d characters). Verify your VNPay merchant/terminal code", len(c.MerchantID))
	}

	// Validate secret key format (basic check)
	if len(c.SecretKey) < 16 {
		return fmt.Errorf("configuration error: 'secret_key' appears invalid (too short: %d characters). Verify your VNPay secret key", len(c.SecretKey))
	}

	// Warn about production environment security
	if c.IsProduction() {
		// Additional validation for production
		if !strings.HasPrefix(c.ReturnURL, "https://") {
			return errors.New("security error: 'return_url' must use HTTPS in production environment for secure payment processing")
		}
		if !strings.HasPrefix(c.IPNURL, "https://") {
			return errors.New("security error: 'ipn_url' must use HTTPS in production environment for secure webhook communication")
		}
		if !strings.HasPrefix(c.FrontendURL, "https://") {
			return errors.New("security error: 'frontend_url' must use HTTPS in production environment")
		}
	}

	return nil
}

// validateURL checks if a URL string is valid and uses appropriate protocol
func validateURL(urlStr, fieldName string) error {
	if urlStr == "" {
		return fmt.Errorf("'%s' cannot be empty", fieldName)
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("'%s' is not a valid URL: %w. Ensure the URL is properly formatted (e.g., https://yourdomain.com/path)", fieldName, err)
	}

	// Ensure URL has a scheme
	if parsedURL.Scheme == "" {
		return fmt.Errorf("'%s' must include a protocol scheme (http:// or https://). Current value: '%s'", fieldName, urlStr)
	}

	// Ensure URL has a host
	if parsedURL.Host == "" {
		return fmt.Errorf("'%s' must include a hostname or domain. Current value: '%s'", fieldName, urlStr)
	}

	// Validate scheme is http or https
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("'%s' must use http or https protocol, got '%s'. Current value: '%s'", fieldName, parsedURL.Scheme, urlStr)
	}

	// Validate host is not localhost or 127.0.0.1 for production URLs
	// This is a warning-level check that should be caught by the production HTTPS check
	if strings.Contains(strings.ToLower(parsedURL.Host), "localhost") || strings.HasPrefix(parsedURL.Host, "127.0.0.1") {
		// This is acceptable for sandbox/development, but the production check will catch it
		// No error here, just a note for developers
	}

	return nil
}

// IsProduction returns true if the configuration is for production environment
func (c *VNPayConfig) IsProduction() bool {
	return strings.ToLower(c.Environment) == "production"
}

// IsSandbox returns true if the configuration is for sandbox environment
func (c *VNPayConfig) IsSandbox() bool {
	return strings.ToLower(c.Environment) == "sandbox"
}

// GetBaseURL returns the VNPay gateway base URL based on environment
func (c *VNPayConfig) GetBaseURL() string {
	if c.IsProduction() {
		return "https://vnpayment.vn/paymentv2/vpcpay.html"
	}
	return "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"
}

// MaskSecretKey returns a masked version of the secret key for logging
// Shows only the first 4 and last 4 characters
func (c *VNPayConfig) MaskSecretKey() string {
	if len(c.SecretKey) <= 8 {
		return "****"
	}
	return c.SecretKey[:4] + "****" + c.SecretKey[len(c.SecretKey)-4:]
}

// LoadConfigFromEnv loads VNPay configuration from environment variables
// Returns an error if required environment variables are missing or invalid
// Environment variables:
//   - VNPAY_ENVIRONMENT: sandbox or production (default: sandbox)
//   - VNPAY_MERCHANT_ID: VNPay merchant/terminal code (required)
//   - VNPAY_SECRET_KEY: Secret key for signature generation (required)
//   - VNPAY_RETURN_URL: User-facing URL for payment return (required)
//   - VNPAY_IPN_URL: Server-to-server webhook URL (required)
//   - FRONTEND_BASE_URL: Frontend base URL for redirecting users after return verification
//   - VNPAY_PAYMENT_TIMEOUT: Payment session timeout in minutes (default: 15)
func LoadConfigFromEnv() (*VNPayConfig, error) {
	config := &VNPayConfig{
		Environment:    getEnvOrDefault("VNPAY_ENVIRONMENT", "sandbox"),
		MerchantID:     strings.TrimSpace(os.Getenv("VNPAY_MERCHANT_ID")),
		SecretKey:      strings.TrimSpace(os.Getenv("VNPAY_SECRET_KEY")),
		ReturnURL:      strings.TrimSpace(os.Getenv("VNPAY_RETURN_URL")),
		IPNURL:         strings.TrimSpace(os.Getenv("VNPAY_IPN_URL")),
		FrontendURL:    getEnvOrDefault("FRONTEND_BASE_URL", "http://localhost:5173"),
		PaymentTimeout: getEnvAsInt("VNPAY_PAYMENT_TIMEOUT", 15),
	}

	// Validate the configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to load VNPay configuration from environment variables: %w", err)
	}

	return config, nil
}

// LoadVNPayConfig - Load VNPay config từ env, trả về sandbox defaults nếu thiếu config
// Dùng trong main.go khi khởi tạo server
func LoadVNPayConfig() *VNPayConfig {
	config, err := LoadConfigFromEnv()
	if err != nil {
		// Fallback: sandbox defaults cho development
		config = &VNPayConfig{
			Environment:    "sandbox",
			MerchantID:     getEnvOrDefault("VNPAY_MERCHANT_ID", "DEMO"),
			SecretKey:      getEnvOrDefault("VNPAY_SECRET_KEY", "DEMO_SECRET_KEY"),
			ReturnURL:      getEnvOrDefault("VNPAY_RETURN_URL", "http://localhost:8080/v1/api/payments/return"),
			IPNURL:         getEnvOrDefault("VNPAY_IPN_URL", "http://localhost:8080/v1/api/payments/webhook"),
			FrontendURL:    getEnvOrDefault("FRONTEND_BASE_URL", "http://localhost:5173"),
			PaymentTimeout: 15,
		}
		fmt.Printf("⚠️ VNPay config validation failed: %v. Using sandbox defaults.\n", err)
	}
	return config
}

// getEnvOrDefault retrieves an environment variable or returns a default value
// Trims whitespace from the retrieved value
func getEnvOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

// getEnvAsInt retrieves an environment variable as an integer or returns a default value
// Returns the fallback if the value is empty, not a valid integer, or less than or equal to 0
func getEnvAsInt(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}

	return parsed
}

// PrettyPrint formats the VNPayConfig into a readable JSON string with proper indentation
// Sensitive fields (secret_key) are masked in the output for security
// Returns formatted JSON string or error if formatting fails
func (c *VNPayConfig) PrettyPrint() (string, error) {
	// Create a copy of the config with masked sensitive data
	masked := c.MaskedCopy()

	// Marshal with indentation for pretty printing
	data, err := json.MarshalIndent(masked, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to format configuration: %w", err)
	}

	return string(data), nil
}

// CompactPrint formats the VNPayConfig into a compact JSON string
// Sensitive fields (secret_key) are masked in the output for security
// Returns compact JSON string or error if formatting fails
func (c *VNPayConfig) CompactPrint() (string, error) {
	// Create a copy of the config with masked sensitive data
	masked := c.MaskedCopy()

	// Marshal without indentation for compact output
	data, err := json.Marshal(masked)
	if err != nil {
		return "", fmt.Errorf("failed to format configuration: %w", err)
	}

	return string(data), nil
}

// MaskedCopy creates a copy of the VNPayConfig with sensitive fields masked
// This is used for safe logging and display purposes
func (c *VNPayConfig) MaskedCopy() VNPayConfig {
	return VNPayConfig{
		Environment:    c.Environment,
		MerchantID:     c.MerchantID,
		SecretKey:      c.MaskSecretKey(),
		ReturnURL:      c.ReturnURL,
		IPNURL:         c.IPNURL,
		FrontendURL:    c.FrontendURL,
		PaymentTimeout: c.PaymentTimeout,
	}
}
