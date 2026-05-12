package payment

// Example usage of LoadConfigFromEnv:
//
// func main() {
//     // Load VNPay configuration from environment variables
//     config, err := payment.LoadConfigFromEnv()
//     if err != nil {
//         log.Fatalf("Failed to load VNPay configuration: %v", err)
//     }
//
//     // Use the configuration
//     log.Printf("VNPay Environment: %s", config.Environment)
//     log.Printf("VNPay Merchant ID: %s", config.MerchantID)
//     log.Printf("VNPay Secret Key: %s", config.MaskSecretKey())
//     log.Printf("VNPay Return URL: %s", config.ReturnURL)
//     log.Printf("VNPay IPN URL: %s", config.IPNURL)
//     log.Printf("VNPay Payment Timeout: %d minutes", config.PaymentTimeout)
//     log.Printf("VNPay Base URL: %s", config.GetBaseURL())
//
//     // Check environment
//     if config.IsProduction() {
//         log.Println("Running in PRODUCTION mode")
//     } else if config.IsSandbox() {
//         log.Println("Running in SANDBOX mode")
//     }
// }
//
// Required environment variables:
// - VNPAY_ENVIRONMENT: sandbox or production (default: sandbox)
// - VNPAY_MERCHANT_ID: VNPay merchant/terminal code (required)
// - VNPAY_SECRET_KEY: Secret key for signature generation (required)
// - VNPAY_RETURN_URL: User-facing URL for payment return (required)
// - VNPAY_IPN_URL: Server-to-server webhook URL (required)
// - VNPAY_PAYMENT_TIMEOUT: Payment session timeout in minutes (default: 15)
