package proxy

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
)

// generateCacheKey создаёт ключ на основе метода, URL и важных заголовков
/*
func (p *Proxy) generateCacheKey(r *http.Request) string {
	// Собираем строку: method + full URL + headers
	var parts []string
	parts = append(parts, r.Method)
	parts = append(parts, r.URL.String())

	// Добавляем заголовки, влияющие на ответ
	headersToInclude := []string{"Accept", "Accept-Encoding", "Accept-Language", "Authorization"}
	for _, h := range headersToInclude {
		if v := r.Header.Get(h); v != "" {
			parts = append(parts, h+":"+v)
		}
	}

	raw := strings.Join(parts, "|")
	hash := md5.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])
}
*/
func (p *Proxy) generateCacheKey(r *http.Request) string {
	raw := r.Method + "|" + r.URL.String()
	hash := md5.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])
}
