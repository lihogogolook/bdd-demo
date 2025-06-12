# BDD 測試框架介紹

## 1. Cucumber 與 Gherkin 語法

**Cucumber** 是一個行為驅動開發 (BDD) 工具，使用自然語言描述測試場景。**Gherkin** 是其描述語法，具有以下關鍵字：

- `Feature`：功能描述
- `Scenario`：測試場景
- `Given`：前置條件
- `When`：執行動作
- `Then`：預期結果
- `And`：連接詞

```gherkin
Feature: 電話號碼風險評估
  Scenario: 評估手機號碼風險
    Given 系統已啟動
    When 我提交手機號碼 "0987654321"
    Then 應該返回風險分數
    And 風險等級應該是 "medium"
```

## 2. Godog 基本使用

**Godog** 是 Go 語言的 Cucumber 實作：

### 安裝

```bash
go get github.com/cucumber/godog/cmd/godog
```

### 執行測試

```bash
godog features/
```

### 步驟定義範例

```go
func InitializeTestSuite(ctx *godog.TestSuiteContext) {
    ctx.BeforeSuite(func() { /* 初始化 */ })
}

func InitializeScenario(ctx *godog.ScenarioContext) {
    ctx.Step(`^我提交手機號碼 "([^"]*)"$`, iSubmitPhoneNumber)
    ctx.Step(`^應該返回風險分數$`, shouldReturnRiskScore)
}
```

## 3. AI 溝通的 Prompt 與 Feature Test

### 與 AI 協作的最佳實務

1. **清晰的需求描述**

   ```
   "請為電話風險評估 API 撰寫 BDD 測試，包含：
   - 有效號碼驗證
   - 風險分數計算
   - 錯誤處理"
   ```

2. **結構化的 Feature 檔案**

   - 使用描述性的場景名稱
   - 明確的步驟定義
   - 涵蓋正常流程和異常情況

3. **可維護的測試設計**
   - 模組化步驟定義
   - 資料驅動測試
   - 清楚的斷言邏輯

### 範例 Prompt

```gherkin
Feature: 電話號碼風險評估
  Scenario: 評估手機號碼風險
    Given 系統已啟動
    When 我提交手機號碼 "0987654321"
    Then 應該返回風險分數
    And 風險等級應該是 "medium"
```

這種方法讓技術需求變得更容易理解，同時提高測試覆蓋率和程式碼品質。
