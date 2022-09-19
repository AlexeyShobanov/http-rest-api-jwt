// определим совой responseWiter, чтобы можно было из него получить информацию о коде возврата (можно еще что-то добавить)
package apiserver

import "net/http"

type responseWiter struct {
	http.ResponseWriter
	code int
}

// statusCode у ResponseWriter записывается в методе WriteHeader, поэтому мы его переопределяем
func (w *responseWiter) WriteHeader(statusCode int) {
	w.code = statusCode // здесь мы запишем статус код в поле code структуры
	// затем передадим обработку стандартному методу WriteHeader из ResponseWriter
	w.ResponseWriter.WriteHeader(statusCode)
}
