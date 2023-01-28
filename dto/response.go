package dto

type (
	OtpResponse struct {
		UserID  uint `json:"user_id"`
		OtpCode int  `json:"otp_code"`
	}
)
