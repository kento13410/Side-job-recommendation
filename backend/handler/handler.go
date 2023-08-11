package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"strconv"
	"strings"
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

type applyRequest struct {
	Name string `json:"name"`
	Mail mail.Address `json:"mail"`
}

type Template struct {
	templates *template.Template
}

type SideJob struct {
	JobName string
	JobDescription string
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
	Skills := []string{"プログラミング", "英語", "デザイン", "マーケティング", "ライティング", "翻訳", "動画編集", "音楽制作", "イラスト", "写真撮影"}
	return c.Render(http.StatusOK, "test.html", Skills)
}

func Calculate(c echo.Context) error {
	var req testRequest
	req.MonthlyIncome, _ = strconv.Atoi(c.FormValue("monthly_income"))
	req.MonthlyWorkHours, _ = strconv.Atoi(c.FormValue("monthly_work_hours"))
	params, _ := c.FormParams()
	req.Name = params["name"]
	for i := range req.Name {
		level, _ := strconv.Atoi(c.FormValue(strconv.Itoa(i)))
		req.Level = append(req.Level, level)
	}
	req.LeaningHours, _ = strconv.Atoi(c.FormValue("learning_hours"))
	fmt.Println(req)

	var skills string
	for i := range req.Name {
		if req.Level[i] > 0 {
			skills += req.Name[i] + "（習熟度：" + fmt.Sprint(req.Level[i]) + "）, "
		}
	}
	message :=
	`以下の条件を満たすおすすめの副業を教えてください。ただし、次のフォーマットに従ってください。

	・条件
	月に欲しい金額：` + fmt.Sprint(req.MonthlyIncome) + `円
	月に働ける時間：` + fmt.Sprint(req.MonthlyWorkHours) + `時間
	現在のスキル(1~5の5段階評価)：` + skills + `
	学習に使うことのできる時間：` + fmt.Sprint(req.LeaningHours) + `時間
	
	・フォーマット説明
	おすすめの副業を1単語で述べ、改行し、補足の説明を加える。
	
	・例
	プログラミング
	あなたにおすすめの副業はプログラミングです。なぜなら..`

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
	return c.JSON(http.StatusOK, req)
}

func RecommendJob(message string) (SideJob, error) {
	client := openai.NewClient("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")
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

	response := resp.Choices[0].Message.Content
	slice := strings.SplitN(response, "\n", 2)
	SideJob := SideJob{JobName: slice[0], JobDescription: slice[1]}
	return SideJob, nil
}