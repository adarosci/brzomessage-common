package dto

import (
	"encoding/json"
)

// ResultMessages retorno mensagens http
type ResultMessages struct {
	Success      bool        `json:"success"`
	Error        error       `json:"error"`
	Message      string      `json:"message"`
	ObjectResult interface{} `json:"result"`
}

// ToJSON converte struct em json para retorno http
func (m ResultMessages) ToJSON() []byte {
	result, _ := json.Marshal(m)
	return []byte(result)
}
