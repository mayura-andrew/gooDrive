package auth

import (
	"encoding/json"
	"os"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expiry       int64  `json:"expiry"`
}

const TokenFilePath = "tokens.json"

// func SaveToken(token *Token) error {
// 	file, err := os.OpenFile(TokenFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	encoder := json.NewEncoder(file)
// 	return encoder.Encode(token)
// }

func LoadToken() (*Token, error) {
	file, err := os.Open(TokenFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var token Token
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&token)
	if err != nil {
		return nil, err
	}

	return &token, nil
}
