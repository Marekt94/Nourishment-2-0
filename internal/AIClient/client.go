package AIClient

import "github.com/invopop/jsonschema"

type AIClientIntf interface {
	ExecutePrompt(p string, s *jsonschema.Schema) (string, bool)
}
