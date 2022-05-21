package utils

import "encoding/hex"

type Signature struct {
	Sign []byte
}

func (signature *Signature) String() string {
	return hex.EncodeToString(signature.Sign)
}
