package people

import "net/http"

func parseQuery(r *http.Request, key string) string {
	value := r.URL.Query().Get(key)
	return value
}
