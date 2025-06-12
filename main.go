package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"bdd-demo/internal/handlers"
	"bdd-demo/internal/services"
)

func main() {
	// 建立服務
	phoneRiskService := services.NewPhoneRiskService()

	// 建立處理器
	phoneHandler := handlers.NewPhoneHandler(phoneRiskService)

	// 設置 Gin 模式
	gin.SetMode(gin.ReleaseMode)

	// 建立路由器
	router := gin.Default()

	// 添加 CORS 中介軟體
	router.Use(func(c *gin.Context) {
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
	router.GET("/health", phoneHandler.HealthCheck)

	// API 路由群組
	api := router.Group("/api")
	{
		phone := api.Group("/phone")
		{
			phone.POST("/risk", phoneHandler.EvaluatePhoneRisk)
		}
	}

	// 根路徑顯示 API 資訊
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"name":        "Phone Risk Assessment API",
			"version":     "1.0.0",
			"description": "電話號碼風險評估 API 服務",
			"endpoints": map[string]string{
				"健康檢查":   "GET /health",
				"電話風險評估": "POST /api/phone/risk",
			},
		})
	})

	// 啟動服務器
	port := ":8080"
	log.Printf("🚀 Phone Risk API 服務啟動在 http://localhost%s", port)
	log.Printf("📋 API 文件: http://localhost%s/", port)
	log.Printf("❤️  健康檢查: http://localhost%s/health", port)

	if err := router.Run(port); err != nil {
		log.Fatal("無法啟動服務器:", err)
	}
}
