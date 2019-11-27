package page

// import (
// 	"encoding/base64"
// 	"encoding/json"
// )

// var B64 = base64.StdEncoding

// type TokenBuilder interface {
// 	ToTokenString(tokenObj interface{}) ([]byte, error)
// 	ParseToken(token []byte, obj interface{}) error
// }

// type TokenAgent struct{}

// func (a TokenAgent) ToTokenString(tokenObj interface{}) ([]byte, error) {
// 	jsonBytes, err := json.Marshal(tokenObj)
// 	if err != nil {
// 		return nil, err
// 	}

// 	token := make([]byte, B64.EncodedLen(len(bs)))
// 	B64.Encode(token, jsonBytes)
// 	return token, nil
// }

// func (a TokenAgent) ParseToken(token []byte, obj interface{}) error {
// 	jsonBytes := make([]byte, B64.DecodedLen(len(token)))
// 	_, err := B64.Decode(jsonBytes, token)
// 	if err != nil {
// 		return err
// 	}

// 	return json.Unmarshal(jsonBytes, obj)
// }
