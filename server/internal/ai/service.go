package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"travel-backend/database"
	"travel-backend/domain"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// InitGeminiClient - Khởi tạo Gemini client
func InitGeminiClient(ctx context.Context) (*genai.Client, error) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GEMINI_API_KEY is not set")
	}
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}
	return client, nil
}

// BuildTourContext - Lấy danh sách tour làm dữ liệu nền
func BuildTourContext() (string, error) {
	var tours []domain.Tour
	if err := database.DB.Where("is_active = ?", true).Select("id, name, location, price_amount, duration, description").Find(&tours).Error; err != nil {
		return "", err
	}

	// Đóng gói thành JSON mỏng để tiết kiệm token
	type MiniTour struct {
		ID       uint   `json:"id"`
		Name     string `json:"name"`
		Loc      string `json:"loc"`
		Price    int64  `json:"price"`
		Duration string `json:"duration"`
		Desc     string `json:"desc"`
	}

	var miniTours []MiniTour
	for _, t := range tours {
		miniTours = append(miniTours, MiniTour{
			ID:       t.ID,
			Name:     t.Name,
			Loc:      t.Location,
			Price:    t.PriceAmount,
			Duration: t.Duration,
			Desc:     t.Description,
		})
	}

	b, _ := json.Marshal(miniTours)
	return string(b), nil
}

// ChatWithBot - Hàm chat chính
func ChatWithBot(ctx context.Context, userMessage string) (string, error) {
	client, err := InitGeminiClient(ctx)
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SetTemperature(0.7)

	tourContext, err := BuildTourContext()
	if err != nil {
		tourContext = "Không thể lấy thông tin tour lúc này."
	}

	systemPrompt := `Bạn là trợ lý ảo thân thiện của công ty du lịch Traveling. Nhiệm vụ của bạn là tư vấn tour cho khách hàng dựa trên dữ liệu sau (định dạng JSON):
` + tourContext + `
Hãy trả lời ngắn gọn, tự nhiên, thân thiện bằng tiếng Việt. Luôn ưu tiên giới thiệu các tour có trong dữ liệu trên. Nếu khách hỏi thông tin không có trong hệ thống, hãy nói rằng hiện tại công ty chưa có tour đó.`

	prompt := []genai.Part{
		genai.Text(systemPrompt),
		genai.Text("\nKhách hàng: " + userMessage),
	}

	resp, err := model.GenerateContent(ctx, prompt...)
	if err != nil {
		return "", err
	}

	var responseText strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			responseText.WriteString(string(text))
		}
	}

	return responseText.String(), nil
}

// AnalyzeBehavior - Phân tích hành vi để gợi ý tour
func AnalyzeBehavior(ctx context.Context, userID uint) ([]uint, error) {
	// Lấy 10 hành động gần nhất
	var activities []domain.UserActivity
	database.DB.Where("user_id = ?", userID).Order("timestamp desc").Limit(10).Find(&activities)

	if len(activities) == 0 {
		return []uint{}, nil
	}

	// Lấy tất cả các tour ID hiện có
	var tours []domain.Tour
	database.DB.Where("is_active = ?", true).Select("id, name, location, price_amount").Find(&tours)

	// Gom dữ liệu
	actBytes, _ := json.Marshal(activities)
	tourBytes, _ := json.Marshal(tours)

	client, err := InitGeminiClient(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.5-flash")
	model.SetTemperature(0.2) // Nhiệt độ thấp để trả về JSON chuẩn xác

	systemPrompt := `Dựa trên lịch sử hành vi (activities) của người dùng và danh sách tour hiện có (tours), hãy chọn ra tối đa 3 tour phù hợp nhất.
Chỉ trả về kết quả là một mảng JSON các ID của tour, KHÔNG CÓ BẤT KỲ TEXT NÀO KHÁC (ví dụ: [1, 5, 8]).
Activities: ` + string(actBytes) + `
Tours: ` + string(tourBytes)

	resp, err := model.GenerateContent(ctx, genai.Text(systemPrompt))
	if err != nil {
		return nil, err
	}

	var rawResp string
	for _, part := range resp.Candidates[0].Content.Parts {
		if text, ok := part.(genai.Text); ok {
			rawResp += string(text)
		}
	}

	// Lọc bỏ markdown code blocks nếu AI trả về (ví dụ ```json [1,2,3] ```)
	rawResp = strings.TrimPrefix(strings.TrimSpace(rawResp), "```json")
	rawResp = strings.TrimPrefix(rawResp, "```")
	rawResp = strings.TrimSuffix(rawResp, "```")
	rawResp = strings.TrimSpace(rawResp)

	var recommendedIDs []uint
	if err := json.Unmarshal([]byte(rawResp), &recommendedIDs); err != nil {
		return nil, fmt.Errorf("failed to parse AI response: %w - %s", err, rawResp)
	}

	return recommendedIDs, nil
}
