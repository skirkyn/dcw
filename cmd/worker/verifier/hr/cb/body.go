package cb

import (
	"github.com/skirkyn/dcw/cmd/common"
	"github.com/skirkyn/dcw/cmd/worker/verifier/hr"
)

var defaultRequestTemplate = `{
  "sms": {
    "token": "%s"
  },
  "constraints": {
    "mode": "ALLOW",
    "types": [
      "NO_2FA",
      "SMS",
      "TOTP",
      "U2F",
      "IDV",
      "RECOVERY_CODE",
      "PUSH",
      "PASSKEY",
      "SECURITY_QUESTION"
    ]
  },
  "action": "web-UnifiedLogin-IdentificationPrompt",
  "bot_token": ""
}`

type BodySupplier struct {
	formatter common.Function[string, []byte]
}

func NewBodySupplier() common.Function[string, []byte] {
	return &BodySupplier{hr.NewFormattingBodySupplier[string](defaultRequestTemplate)}
}

func (sf *BodySupplier) Apply(in string) ([]byte, error) {
	return sf.formatter.Apply(in)
}
