# VNPay Payment Integration Module

## Overview

This module provides VNPay payment gateway integration for the travel booking system, supporting both sandbox (testing) and production (live payment) environments.

## Features

✅ **Environment Switching**: Easy switching between sandbox and production
✅ **Automatic Gateway Selection**: Correct VNPay URL based on environment
✅ **Security Validation**: HTTPS enforcement in production
✅ **Configuration Management**: Environment variable-based configuration
✅ **Comprehensive Testing**: Unit tests for all functionality
✅ **Clear Documentation**: Detailed guides and examples

## Quick Start

### 1. Configure Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
# For sandbox (development/testing)
VNPAY_ENVIRONMENT=sandbox
VNPAY_MERCHANT_ID=your-sandbox-merchant-id
VNPAY_SECRET_KEY=your-sandbox-secret-key
VNPAY_RETURN_URL=http://localhost:5173/payment/return
VNPAY_IPN_URL=http://localhost:8080/api/v1/payments/webhook
VNPAY_PAYMENT_TIMEOUT=15
```

### 2. Load Configuration

```go
import "your-project/server/internal/payment"

config, err := payment.LoadConfigFromEnv()
if err != nil {
    log.Fatalf("Failed to load VNPay configuration: %v", err)
}

// Check environment
if config.IsProduction() {
    log.Println("Running in PRODUCTION mode")
} else if config.IsSandbox() {
    log.Println("Running in SANDBOX mode")
}

// Get gateway URL
gatewayURL := config.GetBaseURL()
```

### 3. Use Configuration

```go
// Configuration is automatically validated
// Gateway URL is automatically selected based on environment
// HTTPS is automatically enforced in production
```

## Module Structure

```
server/internal/payment/
├── config.go                           # Configuration struct and validation
├── config_test.go                      # Configuration unit tests
├── config_roundtrip_test.go            # Round-trip property tests
├── config_prettyprint_test.go          # Pretty print tests
├── environment_switching_test.go       # Environment switching tests
├── config_example.go                   # Usage examples
├── config_example_prettyprint.go       # Pretty print examples
├── README.md                           # This file
├── ENVIRONMENT_SWITCHING.md            # Comprehensive environment guide
├── ENVIRONMENT_QUICK_REFERENCE.md      # Quick reference guide
├── VALIDATION_IMPROVEMENTS.md          # Validation documentation
└── PRETTYPRINT_IMPLEMENTATION.md       # Pretty print documentation
```

## Documentation

### Comprehensive Guides

- **[ENVIRONMENT_SWITCHING.md](./ENVIRONMENT_SWITCHING.md)**: Complete guide to environment switching
  - Sandbox vs Production comparison
  - Configuration methods
  - Security requirements
  - Troubleshooting
  - Best practices

- **[ENVIRONMENT_QUICK_REFERENCE.md](./ENVIRONMENT_QUICK_REFERENCE.md)**: Quick reference for common tasks
  - Quick start commands
  - Common configurations
  - Troubleshooting tips

### Technical Documentation

- **[VALIDATION_IMPROVEMENTS.md](./VALIDATION_IMPROVEMENTS.md)**: Configuration validation details
- **[PRETTYPRINT_IMPLEMENTATION.md](./PRETTYPRINT_IMPLEMENTATION.md)**: Pretty print functionality

## Environment Support

### Sandbox Environment

**Purpose**: Development and testing with test credentials

**Gateway URL**: `https://sandbox.vnpayment.vn/paymentv2/vpcpay.html`

**Features**:
- HTTP URLs allowed for local development
- Test payment cards provided by VNPay
- No real money transactions
- Instant payment confirmation

**Get Credentials**: https://sandbox.vnpayment.vn/devreg

### Production Environment

**Purpose**: Live payment processing with real money

**Gateway URL**: `https://vnpayment.vn/paymentv2/vpcpay.html`

**Requirements**:
- HTTPS required for all URLs
- Valid SSL certificate
- Business verification completed
- Signed merchant agreement

**Get Credentials**: Contact VNPay sales team

## Configuration Methods

### Method 1: IsProduction()

Check if running in production mode:

```go
if config.IsProduction() {
    // Production-specific logic
    enableProductionMonitoring()
}
```

### Method 2: IsSandbox()

Check if running in sandbox mode:

```go
if config.IsSandbox() {
    // Sandbox-specific logic
    enableDebugLogging()
}
```

### Method 3: GetBaseURL()

Get environment-specific gateway URL:

```go
gatewayURL := config.GetBaseURL()
// Sandbox: https://sandbox.vnpayment.vn/paymentv2/vpcpay.html
// Production: https://vnpayment.vn/paymentv2/vpcpay.html
```

## Switching Environments

To switch from sandbox to production:

1. Update `VNPAY_ENVIRONMENT=production`
2. Update `VNPAY_MERCHANT_ID` with production merchant ID
3. Update `VNPAY_SECRET_KEY` with production secret key
4. Change URLs to use HTTPS:
   - `VNPAY_RETURN_URL=https://yourdomain.com/payment/return`
   - `VNPAY_IPN_URL=https://yourdomain.com/api/v1/payments/webhook`
5. Restart application

The system will automatically:
- Select production gateway URL
- Enforce HTTPS validation
- Validate configuration on startup

## Testing

### Run All Tests

```bash
cd server/internal/payment
go test -v
```

### Run Specific Test Suites

```bash
# Configuration validation tests
go test -v -run TestVNPayConfig_Validate

# Environment switching tests
go test -v -run TestEnvironmentSwitching

# Round-trip property tests
go test -v -run TestConfigurationRoundTrip

# Pretty print tests
go test -v -run TestPrettyPrint
```

### Test Coverage

```bash
go test -cover
```

## Security

### Credential Protection

- ✅ Never commit credentials to version control
- ✅ Use environment variables for configuration
- ✅ Use `.env` files for local development (add to `.gitignore`)
- ✅ Use secret management systems for production

### Sensitive Data Masking

```go
// Mask secret key for logging
log.Printf("Secret Key: %s", config.MaskSecretKey())
// Output: Secret Key: ABCD****WXYZ

// Pretty print with masked secrets
output, _ := config.PrettyPrint()
log.Println(output)
```

### HTTPS Enforcement

Production mode automatically enforces HTTPS:

```go
// This will fail in production mode
config := VNPayConfig{
    Environment: "production",
    ReturnURL:   "http://example.com/return",  // ❌ HTTP not allowed
    // ...
}

err := config.Validate()
// Error: "return_url must use HTTPS in production"
```

## Examples

### Basic Usage

```go
package main

import (
    "log"
    "your-project/server/internal/payment"
)

func main() {
    // Load configuration
    config, err := payment.LoadConfigFromEnv()
    if err != nil {
        log.Fatalf("Configuration error: %v", err)
    }

    // Log configuration (with masked secrets)
    log.Printf("Environment: %s", config.Environment)
    log.Printf("Gateway URL: %s", config.GetBaseURL())
    log.Printf("Secret Key: %s", config.MaskSecretKey())

    // Environment-specific logic
    if config.IsProduction() {
        log.Println("🔴 PRODUCTION MODE")
        enableProductionMonitoring()
    } else {
        log.Println("🟡 SANDBOX MODE")
        enableDebugLogging()
    }
}
```

### Environment-Specific Configuration

```go
config, _ := payment.LoadConfigFromEnv()

if config.IsProduction() {
    // Production settings
    setLogLevel("ERROR")
    enableMetrics()
    enableAlerting()
    setStrictValidation()
} else {
    // Sandbox settings
    setLogLevel("DEBUG")
    enableTestMode()
    relaxValidation()
}
```

## Troubleshooting

### Common Issues

**Issue**: "environment must be 'sandbox' or 'production'"
- **Solution**: Check `VNPAY_ENVIRONMENT` value (must be `sandbox` or `production`)

**Issue**: "return_url must use HTTPS in production"
- **Solution**: Change `http://` to `https://` in production environment

**Issue**: Configuration validation fails on startup
- **Solution**: Check all required environment variables are set
- **Solution**: Verify URLs are properly formatted
- **Solution**: Ensure credentials are not empty

### Debug Configuration

```go
config, err := payment.LoadConfigFromEnv()
if err != nil {
    log.Printf("Configuration error: %v", err)
    // Check environment variables
    log.Printf("VNPAY_ENVIRONMENT: %s", os.Getenv("VNPAY_ENVIRONMENT"))
    log.Printf("VNPAY_MERCHANT_ID: %s", os.Getenv("VNPAY_MERCHANT_ID"))
    // ... check other variables
}
```

## API Reference

### VNPayConfig

```go
type VNPayConfig struct {
    Environment    string // "sandbox" or "production"
    MerchantID     string // VNPay merchant/terminal code
    SecretKey      string // Secret key for signature generation
    ReturnURL      string // User-facing payment return URL
    IPNURL         string // Server-to-server webhook URL
    PaymentTimeout int    // Payment session timeout in minutes
}
```

### Methods

- `Validate() error` - Validate configuration
- `IsProduction() bool` - Check if production environment
- `IsSandbox() bool` - Check if sandbox environment
- `GetBaseURL() string` - Get environment-specific gateway URL
- `MaskSecretKey() string` - Get masked secret key for logging
- `PrettyPrint() (string, error)` - Format config as pretty JSON
- `CompactPrint() (string, error)` - Format config as compact JSON
- `MaskedCopy() VNPayConfig` - Create copy with masked secrets

### Functions

- `LoadConfigFromEnv() (*VNPayConfig, error)` - Load config from environment variables

## Support

For questions or issues:

- **VNPay Sandbox**: https://sandbox.vnpayment.vn/devreg
- **VNPay Documentation**: https://vnpay.vn
- **Environment Switching Guide**: [ENVIRONMENT_SWITCHING.md](./ENVIRONMENT_SWITCHING.md)
- **Quick Reference**: [ENVIRONMENT_QUICK_REFERENCE.md](./ENVIRONMENT_QUICK_REFERENCE.md)

## License

This module is part of the travel booking system project.
