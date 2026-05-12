package payment

import (
	"testing"
)

// TestEnvironmentSwitching verifies that environment switching works correctly
// This test ensures that IsProduction(), IsSandbox(), and GetBaseURL() methods
// correctly respond to different environment configurations
func TestEnvironmentSwitching(t *testing.T) {
	tests := []struct {
		name                string
		environment         string
		expectIsProduction  bool
		expectIsSandbox     bool
		expectGatewayURL    string
	}{
		{
			name:               "sandbox environment lowercase",
			environment:        "sandbox",
			expectIsProduction: false,
			expectIsSandbox:    true,
			expectGatewayURL:   "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html",
		},
		{
			name:               "sandbox environment uppercase",
			environment:        "SANDBOX",
			expectIsProduction: false,
			expectIsSandbox:    true,
			expectGatewayURL:   "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html",
		},
		{
			name:               "sandbox environment mixed case",
			environment:        "SandBox",
			expectIsProduction: false,
			expectIsSandbox:    true,
			expectGatewayURL:   "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html",
		},
		{
			name:               "production environment lowercase",
			environment:        "production",
			expectIsProduction: true,
			expectIsSandbox:    false,
			expectGatewayURL:   "https://vnpayment.vn/paymentv2/vpcpay.html",
		},
		{
			name:               "production environment uppercase",
			environment:        "PRODUCTION",
			expectIsProduction: true,
			expectIsSandbox:    false,
			expectGatewayURL:   "https://vnpayment.vn/paymentv2/vpcpay.html",
		},
		{
			name:               "production environment mixed case",
			environment:        "Production",
			expectIsProduction: true,
			expectIsSandbox:    false,
			expectGatewayURL:   "https://vnpayment.vn/paymentv2/vpcpay.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := VNPayConfig{
				Environment: tt.environment,
			}

			// Test IsProduction()
			if got := config.IsProduction(); got != tt.expectIsProduction {
				t.Errorf("IsProduction() = %v, want %v", got, tt.expectIsProduction)
			}

			// Test IsSandbox()
			if got := config.IsSandbox(); got != tt.expectIsSandbox {
				t.Errorf("IsSandbox() = %v, want %v", got, tt.expectIsSandbox)
			}

			// Test GetBaseURL()
			if got := config.GetBaseURL(); got != tt.expectGatewayURL {
				t.Errorf("GetBaseURL() = %v, want %v", got, tt.expectGatewayURL)
			}

			// Verify mutual exclusivity: cannot be both production and sandbox
			if config.IsProduction() && config.IsSandbox() {
				t.Error("Config cannot be both production and sandbox")
			}
		})
	}
}

// TestEnvironmentSwitchingValidation verifies that environment-specific
// validation rules are correctly enforced
func TestEnvironmentSwitchingValidation(t *testing.T) {
	tests := []struct {
		name        string
		config      VNPayConfig
		wantErr     bool
		errContains string
	}{
		{
			name: "sandbox allows http urls",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "http://localhost:5173/payment/return",
				IPNURL:         "http://localhost:8080/api/webhook",
				PaymentTimeout: 15,
			},
			wantErr: false,
		},
		{
			name: "production requires https return url",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "PROD_MERCHANT",
				SecretKey:      "PROD_SECRET_KEY_1234567890",
				ReturnURL:      "http://example.com/payment/return",
				IPNURL:         "https://example.com/api/webhook",
				PaymentTimeout: 15,
			},
			wantErr:     true,
			errContains: "return_url' must use HTTPS in production",
		},
		{
			name: "production requires https ipn url",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "PROD_MERCHANT",
				SecretKey:      "PROD_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "http://example.com/api/webhook",
				PaymentTimeout: 15,
			},
			wantErr:     true,
			errContains: "ipn_url' must use HTTPS in production",
		},
		{
			name: "production accepts https urls",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "PROD_MERCHANT",
				SecretKey:      "PROD_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/webhook",
				PaymentTimeout: 15,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				if err == nil {
					t.Error("Validate() expected error but got nil")
					return
				}
				if tt.errContains != "" && !contains(err.Error(), tt.errContains) {
					t.Errorf("Validate() error = %v, want error containing %v", err, tt.errContains)
				}
			} else {
				if err != nil {
					t.Errorf("Validate() unexpected error = %v", err)
				}
			}
		})
	}
}

// TestEnvironmentSwitchingConsistency verifies that environment methods
// are consistent with each other
func TestEnvironmentSwitchingConsistency(t *testing.T) {
	environments := []string{
		"sandbox", "SANDBOX", "Sandbox",
		"production", "PRODUCTION", "Production",
		"invalid", "test", "dev", "",
	}

	for _, env := range environments {
		t.Run("environment_"+env, func(t *testing.T) {
			config := VNPayConfig{Environment: env}

			isProduction := config.IsProduction()
			isSandbox := config.IsSandbox()

			// If it's production, it cannot be sandbox
			if isProduction && isSandbox {
				t.Errorf("Environment %q: cannot be both production and sandbox", env)
			}

			// If it's sandbox, it cannot be production
			if isSandbox && isProduction {
				t.Errorf("Environment %q: cannot be both sandbox and production", env)
			}

			// GetBaseURL should return appropriate URL
			baseURL := config.GetBaseURL()
			if isProduction {
				expectedURL := "https://vnpayment.vn/paymentv2/vpcpay.html"
				if baseURL != expectedURL {
					t.Errorf("Production environment %q: GetBaseURL() = %v, want %v", env, baseURL, expectedURL)
				}
			} else if isSandbox {
				expectedURL := "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"
				if baseURL != expectedURL {
					t.Errorf("Sandbox environment %q: GetBaseURL() = %v, want %v", env, baseURL, expectedURL)
				}
			}
		})
	}
}

// TestEnvironmentSwitchingGatewayURLs verifies that correct gateway URLs
// are returned for each environment
func TestEnvironmentSwitchingGatewayURLs(t *testing.T) {
	const (
		sandboxGatewayURL    = "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html"
		productionGatewayURL = "https://vnpayment.vn/paymentv2/vpcpay.html"
	)

	tests := []struct {
		name        string
		environment string
		wantURL     string
	}{
		{
			name:        "sandbox returns sandbox gateway",
			environment: "sandbox",
			wantURL:     sandboxGatewayURL,
		},
		{
			name:        "production returns production gateway",
			environment: "production",
			wantURL:     productionGatewayURL,
		},
		{
			name:        "SANDBOX uppercase returns sandbox gateway",
			environment: "SANDBOX",
			wantURL:     sandboxGatewayURL,
		},
		{
			name:        "PRODUCTION uppercase returns production gateway",
			environment: "PRODUCTION",
			wantURL:     productionGatewayURL,
		},
		{
			name:        "invalid environment defaults to sandbox gateway",
			environment: "invalid",
			wantURL:     sandboxGatewayURL,
		},
		{
			name:        "empty environment defaults to sandbox gateway",
			environment: "",
			wantURL:     sandboxGatewayURL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := VNPayConfig{Environment: tt.environment}
			gotURL := config.GetBaseURL()

			if gotURL != tt.wantURL {
				t.Errorf("GetBaseURL() = %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
