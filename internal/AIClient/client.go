package AIClient

type AIClientIntf interface {
	ExecutePrompt(p string, s AIResponseSchemaIntf) (string, bool)
}

type AIResponseSchemaIntf interface {
	GetAIResponseSchema() map[string]interface{}
	MarshalJSON() ([]byte, error)
}
