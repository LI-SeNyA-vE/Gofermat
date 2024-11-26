package middleware

import (
	"compress/gzip"
	"github.com/LI-SeNyA-vE/Gofermat/internal/global"
	"io"
	"net/http"
)

type (
	gzipWriter struct {
		http.ResponseWriter
		Writer io.Writer
	}
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Encoding") == "gzip" {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, "Ошибка при создании gzip.Reader", http.StatusInternalServerError)
				return
			}
			defer gz.Close()
			// Замена r.Body на распакованный stream
			r.Body = io.NopCloser(gz)
		}
		next.ServeHTTP(w, r)
	})
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func UnGzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		global.Logger.Info("Провалились в функцию UnGzipMiddleware")
		if !(r.Header.Get("Accept-Encoding") == "gzip") {
			global.Logger.Info("Accept-Encoding не равен gzip")
			next.ServeHTTP(w, r)
			return
		}
		global.Logger.Info("Accept-Encoding равен gzip")
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			global.Logger.Info("Ошибка при gzip.NewWriterLevel(w, gzip.BestSpeed)")
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()
		global.Logger.Info("Нет ошибки при gzip.NewWriterLevel(w, gzip.BestSpeed)")
		w.Header().Set("Content-Encoding", "gzip")
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)

	})
}
