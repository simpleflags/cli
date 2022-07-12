package config

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
)

func ExtractClaimsFromJWT(token string) map[string]interface{} {
	parts := strings.Split(token, ".")
	bytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		log.Printf("error decoding token with err: %v", err)
		return nil
	}

	result := make(map[string]interface{})
	err = json.Unmarshal(bytes, &result)
	if err != nil {
		log.Printf("error unmarshalling token with err: %v", err)
		return result
	}

	return result
}
