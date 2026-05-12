# VNPay Environment Switching Guide

## Overview

The VNPay payment integration supports two environments:
- **Sandbox**: For development and testing with test credentials
- **Production**: For live payment processing with real money

This guide explains how to configure, switch between, and use these environments.

## Environment Configuration

### Configuration Methods

The system loads VNPay configuration from environment variables using the `LoadConfigFromEnv()` function. The environment is determined by the `VNPAY_ENVIRONMENT` variable.

### Required Environment Variables

```bash
# Environment: sandbox or production (default: sandbox)
VNPAY_ENVIRONMENT=sandbox

# Merchant credentials (different for sandbox vs production)
VNPAY_MERCHANT_ID=your-merchant-id
VNPAY_SECRET_KEY=your-secret-key

# Callback URLs (must use HTTPS in production)
VNPAY_RETURN_URL=http://localhost:5173/payment/return
VNPAY_IPN_URL=http://localhost:8080/api/v1/payments/webhook

# Payment session timeout in minutes (default: 15)
VNPAY_PAYMENT_TIMEOUT=15
```

## Sandbox Environment

### Purpose
- Development and testing
- No real money transactions
- Test payment scenarios without financial risk
- Validate integration before production deployment

### Sandbox Configuration

```bash
# .env file for sandbox
VNPAY_ENVIRONMENT=sandbox
VNPAY_MERCHANT_ID=your-sandbox-merchant-id
VNPAY_SECRET_KEY=your-sandbox-secret-key
VNPAY_RETURN_URL=http://localhost:5173/payment/return
VNPAY_IPN_URL=http://localhost:8080/api/v1/payments/webhook
VNPAY_PAYMENT_TIMEOUT=15
```

### Sandbox Gateway URL
```
https://sandbox.vnpayment.vn/paymentv2/vpcpay.html
```

### Getting Sandbox Credentials

1. Visit [VNPay Sandbox Registration](https://sandbox.vnpayment.vn/devreg)
2. Register for a sandbox account
3. Complete the registration form
4. Receive your sandbox credentials:
   - **Merchant ID** (TmnCode)
   - **Secret Key** (HashSecret)
5. Update your `.env` file with these credentials

### Sandbox Features

- ✅ HTTP URLs allowed for local development
- ✅ Test payment cards provided by VNPay
- ✅ Simulated payment responses
- ✅ No SSL certificate required for testing
- ✅ Instant payment confirmation
- ✅ Full API feature parity with production

### Sandbox Test Cards

VNPay provides test cards for sandbox testing. Refer to VNPay sandbox documentation for:
- Test card numbers
- Test bank accounts
- Test QR codes
- Simulated payment scenarios (success, failure, timeout)

## Production Environment

### Purpose
- Live payment processing
- Real money transactions
- Customer-facing payment gateway
- Production-grade security and reliability

### Production Configuration

```bash
# .env file for production
VNPAY_ENVIRONMENT=production
VNPAY_MERCHANT_ID=your-production-merchant-id
VNPAY_SECRET_KEY=your-production-secret-key
VNPAY_RETURN_URL=https://yourdomain.com/payment/return
VNPAY_IPN_URL=https://yourdomain.com/api/v1/payments/webhook
VNPAY_PAYMENT_TIMEOUT=15
```

### Production Gateway URL
```
https://vnpayment.vn/paymentv2/vpcpay.html
```

### Getting Production Credentials

1. Contact VNPay sales team
2. Complete business verification process
3. Sign merchant agreement and contract
4. Provide required business documents:
   - Business registration certificate
   - Tax identification number
   - Bank account information
   - Business license
5. Complete integration testing in sandbox
6. Pass VNPay's production readiness review
7. Receive production credentials:
   - **Production Merchant ID**
   - **Production Secret Key**

### Production Requirements

- ⚠️ **HTTPS Required**: All URLs must use HTTPS protocol
- ⚠️ **Valid SSL Certificate**: Domain must have valid SSL certificate
- ⚠️ **Security Validation**: System validates HTTPS in production mode
- ⚠️ **Business Verification**: Must complete VNPay's verification process
- ⚠️ **Contract Required**: Signed merchant agreement with VNPay

### Production Security Checks

The system automatically enforces security requirements in production:

```go
if c.IsProduction() {
    // Validate HTTPS for return URL
    if !strings.HasPrefix(c.ReturnURL, "https://") {
        return errors.New("security error: 'return_url' must use HTTPS in production")
    }
    
    // Validate HTTPS for IPN URL
    if !strings.HasPrefix(c.IPNURL, "https://") {
        return errors.New("security error: 'ipn_url' must use HTTPS in production")
    }
}
```

## Switching Between Environments

### Method 1: Environment Variable (Recommended)

Simply change the `VNPAY_ENVIRONMENT` variable:

```bash
# Switch to sandbox
export VNPAY_ENVIRONMENT=sandbox

# Switch to production
export VNPAY_ENVIRONMENT=production
```

### Method 2: Update .env File

Edit your `.env` file:

```bash
# For sandbox
VNPAY_ENVIRONMENT=sandbox
VNPAY_MERCHANT_ID=sandbox-merchant-id
VNPAY_SECRET_KEY=sandbox-secret-key
VNPAY_RETURN_URL=http://localhost:5173/payment/return
VNPAY_IPN_URL=http://localhost:8080/api/v1/payments/webhook

# For production
VNPAY_ENVIRONMENT=production
VNPAY_MERCHANT_ID=production-merchant-id
VNPAY_SECRET_KEY=production-secret-key
VNPAY_RETURN_URL=https://yourdomain.com/payment/return
VNPAY_IPN_URL=https://yourdomain.com/api/v1/payments/webhook
```

**Important**: When switching environments, you must also update:
- Merchant ID (different for sandbox vs production)
- Secret Key (different for sandbox vs production)
- URLs (HTTP for sandbox, HTTPS for production)

### Method 3: Separate Configuration Files

Maintain separate configuration files:

```bash
# .env.sandbox
VNPAY_ENVIRONMENT=sandbox
VNPAY_MERCHANT_ID=sandbox-merchant-id
VNPAY_SECRET_KEY=sandbox-secret-key
VNPAY_RETURN_URL=http://localhost:5173/payment/return
VNPAY_IPN_URL=http://localhost:8080/api/v1/payments/webhook

# .env.production
VNPAY_ENVIRONMENT=production
VNPAY_MERCHANT_ID=production-merchant-id
VNPAY_SECRET_KEY=production-secret-key
VNPAY_RETURN_URL=https://yourdomain.com/payment/return
VNPAY_IPN_URL=https://yourdomain.com/api/v1/payments/webhook
```

Load the appropriate file:

```bash
# Use sandbox
cp .env.sandbox .env

# Use production
cp .env.production .env
```

## Using Environment Methods in Code

### Check Current Environment

```go
config, err := payment.LoadConfigFromEnv()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

// Check if running in production
if config.IsProduction() {
    log.Println("Running in PRODUCTION mode - real payments enabled")
    // Enable production monitoring
    // Enable production logging
    // Enable production alerting
}

// Check if running in sandbox
if config.IsSandbox() {
    log.Println("Running in SANDBOX mode - test payments only")
    // Enable debug logging
    // Disable production alerts
}
```

### Get Environment-Specific Gateway URL

```go
config, err := payment.LoadConfigFromEnv()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

// Automatically returns correct URL based on environment
gatewayURL := config.GetBaseURL()

// Sandbox: https://sandbox.vnpayment.vn/paymentv2/vpcpay.html
// Production: https://vnpayment.vn/paymentv2/vpcpay.html
```

### Environment-Specific Logic

```go
config, err := payment.LoadConfigFromEnv()
if err != nil {
    log.Fatalf("Failed to load config: %v", err)
}

if config.IsProduction() {
    // Production-specific logic
    enableProductionMonitoring()
    enableSecurityAuditing()
    setStrictErrorHandling()
    enablePaymentAlerts()
} else {
    // Sandbox-specific logic
    enableDebugLogging()
    enableTestMode()
    relaxValidation()
}
```

## Environment Validation

### Automatic Validation

The system automatically validates configuration when loading:

```go
config, err := payment.LoadConfigFromEnv()
if err != nil {
    // Configuration validation failed
    log.Fatalf("Configuration error: %v", err)
}
```

### Validation Rules

#### Sandbox Environment
- ✅ Environment must be "sandbox" (case-insensitive)
- ✅ HTTP URLs allowed for development
- ✅ Localhost URLs allowed
- ✅ Merchant ID must be at least 4 characters
- ✅ Secret Key must be at least 16 characters

#### Production Environment
- ✅ Environment must be "production" (case-insensitive)
- ⚠️ **HTTPS Required**: Return URL must use HTTPS
- ⚠️ **HTTPS Required**: IPN URL must use HTTPS
- ✅ Merchant ID must be at least 4 characters
- ✅ Secret Key must be at least 16 characters
- ⚠️ Localhost URLs not recommended (but not blocked)

### Validation Error Messages

The system provides clear, actionable error messages:

```
❌ Invalid environment:
"configuration error: 'environment' must be 'sandbox' or 'production', got 'invalid'. 
Use 'sandbox' for testing or 'production' for live payments"

❌ Production HTTPS requirement:
"security error: 'return_url' must use HTTPS in production environment 
for secure payment processing"

❌ Missing configuration:
"configuration error: 'merchant_id' is required. 
Set VNPAY_MERCHANT_ID to your VNPay merchant/terminal code"
```

## Environment Differences Summary

| Feature | Sandbox | Production |
|---------|---------|------------|
| **Gateway URL** | `sandbox.vnpayment.vn` | `vnpayment.vn` |
| **Real Money** | ❌ No | ✅ Yes |
| **HTTPS Required** | ❌ No | ✅ Yes |
| **SSL Certificate** | ❌ Not required | ✅ Required |
| **Test Cards** | ✅ Available | ❌ Not available |
| **Credentials** | Sandbox credentials | Production credentials |
| **Localhost URLs** | ✅ Allowed | ⚠️ Not recommended |
| **HTTP URLs** | ✅ Allowed | ❌ Blocked |
| **Business Verification** | ❌ Not required | ✅ Required |
| **Contract** | ❌ Not required | ✅ Required |
| **Payment Speed** | Instant | 1-3 minutes |
| **Monitoring** | Optional | Recommended |
| **Alerting** | Optional | Required |

## Best Practices

### Development Workflow

1. **Start with Sandbox**
   - Develop and test all payment features in sandbox
   - Use sandbox credentials for local development
   - Test all payment scenarios (success, failure, timeout)

2. **Test Thoroughly**
   - Test payment initiation
   - Test payment success flow
   - Test payment failure handling
   - Test payment timeout scenarios
   - Test webhook processing
   - Test signature validation

3. **Prepare for Production**
   - Obtain production credentials from VNPay
   - Set up HTTPS with valid SSL certificate
   - Update environment variables
   - Configure production monitoring
   - Set up payment alerting

4. **Deploy to Production**
   - Switch to production environment
   - Verify configuration validation passes
   - Test with small transaction first
   - Monitor payment processing
   - Have rollback plan ready

### Security Best Practices

1. **Protect Credentials**
   ```bash
   # Never commit credentials to version control
   echo ".env" >> .gitignore
   
   # Use environment variables in production
   # Store secrets in secure secret management system
   ```

2. **Separate Environments**
   ```bash
   # Use different credentials for each environment
   # Never use production credentials in sandbox
   # Never use sandbox credentials in production
   ```

3. **Validate Environment**
   ```go
   // Always validate configuration on startup
   config, err := payment.LoadConfigFromEnv()
   if err != nil {
       log.Fatalf("Configuration validation failed: %v", err)
   }
   
   // Log environment (but not credentials)
   log.Printf("VNPay Environment: %s", config.Environment)
   log.Printf("VNPay Gateway: %s", config.GetBaseURL())
   ```

4. **Monitor Production**
   ```go
   if config.IsProduction() {
       // Enable production monitoring
       enablePaymentMetrics()
       enableSecurityAuditing()
       enableAlertingSystem()
   }
   ```

### Configuration Management

1. **Use Environment Variables**
   - Store configuration in environment variables
   - Never hardcode credentials in source code
   - Use `.env` files for local development
   - Use secret management systems for production

2. **Validate on Startup**
   - Load and validate configuration when application starts
   - Fail fast if configuration is invalid
   - Log configuration errors clearly

3. **Mask Sensitive Data**
   ```go
   // Use MaskSecretKey() for logging
   log.Printf("Secret Key: %s", config.MaskSecretKey())
   // Output: Secret Key: ABCD****WXYZ
   
   // Use PrettyPrint() for safe configuration display
   output, _ := config.PrettyPrint()
   log.Println(output)
   ```

## Troubleshooting

### Common Issues

#### Issue: "environment must be 'sandbox' or 'production'"
**Solution**: Check `VNPAY_ENVIRONMENT` variable spelling and value

```bash
# Correct values (case-insensitive)
VNPAY_ENVIRONMENT=sandbox
VNPAY_ENVIRONMENT=production
VNPAY_ENVIRONMENT=SANDBOX
VNPAY_ENVIRONMENT=Production

# Incorrect values
VNPAY_ENVIRONMENT=test      # ❌ Invalid
VNPAY_ENVIRONMENT=dev       # ❌ Invalid
VNPAY_ENVIRONMENT=staging   # ❌ Invalid
```

#### Issue: "return_url must use HTTPS in production"
**Solution**: Update URLs to use HTTPS when in production mode

```bash
# Sandbox - HTTP allowed
VNPAY_ENVIRONMENT=sandbox
VNPAY_RETURN_URL=http://localhost:5173/payment/return  # ✅ OK

# Production - HTTPS required
VNPAY_ENVIRONMENT=production
VNPAY_RETURN_URL=https://yourdomain.com/payment/return  # ✅ OK
VNPAY_RETURN_URL=http://yourdomain.com/payment/return   # ❌ Error
```

#### Issue: Wrong gateway URL being used
**Solution**: Verify environment variable is set correctly

```go
config, _ := payment.LoadConfigFromEnv()
log.Printf("Environment: %s", config.Environment)
log.Printf("Gateway URL: %s", config.GetBaseURL())

// Expected output for sandbox:
// Environment: sandbox
// Gateway URL: https://sandbox.vnpayment.vn/paymentv2/vpcpay.html

// Expected output for production:
// Environment: production
// Gateway URL: https://vnpayment.vn/paymentv2/vpcpay.html
```

#### Issue: Payment works in sandbox but fails in production
**Checklist**:
- ✅ Using production credentials (not sandbox credentials)
- ✅ URLs use HTTPS protocol
- ✅ SSL certificate is valid
- ✅ Domain is accessible from VNPay servers
- ✅ Webhook URL is publicly accessible
- ✅ Firewall allows VNPay IP addresses
- ✅ Production merchant account is active

## Testing Environment Switching

### Unit Tests

The system includes comprehensive tests for environment switching:

```bash
# Run all config tests
go test ./server/internal/payment -v -run TestVNPayConfig

# Test environment detection
go test ./server/internal/payment -v -run TestVNPayConfig_IsProduction
go test ./server/internal/payment -v -run TestVNPayConfig_IsSandbox

# Test gateway URL generation
go test ./server/internal/payment -v -run TestVNPayConfig_GetBaseURL
```

### Manual Testing

1. **Test Sandbox Configuration**
   ```bash
   # Set sandbox environment
   export VNPAY_ENVIRONMENT=sandbox
   export VNPAY_MERCHANT_ID=sandbox-merchant
   export VNPAY_SECRET_KEY=sandbox-secret-key-1234567890
   export VNPAY_RETURN_URL=http://localhost:5173/payment/return
   export VNPAY_IPN_URL=http://localhost:8080/api/v1/payments/webhook
   
   # Run application
   go run main.go
   
   # Verify sandbox gateway URL is used
   # Verify HTTP URLs are accepted
   ```

2. **Test Production Configuration**
   ```bash
   # Set production environment
   export VNPAY_ENVIRONMENT=production
   export VNPAY_MERCHANT_ID=production-merchant
   export VNPAY_SECRET_KEY=production-secret-key-1234567890
   export VNPAY_RETURN_URL=https://yourdomain.com/payment/return
   export VNPAY_IPN_URL=https://yourdomain.com/api/v1/payments/webhook
   
   # Run application
   go run main.go
   
   # Verify production gateway URL is used
   # Verify HTTPS validation is enforced
   ```

3. **Test Environment Switching**
   ```bash
   # Start with sandbox
   export VNPAY_ENVIRONMENT=sandbox
   go run main.go
   # Verify sandbox mode
   
   # Switch to production (restart required)
   export VNPAY_ENVIRONMENT=production
   go run main.go
   # Verify production mode
   ```

## Example Code

### Complete Example

```go
package main

import (
    "log"
    "your-project/server/internal/payment"
)

func main() {
    // Load configuration from environment variables
    config, err := payment.LoadConfigFromEnv()
    if err != nil {
        log.Fatalf("Failed to load VNPay configuration: %v", err)
    }

    // Log configuration (with masked secrets)
    log.Printf("VNPay Configuration:")
    log.Printf("  Environment: %s", config.Environment)
    log.Printf("  Merchant ID: %s", config.MerchantID)
    log.Printf("  Secret Key: %s", config.MaskSecretKey())
    log.Printf("  Gateway URL: %s", config.GetBaseURL())
    log.Printf("  Return URL: %s", config.ReturnURL)
    log.Printf("  IPN URL: %s", config.IPNURL)
    log.Printf("  Payment Timeout: %d minutes", config.PaymentTimeout)

    // Check environment and configure accordingly
    if config.IsProduction() {
        log.Println("🔴 PRODUCTION MODE - Real payments enabled")
        log.Println("⚠️  Ensure monitoring and alerting are configured")
        
        // Production-specific setup
        enableProductionMonitoring()
        enableSecurityAuditing()
        enablePaymentAlerts()
    } else if config.IsSandbox() {
        log.Println("🟡 SANDBOX MODE - Test payments only")
        log.Println("ℹ️  Using test credentials and sandbox gateway")
        
        // Sandbox-specific setup
        enableDebugLogging()
        enableTestMode()
    }

    // Use the configuration for payment processing
    // ... rest of your application code
}

func enableProductionMonitoring() {
    log.Println("✅ Production monitoring enabled")
}

func enableSecurityAuditing() {
    log.Println("✅ Security auditing enabled")
}

func enablePaymentAlerts() {
    log.Println("✅ Payment alerting enabled")
}

func enableDebugLogging() {
    log.Println("✅ Debug logging enabled")
}

func enableTestMode() {
    log.Println("✅ Test mode enabled")
}
```

## Summary

The VNPay payment integration provides complete support for environment switching:

✅ **Easy Configuration**: Simple environment variable controls environment
✅ **Automatic Gateway Selection**: `GetBaseURL()` returns correct URL
✅ **Environment Detection**: `IsProduction()` and `IsSandbox()` methods
✅ **Security Validation**: Automatic HTTPS enforcement in production
✅ **Clear Documentation**: Comprehensive guides and examples
✅ **Thorough Testing**: Unit tests cover all environment scenarios

To switch environments:
1. Update `VNPAY_ENVIRONMENT` variable
2. Update credentials (merchant ID and secret key)
3. Update URLs (HTTP for sandbox, HTTPS for production)
4. Restart application
5. Verify correct gateway URL is used

For questions or issues, refer to:
- VNPay Sandbox: https://sandbox.vnpayment.vn/devreg
- VNPay Documentation: https://vnpay.vn
- This guide: `server/internal/payment/ENVIRONMENT_SWITCHING.md`
