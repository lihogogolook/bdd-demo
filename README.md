# Phone Risk Assessment API

這是一個使用 Go 語言和 Gin 框架建立的電話號碼風險評估 API 服務。

## 功能特色

- 🔍 **電話號碼驗證**：支援台灣手機號碼和市話格式驗證
- 📊 **風險評估**：根據電話號碼特徵計算風險分數 (0-100)
- 🎯 **風險等級**：提供三級風險分類 (低/中/高)
- 🌐 **RESTful API**：標準的 REST API 設計
- ❤️ **健康檢查**：服務狀態監控端點

## 專案結構

```
bdd-demo/
├── go.mod
├── main.go                    # 主程式入口
├── README.md
├── api_test.md               # API 測試範例
├── .gitignore
└── internal/
    ├── models/
    │   └── phone.go          # 資料模型
    ├── services/
    │   └── phone_risk.go     # 風險評估服務
    └── handlers/
        └── phone.go          # API 處理器
```

## 快速開始

### 1. 安裝相依套件

```bash
go mod tidy
```

### 2. 啟動服務

```bash
go run main.go
```

服務將在 `http://localhost:8080` 啟動

### 3. 測試 API

#### 健康檢查

```bash
curl -X GET http://localhost:8080/health
```

#### 評估電話號碼風險

```bash
curl -X POST http://localhost:8080/api/phone/risk \
  -H "Content-Type: application/json" \
  -d '{"phone_number": "0987654321"}'
```

## API 端點

| 方法 | 路徑              | 描述             |
| ---- | ----------------- | ---------------- |
| GET  | `/`               | API 服務資訊     |
| GET  | `/health`         | 健康檢查         |
| POST | `/api/phone/risk` | 電話號碼風險評估 |

## 請求/回應格式

### 風險評估請求

```json
{
  "phone_number": "0987654321"
}
```

### 成功回應

```json
{
  "phone_number": "0987654321",
  "risk_score": 75.5,
  "risk_level": "medium",
  "message": "電話號碼風險評估完成"
}
```

### 錯誤回應

```json
{
  "error": "無效的電話號碼格式",
  "message": "請提供有效的台灣電話號碼（手機：09XXXXXXXX，市話：0X-XXXXXXXX）"
}
```

## 支援的電話號碼格式

- **手機號碼**：`09XXXXXXXX` (10 位數字)
- **市話號碼**：`0X-XXXXXXXX` 或 `0XX-XXXXXXX`

## 風險等級

- **低風險 (low)**：0 - 29.9 分
- **中等風險 (medium)**：30 - 69.9 分
- **高風險 (high)**：70 - 100 分

## 開發

### 建置

```bash
go build -o phone-risk-api
```

### 測試

```bash
go test ./...
```

### 查看詳細的 API 測試範例

請參考 [api_test.md](./api_test.md) 檔案
