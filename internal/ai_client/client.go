package aiclient

type Client struct {
	ApiKey string
	Model  string
}

type ClientIntf interface {
	ExecutePrompt(p string)
}
