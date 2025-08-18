package http_connect

import (
	"crypto/tls"
	"github.com/ManyakRus/starter/log"
	"io"
	"net"
	"net/http"
	"time"
)

// Client - клиент для http
var Client *http.Client

// Authentication - ненужная
// функция для аутентификации
func Authentication(URL, login, password string) error {
	var err error

	req, err := http.NewRequest(http.MethodGet, URL, http.NoBody)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(login, password)

	res, err := Client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	print(resBody)

	//
	return err
}

// CreateClient - создаёт клиент http
// для работы без ошибки сертификатов
func CreateClient() {
	Client = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			TLSClientConfig: &tls.Config{
				// See comment above.
				// UNSAFE!
				// DON'T USE IN PRODUCTION!
				InsecureSkipVerify: true,
			},
		},
	}
}
