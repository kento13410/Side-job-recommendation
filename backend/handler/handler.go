package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"

	"github.com/labstack/echo/v4"
)

type testRequest struct {
	MonthlyIncome int `json:"monthly_income"`
	MonthlyWorkHours int `json:"monthly_work_hours"`
	Name []string `json:"name"`
	Level []int `json:"level"`
	LeaningHours int `json:"leaning_hours"`
}

type ChatGPTRequest struct {
	Model    string `json:"model"`
	Query    string `json:"query"`
	MaxTokens int `json:"max_tokens"`
}

type ChatGPTResponse struct {
    Response      string  `json:"response"`
    Score         float64 `json:"score"`
    ConversationID string `json:"conversationId"`
    UserID        string `json:"userId"`
    BotID         string `json:"botId"`
    SessionID     string `json:"sessionId"`
}

type blob []byte

type applyRequest struct {
	Name string `json:"name"`
	Mail mail.Address `json:"mail"`
	Artifact blob `json:"artifact"`
}

func Init(e *echo.Echo) {
	e.GET("/", Root)
	e.GET("/test", Test)
	e.POST("/calculate", Calculate)
	e.GET("/share", Share)
	e.POST("/apply", Apply)
}

func Root(c echo.Context) error {
	return c.File("../frontend/templates/index.html")
}

func Test(c echo.Context) error {
	return c.Render(http.StatusOK, "../frontend/templates/test.html", nil)
}

func Calculate(c echo.Context) error {
	var req testRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	fmt.Println(req)

	var skills string
	for i := range req.Name {
		skills += req.Name[i] + "（習熟度：" + fmt.Sprint(req.Level[i]) + "）, "
	}
	message := "以下の条件を満たすおすすめの副業を教えてください。月に欲しい金額：" + fmt.Sprint(req.MonthlyIncome) + "円, 月に働ける時間：" + fmt.Sprint(req.MonthlyWorkHours) + "時間, 現在のスキル：" + skills + "学習に使うことのできる時間：" + fmt.Sprint(req.LeaningHours) + "時間"

	recommendation, err := RecommendJob(message)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.Render(http.StatusOK, "../../frontend/templates/test.html", recommendation)
}

func Share(c echo.Context) error {
	return nil
}

func Apply(c echo.Context) error {
	var req applyRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	return c.JSON(http.StatusOK, req)
}

// RecommendJob は、ChatGPT APIを使っておすすめの副業を出力する関数です。
// 引数には、現在所有しているスキルとその習熟度、月に欲しい金額、月に働ける時間、学習に使うことのできる時間などの情報が含まれるメッセージを渡します。
// 返り値には、おすすめの副業に関する情報やアドバイスなどが含まれるレスポンスを返します。
func RecommendJob(message string) (string, error) {
    // ChatGPT APIのキーとエンドポイントを設定する
    key, ok := os.LookupEnv("OPEN_AI_SECRET")
	if !ok {
		panic("open-api-secret is empty")
	}
    endpoint := "https://chatgpt.cognitiveservices.azure.com/generateAnswer" // ここに自分のエンドポイントを入力する

    // リクエスト用のJSONデータを作成する
    request := ChatGPTRequest{
		Model:         "gpt-3.5-turbo",
        Query:         message,
		MaxTokens:     500,
    }
    requestBody, err := json.Marshal(request)
    if err != nil {
        return "", err
    }

    // HTTPクライアントを作成する
    client := &http.Client{}

    // HTTPリクエストを作成する
    req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
    if err != nil {
        return "", err
    }

    // ヘッダーにキーとコンテントタイプを設定する
    req.Header.Add("Ocp-Apim-Subscription-Key", key)
    req.Header.Add("Content-Type", "application/json")

    // HTTPリクエストを送信する
    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    
    // レスポンスボディを読み込む
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    // レスポンスボディを閉じる
    defer resp.Body.Close()

    // レスポンス用のJSONデータをパースする
    var response ChatGPTResponse
    err = json.Unmarshal(responseBody, &response)
    if err != nil {
        return "", err
    }

    // レスポンスからおすすめの副業を取得する
    job := response.Response

    // おすすめの副業を返す
    return job, nil

}