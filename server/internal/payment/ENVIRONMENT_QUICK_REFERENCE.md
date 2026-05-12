# VNPay Environment Switching - Quick Reference

## Quick Start

### Sandbox (Development/Testing)
```bash
VNPAY_ENVIRONMENT=sandbox
VNPAY_MERCHANT_ID=your-sandbox-merchant-id
VNPAY_SECRET_KEY=your-sandbox-secret-key
VNPAY_RETURN_URL=http://localhost:5173/payment/return
VNPAY_IPN_URL=http://localhost:8080/api/v1/payments/webhook
```

### Production (Live Payments)
```bash
VNPAY_ENVIRONMENT=production
VNPAY_MERCHANT_ID=your-production-merchant-id
VNPAY_SECRET_KEY=your-production-secret-key
VNPAY_RETURN_URL=https://yourdomain.com/payment/return
VNPAY_IPN_URL=https://yourdomain.com/api/v1/payments/webhook
```

## Gateway URLs

| Environment | Gateway URL |
|-------------|-------------|
| Sandbox | `https://sandbox.vnpayment.vn/paymentv2/vpcpay.html` |
| Production | `https://vnpayment.vn/paymentv2/vpcpay.html` |

## Code Examples

### Check Environment
```go
config, _ := payment.LoadConfigFromEnv()

if config.IsProduction() {
    // Production mode
}

if config.IsSandbox() {
    // Sandbox mode
}
```

### Get Gateway URL
```go
config, _ := payment.LoadConfigFromEnv()
gatewayURL := config.GetBaseURL()
// Returns correct URL based on environment
```

## Key Differences

| Feature | Sandbox | Production |
|---------|---------|------------|
| Real Money | ❌ | ✅ |
| HTTPS Required | ❌ | ✅ |
| Test Cards | ✅ | ❌ |
| Business Verification | ❌ | ✅ |

## Common Commands

```bash
# Switch to sandbox
export VNPAY_ENVIRONMENT=sandbox

# Switch to production
export VNPAY_ENVIRONMENT=production

# Run tests
go test ./server/internal/payment -v

# Check configuration
go run main.go
```

## Troubleshooting

### Error: "environment must be 'sandbox' or 'production'"
✅ Fix: Use `sandbox` or `production` (case-insensitive)

### Error: "return_url must use HTTPS in production"
✅ Fix: Change `http://` to `https://` in production

### Payment works in sandbox but not production
✅ Check: Using production credentials (not sandbox)
✅ Check: URLs use HTTPS
✅ Check: SSL certificate is valid

## Get Credentials

**Sandbox**: https://sandbox.vnpayment.vn/devreg
**Production**: Contact VNPay sales team

## Full Documentation

See `ENVIRONMENT_SWITCHING.md` for complete guide.
