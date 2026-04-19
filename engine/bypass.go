package engine

import (
	"crypto/tls"
	"net/http"
	"time"
)

// CreateCustomClient bikin client yang sidik jarinya lebih "manusiawi"
func CreateCustomClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			// Bypass pengecekan sertifikat (biar kenceng)
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
				MinVersion:         tls.VersionTLS12,
				CurvePreferences:   []tls.CurveID{tls.CurveP256, tls.X25519},
			},
			ForceAttemptHTTP2: true, // CF suka HTTP/2
			MaxIdleConns:      100,
			IdleConnTimeout:   90 * time.Second,
		},
		Timeout: 5 * time.Second,
	}
}

// BuildCFRequest nyusun header biar mirip Browser asli
func BuildCFRequest(target string) *http.Request {
	req, _ := http.NewRequest("GET", target, nil)

	// Urutan dan isi Header harus konsisten
	req.Header.Set("authority", "target.com")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("sec-ch-ua", `"Chromium";v="124", "Google Chrome";v="124", "Not-A.Brand";v="99"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", GetRandomUA()) // Panggil dari headers.go
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("accept-language", "en-US,en;q=0.9,id;q=0.8")

	return req
}