package page

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

var B64 = base64.StdEncoding

//SkipTakeTokenData is a standard struct to use for TokenData
type SkipTakeTokenData struct {
	Skip  int64           `json:"skip"`
	Query json.RawMessage `json:"query,omitempty"`
}

func (td SkipTakeTokenData) ParseQuery(q interface{}) error {
	return json.Unmarshal([]byte(td.Query), q)
}

func (td *SkipTakeTokenData) SetQuery(q interface{}) {
	bs, err := json.Marshal(q)
	if err != nil {
		//this should only be used with structs that are
		//universally marshallable, so this should never
		//happen
		panic(err)
	}

	td.Query = json.RawMessage(bs)
}

//Tokener converts tokens to structured data and back
//go:generate counterfeiter . Tokener
type Tokener interface {
	ToTokenString(tokenData interface{}) ([]byte, error)
	ParseToken(token string, obj interface{}) error
}

//Base64JSONTokener converts nbetween structured token objects and
//base64-encoded JSON string.
type Base64JSONTokener struct{}

//ToTokenString converts a struct to a base64-JSON byte string
func (a Base64JSONTokener) ToTokenString(tokenData interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(tokenData)
	if err != nil {
		return nil, err
	}

	token := make([]byte, B64.EncodedLen(len(jsonBytes)))
	B64.Encode(token, jsonBytes)
	return token, nil
}

//TODO log these errors but don't pass them back to the frontend

//ParseToken parses a base64-JSON byte string back into a struct
func (a Base64JSONTokener) ParseToken(token string, obj interface{}) error {
	jsonBytes, err := B64.DecodeString(token)
	if err != nil {
		return fmt.Errorf("failed decoding base64: %w", err)
	}

	err = json.Unmarshal(jsonBytes, obj)
	if err != nil {
		return fmt.Errorf("failed unmarshaling json: %w", err)
	}
	return nil
}
