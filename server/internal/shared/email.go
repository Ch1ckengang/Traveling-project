package shared

import (
	"fmt"
	"log"
	"net/smtp"
	"os"
)

// EmailConfig - Cấu hình email service
type EmailConfig struct {
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	Enabled      bool
}

var emailConfig *EmailConfig

// InitEmailService - Khởi tạo email service từ environment variables
func InitEmailService() {
	emailConfig = &EmailConfig{
		SMTPHost:     os.Getenv("SMTP_HOST"),
		SMTPPort:     os.Getenv("SMTP_PORT"),
		SMTPUsername: os.Getenv("SMTP_USERNAME"),
		SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail:    os.Getenv("SMTP_FROM_EMAIL"),
		FromName:     os.Getenv("SMTP_FROM_NAME"),
		Enabled:      os.Getenv("EMAIL_ENABLED") == "true",
	}

	if !emailConfig.Enabled {
		log.Println("📧 Email service DISABLED - running in dev mode")
		return
	}

	if emailConfig.SMTPHost == "" || emailConfig.SMTPPort == "" {
		log.Println("⚠️ Email service configured but missing SMTP credentials")
		emailConfig.Enabled = false
		return
	}

	log.Printf("✅ Email service ENABLED (%s:%s)", emailConfig.SMTPHost, emailConfig.SMTPPort)
}

// SendEmail - Gửi email với SMTP
func SendEmail(to, subject, body string) error {
	if emailConfig == nil {
		InitEmailService()
	}

	// Dev mode: chỉ log ra console
	if !emailConfig.Enabled {
		log.Printf("📧 [DEV MODE] Email to %s\nSubject: %s\nBody: %s", to, subject, body)
		return nil
	}

	// Production mode: gửi email thật
	from := emailConfig.FromEmail
	if emailConfig.FromName != "" {
		from = fmt.Sprintf("%s <%s>", emailConfig.FromName, emailConfig.FromEmail)
	}

	msg := []byte(fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s\r\n",
		from, to, subject, body,
	))

	auth := smtp.PlainAuth("", emailConfig.SMTPUsername, emailConfig.SMTPPassword, emailConfig.SMTPHost)
	addr := fmt.Sprintf("%s:%s", emailConfig.SMTPHost, emailConfig.SMTPPort)

	err := smtp.SendMail(addr, auth, emailConfig.FromEmail, []string{to}, msg)
	if err != nil {
		log.Printf("❌ Failed to send email to %s: %v", to, err)
		return fmt.Errorf("không thể gửi email: %w", err)
	}

	log.Printf("✅ Email sent successfully to %s", to)
	return nil
}

// SendOTPEmail - Gửi email chứa mã OTP
func SendOTPEmail(to, code string) error {
	subject := "Mã xác thực OTP - Traveling"
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 10px;">
				<h2 style="color: #2563eb;">Xác thực tài khoản Traveling</h2>
				<p>Xin chào,</p>
				<p>Mã OTP của bạn là:</p>
				<div style="background-color: #f3f4f6; padding: 20px; text-align: center; border-radius: 5px; margin: 20px 0;">
					<h1 style="color: #2563eb; font-size: 36px; margin: 0; letter-spacing: 5px;">%s</h1>
				</div>
				<p>Mã này sẽ hết hạn sau <strong>3 phút</strong>.</p>
				<p>Nếu bạn không yêu cầu mã này, vui lòng bỏ qua email này.</p>
				<hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">
				<p style="font-size: 12px; color: #666;">
					Email này được gửi tự động từ hệ thống Traveling.<br>
					Vui lòng không trả lời email này.
				</p>
			</div>
		</body>
		</html>
	`, code)

	return SendEmail(to, subject, body)
}

// SendPasswordResetEmail - Gửi email reset mật khẩu
func SendPasswordResetEmail(to, resetToken string) error {
	subject := "Đặt lại mật khẩu - Traveling"
	resetURL := fmt.Sprintf("http://localhost:5173/reset-password?token=%s", resetToken)

	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 10px;">
				<h2 style="color: #2563eb;">Đặt lại mật khẩu</h2>
				<p>Xin chào,</p>
				<p>Bạn đã yêu cầu đặt lại mật khẩu cho tài khoản Traveling của mình.</p>
				<p>Nhấn vào nút bên dưới để đặt lại mật khẩu:</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="%s" style="background-color: #2563eb; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">
						Đặt lại mật khẩu
					</a>
				</div>
				<p>Hoặc copy link sau vào trình duyệt:</p>
				<p style="background-color: #f3f4f6; padding: 10px; border-radius: 5px; word-break: break-all;">
					%s
				</p>
				<p>Link này sẽ hết hạn sau <strong>1 giờ</strong>.</p>
				<p>Nếu bạn không yêu cầu đặt lại mật khẩu, vui lòng bỏ qua email này.</p>
				<hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">
				<p style="font-size: 12px; color: #666;">
					Email này được gửi tự động từ hệ thống Traveling.<br>
					Vui lòng không trả lời email này.
				</p>
			</div>
		</body>
		</html>
	`, resetURL, resetURL)

	return SendEmail(to, subject, body)
}

// SendBookingConfirmationEmail - Gửi email xác nhận booking
func SendBookingConfirmationEmail(to, bookingCode, tourName, travelDate string) error {
	subject := fmt.Sprintf("Xác nhận đặt tour #%s - Traveling", bookingCode)
	body := fmt.Sprintf(`
		<html>
		<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333;">
			<div style="max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 10px;">
				<h2 style="color: #10b981;">Đặt tour thành công! 🎉</h2>
				<p>Xin chào,</p>
				<p>Cảm ơn bạn đã đặt tour tại Traveling. Đơn đặt tour của bạn đã được xác nhận.</p>
				<div style="background-color: #f3f4f6; padding: 20px; border-radius: 5px; margin: 20px 0;">
					<p style="margin: 5px 0;"><strong>Mã đặt tour:</strong> %s</p>
					<p style="margin: 5px 0;"><strong>Tour:</strong> %s</p>
					<p style="margin: 5px 0;"><strong>Ngày khởi hành:</strong> %s</p>
				</div>
				<p>Vui lòng thanh toán trong vòng <strong>24 giờ</strong> để giữ chỗ.</p>
				<div style="text-align: center; margin: 30px 0;">
					<a href="http://localhost:5173/account/bookings" style="background-color: #2563eb; color: white; padding: 12px 30px; text-decoration: none; border-radius: 5px; display: inline-block;">
						Xem chi tiết đơn hàng
					</a>
				</div>
				<p>Nếu có bất kỳ thắc mắc nào, vui lòng liên hệ với chúng tôi.</p>
				<hr style="border: none; border-top: 1px solid #ddd; margin: 20px 0;">
				<p style="font-size: 12px; color: #666;">
					Email này được gửi tự động từ hệ thống Traveling.<br>
					Vui lòng không trả lời email này.
				</p>
			</div>
		</body>
		</html>
	`, bookingCode, tourName, travelDate)

	return SendEmail(to, subject, body)
}
