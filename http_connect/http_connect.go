package http_connect

// Authentication - ненужная
// функция для аутентификации
func Authentication() error {
	var err error

	URL := config.Settings.LOKI_URL + config.Settings.LOKI_API_PATH //+ "/login"

	//client := http.Client{Timeout: 60 * time.Second}

	req, err := http.NewRequest(http.MethodGet, URL, http.NoBody)
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(config.Settings.LOKI_LOGIN, config.Settings.LOKI_PASSWORD)

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
