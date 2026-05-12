# Design Document: VNPay Payment Integration

## Overview

This design document outlines the technical implementation for integrating VNPay payment gateway into the existing travel booking system. The solution extends the current booking module with secure payment processing capabilities, maintaining the existing monolithic architecture while adding new payment-specific components.

## Architecture Overview

### System Integration Points

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   React Client  │    │   Go Backend     │    │     VNPay       │
│                 │    │                  │    │    Gateway      │
│  Payment UI     │◄──►│  Payment Module  │◄──►│                 │
│  Status Pages   │    │  Webhook Handler │    │  Sandbox/Prod   │
│  Error Handling │    │  Audit Logger    │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │   PostgreSQL     │
                       │                  │
                       │  Payment Tables  │
                       │  Audit Logs      │
                       └──────────────────┘
```

### Module Structure

Following the existing monolithic pattern:

```
server/internal/payment/
├── handler.go          # HTTP endpoints and request/response handling
├── service.go          # Business logic and VNPay integration
├── repository.go       # Database operations and transactions
├── vnpay_client.go     # VNPay API client and signature handling
├── config.go           # Configuration parsing and validation
└── audit.go            # Audit logging functionality

server/domain/
├── payment.go          # Payment entity models
├── payment_dto.go      # Request/response DTOs
└── payment_audit.go    # Audit log models
```

## Database Design

### New Tables

#### payments table
```sql
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL REFERENCES bookings(id),
    vnpay_transaction_id VARCHAR(100) UNIQUE,
    transaction_reference VARCHAR(50) UNIQUE NOT NULL,
    amount DECIMAL(12,2) NOT NULL,
    currency VARCHAR(3) DEFAULT 'VND',
    payment_method VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    vnpay_response_code VARCHAR(10),
    vnpay_message TEXT,
    session_expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_payments_booking_id ON payments(booking_id);
CREATE INDEX idx_payments_transaction_ref ON payments(transaction_reference);
CREATE INDEX idx_payments_vnpay_txn_id ON payments(vnpay_transaction_id);
CREATE INDEX idx_payments_status ON payments(status);
```

#### payment_audit_logs table
```sql
CREATE TABLE payment_audit_logs (
    id SERIAL PRIMARY KEY,
    payment_id INTEGER REFERENCES payments(id),
    booking_id INTEGER REFERENCES bookings(id),
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB,
    user_id INTEGER REFERENCES users(id),
    ip_address INET,
    user_agent TEXT,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_payment_id ON payment_audit_logs(payment_id);
CREATE INDEX idx_audit_event_type ON payment_audit_logs(event_type);
CREATE INDEX idx_audit_timestamp ON payment_audit_logs(timestamp);
```

### Modified Tables

#### bookings table (add payment tracking)
```sql
ALTER TABLE bookings 
ADD COLUMN payment_id INTEGER REFERENCES payments(id),
ADD COLUMN payment_deadline TIMESTAMP;

-- Update existing payment_status enum if needed
ALTER TABLE bookings 
ALTER COLUMN payment_status TYPE VARCHAR(20);
```

## API Design

### New Endpoints

#### POST /v1/api/payments/initiate
**Purpose:** Create payment session and redirect to VNPay

**Request:**
```json
{
  "booking_id": 123
}
```

**Response:**
```json
{
  "success": true,
  "payment_url": "https://sandbox.vnpayment.vn/paymentv2/vpcpay.html?...",
  "transaction_reference": "PAY20240430001",
  "expires_at": "2024-04-30T16:00:00Z"
}
```

#### GET /v1/api/payments/return
**Purpose:** Handle VNPay return redirect

**Query Parameters:** VNPay standard return parameters
**Response:** Redirect to frontend success/failure page

#### POST /v1/api/payments/webhook
**Purpose:** Process VNPay IPN notifications

**Request:** VNPay IPN payload with signature
**Response:** HTTP 200 OK or error status

#### GET /v1/api/payments/status/:transaction_reference
**Purpose:** Check payment status

**Response:**
```json
{
  "success": true,
  "payment": {
    "transaction_reference": "PAY20240430001",
    "status": "paid",
    "amount": 2500000,
    "booking_id": 123,
    "created_at": "2024-04-30T15:30:00Z"
  }
}
```

#### POST /v1/api/payments/refund
**Purpose:** Process payment refunds (admin only)

**Request:**
```json
{
  "payment_id": 456,
  "refund_amount": 2500000,
  "reason": "Customer cancellation"
}
```

## VNPay Integration

### Configuration Management

```go
type VNPayConfig struct {
    Environment    string `json:"environment"`     // sandbox/production
    MerchantID     string `json:"merchant_id"`
    SecretKey      string `json:"secret_key"`
    ReturnURL      string `json:"return_url"`
    IPNURL         string `json:"ipn_url"`
    PaymentTimeout int    `json:"payment_timeout"` // minutes
}

func (c *VNPayConfig) Validate() error {
    // Validation logic for required fields and URL formats
}
```

### Signature Generation and Validation

```go
func GenerateSignature(params map[string]string, secretKey string) string {
    // Sort parameters, create query string, HMAC-SHA512 hash
}

func ValidateSignature(params map[string]string, signature, secretKey string) bool {
    // Regenerate signature and compare
}
```

### Payment URL Generation

```go
func (c *VNPayClient) CreatePaymentURL(payment *Payment) (string, error) {
    params := map[string]string{
        "vnp_Version":    "2.1.0",
        "vnp_Command":    "pay",
        "vnp_TmnCode":    c.config.MerchantID,
        "vnp_Amount":     strconv.FormatInt(payment.Amount*100, 10), // VND cents
        "vnp_CurrCode":   "VND",
        "vnp_TxnRef":     payment.TransactionReference,
        "vnp_OrderInfo":  fmt.Sprintf("Thanh toan tour %s", payment.BookingCode),
        "vnp_ReturnUrl":  c.config.ReturnURL,
        "vnp_IpAddr":     payment.ClientIP,
        "vnp_CreateDate": time.Now().Format("20060102150405"),
        "vnp_ExpireDate": payment.ExpiresAt.Format("20060102150405"),
    }
    
    signature := GenerateSignature(params, c.config.SecretKey)
    params["vnp_SecureHash"] = signature
    
    return c.config.BaseURL + "?" + buildQueryString(params), nil
}
```

## Business Logic Implementation

### Payment Flow State Machine

```
┌─────────┐    initiate    ┌────────────┐    vnpay_redirect    ┌─────────┐
│ pending │──────────────► │ processing │────────────────────► │ success │
└─────────┘                └────────────┘                      └─────────┘
     │                           │                                   │
     │ timeout/expire            │ failure/cancel                    │
     ▼                           ▼                                   ▼
┌─────────┐                ┌─────────┐                        ┌───────────┐
│ expired │                │ failed  │                        │ confirmed │
└─────────┘                └─────────┘                        └───────────┘
                                 │                                   │
                                 │ retry                             │ refund
                                 ▼                                   ▼
                           ┌────────────┐                      ┌──────────┐
                           │ processing │                      │ refunded │
                           └────────────┘                      └──────────┘
```

### Service Layer Methods

```go
type PaymentService struct {
    repo      PaymentRepository
    vnpay     VNPayClient
    audit     AuditLogger
    booking   BookingService
}

func (s *PaymentService) InitiatePayment(userID, bookingID uint) (*PaymentResponse, error) {
    // 1. Validate booking exists and belongs to user
    // 2. Check booking is in "booked" status and "unpaid"
    // 3. Create payment record with unique transaction reference
    // 4. Generate VNPay payment URL
    // 5. Log payment initiation
    // 6. Return payment URL and reference
}

func (s *PaymentService) ProcessWebhook(params map[string]string) error {
    // 1. Validate signature
    // 2. Extract transaction reference and status
    // 3. Update payment status in transaction
    // 4. Update booking status if payment successful
    // 5. Send confirmation email
    // 6. Log webhook processing
}

func (s *PaymentService) HandleReturn(params map[string]string) (*PaymentStatus, error) {
    // 1. Validate signature
    // 2. Extract payment result
    // 3. Return status for frontend display
    // Note: Don't update status here, rely on webhook for authoritative updates
}
```

## Security Implementation

### Rate Limiting

```go
func PaymentRateLimit() gin.HandlerFunc {
    limiter := rate.NewLimiter(rate.Every(12*time.Second), 5) // 5 requests per minute
    return gin_rate_limit.RateLimiter(limiter, &gin_rate_limit.Options{
        ErrorHandler: func(c *gin.Context, info gin_rate_limit.Info) {
            c.JSON(429, gin.H{"error": "Too many payment requests"})
        },
    })
}
```

### Input Validation and Sanitization

```go
func ValidatePaymentRequest(req *InitiatePaymentRequest) error {
    if req.BookingID <= 0 {
        return errors.New("invalid booking ID")
    }
    // Additional validation
}

func SanitizeVNPayParams(params map[string]string) map[string]string {
    sanitized := make(map[string]string)
    for k, v := range params {
        // Remove potentially dangerous characters
        sanitized[k] = html.EscapeString(strings.TrimSpace(v))
    }
    return sanitized
}
```

### CSRF Protection

```go
func CSRFProtection() gin.HandlerFunc {
    return csrf.Middleware(csrf.Options{
        Secret: os.Getenv("CSRF_SECRET"),
        ErrorFunc: func(c *gin.Context) {
            c.JSON(403, gin.H{"error": "CSRF token invalid"})
        },
    })
}
```

## Frontend Implementation

### Payment Components

#### PaymentButton Component
```jsx
const PaymentButton = ({ booking, onPaymentStart, onPaymentComplete }) => {
  const [loading, setLoading] = useState(false);
  
  const handlePayment = async () => {
    setLoading(true);
    try {
      const response = await paymentService.initiatePayment(booking.id);
      onPaymentStart();
      window.location.href = response.payment_url;
    } catch (error) {
      setLoading(false);
      // Handle error
    }
  };

  return (
    <Button 
      onClick={handlePayment} 
      loading={loading}
      disabled={booking.payment_status !== 'unpaid'}
    >
      Thanh toán ngay - {formatCurrency(booking.total_amount)}
    </Button>
  );
};
```

#### PaymentStatus Component
```jsx
const PaymentStatus = ({ transactionReference }) => {
  const [status, setStatus] = useState('checking');
  
  useEffect(() => {
    const checkStatus = async () => {
      try {
        const response = await paymentService.getPaymentStatus(transactionReference);
        setStatus(response.payment.status);
      } catch (error) {
        setStatus('error');
      }
    };
    
    const interval = setInterval(checkStatus, 3000);
    return () => clearInterval(interval);
  }, [transactionReference]);

  return (
    <div className="payment-status">
      {status === 'checking' && <Spinner />}
      {status === 'paid' && <SuccessMessage />}
      {status === 'failed' && <ErrorMessage />}
    </div>
  );
};
```

### Payment Service

```javascript
class PaymentService {
  async initiatePayment(bookingId) {
    const response = await httpClient.post('/v1/api/payments/initiate', {
      booking_id: bookingId
    });
    return response.data;
  }

  async getPaymentStatus(transactionReference) {
    const response = await httpClient.get(`/v1/api/payments/status/${transactionReference}`);
    return response.data;
  }

  async retryPayment(paymentId) {
    const response = await httpClient.post(`/v1/api/payments/${paymentId}/retry`);
    return response.data;
  }
}
```

## Error Handling Strategy

### Error Categories and Responses

1. **Validation Errors (400)**
   - Invalid booking ID
   - Booking already paid
   - Invalid payment amount

2. **Authentication Errors (401)**
   - Missing or invalid JWT token
   - User not authorized for booking

3. **Business Logic Errors (409)**
   - Booking expired
   - Insufficient tour slots
   - Payment already processed

4. **External Service Errors (502/503)**
   - VNPay service unavailable
   - Network timeout
   - Invalid VNPay response

5. **Security Errors (403)**
   - Invalid signature
   - Rate limit exceeded
   - CSRF token invalid

### Error Recovery Mechanisms

```go
func (s *PaymentService) InitiatePaymentWithRetry(userID, bookingID uint) (*PaymentResponse, error) {
    var lastErr error
    
    for attempt := 1; attempt <= 3; attempt++ {
        result, err := s.InitiatePayment(userID, bookingID)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        if !isRetryableError(err) {
            break
        }
        
        time.Sleep(time.Duration(attempt) * time.Second)
    }
    
    return nil, fmt.Errorf("payment initiation failed after 3 attempts: %w", lastErr)
}
```

## Audit and Logging

### Audit Event Types

```go
const (
    EventPaymentInitiated   = "payment_initiated"
    EventPaymentCompleted   = "payment_completed"
    EventPaymentFailed      = "payment_failed"
    EventWebhookReceived    = "webhook_received"
    EventRefundProcessed    = "refund_processed"
    EventSecurityViolation  = "security_violation"
)
```

### Audit Logger Implementation

```go
type AuditLogger struct {
    repo AuditRepository
}

func (a *AuditLogger) LogPaymentEvent(eventType string, paymentID, userID uint, data interface{}) {
    log := &PaymentAuditLog{
        PaymentID: paymentID,
        UserID:    userID,
        EventType: eventType,
        EventData: data,
        Timestamp: time.Now(),
        IPAddress: getCurrentIP(),
        UserAgent: getCurrentUserAgent(),
    }
    
    a.repo.CreateAuditLog(log)
}
```

## Testing Strategy

### Unit Tests

Property-based tests for critical components:

```go
func TestPaymentSignatureValidation(t *testing.T) {
    quick.Check(func(params map[string]string, secretKey string) bool {
        // Generate signature
        signature := GenerateSignature(params, secretKey)
        
        // Validate signature
        return ValidateSignature(params, signature, secretKey)
    }, nil)
}

func TestConfigurationRoundTrip(t *testing.T) {
    quick.Check(func(config VNPayConfig) bool {
        // Serialize config
        data, err := json.Marshal(config)
        if err != nil {
            return false
        }
        
        // Parse config
        var parsed VNPayConfig
        err = json.Unmarshal(data, &parsed)
        if err != nil {
            return false
        }
        
        // Compare configs
        return reflect.DeepEqual(config, parsed)
    }, nil)
}
```

### Integration Tests

```go
func TestVNPaySandboxIntegration(t *testing.T) {
    // Test with VNPay sandbox environment
    client := NewVNPayClient(sandboxConfig)
    
    payment := &Payment{
        Amount: 100000,
        TransactionReference: "TEST001",
    }
    
    url, err := client.CreatePaymentURL(payment)
    assert.NoError(t, err)
    assert.Contains(t, url, "sandbox.vnpayment.vn")
}
```

### End-to-End Tests

```javascript
describe('Payment Flow', () => {
  it('should complete payment successfully', async () => {
    // 1. Create booking
    const booking = await createTestBooking();
    
    // 2. Initiate payment
    const paymentResponse = await paymentService.initiatePayment(booking.id);
    expect(paymentResponse.payment_url).toContain('vnpayment.vn');
    
    // 3. Simulate VNPay callback
    await simulateVNPayCallback(paymentResponse.transaction_reference, 'success');
    
    // 4. Verify booking status updated
    const updatedBooking = await getBooking(booking.id);
    expect(updatedBooking.payment_status).toBe('paid');
    expect(updatedBooking.status).toBe('confirmed');
  });
});
```

## Deployment Considerations

### Environment Configuration

```yaml
# .env.production
VNPAY_ENVIRONMENT=production
VNPAY_MERCHANT_ID=${VNPAY_MERCHANT_ID}
VNPAY_SECRET_KEY=${VNPAY_SECRET_KEY}
VNPAY_RETURN_URL=https://yourdomain.com/payment/return
VNPAY_IPN_URL=https://yourdomain.com/api/v1/payments/webhook

# .env.staging
VNPAY_ENVIRONMENT=sandbox
VNPAY_MERCHANT_ID=sandbox_merchant
VNPAY_SECRET_KEY=sandbox_secret
VNPAY_RETURN_URL=https://staging.yourdomain.com/payment/return
VNPAY_IPN_URL=https://staging.yourdomain.com/api/v1/payments/webhook
```

### Database Migration

```sql
-- Migration: 001_add_payment_tables.up.sql
CREATE TABLE payments (...);
CREATE TABLE payment_audit_logs (...);
ALTER TABLE bookings ADD COLUMN payment_id INTEGER REFERENCES payments(id);

-- Migration: 001_add_payment_tables.down.sql
ALTER TABLE bookings DROP COLUMN payment_id;
DROP TABLE payment_audit_logs;
DROP TABLE payments;
```

### Monitoring and Alerting

```go
// Metrics collection
var (
    paymentInitiations = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "payment_initiations_total",
            Help: "Total number of payment initiations",
        },
        []string{"status"},
    )
    
    paymentDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "payment_processing_duration_seconds",
            Help: "Payment processing duration",
        },
        []string{"operation"},
    )
)
```

## Performance Considerations

### Database Optimization

1. **Indexing Strategy**
   - Index on booking_id for fast payment lookups
   - Index on transaction_reference for webhook processing
   - Composite index on (status, created_at) for reporting

2. **Connection Pooling**
   ```go
   db.SetMaxOpenConns(25)
   db.SetMaxIdleConns(5)
   db.SetConnMaxLifetime(5 * time.Minute)
   ```

3. **Query Optimization**
   - Use prepared statements for frequent queries
   - Implement pagination for payment history
   - Use database transactions for atomic updates

### Caching Strategy

```go
// Cache payment status for frequent checks
func (s *PaymentService) GetPaymentStatus(txnRef string) (*Payment, error) {
    // Check cache first
    if cached := s.cache.Get("payment:" + txnRef); cached != nil {
        return cached.(*Payment), nil
    }
    
    // Fetch from database
    payment, err := s.repo.GetByTransactionReference(txnRef)
    if err != nil {
        return nil, err
    }
    
    // Cache for 5 minutes
    s.cache.Set("payment:"+txnRef, payment, 5*time.Minute)
    return payment, nil
}
```

## Security Checklist

- [ ] All VNPay communications use HTTPS
- [ ] Payment signatures validated using HMAC-SHA512
- [ ] Rate limiting implemented on payment endpoints
- [ ] Input validation and sanitization applied
- [ ] CSRF protection enabled for payment initiation
- [ ] Sensitive data (secret keys) stored in environment variables
- [ ] Payment amounts validated against booking totals
- [ ] Audit logging captures all payment events
- [ ] Database queries use parameterized statements
- [ ] Error messages don't expose sensitive information

## Correctness Properties

Based on the prework analysis, here are the key correctness properties to implement:

### Property 1: Configuration Round-Trip Integrity
```go
// For all valid VNPay configurations, parsing then serializing then parsing 
// should produce an equivalent configuration
func TestConfigurationRoundTrip(t *testing.T) {
    quick.Check(func(config VNPayConfig) bool {
        if !config.IsValid() {
            return true // Skip invalid configs
        }
        
        serialized := config.Serialize()
        parsed, err := ParseVNPayConfig(serialized)
        
        return err == nil && config.Equals(parsed)
    }, nil)
}
```

### Property 2: Payment Signature Validation Invariant
```go
// For all valid parameter sets and secret keys, generating a signature 
// and then validating it should always return true
func TestSignatureValidationInvariant(t *testing.T) {
    quick.Check(func(params map[string]string, secretKey string) bool {
        if secretKey == "" {
            return true // Skip empty secret keys
        }
        
        signature := GenerateSignature(params, secretKey)
        return ValidateSignature(params, signature, secretKey)
    }, nil)
}
```

### Property 3: Payment State Machine Consistency
```go
// Payment status transitions should always be valid according to business rules
func TestPaymentStateTransitions(t *testing.T) {
    quick.Check(func(initialStatus, newStatus string) bool {
        return IsValidStatusTransition(initialStatus, newStatus)
    }, nil)
}
```

### Property 4: Audit Log Completeness
```go
// Every payment operation should generate corresponding audit log entries
func TestAuditLogCompleteness(t *testing.T) {
    quick.Check(func(operation PaymentOperation) bool {
        initialLogCount := getAuditLogCount()
        
        executePaymentOperation(operation)
        
        finalLogCount := getAuditLogCount()
        return finalLogCount > initialLogCount
    }, nil)
}
```

### Property 5: Webhook Idempotency
```go
// Processing the same webhook multiple times should not change the final state
func TestWebhookIdempotency(t *testing.T) {
    quick.Check(func(webhookPayload VNPayWebhook) bool {
        if !webhookPayload.IsValid() {
            return true
        }
        
        // Process webhook first time
        result1 := processWebhook(webhookPayload)
        state1 := getPaymentState(webhookPayload.TransactionRef)
        
        // Process same webhook again
        result2 := processWebhook(webhookPayload)
        state2 := getPaymentState(webhookPayload.TransactionRef)
        
        return result1.Equals(result2) && state1.Equals(state2)
    }, nil)
}
```

This design provides a comprehensive technical foundation for implementing the VNPay payment integration while maintaining security, reliability, and testability standards.