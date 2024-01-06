package accessToken

import "net/http"

func ReadAccessToken(r *http.Request, accessKey string) string {
	return r.Header.Get(accessKey)
}
