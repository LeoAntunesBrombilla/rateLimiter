package accessToken

import (
	"net/http"
	"os"
	"strings"
)

func getAPIKeyNames() []string {
	apiKeyNameList, exists := os.LookupEnv("API_KEY_NAME_LIST")
	if !exists {
		return []string{"default_api_key"}
	}
	return strings.Split(apiKeyNameList, ",")
}

func ReadAccessToken(r *http.Request, accessKey string) (string, string) {
	apiKeyNames := getAPIKeyNames()
	for _, apiKeyName := range apiKeyNames {
		if value := r.Header.Get(apiKeyName); value != "" {
			return apiKeyName, value
		}
	}

	return "", ""
}
