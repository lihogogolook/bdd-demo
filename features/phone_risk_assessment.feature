Feature: 電話號碼風險評估 API
  作為一個 API 用戶
  我想要能夠評估電話號碼的風險等級
  以便做出相應的商業決策

  Background:
    Given API 服務已經啟動
    And API 端點為 "http://localhost:8080"

  Scenario: 健康檢查端點正常運作
    When 我發送 GET 請求到 "/health"
    Then 我應該收到狀態碼 200
    And 回應應該包含系統狀態資訊

  Scenario: 取得 API 基本資訊
    When 我發送 GET 請求到 "/"
    Then 我應該收到狀態碼 200
    And 回應應該包含 API 基本資訊

  Scenario: 評估有效的手機號碼風險
    Given 我有一個有效的手機號碼 "0987654321"
    When 我發送 POST 請求到 "/api/phone/risk" 包含電話號碼
    Then 我應該收到狀態碼 200
    And 回應應該包含 "phone_number" 為 "0987654321"
    And 回應應該包含 "risk_score" 為數字
    And 回應應該包含 "risk_level" 為 "low", "medium", 或 "high" 之一
    And 回應應該包含 "message" 為 "電話號碼風險評估完成"

  Scenario: 評估有效的市話號碼風險
    Given 我有一個有效的市話號碼 "02-12345678"
    When 我發送 POST 請求到 "/api/phone/risk" 包含電話號碼
    Then 我應該收到狀態碼 200
    And 回應應該包含 "phone_number" 為 "02-12345678"
    And 回應應該包含 "risk_score" 為數字
    And 回應應該包含 "risk_level" 為 "low", "medium", 或 "high" 之一
    And 回應應該包含 "message" 為 "電話號碼風險評估完成"

  Scenario Outline: 測試不同風險等級的電話號碼
    Given 我有一個電話號碼 "<phone_number>"
    When 我發送 POST 請求到 "/api/phone/risk" 包含電話號碼
    Then 我應該收到狀態碼 200
    And 回應的 "risk_level" 應該為 "<expected_risk_level>"

    Examples:
      | phone_number | expected_risk_level |
      |   0911111111 | low                 |
      |   0922222222 | medium              |
      |   0933333333 | high                |

  Scenario: 測試無效的電話號碼格式
    Given 我有一個無效的電話號碼 "1234567890"
    When 我發送 POST 請求到 "/api/phone/risk" 包含電話號碼
    Then 我應該收到狀態碼 400
    And 回應應該包含 "error" 為 "無效的電話號碼格式"
    And 回應應該包含 "message" 包含 "請提供有效的台灣電話號碼"

  Scenario: 測試空的電話號碼
    Given 我有一個空的電話號碼 ""
    When 我發送 POST 請求到 "/api/phone/risk" 包含電話號碼
    Then 我應該收到狀態碼 400
    And 回應應該包含錯誤訊息

  Scenario: 測試缺少電話號碼欄位的請求
    When 我發送 POST 請求到 "/api/phone/risk" 不包含電話號碼欄位
    Then 我應該收到狀態碼 400
    And 回應應該包含錯誤訊息
