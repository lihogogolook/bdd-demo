package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"bdd-demo/internal/handlers"
	"bdd-demo/internal/models"
	"bdd-demo/internal/services"

	"github.com/cucumber/godog"
	"github.com/gin-gonic/gin"
)

type apiFeatureContext struct {
	server       *httptest.Server
	response     *http.Response
	responseBody []byte
	phoneNumber  string
	router       *gin.Engine
}

func (ctx *apiFeatureContext) aPIService() error {
	// 設置 Gin 為測試模式
	gin.SetMode(gin.TestMode)

	// 建立服務
	phoneRiskService := services.NewPhoneRiskService()

	// 建立處理器
	phoneHandler := handlers.NewPhoneHandler(phoneRiskService)

	// 建立路由器
	ctx.router = gin.New()

	// 添加 CORS 中介軟體
	ctx.router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// 健康檢查端點
	ctx.router.GET("/health", phoneHandler.HealthCheck)

	// 根端點
	ctx.router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": "Phone Risk Assessment API",
			"version": "1.0.0",
			"status":  "running",
		})
	})

	// API 路由群組
	api := ctx.router.Group("/api")
	{
		phone := api.Group("/phone")
		{
			phone.POST("/risk", phoneHandler.EvaluatePhoneRisk)
		}
	}

	// 建立測試服務器
	ctx.server = httptest.NewServer(ctx.router)

	return nil
}

func (ctx *apiFeatureContext) aPIEndpoint(endpoint string) error {
	// 在實際場景中，這裡可以設置不同的端點
	// 目前使用測試服務器的 URL
	return nil
}

func (ctx *apiFeatureContext) iSendGETRequestTo(path string) error {
	url := ctx.server.URL + path
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	ctx.response = resp
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	ctx.responseBody = body
	resp.Body.Close()

	return nil
}

func (ctx *apiFeatureContext) iShouldReceiveStatusCode(expectedCode int) error {
	if ctx.response.StatusCode != expectedCode {
		return fmt.Errorf("expected status code %d, but got %d", expectedCode, ctx.response.StatusCode)
	}
	return nil
}

func (ctx *apiFeatureContext) responseShouldContainSystemStatusInfo() error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(ctx.responseBody, &responseData); err != nil {
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	// 檢查是否包含狀態資訊
	if _, exists := responseData["status"]; !exists {
		return fmt.Errorf("response does not contain status information")
	}

	return nil
}

func (ctx *apiFeatureContext) responseShouldContainAPIBasicInfo() error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(ctx.responseBody, &responseData); err != nil {
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	// 檢查是否包含基本 API 資訊
	requiredFields := []string{"service", "version", "status"}
	for _, field := range requiredFields {
		if _, exists := responseData[field]; !exists {
			return fmt.Errorf("response does not contain required field: %s", field)
		}
	}

	return nil
}

func (ctx *apiFeatureContext) iHaveAValidMobileNumber(phoneNumber string) error {
	ctx.phoneNumber = phoneNumber
	return nil
}

func (ctx *apiFeatureContext) iHaveAValidLandlineNumber(phoneNumber string) error {
	ctx.phoneNumber = phoneNumber
	return nil
}

func (ctx *apiFeatureContext) iHaveAPhoneNumber(phoneNumber string) error {
	ctx.phoneNumber = phoneNumber
	return nil
}

func (ctx *apiFeatureContext) iHaveAnInvalidPhoneNumber(phoneNumber string) error {
	ctx.phoneNumber = phoneNumber
	return nil
}

func (ctx *apiFeatureContext) iHaveAnEmptyPhoneNumber(phoneNumber string) error {
	ctx.phoneNumber = phoneNumber
	return nil
}

func (ctx *apiFeatureContext) iSendPOSTRequestToWithPhoneNumber(path string) error {
	url := ctx.server.URL + path

	requestBody := models.PhoneRiskRequest{
		PhoneNumber: ctx.phoneNumber,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	ctx.response = resp
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	ctx.responseBody = body
	resp.Body.Close()

	return nil
}

func (ctx *apiFeatureContext) iSendPOSTRequestToWithoutPhoneNumberField(path string) error {
	url := ctx.server.URL + path

	// 發送空的 JSON 物件
	emptyJSON := []byte("{}")

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(emptyJSON))
	if err != nil {
		return err
	}

	ctx.response = resp
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	ctx.responseBody = body
	resp.Body.Close()

	return nil
}

func (ctx *apiFeatureContext) responseShouldContainPhoneNumberAs(field, expectedValue string) error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(ctx.responseBody, &responseData); err != nil {
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	actualValue, exists := responseData[field]
	if !exists {
		return fmt.Errorf("field %s not found in response", field)
	}

	if actualValue != expectedValue {
		return fmt.Errorf("expected %s to be %s, but got %v", field, expectedValue, actualValue)
	}

	return nil
}

func (ctx *apiFeatureContext) responseShouldContainAsNumber(field string) error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(ctx.responseBody, &responseData); err != nil {
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	value, exists := responseData[field]
	if !exists {
		return fmt.Errorf("field %s not found in response", field)
	}

	// 檢查是否為數字類型
	switch value.(type) {
	case float64, int, int64:
		return nil
	default:
		return fmt.Errorf("field %s is not a number, got %T", field, value)
	}
}

func (ctx *apiFeatureContext) responseShouldContainAsOneOf(field, allowedValues string) error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(ctx.responseBody, &responseData); err != nil {
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	actualValue, exists := responseData[field]
	if !exists {
		return fmt.Errorf("field %s not found in response", field)
	}

	// 解析允許的值
	allowed := strings.Split(strings.ReplaceAll(allowedValues, " ", ""), ",")
	actualStr := fmt.Sprintf("%v", actualValue)

	for _, allowedValue := range allowed {
		cleanValue := strings.Trim(allowedValue, `"`)
		if actualStr == cleanValue {
			return nil
		}
	}

	return fmt.Errorf("field %s value %v is not one of allowed values: %s", field, actualValue, allowedValues)
}

func (ctx *apiFeatureContext) responseShouldContainMessageAs(field, expectedMessage string) error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(ctx.responseBody, &responseData); err != nil {
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	actualValue, exists := responseData[field]
	if !exists {
		return fmt.Errorf("field %s not found in response", field)
	}

	if actualValue != expectedMessage {
		return fmt.Errorf("expected %s to be %s, but got %v", field, expectedMessage, actualValue)
	}

	return nil
}

func (ctx *apiFeatureContext) responseRiskLevelShouldBe(field, expectedRiskLevel string) error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(ctx.responseBody, &responseData); err != nil {
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	actualRiskLevel, exists := responseData[field]
	if !exists {
		return fmt.Errorf("field %s not found in response", field)
	}

	if actualRiskLevel != expectedRiskLevel {
		return fmt.Errorf("expected %s to be %s, but got %v", field, expectedRiskLevel, actualRiskLevel)
	}

	return nil
}

func (ctx *apiFeatureContext) responseShouldContainErrorAs(field, expectedError string) error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(ctx.responseBody, &responseData); err != nil {
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	actualValue, exists := responseData[field]
	if !exists {
		return fmt.Errorf("field %s not found in response", field)
	}

	if actualValue != expectedError {
		return fmt.Errorf("expected %s to be %s, but got %v", field, expectedError, actualValue)
	}

	return nil
}

func (ctx *apiFeatureContext) responseShouldContainMessageContaining(field, expectedSubstring string) error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(ctx.responseBody, &responseData); err != nil {
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	actualValue, exists := responseData[field]
	if !exists {
		return fmt.Errorf("field %s not found in response", field)
	}

	actualStr := fmt.Sprintf("%v", actualValue)
	if !strings.Contains(actualStr, expectedSubstring) {
		return fmt.Errorf("expected %s to contain %s, but got %v", field, expectedSubstring, actualValue)
	}

	return nil
}

func (ctx *apiFeatureContext) responseShouldContainErrorMessage() error {
	var responseData map[string]interface{}
	if err := json.Unmarshal(ctx.responseBody, &responseData); err != nil {
		return fmt.Errorf("failed to parse response JSON: %v", err)
	}

	// 檢查是否包含錯誤訊息 (error 或 message 欄位)
	if _, hasError := responseData["error"]; hasError {
		return nil
	}
	if _, hasMessage := responseData["message"]; hasMessage {
		return nil
	}

	return fmt.Errorf("response does not contain error message")
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	apiCtx := &apiFeatureContext{}

	// Background 步驟
	ctx.Step(`^API 服務已經啟動$`, apiCtx.aPIService)
	ctx.Step(`^API 端點為 "([^"]*)"$`, apiCtx.aPIEndpoint)

	// When 步驟
	ctx.Step(`^我發送 GET 請求到 "([^"]*)"$`, apiCtx.iSendGETRequestTo)
	ctx.Step(`^我發送 POST 請求到 "([^"]*)" 包含電話號碼$`, apiCtx.iSendPOSTRequestToWithPhoneNumber)
	ctx.Step(`^我發送 POST 請求到 "([^"]*)" 不包含電話號碼欄位$`, apiCtx.iSendPOSTRequestToWithoutPhoneNumberField)

	// Given 步驟
	ctx.Step(`^我有一個有效的手機號碼 "([^"]*)"$`, apiCtx.iHaveAValidMobileNumber)
	ctx.Step(`^我有一個有效的市話號碼 "([^"]*)"$`, apiCtx.iHaveAValidLandlineNumber)
	ctx.Step(`^我有一個電話號碼 "([^"]*)"$`, apiCtx.iHaveAPhoneNumber)
	ctx.Step(`^我有一個無效的電話號碼 "([^"]*)"$`, apiCtx.iHaveAnInvalidPhoneNumber)
	ctx.Step(`^我有一個空的電話號碼 "([^"]*)"$`, apiCtx.iHaveAnEmptyPhoneNumber)

	// Then 步驟
	ctx.Step(`^我應該收到狀態碼 (\d+)$`, apiCtx.iShouldReceiveStatusCode)
	ctx.Step(`^回應應該包含系統狀態資訊$`, apiCtx.responseShouldContainSystemStatusInfo)
	ctx.Step(`^回應應該包含 API 基本資訊$`, apiCtx.responseShouldContainAPIBasicInfo)
	ctx.Step(`^回應應該包含 "([^"]*)" 為 "([^"]*)"$`, apiCtx.responseShouldContainPhoneNumberAs)
	ctx.Step(`^回應應該包含 "([^"]*)" 為數字$`, apiCtx.responseShouldContainAsNumber)
	ctx.Step(`^回應應該包含 "([^"]*)" 為 "([^"]*)", "([^"]*)", 或 "([^"]*)" 之一$`, func(field, val1, val2, val3 string) error {
		allowedValues := fmt.Sprintf(`"%s", "%s", "%s"`, val1, val2, val3)
		return apiCtx.responseShouldContainAsOneOf(field, allowedValues)
	})
	ctx.Step(`^回應的 "([^"]*)" 應該為 "([^"]*)"$`, apiCtx.responseRiskLevelShouldBe)
	ctx.Step(`^回應應該包含 "([^"]*)" 為 "([^"]*)"$`, apiCtx.responseShouldContainMessageAs)
	ctx.Step(`^回應應該包含 "([^"]*)" 包含 "([^"]*)"$`, apiCtx.responseShouldContainMessageContaining)
	ctx.Step(`^回應應該包含錯誤訊息$`, apiCtx.responseShouldContainErrorMessage)

	// 場景後清理
	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if apiCtx.server != nil {
			apiCtx.server.Close()
		}
		return ctx, nil
	})
}
