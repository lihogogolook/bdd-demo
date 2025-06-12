package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"bdd-demo/internal/models"
	"bdd-demo/internal/services"
)

// PhoneHandler 電話號碼相關的 API 處理器
type PhoneHandler struct {
	phoneRiskService *services.PhoneRiskService
}

// NewPhoneHandler 建立新的電話號碼處理器
func NewPhoneHandler(phoneRiskService *services.PhoneRiskService) *PhoneHandler {
	return &PhoneHandler{
		phoneRiskService: phoneRiskService,
	}
}

// EvaluatePhoneRisk 評估電話號碼風險
// @Summary 評估電話號碼風險
// @Description 輸入電話號碼，回傳風險分數和風險等級
// @Tags phone
// @Accept json
// @Produce json
// @Param request body models.PhoneRiskRequest true "電話號碼風險評估請求"
// @Success 200 {object} models.PhoneRiskResponse
// @Failure 400 {object} models.ErrorResponse
// @Router /api/phone/risk [post]
func (h *PhoneHandler) EvaluatePhoneRisk(c *gin.Context) {
	var req models.PhoneRiskRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "請求格式錯誤",
			Message: "請提供有效的 JSON 格式和必要的 phone_number 欄位",
		})
		return
	}

	response, err := h.phoneRiskService.EvaluatePhoneRisk(req.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   err.Error(),
			Message: "請提供有效的台灣電話號碼（手機：09XXXXXXXX，市話：0X-XXXXXXXX）",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// HealthCheck 健康檢查
// @Summary 健康檢查
// @Description 檢查 API 服務狀態
// @Tags system
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *PhoneHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "Phone Risk API 服務正常運行",
	})
}
