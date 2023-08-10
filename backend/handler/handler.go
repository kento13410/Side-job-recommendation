package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strconv"
	"text/template"

	"github.com/labstack/echo/v4"
	"github.com/sashabaranov/go-openai"
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

type blob []byte

type applyRequest struct {
	Name string `json:"name"`
	Mail mail.Address `json:"mail"`
	Artifact blob `json:"artifact"`
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func Init(e *echo.Echo) {
	e.GET("/", Root)
	e.GET("/test", Test)
	e.POST("/calculate", Calculate)
	e.GET("/share", Share)
	e.POST("/apply", Apply)

	t := &Template{
		templates: template.Must(template.ParseGlob("../frontend/templates/*.html")),
	}
	e.Renderer = t
}

func Root(c echo.Context) error {
	return c.File("../frontend/templates/index.html")
}

func Test(c echo.Context) error {
	return c.Render(http.StatusOK, "test.html", nil)
}

func Calculate(c echo.Context) error {
	var req testRequest
	req.MonthlyIncome, _ = strconv.Atoi(c.FormValue("monthly_income"))
	req.MonthlyWorkHours, _ = strconv.Atoi(c.FormValue("monthly_work_hours"))
	params, _ := c.FormParams()
	req.Name = params["name"]
	// for i := range req.Name {
	// 	req.Level[i], _ = strconv.Atoi(params["level"][i])
	// }
	req.Level = make([]int, len(req.Name))
	req.LeaningHours, _ = strconv.Atoi(c.FormValue("learning_hours"))

	var skills string
	for i := range req.Name {
		skills += req.Name[i] + "（習熟度：" + fmt.Sprint(req.Level[i]) + "）, "
	}
	message := "以下の条件を満たすおすすめの副業を教えてください。月に欲しい金額：" + fmt.Sprint(req.MonthlyIncome) + "円, 月に働ける時間：" + fmt.Sprint(req.MonthlyWorkHours) + "時間, 現在のスキル：" + skills + "学習に使うことのできる時間：" + fmt.Sprint(req.LeaningHours) + "時間"

	Recommendation, err := RecommendJob(message)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.Render(http.StatusOK, "result.html", Recommendation)
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

func RecommendJob(message string) (string, error) {
	client := openai.NewClient("OPENAI_API_KEY")
	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
			MaxTokens: 500, // 出力トークン数の上限
			Temperature: 0, // 出力のランダム性を低くする
		},
	)

	if err != nil {
		panic(err)
	}

	return resp.Choices[0].Message.Content, nil
}
