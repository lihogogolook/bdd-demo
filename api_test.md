# Phone Risk Assessment API 測試

## 測試請求範例

### 1. 健康檢查

```bash
curl -X GET http://localhost:8080/health
```

### 2. 取得 API 資訊

```bash
curl -X GET http://localhost:8080/
```

### 3. 評估手機號碼風險

```bash
curl -X POST http://localhost:8080/api/phone/risk \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "0987654321"
  }'
```

### 4. 評估市話號碼風險

```bash
curl -X POST http://localhost:8080/api/phone/risk \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "02-12345678"
  }'
```

### 5. 測試無效號碼

```bash
curl -X POST http://localhost:8080/api/phone/risk \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "1234567890"
  }'
```

## 回應格式

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

## 風險等級說明

- **low** (0-29.9): 低風險
- **medium** (30-69.9): 中等風險
- **high** (70-100): 高風險
