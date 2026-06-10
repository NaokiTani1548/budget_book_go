package ocr

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"budget-book-go/internal/application/dto"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type AnalyzeReceiptUseCase struct {
	apiKey string
}

func NewAnalyzeReceiptUseCase(apiKey string) *AnalyzeReceiptUseCase {
	return &AnalyzeReceiptUseCase{apiKey: apiKey}
}

func (uc *AnalyzeReceiptUseCase) Execute(ctx context.Context, imageData []byte, mimeType string) (*dto.OCRResult, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(uc.apiKey))
	if err != nil {
		return nil, fmt.Errorf("Geminiクライアントの作成に失敗しました: %w", err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-3.1-flash-lite")
	model.SetTemperature(0.1)

	prompt := `この画像は決済履歴のスクリーンショットまたは領収書です。
画像から支出情報を抽出し、以下のJSON形式で返してください。
複数の支出がある場合はすべて抽出してください。

必ず以下のJSON形式のみを返してください。説明文やマークダウンは不要です。

{
  "items": [
    {
      "description": "店名や商品名",
      "amount": 金額（数値）,
      "expenseDate": "YYYY-MM-DD形式の日付（画像から読み取れない場合は今日の日付）",
      "paymentMethod": "CASH|CREDIT_CARD（画像から判断、不明ならCREDIT_CARD）"
    }
  ]
}`

	resp, err := model.GenerateContent(ctx,
		genai.Text(prompt),
		genai.Blob{MIMEType: mimeType, Data: imageData},
	)
	if err != nil {
		return nil, fmt.Errorf("Gemini APIの呼び出しに失敗しました: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("Geminiからの応答が空です")
	}

	text := fmt.Sprintf("%v", resp.Candidates[0].Content.Parts[0])

	// マークダウンのコードブロックを除去
	text = strings.TrimPrefix(text, "```json")
	text = strings.TrimPrefix(text, "```")
	text = strings.TrimSuffix(text, "```")
	text = strings.TrimSpace(text)

	var result dto.OCRResult
	if err := json.Unmarshal([]byte(text), &result); err != nil {
		return nil, fmt.Errorf("Geminiの応答をパースできません: %w\n応答: %s", err, text)
	}

	return &result, nil
}