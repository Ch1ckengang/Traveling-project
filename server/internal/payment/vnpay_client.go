package payment

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
	"travel-backend/domain"
)

// VNPayClient - Client xử lý giao tiếp với VNPay payment gateway
type VNPayClient struct {
	config *VNPayConfig
}

// NewVNPayClient - Tạo VNPay client mới từ config
func NewVNPayClient(config *VNPayConfig) *VNPayClient {
	return &VNPayClient{config: config}
}

// GeneratePaymentURL - Tạo URL thanh toán VNPay để redirect user
func (c *VNPayClient) GeneratePaymentURL(req *domain.VNPayPaymentRequest) (string, error) {
	if req == nil {
		return "", fmt.Errorf("payment request cannot be nil")
	}

	now := time.Now()
	createDate := now.Format("20060102150405")
	expireDate := req.ExpiresAt.Format("20060102150405")

	params := map[string]string{
		"vnp_Version":    "2.1.0",
		"vnp_Command":    "pay",
		"vnp_TmnCode":    c.config.MerchantID,
		"vnp_Locale":     "vn",
		"vnp_CurrCode":   "VND",
		"vnp_TxnRef":     req.TransactionReference,
		"vnp_OrderInfo":  req.OrderInfo,
		"vnp_OrderType":  "other",
		"vnp_Amount":     fmt.Sprintf("%d", req.Amount*100), // VNPay yêu cầu nhân 100
		"vnp_ReturnUrl":  req.ReturnURL,
		"vnp_IpAddr":     req.ClientIP,
		"vnp_CreateDate": createDate,
		"vnp_ExpireDate": expireDate,
	}

	// Tạo secure hash
	hashData := c.buildHashData(params)
	secureHash := c.generateHMACSHA512(hashData)

	// Build query string
	queryParams := url.Values{}
	for key, value := range params {
		queryParams.Set(key, value)
	}
	queryParams.Set("vnp_SecureHash", secureHash)

	paymentURL := fmt.Sprintf("%s?%s", c.config.GetBaseURL(), queryParams.Encode())
	return paymentURL, nil
}

// ValidateSignature - Kiểm tra chữ ký từ VNPay callback/webhook
func (c *VNPayClient) ValidateSignature(params map[string]string, receivedHash string) bool {
	hashData := c.buildHashData(params)
	expectedHash := c.generateHMACSHA512(hashData)
	return strings.EqualFold(expectedHash, receivedHash)
}

// buildHashData - Xây dựng chuỗi hash data từ params (sorted by key)
func (c *VNPayClient) buildHashData(params map[string]string) string {
	// Sort keys alphabetically
	keys := make([]string, 0, len(params))
	for k := range params {
		if k != "vnp_SecureHash" && k != "vnp_SecureHashType" && params[k] != "" {
			keys = append(keys, k)
		}
	}
	sort.Strings(keys)

	// Build hash data string
	var parts []string
	for _, k := range keys {
		value := url.QueryEscape(params[k])
		parts = append(parts, fmt.Sprintf("%s=%s", k, value))
	}

	return strings.Join(parts, "&")
}

// generateHMACSHA512 - Tạo HMAC-SHA512 signature
func (c *VNPayClient) generateHMACSHA512(data string) string {
	h := hmac.New(sha512.New, []byte(c.config.SecretKey))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

// ParseVNPayResponseCode - Giải mã VNPay response code
func (c *VNPayClient) ParseVNPayResponseCode(code string) (bool, string) {
	isSuccess := code == "00"
	message := domain.GetVNPayErrorMessage(code)
	return isSuccess, message
}
