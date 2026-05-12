package payment

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestVNPayConfig_PrettyPrint tests the pretty print functionality
func TestVNPayConfig_PrettyPrint(t *testing.T) {
	tests := []struct {
		name       string
		config     VNPayConfig
		wantErr    bool
		checkFunc  func(*testing.T, string)
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
			checkFunc: func(t *testing.T, output string) {
				// Check that output is valid JSON
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(output), &parsed); err != nil {
					t.Errorf("PrettyPrint() output is not valid JSON: %v", err)
				}

				// Check that secret key is masked
				if strings.Contains(output, "TEST_SECRET_KEY_1234567890") {
					t.Errorf("PrettyPrint() should mask secret key, but found unmasked value")
				}

				// Check that masked secret key is present
				if !strings.Contains(output, "TEST****7890") {
					t.Errorf("PrettyPrint() should contain masked secret key 'TEST****7890', got: %s", output)
				}

				// Check that other fields are present
				if !strings.Contains(output, "TEST_MERCHANT") {
					t.Errorf("PrettyPrint() should contain merchant ID")
				}
				if !strings.Contains(output, "sandbox") {
					t.Errorf("PrettyPrint() should contain environment")
				}
				if !strings.Contains(output, "https://example.com/payment/return") {
					t.Errorf("PrettyPrint() should contain return URL")
				}

				// Check that output has proper indentation (pretty format)
				if !strings.Contains(output, "\n") {
					t.Errorf("PrettyPrint() should have newlines for pretty formatting")
				}
				if !strings.Contains(output, "  ") {
					t.Errorf("PrettyPrint() should have indentation")
				}
			},
		},
		{
			name: "production configuration with long secret key",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "PROD_MERCHANT_123",
				SecretKey:      "VERY_LONG_PRODUCTION_SECRET_KEY_ABCDEFGHIJKLMNOPQRSTUVWXYZ",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 30,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, output string) {
				// Check that full secret key is not exposed
				if strings.Contains(output, "VERY_LONG_PRODUCTION_SECRET_KEY_ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
					t.Errorf("PrettyPrint() should mask secret key")
				}

				// Check that masked version is present
				if !strings.Contains(output, "VERY****WXYZ") {
					t.Errorf("PrettyPrint() should contain masked secret key 'VERY****WXYZ'")
				}

				// Check production environment is present
				if !strings.Contains(output, "production") {
					t.Errorf("PrettyPrint() should contain production environment")
				}
			},
		},
		{
			name: "configuration with short secret key",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "SHORT",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, output string) {
				// Short secret keys should be completely masked
				if strings.Contains(output, "SHORT") {
					t.Errorf("PrettyPrint() should mask short secret key")
				}

				// Check that masked version is present
				if !strings.Contains(output, "****") {
					t.Errorf("PrettyPrint() should contain masked secret key '****'")
				}
			},
		},
		{
			name: "configuration with special characters in URLs",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return?source=app&version=1.0",
				IPNURL:         "https://example.com/api/payment/webhook?token=abc123",
				PaymentTimeout: 15,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, output string) {
				// Check that URLs with special characters are properly escaped
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(output), &parsed); err != nil {
					t.Errorf("PrettyPrint() output is not valid JSON: %v", err)
				}

				// URLs should be present and properly formatted
				if !strings.Contains(output, "return_url") {
					t.Errorf("PrettyPrint() should contain return_url field")
				}
				if !strings.Contains(output, "ipn_url") {
					t.Errorf("PrettyPrint() should contain ipn_url field")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.config.PrettyPrint()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("PrettyPrint() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("PrettyPrint() unexpected error = %v", err)
				return
			}

			if output == "" {
				t.Errorf("PrettyPrint() returned empty string")
				return
			}

			// Run custom checks
			if tt.checkFunc != nil {
				tt.checkFunc(t, output)
			}
		})
	}
}

// TestVNPayConfig_CompactPrint tests the compact print functionality
func TestVNPayConfig_CompactPrint(t *testing.T) {
	tests := []struct {
		name       string
		config     VNPayConfig
		wantErr    bool
		checkFunc  func(*testing.T, string)
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
			checkFunc: func(t *testing.T, output string) {
				// Check that output is valid JSON
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(output), &parsed); err != nil {
					t.Errorf("CompactPrint() output is not valid JSON: %v", err)
				}

				// Check that secret key is masked
				if strings.Contains(output, "TEST_SECRET_KEY_1234567890") {
					t.Errorf("CompactPrint() should mask secret key")
				}

				// Check that masked secret key is present
				if !strings.Contains(output, "TEST****7890") {
					t.Errorf("CompactPrint() should contain masked secret key")
				}

				// Check that output is compact (no extra whitespace)
				if strings.Contains(output, "\n") {
					t.Errorf("CompactPrint() should not have newlines")
				}
				if strings.Contains(output, "  ") {
					t.Errorf("CompactPrint() should not have extra indentation")
				}

				// Check that other fields are present
				if !strings.Contains(output, "TEST_MERCHANT") {
					t.Errorf("CompactPrint() should contain merchant ID")
				}
			},
		},
		{
			name: "production configuration",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "PROD_MERCHANT",
				SecretKey:      "PROD_SECRET_KEY_ABCDEFGHIJKLMNOP",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 30,
			},
			wantErr: false,
			checkFunc: func(t *testing.T, output string) {
				// Check that secret key is masked
				if strings.Contains(output, "PROD_SECRET_KEY_ABCDEFGHIJKLMNOP") {
					t.Errorf("CompactPrint() should mask secret key")
				}

				// Check compact format
				if strings.Contains(output, "\n") {
					t.Errorf("CompactPrint() should be on a single line")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output, err := tt.config.CompactPrint()
			
			if tt.wantErr {
				if err == nil {
					t.Errorf("CompactPrint() expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("CompactPrint() unexpected error = %v", err)
				return
			}

			if output == "" {
				t.Errorf("CompactPrint() returned empty string")
				return
			}

			// Run custom checks
			if tt.checkFunc != nil {
				tt.checkFunc(t, output)
			}
		})
	}
}

// TestVNPayConfig_MaskedCopy tests the masked copy functionality
func TestVNPayConfig_MaskedCopy(t *testing.T) {
	tests := []struct {
		name   string
		config VNPayConfig
		check  func(*testing.T, VNPayConfig, VNPayConfig)
	}{
		{
			name: "masked copy preserves all fields except secret key",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			check: func(t *testing.T, original, masked VNPayConfig) {
				// Check that non-sensitive fields are preserved
				if masked.Environment != original.Environment {
					t.Errorf("MaskedCopy() Environment = %v, want %v", masked.Environment, original.Environment)
				}
				if masked.MerchantID != original.MerchantID {
					t.Errorf("MaskedCopy() MerchantID = %v, want %v", masked.MerchantID, original.MerchantID)
				}
				if masked.ReturnURL != original.ReturnURL {
					t.Errorf("MaskedCopy() ReturnURL = %v, want %v", masked.ReturnURL, original.ReturnURL)
				}
				if masked.IPNURL != original.IPNURL {
					t.Errorf("MaskedCopy() IPNURL = %v, want %v", masked.IPNURL, original.IPNURL)
				}
				if masked.PaymentTimeout != original.PaymentTimeout {
					t.Errorf("MaskedCopy() PaymentTimeout = %v, want %v", masked.PaymentTimeout, original.PaymentTimeout)
				}

				// Check that secret key is masked
				if masked.SecretKey == original.SecretKey {
					t.Errorf("MaskedCopy() should mask secret key")
				}
				if masked.SecretKey != "TEST****7890" {
					t.Errorf("MaskedCopy() SecretKey = %v, want TEST****7890", masked.SecretKey)
				}
			},
		},
		{
			name: "masked copy with short secret key",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "SHORT",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
			check: func(t *testing.T, original, masked VNPayConfig) {
				// Short secret keys should be completely masked
				if masked.SecretKey != "****" {
					t.Errorf("MaskedCopy() SecretKey = %v, want ****", masked.SecretKey)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			masked := tt.config.MaskedCopy()
			
			if tt.check != nil {
				tt.check(t, tt.config, masked)
			}
		})
	}
}

// **Validates: Requirements 12.5**
// Property: Round-trip property for pretty printer
// For all valid VNPay configuration objects, parsing then printing then parsing
// should produce equivalent configuration (with masked secret key)
func TestPrettyPrintRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		config VNPayConfig
	}{
		{
			name: "sandbox configuration",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      "TEST_SECRET_KEY_1234567890",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
		},
		{
			name: "production configuration",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "PROD_MERCHANT",
				SecretKey:      "PROD_SECRET_KEY_ABCDEFGHIJKLMNOP",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Step 1: Pretty print the configuration
			prettyOutput, err := tt.config.PrettyPrint()
			if err != nil {
				t.Fatalf("PrettyPrint() error = %v", err)
			}

			// Step 2: Parse the pretty printed output
			var parsed VNPayConfig
			err = json.Unmarshal([]byte(prettyOutput), &parsed)
			if err != nil {
				t.Fatalf("Failed to parse pretty printed output: %v", err)
			}

			// Step 3: Pretty print the parsed configuration
			prettyOutput2, err := parsed.PrettyPrint()
			if err != nil {
				t.Fatalf("PrettyPrint() on parsed config error = %v", err)
			}

			// Step 4: Parse again
			var parsed2 VNPayConfig
			err = json.Unmarshal([]byte(prettyOutput2), &parsed2)
			if err != nil {
				t.Fatalf("Failed to parse second pretty printed output: %v", err)
			}

			// Step 5: Verify that the two parsed configs are equivalent
			// Note: Secret keys will be masked, so we compare the masked versions
			if parsed.Environment != parsed2.Environment {
				t.Errorf("Round-trip failed: Environment differs")
			}
			if parsed.MerchantID != parsed2.MerchantID {
				t.Errorf("Round-trip failed: MerchantID differs")
			}
			if parsed.SecretKey != parsed2.SecretKey {
				t.Errorf("Round-trip failed: SecretKey differs")
			}
			if parsed.ReturnURL != parsed2.ReturnURL {
				t.Errorf("Round-trip failed: ReturnURL differs")
			}
			if parsed.IPNURL != parsed2.IPNURL {
				t.Errorf("Round-trip failed: IPNURL differs")
			}
			if parsed.PaymentTimeout != parsed2.PaymentTimeout {
				t.Errorf("Round-trip failed: PaymentTimeout differs")
			}

			// Verify that the pretty printed outputs are identical
			if prettyOutput != prettyOutput2 {
				t.Errorf("Round-trip failed: pretty printed outputs differ\nFirst:\n%s\nSecond:\n%s", prettyOutput, prettyOutput2)
			}
		})
	}
}

// TestCompactPrintRoundTrip tests round-trip for compact print
func TestCompactPrintRoundTrip(t *testing.T) {
	config := VNPayConfig{
		Environment:    "sandbox",
		MerchantID:     "TEST_MERCHANT",
		SecretKey:      "TEST_SECRET_KEY_1234567890",
		ReturnURL:      "https://example.com/payment/return",
		IPNURL:         "https://example.com/api/payment/webhook",
		PaymentTimeout: 15,
	}

	// Step 1: Compact print
	compactOutput, err := config.CompactPrint()
	if err != nil {
		t.Fatalf("CompactPrint() error = %v", err)
	}

	// Step 2: Parse
	var parsed VNPayConfig
	err = json.Unmarshal([]byte(compactOutput), &parsed)
	if err != nil {
		t.Fatalf("Failed to parse compact output: %v", err)
	}

	// Step 3: Compact print again
	compactOutput2, err := parsed.CompactPrint()
	if err != nil {
		t.Fatalf("CompactPrint() on parsed config error = %v", err)
	}

	// Step 4: Verify outputs are identical
	if compactOutput != compactOutput2 {
		t.Errorf("Round-trip failed: compact outputs differ\nFirst:  %s\nSecond: %s", compactOutput, compactOutput2)
	}
}

// TestPrettyPrintVsCompactPrint verifies that both formats produce equivalent data
func TestPrettyPrintVsCompactPrint(t *testing.T) {
	config := VNPayConfig{
		Environment:    "sandbox",
		MerchantID:     "TEST_MERCHANT",
		SecretKey:      "TEST_SECRET_KEY_1234567890",
		ReturnURL:      "https://example.com/payment/return",
		IPNURL:         "https://example.com/api/payment/webhook",
		PaymentTimeout: 15,
	}

	// Get both formats
	prettyOutput, err := config.PrettyPrint()
	if err != nil {
		t.Fatalf("PrettyPrint() error = %v", err)
	}

	compactOutput, err := config.CompactPrint()
	if err != nil {
		t.Fatalf("CompactPrint() error = %v", err)
	}

	// Parse both
	var prettyParsed, compactParsed VNPayConfig
	
	err = json.Unmarshal([]byte(prettyOutput), &prettyParsed)
	if err != nil {
		t.Fatalf("Failed to parse pretty output: %v", err)
	}

	err = json.Unmarshal([]byte(compactOutput), &compactParsed)
	if err != nil {
		t.Fatalf("Failed to parse compact output: %v", err)
	}

	// Verify both produce the same data
	if prettyParsed.Environment != compactParsed.Environment {
		t.Errorf("Environment differs between formats")
	}
	if prettyParsed.MerchantID != compactParsed.MerchantID {
		t.Errorf("MerchantID differs between formats")
	}
	if prettyParsed.SecretKey != compactParsed.SecretKey {
		t.Errorf("SecretKey differs between formats")
	}
	if prettyParsed.ReturnURL != compactParsed.ReturnURL {
		t.Errorf("ReturnURL differs between formats")
	}
	if prettyParsed.IPNURL != compactParsed.IPNURL {
		t.Errorf("IPNURL differs between formats")
	}
	if prettyParsed.PaymentTimeout != compactParsed.PaymentTimeout {
		t.Errorf("PaymentTimeout differs between formats")
	}
}

// TestPrettyPrintSecurityMasking verifies that sensitive data is always masked
func TestPrettyPrintSecurityMasking(t *testing.T) {
	sensitiveSecrets := []string{
		"SUPER_SECRET_KEY_DO_NOT_EXPOSE_1234567890",
		"PRODUCTION_KEY_ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"API_SECRET_TOKEN_9876543210",
		"CONFIDENTIAL_KEY_ZYXWVUTSRQPONMLKJIHGFEDCBA",
	}

	for _, secret := range sensitiveSecrets {
		t.Run("masking_"+secret[:10], func(t *testing.T) {
			config := VNPayConfig{
				Environment:    "production",
				MerchantID:     "TEST_MERCHANT",
				SecretKey:      secret,
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			}

			// Test pretty print
			prettyOutput, err := config.PrettyPrint()
			if err != nil {
				t.Fatalf("PrettyPrint() error = %v", err)
			}

			if strings.Contains(prettyOutput, secret) {
				t.Errorf("PrettyPrint() exposed secret key: %s", secret)
			}

			// Test compact print
			compactOutput, err := config.CompactPrint()
			if err != nil {
				t.Fatalf("CompactPrint() error = %v", err)
			}

			if strings.Contains(compactOutput, secret) {
				t.Errorf("CompactPrint() exposed secret key: %s", secret)
			}
		})
	}
}
