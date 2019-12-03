// package page

// import (
// 	"encoding/base64"
// 	"encoding/json"
// )

// var B64 = base64.StdEncoding

// //SkipTakeTokenData is a standard struct to use for TokenData
// type SkipTakeTokenData struct {
// 	Skip      int64       `json:"skip"`
// 	QueryMeta interface{} `json:"query,omitempty"`
// }

// type Tokener interface {
// 	ToTokenString(tokenData interface{}) ([]byte, error)
// 	ParseToken(token []byte, obj interface{}) error
// }

// type StandardTokener struct{}

// func (a StandardTokener) ToTokenString(tokenData interface{}) ([]byte, error) {
// 	jsonBytes, err := json.Marshal(tokenData)
// 	if err != nil {
// 		return nil, err
// 	}

// 	token := make([]byte, B64.EncodedLen(len(jsonBytes)))
// 	B64.Encode(token, jsonBytes)
// 	return token, nil
// }

// func (a StandardTokener) ParseToken(token []byte, obj interface{}) error {
// 	jsonBytes := make([]byte, B64.DecodedLen(len(token)))
// 	_, err := B64.Decode(jsonBytes, token)
// 	if err != nil {
// 		return err
// 	}

// 	return json.Unmarshal(jsonBytes, obj)
// }
