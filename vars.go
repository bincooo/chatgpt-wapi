package chatgpt

const (
	BU        = "https://chat.openai.com/backend-api"
	Gpt3Model = "text-davinci-002-render-sha"
)

var H = map[string]string{
	"User-Agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0",
	"Content-Type": "application/json; charset=utf-8",
	"Accept":       "text/event-stream",
}
