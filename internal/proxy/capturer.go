package proxy

import (
	"bytes"
	"net/http"
)

type responseCapturer struct {
	http.ResponseWriter
	buffer     *bytes.Buffer
	statusCode int
	header     http.Header
}

func (r *responseCapturer) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	// Сохраняем копию заголовков на момент вызова WriteHeader
	r.header = r.ResponseWriter.Header().Clone()
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseCapturer) Write(data []byte) (int, error) {
	// Если статус ещё не установлен, считаем, что это 200 OK
	if r.statusCode == 0 {
		r.WriteHeader(http.StatusOK)
	}
	r.buffer.Write(data)
	return r.ResponseWriter.Write(data)
}

// Header() уже правильно возвращает заголовки оригинального ResponseWriter
func (r *responseCapturer) Header() http.Header {
	return r.ResponseWriter.Header()
}
