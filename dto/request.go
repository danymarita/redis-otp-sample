package dto

type (
	OtpRequest struct {
		UserID uint `json:"user_id"`
	}

	OtpValidateRequest struct {
		UserID  uint `json:"user_id"`
		OtpCode int  `json:"otp_code"`
	}
)
