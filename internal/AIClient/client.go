package AIClient

type AIClientIntf interface {
	ExecutePrompt(p string) (string, bool)
}
