package payment

import "travel-backend/domain"

// PaymentGateway defines the provider boundary for payment integrations.
// VNPay is the first implementation; additional providers can implement this
// interface without changing booking/payment orchestration.
type PaymentGateway interface {
	ProviderName() string
	GeneratePaymentURL(req *domain.VNPayPaymentRequest) (string, error)
	ValidateSignature(params map[string]string, receivedHash string) bool
	ParseResponseCode(code string) (bool, string)
}
