package payment

import (
	"net/url"
	"testing"
	"time"
	"travel-backend/domain"
)

func testVNPayConfig() *VNPayConfig {
	return &VNPayConfig{
		Environment:    "sandbox",
		MerchantID:     "DEMO1234",
		SecretKey:      "DEMO_SECRET_KEY_123456789",
		ReturnURL:      "http://localhost:8080/v1/api/payments/return",
		IPNURL:         "http://localhost:8080/v1/api/payments/webhook",
		FrontendURL:    "http://localhost:5173",
		PaymentTimeout: 15,
	}
}

func TestGeneratePaymentURLSendsVNPayAmountMultiplier(t *testing.T) {
	client := NewVNPayClient(testVNPayConfig())

	paymentURL, err := client.GeneratePaymentURL(&domain.VNPayPaymentRequest{
		TransactionReference: "PAY202605170001",
		Amount:               2500000,
		OrderInfo:            "Thanh toan booking TOUR-1",
		ClientIP:             "127.0.0.1",
		ReturnURL:            "http://localhost:8080/v1/api/payments/return",
		ExpiresAt:            time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		t.Fatalf("GeneratePaymentURL() unexpected error = %v", err)
	}

	parsed, err := url.Parse(paymentURL)
	if err != nil {
		t.Fatalf("url.Parse() unexpected error = %v", err)
	}

	if got := parsed.Query().Get("vnp_Amount"); got != "250000000" {
		t.Fatalf("vnp_Amount = %s, want 250000000", got)
	}
}

func TestVNPayConfigRequiresFrontendURL(t *testing.T) {
	config := testVNPayConfig()
	config.FrontendURL = ""

	if err := config.Validate(); err == nil {
		t.Fatal("Validate() expected error for missing frontend URL")
	}
}

func TestBuildPaymentResultURL(t *testing.T) {
	got := buildPaymentResultURL("https://traveling.example", "success", "PAY123", "Giao dịch thành công")

	parsed, err := url.Parse(got)
	if err != nil {
		t.Fatalf("url.Parse() unexpected error = %v", err)
	}

	if parsed.Scheme != "https" || parsed.Host != "traveling.example" || parsed.Path != "/payment/result" {
		t.Fatalf("result URL = %s, want https://traveling.example/payment/result", got)
	}
	if parsed.Query().Get("status") != "success" || parsed.Query().Get("ref") != "PAY123" {
		t.Fatalf("query = %s, want status/ref params", parsed.RawQuery)
	}
}
