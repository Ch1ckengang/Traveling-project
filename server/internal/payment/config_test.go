package payment

import (
	"os"
	"strings"
	"testing"
)

func TestVNPayConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  VNPayConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid sandbox configuration",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: false,
		},
		{
			name: "valid production configuration",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "PROD_MERCHANT",
				SecretKey:      "PROD_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: false,
		},
		{
			name: "missing environment",
			config: VNPayConfig{
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'environment' is required",
		},
		{
			name: "missing merchant_id",
			config: VNPayConfig{
				Environment:    "sandbox",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'merchant_id' is required",
		},
		{
			name: "missing secret_key",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'secret_key' is required",
		},
		{
			name: "missing return_url",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'return_url' is required",
		},
		{
			name: "missing ipn_url",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'ipn_url' is required",
		},
		{
			name: "invalid environment value",
			config: VNPayConfig{
				Environment:    "invalid",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'environment' must be 'sandbox' or 'production'",
		},
		{
			name: "invalid return_url format",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "not-a-valid-url",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'return_url' must include a protocol scheme",
		},
		{
			name: "invalid ipn_url format",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "ftp://example.com/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'ipn_url' must use http or https protocol",
		},
		{
			name: "zero payment timeout",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 0,
			},
			wantErr: true,
			errMsg:  "'payment_timeout' must be greater than 0",
		},
		{
			name: "negative payment timeout",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: -5,
			},
			wantErr: true,
			errMsg:  "'payment_timeout' must be greater than 0",
		},
		{
			name: "http url allowed for sandbox",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "http://localhost:3000/payment/return",
				IPNURL:         "http://localhost:8080/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: false,
		},
		{
			name: "case insensitive environment - SANDBOX",
			config: VNPayConfig{
				Environment:    "SANDBOX",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: false,
		},
		{
			name: "case insensitive environment - Production",
			config: VNPayConfig{
				Environment:    "Production",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: false,
		},
		{
			name: "merchant_id too short",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "ABC",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'merchant_id' appears invalid (too short",
		},
		{
			name: "secret_key too short",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "SHORT",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'secret_key' appears invalid (too short",
		},
		{
			name: "production requires https return_url",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "PROD_MERCHANT",
				SecretKey:      "PROD_SECRET_KEY_1234567890",
				ReturnURL:      "http://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'return_url' must use HTTPS in production",
		},
		{
			name: "production requires https ipn_url",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "PROD_MERCHANT",
				SecretKey:      "PROD_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "http://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: true,
			errMsg:  "'ipn_url' must use HTTPS in production",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				if err == nil {
					t.Errorf("Validate() expected error but got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestVNPayConfig_IsProduction(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		want        bool
	}{
		{
			name:        "production lowercase",
			environment: "production",
			want:        true,
		},
		{
			name:        "production uppercase",
			environment: "PRODUCTION",
			want:        true,
		},
		{
			name:        "production mixed case",
			environment: "Production",
			want:        true,
		},
		{
			name:        "sandbox",
			environment: "sandbox",
			want:        false,
		},
		{
			name:        "invalid",
			environment: "invalid",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := VNPayConfig{Environment: tt.environment}
			if got := config.IsProduction(); got != tt.want {
				t.Errorf("IsProduction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVNPayConfig_IsSandbox(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		want        bool
	}{
		{
			name:        "sandbox lowercase",
			environment: "sandbox",
			want:        true,
		},
		{
			name:        "sandbox uppercase",
			environment: "SANDBOX",
			want:        true,
		},
		{
			name:        "sandbox mixed case",
			environment: "Sandbox",
			want:        true,
		},
		{
			name:        "production",
			environment: "production",
			want:        false,
		},
		{
			name:        "invalid",
			environment: "invalid",
			want:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := VNPayConfig{Environment: tt.environment}
			if got := config.IsSandbox(); got != tt.want {
				t.Errorf("IsSandbox() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVNPayConfig_GetBaseURL(t *testing.T) {
	tests := []struct {
		name        string
		environment string
		want        string
	}{
		{
			name:        "production environment",
			environment: "production",
			want:        "https://vnpayment.vn/paymentv2/vpcpay.html",
		},
		{
			name:        "sandbox environment",
			environment: "sandbox",
			want:        "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html",
		},
		{
			name:        "production uppercase",
			environment: "PRODUCTION",
			want:        "https://vnpayment.vn/paymentv2/vpcpay.html",
		},
		{
			name:        "sandbox uppercase",
			environment: "SANDBOX",
			want:        "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := VNPayConfig{Environment: tt.environment}
			if got := config.GetBaseURL(); got != tt.want {
				t.Errorf("GetBaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVNPayConfig_MaskSecretKey(t *testing.T) {
	tests := []struct {
		name      string
		secretKey string
		want      string
	}{
		{
			name:      "long secret key",
			secretKey: "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
			want:      "ABCD****WXYZ",
		},
		{
			name:      "short secret key",
			secretKey: "SHORT",
			want:      "****",
		},
		{
			name:      "exactly 8 characters",
			secretKey: "12345678",
			want:      "****",
		},
		{
			name:      "9 characters",
			secretKey: "123456789",
			want:      "1234****6789",
		},
		{
			name:      "empty secret key",
			secretKey: "",
			want:      "****",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := VNPayConfig{SecretKey: tt.secretKey}
			if got := config.MaskSecretKey(); got != tt.want {
				t.Errorf("MaskSecretKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	tests := []struct {
		name      string
		urlStr    string
		fieldName string
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid https url",
			urlStr:    "https://example.com/path",
			fieldName: "test_url",
			wantErr:   false,
		},
		{
			name:      "valid http url",
			urlStr:    "http://localhost:8080/path",
			fieldName: "test_url",
			wantErr:   false,
		},
		{
			name:      "empty url",
			urlStr:    "",
			fieldName: "test_url",
			wantErr:   true,
			errMsg:    "'test_url' cannot be empty",
		},
		{
			name:      "url without scheme",
			urlStr:    "example.com/path",
			fieldName: "test_url",
			wantErr:   true,
			errMsg:    "'test_url' must include a protocol scheme",
		},
		{
			name:      "url with invalid scheme",
			urlStr:    "ftp://example.com/path",
			fieldName: "test_url",
			wantErr:   true,
			errMsg:    "'test_url' must use http or https protocol",
		},
		{
			name:      "url without host",
			urlStr:    "https:///path",
			fieldName: "test_url",
			wantErr:   true,
			errMsg:    "'test_url' must include a hostname or domain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURL(tt.urlStr, tt.fieldName)
			if tt.wantErr {
				if err == nil {
					t.Errorf("validateURL() expected error but got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("validateURL() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("validateURL() unexpected error = %v", err)
				}
			}
		})
	}
}

func TestLoadConfigFromEnv(t *testing.T) {
	// Save original environment variables
	originalEnv := map[string]string{
		"VNPAY_ENVIRONMENT":     os.Getenv("VNPAY_ENVIRONMENT"),
		"VNPAY_MERCHANT_ID":     os.Getenv("VNPAY_MERCHANT_ID"),
		"VNPAY_SECRET_KEY":      os.Getenv("VNPAY_SECRET_KEY"),
		"VNPAY_RETURN_URL":      os.Getenv("VNPAY_RETURN_URL"),
		"VNPAY_IPN_URL":         os.Getenv("VNPAY_IPN_URL"),
		"VNPAY_PAYMENT_TIMEOUT": os.Getenv("VNPAY_PAYMENT_TIMEOUT"),
	}

	// Restore environment variables after test
	defer func() {
		for key, value := range originalEnv {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
		errMsg  string
		check   func(*testing.T, *VNPayConfig)
	}{
		{
			name: "valid configuration with all env vars",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT":     "production",
				"VNPAY_MERCHANT_ID":     "PROD_MERCHANT_123",
				"VNPAY_SECRET_KEY":      "PROD_SECRET_KEY_XYZ",
				"VNPAY_RETURN_URL":      "https://example.com/payment/return",
				"VNPAY_IPN_URL":         "https://example.com/api/payment/webhook",
				"VNPAY_PAYMENT_TIMEOUT": "20",
			},
			wantErr: false,
			check: func(t *testing.T, config *VNPayConfig) {
				if config.Environment != "production" {
					t.Errorf("Environment = %v, want production", config.Environment)
				}
				if config.MerchantID != "PROD_MERCHANT_123" {
					t.Errorf("MerchantID = %v, want PROD_MERCHANT_123", config.MerchantID)
				}
				if config.SecretKey != "PROD_SECRET_KEY_XYZ" {
					t.Errorf("SecretKey = %v, want PROD_SECRET_KEY_XYZ", config.SecretKey)
				}
				if config.ReturnURL != "https://example.com/payment/return" {
					t.Errorf("ReturnURL = %v, want https://example.com/payment/return", config.ReturnURL)
				}
				if config.IPNURL != "https://example.com/api/payment/webhook" {
					t.Errorf("IPNURL = %v, want https://example.com/api/payment/webhook", config.IPNURL)
				}
				if config.PaymentTimeout != 20 {
					t.Errorf("PaymentTimeout = %v, want 20", config.PaymentTimeout)
				}
			},
		},
		{
			name: "valid configuration with default environment",
			envVars: map[string]string{
				"VNPAY_MERCHANT_ID": "TEST_MERCHANT",
				"VNPAY_SECRET_KEY":  "TEST_SECRET_KEY_1234567890",
				"VNPAY_RETURN_URL":  "https://example.com/return",
				"VNPAY_IPN_URL":     "https://example.com/webhook",
			},
			wantErr: false,
			check: func(t *testing.T, config *VNPayConfig) {
				if config.Environment != "sandbox" {
					t.Errorf("Environment = %v, want sandbox (default)", config.Environment)
				}
				if config.PaymentTimeout != 15 {
					t.Errorf("PaymentTimeout = %v, want 15 (default)", config.PaymentTimeout)
				}
			},
		},
		{
			name: "valid configuration with whitespace trimming",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT": "  sandbox  ",
				"VNPAY_MERCHANT_ID": "  TEST_MERCHANT  ",
				"VNPAY_SECRET_KEY":  "  TEST_SECRET_KEY_1234567890  ",
				"VNPAY_RETURN_URL":  "  https://example.com/return  ",
				"VNPAY_IPN_URL":     "  https://example.com/webhook  ",
			},
			wantErr: false,
			check: func(t *testing.T, config *VNPayConfig) {
				if config.Environment != "sandbox" {
					t.Errorf("Environment = %v, want sandbox (trimmed)", config.Environment)
				}
				if config.MerchantID != "TEST_MERCHANT" {
					t.Errorf("MerchantID = %v, want TEST_MERCHANT (trimmed)", config.MerchantID)
				}
			},
		},
		{
			name: "missing merchant_id",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT": "sandbox",
				"VNPAY_SECRET_KEY":  "TEST_SECRET_KEY_1234567890",
				"VNPAY_RETURN_URL":  "https://example.com/return",
				"VNPAY_IPN_URL":     "https://example.com/webhook",
			},
			wantErr: true,
			errMsg:  "'merchant_id' is required",
		},
		{
			name: "missing secret_key",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT": "sandbox",
				"VNPAY_MERCHANT_ID": "TEST_MERCHANT",
				"VNPAY_RETURN_URL":  "https://example.com/return",
				"VNPAY_IPN_URL":     "https://example.com/webhook",
			},
			wantErr: true,
			errMsg:  "'secret_key' is required",
		},
		{
			name: "missing return_url",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT": "sandbox",
				"VNPAY_MERCHANT_ID": "TEST_MERCHANT",
				"VNPAY_SECRET_KEY":  "TEST_SECRET_KEY_1234567890",
				"VNPAY_IPN_URL":     "https://example.com/webhook",
			},
			wantErr: true,
			errMsg:  "'return_url' is required",
		},
		{
			name: "missing ipn_url",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT": "sandbox",
				"VNPAY_MERCHANT_ID": "TEST_MERCHANT",
				"VNPAY_SECRET_KEY":  "TEST_SECRET_KEY_1234567890",
				"VNPAY_RETURN_URL":  "https://example.com/return",
			},
			wantErr: true,
			errMsg:  "'ipn_url' is required",
		},
		{
			name: "invalid environment value",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT": "invalid",
				"VNPAY_MERCHANT_ID": "TEST_MERCHANT",
				"VNPAY_SECRET_KEY":  "TEST_SECRET_KEY_1234567890",
				"VNPAY_RETURN_URL":  "https://example.com/return",
				"VNPAY_IPN_URL":     "https://example.com/webhook",
			},
			wantErr: true,
			errMsg:  "'environment' must be 'sandbox' or 'production'",
		},
		{
			name: "invalid return_url format",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT": "sandbox",
				"VNPAY_MERCHANT_ID": "TEST_MERCHANT",
				"VNPAY_SECRET_KEY":  "TEST_SECRET_KEY_1234567890",
				"VNPAY_RETURN_URL":  "not-a-valid-url",
				"VNPAY_IPN_URL":     "https://example.com/webhook",
			},
			wantErr: true,
			errMsg:  "'return_url' must include a protocol scheme",
		},
		{
			name: "invalid payment timeout - not a number",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT":     "sandbox",
				"VNPAY_MERCHANT_ID":     "TEST_MERCHANT",
				"VNPAY_SECRET_KEY":      "TEST_SECRET_KEY_1234567890",
				"VNPAY_RETURN_URL":      "https://example.com/return",
				"VNPAY_IPN_URL":         "https://example.com/webhook",
				"VNPAY_PAYMENT_TIMEOUT": "not-a-number",
			},
			wantErr: false, // Falls back to default
			check: func(t *testing.T, config *VNPayConfig) {
				if config.PaymentTimeout != 15 {
					t.Errorf("PaymentTimeout = %v, want 15 (default fallback)", config.PaymentTimeout)
				}
			},
		},
		{
			name: "invalid payment timeout - zero",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT":     "sandbox",
				"VNPAY_MERCHANT_ID":     "TEST_MERCHANT",
				"VNPAY_SECRET_KEY":      "TEST_SECRET_KEY_1234567890",
				"VNPAY_RETURN_URL":      "https://example.com/return",
				"VNPAY_IPN_URL":         "https://example.com/webhook",
				"VNPAY_PAYMENT_TIMEOUT": "0",
			},
			wantErr: false, // Falls back to default
			check: func(t *testing.T, config *VNPayConfig) {
				if config.PaymentTimeout != 15 {
					t.Errorf("PaymentTimeout = %v, want 15 (default fallback)", config.PaymentTimeout)
				}
			},
		},
		{
			name: "invalid payment timeout - negative",
			envVars: map[string]string{
				"VNPAY_ENVIRONMENT":     "sandbox",
				"VNPAY_MERCHANT_ID":     "TEST_MERCHANT",
				"VNPAY_SECRET_KEY":      "TEST_SECRET_KEY_1234567890",
				"VNPAY_RETURN_URL":      "https://example.com/return",
				"VNPAY_IPN_URL":         "https://example.com/webhook",
				"VNPAY_PAYMENT_TIMEOUT": "-5",
			},
			wantErr: false, // Falls back to default
			check: func(t *testing.T, config *VNPayConfig) {
				if config.PaymentTimeout != 15 {
					t.Errorf("PaymentTimeout = %v, want 15 (default fallback)", config.PaymentTimeout)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear all environment variables first
			os.Unsetenv("VNPAY_ENVIRONMENT")
			os.Unsetenv("VNPAY_MERCHANT_ID")
			os.Unsetenv("VNPAY_SECRET_KEY")
			os.Unsetenv("VNPAY_RETURN_URL")
			os.Unsetenv("VNPAY_IPN_URL")
			os.Unsetenv("VNPAY_PAYMENT_TIMEOUT")

			// Set test environment variables
			for key, value := range tt.envVars {
				os.Setenv(key, value)
			}

			// Load configuration
			config, err := LoadConfigFromEnv()

			// Check error expectation
			if tt.wantErr {
				if err == nil {
					t.Errorf("LoadConfigFromEnv() expected error but got nil")
					return
				}
				if tt.errMsg != "" && !strings.Contains(err.Error(), tt.errMsg) {
					t.Errorf("LoadConfigFromEnv() error = %v, want error containing %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("LoadConfigFromEnv() unexpected error = %v", err)
					return
				}
				if config == nil {
					t.Errorf("LoadConfigFromEnv() returned nil config")
					return
				}
				// Run additional checks if provided
				if tt.check != nil {
					tt.check(t, config)
				}
			}
		})
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback string
		want     string
	}{
		{
			name:     "environment variable exists",
			key:      "TEST_VAR_1",
			value:    "test_value",
			fallback: "default",
			want:     "test_value",
		},
		{
			name:     "environment variable empty - use fallback",
			key:      "TEST_VAR_2",
			value:    "",
			fallback: "default",
			want:     "default",
		},
		{
			name:     "environment variable with whitespace",
			key:      "TEST_VAR_3",
			value:    "  test_value  ",
			fallback: "default",
			want:     "test_value",
		},
		{
			name:     "environment variable only whitespace - use fallback",
			key:      "TEST_VAR_4",
			value:    "   ",
			fallback: "default",
			want:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
			} else {
				os.Unsetenv(tt.key)
			}
			defer os.Unsetenv(tt.key)

			got := getEnvOrDefault(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("getEnvOrDefault() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnvAsInt(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    string
		fallback int
		want     int
	}{
		{
			name:     "valid integer",
			key:      "TEST_INT_1",
			value:    "42",
			fallback: 10,
			want:     42,
		},
		{
			name:     "empty value - use fallback",
			key:      "TEST_INT_2",
			value:    "",
			fallback: 10,
			want:     10,
		},
		{
			name:     "invalid integer - use fallback",
			key:      "TEST_INT_3",
			value:    "not-a-number",
			fallback: 10,
			want:     10,
		},
		{
			name:     "zero value - use fallback",
			key:      "TEST_INT_4",
			value:    "0",
			fallback: 10,
			want:     10,
		},
		{
			name:     "negative value - use fallback",
			key:      "TEST_INT_5",
			value:    "-5",
			fallback: 10,
			want:     10,
		},
		{
			name:     "value with whitespace",
			key:      "TEST_INT_6",
			value:    "  25  ",
			fallback: 10,
			want:     25,
		},
		{
			name:     "large integer",
			key:      "TEST_INT_7",
			value:    "999999",
			fallback: 10,
			want:     999999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environment variable
			if tt.value != "" {
				os.Setenv(tt.key, tt.value)
			} else {
				os.Unsetenv(tt.key)
			}
			defer os.Unsetenv(tt.key)

			got := getEnvAsInt(tt.key, tt.fallback)
			if got != tt.want {
				t.Errorf("getEnvAsInt() = %v, want %v", got, tt.want)
			}
		})
	}
}
