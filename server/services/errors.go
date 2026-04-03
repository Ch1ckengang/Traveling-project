package services

import "errors"

var (
	ErrInvalidAuthPayload     = errors.New("Dữ liệu đăng nhập không hợp lệ")
	ErrInvalidCredentials     = errors.New("Email hoặc mật khẩu không đúng")
	ErrLoginProcessFailed     = errors.New("Không thể xử lý đăng nhập")
	ErrInvalidName            = errors.New("Họ tên không được để trống")
	ErrInvalidRegisterEmail   = errors.New("Email không hợp lệ")
	ErrWeakPassword           = errors.New("Mật khẩu phải có ít nhất 8 ký tự")
	ErrEmailAlreadyRegistered = errors.New("Email đã được đăng ký")
	ErrEmailCheckFailed       = errors.New("Không thể xác thực email")
	ErrPasswordHashFailed     = errors.New("Không thể bảo mật mật khẩu")
	ErrCreateAccountFailed    = errors.New("Không thể tạo tài khoản")

	ErrInvalidBookingPayload = errors.New("Dữ liệu đặt tour không hợp lệ")
	ErrInvalidFullName       = errors.New("Họ tên không được để trống")
	ErrInvalidPhone          = errors.New("Số điện thoại không hợp lệ")
	ErrInvalidEmail          = errors.New("Email không hợp lệ")
	ErrInvalidQuantity       = errors.New("Số lượng khách phải lớn hơn 0")
	ErrInvalidTravelDate     = errors.New("Ngày đi không hợp lệ")
	ErrTravelDateInPast      = errors.New("Ngày đi không được ở quá khứ")
	ErrTourNotFound          = errors.New("Tour không tồn tại")
	ErrCreateBookingFailed   = errors.New("Không thể tạo đơn đặt tour")
)
