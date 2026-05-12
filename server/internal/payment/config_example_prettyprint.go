package payment

import "fmt"

// ExamplePrettyPrint demonstrates the pretty print functionality
func ExamplePrettyPrint() {
	config := VNPayConfig{
		Environment:    "sandbox",
		MerchantID:     "TEST_MERCHANT_123",
		SecretKey:      "SUPER_SECRET_KEY_DO_NOT_EXPOSE_1234567890",
		ReturnURL:      "https://example.com/payment/return",
		IPNURL:         "https://example.com/api/payment/webhook",
		PaymentTimeout: 15,
	}

	// Pretty print with indentation
	prettyOutput, err := config.PrettyPrint()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Pretty Print Output:")
	fmt.Println(prettyOutput)
	fmt.Println()

	// Compact print without indentation
	compactOutput, err := config.CompactPrint()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("Compact Print Output:")
	fmt.Println(compactOutput)
}

// ExampleMaskedCopy demonstrates the masked copy functionality
func ExampleMaskedCopy() {
	config := VNPayConfig{
		Environment:    "production",
		MerchantID:     "PROD_MERCHANT_456",
		SecretKey:      "PRODUCTION_SECRET_KEY_ABCDEFGHIJKLMNOP",
		ReturnURL:      "https://myapp.com/payment/return",
		IPNURL:         "https://myapp.com/api/payment/webhook",
		PaymentTimeout: 30,
	}

	// Create a masked copy for safe logging
	masked := config.MaskedCopy()

	fmt.Printf("Original Secret Key: %s\n", config.SecretKey)
	fmt.Printf("Masked Secret Key: %s\n", masked.SecretKey)
	fmt.Printf("Other fields preserved: Environment=%s, MerchantID=%s\n",
		masked.Environment, masked.MerchantID)
}
