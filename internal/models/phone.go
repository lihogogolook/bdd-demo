package models

// PhoneRiskRequest 電話號碼風險評估請求
type PhoneRiskRequest struct {
	PhoneNumber string `json:"phone_number" binding:"required" example:"0987654321"`
}

// PhoneRiskResponse 電話號碼風險評估回應
type PhoneRiskResponse struct {
	PhoneNumber string  `json:"phone_number" example:"0987654321"`
	RiskScore   float64 `json:"risk_score" example:"75.5"`
	RiskLevel   string  `json:"risk_level" example:"medium"`
	Message     string  `json:"message" example:"電話號碼風險評估完成"`
}

// ErrorResponse 錯誤回應
type ErrorResponse struct {
	Error   string `json:"error" example:"無效的電話號碼格式"`
	Message string `json:"message" example:"請提供有效的台灣電話號碼"`
}
