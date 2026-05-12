package payment

import (
	"encoding/json"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"
)

// **Validates: Requirements 11.5**
// Property: Configuration Round-Trip Integrity
// For all valid VNPay configurations, parsing then serializing then parsing
// should produce an equivalent configuration.
func TestConfigurationRoundTripIntegrity(t *testing.T) {
	// Property-based test using testing/quick
	f := func(config VNPayConfig) bool {
		// Skip invalid configurations - we only test round-trip for valid configs
		if err := config.Validate(); err != nil {
			return true // Skip invalid configs
		}

		// Step 1: Serialize the configuration to JSON
		serialized, err := json.Marshal(config)
		if err != nil {
			t.Logf("Failed to serialize config: %v", err)
			return false
		}

		// Step 2: Parse the serialized JSON back to a config object
		var parsed VNPayConfig
		err = json.Unmarshal(serialized, &parsed)
		if err != nil {
			t.Logf("Failed to parse serialized config: %v", err)
			return false
		}

		// Step 3: Serialize the parsed config again
		reserialized, err := json.Marshal(parsed)
		if err != nil {
			t.Logf("Failed to re-serialize parsed config: %v", err)
			return false
		}

		// Step 4: Parse the re-serialized JSON
		var reparsed VNPayConfig
		err = json.Unmarshal(reserialized, &reparsed)
		if err != nil {
			t.Logf("Failed to re-parse config: %v", err)
			return false
		}

		// Step 5: Verify that original config equals the final parsed config
		if !reflect.DeepEqual(config, reparsed) {
			t.Logf("Round-trip failed: original != reparsed\nOriginal: %+v\nReparsed: %+v", config, reparsed)
			return false
		}

		// Also verify that the serialized JSON strings are identical
		if string(serialized) != string(reserialized) {
			t.Logf("Serialized JSON differs:\nFirst:  %s\nSecond: %s", serialized, reserialized)
			return false
		}

		return true
	}

	// Custom configuration for quick.Check to generate more test cases
	config := &quick.Config{
		MaxCount: 100, // Run 100 random test cases
		Values: func(args []reflect.Value, r *rand.Rand) {
			// Generate a random VNPayConfig
			args[0] = reflect.ValueOf(generateValidConfig(r))
		},
	}

	if err := quick.Check(f, config); err != nil {
		t.Errorf("Configuration round-trip property failed: %v", err)
	}
}

// generateValidConfig creates a valid VNPayConfig for property-based testing
func generateValidConfig(r *rand.Rand) VNPayConfig {
	// Generate random but valid configuration values
	environments := []string{"sandbox", "production", "SANDBOX", "PRODUCTION", "Sandbox", "Production"}
	protocols := []string{"http://", "https://"}
	domains := []string{"example.com", "test.org", "localhost:3000", "api.example.com"}
	paths := []string{"/payment/return", "/api/webhook", "/callback", "/ipn"}

	// Select random environment
	env := environments[r.Intn(len(environments))]

	// Generate merchant ID (at least 4 characters)
	merchantID := randomString(r, 8, 20)

	// Generate secret key (at least 16 characters)
	secretKey := randomString(r, 16, 32)

	// Generate URLs
	// For production, always use HTTPS
	var returnProtocol, ipnProtocol string
	if env == "production" || env == "PRODUCTION" || env == "Production" {
		returnProtocol = "https://"
		ipnProtocol = "https://"
	} else {
		// For sandbox, can use either HTTP or HTTPS
		returnProtocol = protocols[r.Intn(len(protocols))]
		ipnProtocol = protocols[r.Intn(len(protocols))]
	}

	returnURL := returnProtocol + domains[r.Intn(len(domains))] + paths[r.Intn(len(paths))]
	ipnURL := ipnProtocol + domains[r.Intn(len(domains))] + paths[r.Intn(len(paths))]

	// Generate payment timeout (1 to 60 minutes)
	paymentTimeout := r.Intn(60) + 1

	return VNPayConfig{
		Environment:    env,
		MerchantID:     merchantID,
		SecretKey:      secretKey,
		ReturnURL:      returnURL,
		IPNURL:         ipnURL,
		PaymentTimeout: paymentTimeout,
	}
}

// randomString generates a random alphanumeric string of length between min and max
func randomString(r *rand.Rand, min, max int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-"
	length := min + r.Intn(max-min+1)
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

// TestConfigurationRoundTripWithSpecificExamples tests round-trip with known examples
func TestConfigurationRoundTripWithSpecificExamples(t *testing.T) {
	testCases := []struct {
		name   string
		config VNPayConfig
	}{
		{
			name: "sandbox configuration",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "TEST_MERCHANT_123",
				SecretKey:      "TEST_SECRET_KEY_ABCDEFGHIJKLMNOP",
				ReturnURL:      "https://example.com/payment/return",
				IPNURL:         "https://example.com/api/payment/webhook",
				PaymentTimeout: 15,
			},
		},
		{
			name: "production configuration",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "PROD_MERCHANT_456",
				SecretKey:      "PROD_SECRET_KEY_QRSTUVWXYZ123456",
				ReturnURL:      "https://myapp.com/payment/callback",
				IPNURL:         "https://myapp.com/api/v1/payments/ipn",
				PaymentTimeout: 30,
			},
		},
		{
			name: "uppercase environment",
			config: VNPayConfig{
				Environment:    "SANDBOX",
				MerchantID:     "UPPER_MERCHANT",
				SecretKey:      "UPPER_SECRET_KEY_1234567890ABCD",
				ReturnURL:      "http://localhost:3000/return",
				IPNURL:         "http://localhost:8080/webhook",
				PaymentTimeout: 10,
			},
		},
		{
			name: "mixed case environment",
			config: VNPayConfig{
				Environment:    "Production",
				MerchantID:     "Mixed_Merchant_789",
				SecretKey:      "Mixed_Secret_Key_ZYXWVUTSRQPONM",
				ReturnURL:      "https://secure.example.org/payments/return",
				IPNURL:         "https://secure.example.org/api/payments/ipn",
				PaymentTimeout: 45,
			},
		},
		{
			name: "minimum valid lengths",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "ABCD",
				SecretKey:      "1234567890ABCDEF",
				ReturnURL:      "https://a.co/r",
				IPNURL:         "https://a.co/i",
				PaymentTimeout: 1,
			},
		},
		{
			name: "long values",
			config: VNPayConfig{
				Environment:    "production",
				MerchantID:     "VERY_LONG_MERCHANT_ID_WITH_MANY_CHARACTERS_12345678901234567890",
				SecretKey:      "VERY_LONG_SECRET_KEY_WITH_MANY_CHARACTERS_ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890",
				ReturnURL:      "https://very-long-domain-name-for-testing.example.com/api/v1/payments/return/callback",
				IPNURL:         "https://very-long-domain-name-for-testing.example.com/api/v1/payments/ipn/webhook",
				PaymentTimeout: 60,
			},
		},
		{
			name: "special characters in URLs",
			config: VNPayConfig{
				Environment:    "sandbox",
				MerchantID:     "SPECIAL_MERCHANT",
				SecretKey:      "SPECIAL_SECRET_KEY_123456789",
				ReturnURL:      "https://example.com/payment/return?source=app&version=1.0",
				IPNURL:         "https://example.com/api/payment/webhook?token=abc123",
				PaymentTimeout: 20,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Validate the config first
			if err := tc.config.Validate(); err != nil {
				t.Fatalf("Test config is invalid: %v", err)
			}

			// Step 1: Serialize
			serialized, err := json.Marshal(tc.config)
			if err != nil {
				t.Fatalf("Failed to serialize config: %v", err)
			}

			// Step 2: Parse
			var parsed VNPayConfig
			err = json.Unmarshal(serialized, &parsed)
			if err != nil {
				t.Fatalf("Failed to parse serialized config: %v", err)
			}

			// Step 3: Re-serialize
			reserialized, err := json.Marshal(parsed)
			if err != nil {
				t.Fatalf("Failed to re-serialize parsed config: %v", err)
			}

			// Step 4: Re-parse
			var reparsed VNPayConfig
			err = json.Unmarshal(reserialized, &reparsed)
			if err != nil {
				t.Fatalf("Failed to re-parse config: %v", err)
			}

			// Verify equality
			if !reflect.DeepEqual(tc.config, reparsed) {
				t.Errorf("Round-trip failed: configs are not equal\nOriginal: %+v\nReparsed: %+v", tc.config, reparsed)
			}

			// Verify JSON equality
			if string(serialized) != string(reserialized) {
				t.Errorf("Serialized JSON differs:\nFirst:  %s\nSecond: %s", serialized, reserialized)
			}

			// Verify the reparsed config is still valid
			if err := reparsed.Validate(); err != nil {
				t.Errorf("Reparsed config failed validation: %v", err)
			}
		})
	}
}

// TestConfigurationRoundTripPreservesValidation tests that validation status is preserved
func TestConfigurationRoundTripPreservesValidation(t *testing.T) {
	t.Run("valid config remains valid after round-trip", func(t *testing.T) {
		config := VNPayConfig{
			Environment:    "sandbox",
			MerchantID:     "TEST_MERCHANT",
			SecretKey:      "TEST_SECRET_KEY_1234567890",
			ReturnURL:      "https://example.com/return",
			IPNURL:         "https://example.com/webhook",
			PaymentTimeout: 15,
		}

		// Verify original is valid
		if err := config.Validate(); err != nil {
			t.Fatalf("Original config should be valid: %v", err)
		}

		// Round-trip
		serialized, _ := json.Marshal(config)
		var parsed VNPayConfig
		json.Unmarshal(serialized, &parsed)

		// Verify parsed is still valid
		if err := parsed.Validate(); err != nil {
			t.Errorf("Parsed config should remain valid: %v", err)
		}
	})
}
