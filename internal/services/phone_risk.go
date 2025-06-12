package services

import (
	"crypto/md5"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"bdd-demo/internal/models"
)

// PhoneRiskService 電話號碼風險評估服務
type PhoneRiskService struct{}

// NewPhoneRiskService 建立新的電話號碼風險評估服務
func NewPhoneRiskService() *PhoneRiskService {
	return &PhoneRiskService{}
}

// ValidatePhoneNumber 驗證電話號碼格式
func (s *PhoneRiskService) ValidatePhoneNumber(phoneNumber string) bool {
	// 移除所有空格和破折號
	cleaned := strings.ReplaceAll(strings.ReplaceAll(phoneNumber, " ", ""), "-", "")

	// 台灣手機號碼格式：09XXXXXXXX (10位數字)
	mobilePattern := `^09\d{8}$`

	// 台灣市話格式：0X-XXXXXXXX 或 0XX-XXXXXXX
	landlinePattern := `^0\d{1,2}\d{7,8}$`

	mobileRegex := regexp.MustCompile(mobilePattern)
	landlineRegex := regexp.MustCompile(landlinePattern)

	return mobileRegex.MatchString(cleaned) || landlineRegex.MatchString(cleaned)
}

// CalculateRiskScore 計算風險分數
func (s *PhoneRiskService) CalculateRiskScore(phoneNumber string) float64 {
	// 使用電話號碼生成一個偽隨機但一致的風險分數
	hash := md5.Sum([]byte(phoneNumber))
	hashString := fmt.Sprintf("%x", hash)

	// 取前8位並轉換為數字
	hashInt, _ := strconv.ParseInt(hashString[:8], 16, 64)

	// 計算0-100的風險分數
	riskScore := math.Abs(float64(hashInt%10000)) / 100.0

	// 加入一些基於號碼特徵的調整
	cleaned := strings.ReplaceAll(strings.ReplaceAll(phoneNumber, " ", ""), "-", "")

	// 特殊號碼模式調整
	if strings.Contains(cleaned, "1234") || strings.Contains(cleaned, "5678") {
		riskScore += 15.0 // 連續數字被認為風險較高
	}

	if strings.Count(cleaned, cleaned[2:3]) >= 4 {
		riskScore += 10.0 // 重複數字較多
	}

	// 確保分數在0-100範圍內
	if riskScore > 100 {
		riskScore = 100
	}

	return math.Round(riskScore*10) / 10 // 四捨五入到小數點後一位
}

// GetRiskLevel 根據風險分數取得風險等級
func (s *PhoneRiskService) GetRiskLevel(riskScore float64) string {
	switch {
	case riskScore < 30:
		return "low"
	case riskScore < 70:
		return "medium"
	default:
		return "high"
	}
}

// EvaluatePhoneRisk 評估電話號碼風險
func (s *PhoneRiskService) EvaluatePhoneRisk(phoneNumber string) (*models.PhoneRiskResponse, error) {
	if !s.ValidatePhoneNumber(phoneNumber) {
		return nil, fmt.Errorf("無效的電話號碼格式")
	}

	riskScore := s.CalculateRiskScore(phoneNumber)
	riskLevel := s.GetRiskLevel(riskScore)

	return &models.PhoneRiskResponse{
		PhoneNumber: phoneNumber,
		RiskScore:   riskScore,
		RiskLevel:   riskLevel,
		Message:     "電話號碼風險評估完成",
	}, nil
}
