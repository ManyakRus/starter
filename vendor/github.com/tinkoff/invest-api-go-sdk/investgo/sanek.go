package investgo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"

	pb "github.com/tinkoff/invest-api-go-sdk/proto"
	"github.com/tinkoff/invest-api-go-sdk/retry"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/metadata"
)

// NewClient_WithCertificate - создание клиента для API Тинькофф инвестиций с поддержкой кастомного сертификата
func NewClient_WithCertificate(ctx context.Context, conf Config, l Logger, certPath ...string) (*Client, error) {
	setDefaultConfig(&conf)

	var authKey ctxKey = "authorization"
	ctx = context.WithValue(ctx, authKey, fmt.Sprintf("Bearer %s", conf.Token))
	ctx = metadata.AppendToOutgoingContext(ctx, "x-app-name", conf.AppName)

	opts := []retry.CallOption{
		retry.WithCodes(codes.Unavailable, codes.Internal),
		retry.WithBackoff(retry.BackoffLinear(WAIT_BETWEEN)),
		retry.WithMax(conf.MaxRetries),
	}

	// при исчерпывании лимита запросов в минуту, нужно ждать дольше
	exhaustedOpts := []retry.CallOption{
		retry.WithCodes(codes.ResourceExhausted),
		retry.WithMax(conf.MaxRetries),
		retry.WithOnRetryCallback(func(ctx context.Context, attempt uint, err error) {
			l.Infof("Resource Exhausted, sleep for %vs...", attempt)
		}),
	}

	streamInterceptors := []grpc.StreamClientInterceptor{
		retry.StreamClientInterceptor(opts...),
	}

	var unaryInterceptors []grpc.UnaryClientInterceptor
	if conf.DisableResourceExhaustedRetry {
		unaryInterceptors = []grpc.UnaryClientInterceptor{
			retry.UnaryClientInterceptor(opts...),
		}
	} else {
		unaryInterceptors = []grpc.UnaryClientInterceptor{
			retry.UnaryClientInterceptor(opts...),
			retry.UnaryClientInterceptorRE(exhaustedOpts...),
		}
	}

	// Создаём TLS конфигурацию с кастомным сертификатом
	tlsConfig, err := createTLSConfig(certPath...)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS config: %w", err)
	}

	conn, err := grpc.Dial(conf.EndPoint,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithPerRPCCredentials(oauth.TokenSource{
			TokenSource: oauth2.StaticTokenSource(&oauth2.Token{AccessToken: conf.Token}),
		}),
		grpc.WithChainUnaryInterceptor(unaryInterceptors...),
		grpc.WithChainStreamInterceptor(streamInterceptors...))
	if err != nil {
		return nil, err
	}

	client := &Client{
		conn:   conn,
		Config: conf,
		Logger: l,
		ctx:    ctx,
	}

	if conf.AccountId == "" {
		s := client.NewSandboxServiceClient()
		accountsResp, err := s.GetSandboxAccounts()
		if err != nil {
			return nil, err
		}
		accs := accountsResp.GetAccounts()
		if len(accs) < 1 {
			resp, err := s.OpenSandboxAccount()
			if err != nil {
				return nil, err
			}
			client.Config.AccountId = resp.GetAccountId()
		} else {
			for _, acc := range accs {
				if acc.GetStatus() == pb.AccountStatus_ACCOUNT_STATUS_OPEN {
					client.Config.AccountId = acc.GetId()
					break
				}
			}
		}
	}

	return client, nil
}

// createTLSConfig создаёт TLS конфигурацию с опциональным кастомным сертификатом
func createTLSConfig(certPaths ...string) (*tls.Config, error) {
	var certPool *x509.CertPool

	if len(certPaths) > 0 && certPaths[0] != "" {
		// Используем кастомный пул сертификатов
		customPool, err := loadCertificates(certPaths...)
		if err != nil {
			return nil, fmt.Errorf("load custom certificates: %w", err)
		}
		certPool = customPool
	} else {
		// Используем системные сертификаты
		systemPool, err := x509.SystemCertPool()
		if err != nil {
			// Если системный пул недоступен, создаём пустой
			systemPool = x509.NewCertPool()
		}
		certPool = systemPool
	}

	return &tls.Config{
		RootCAs:      certPool,
		NextProtos:   []string{"h2"}, // Обязательно для gRPC 1.67+
		MinVersion:   tls.VersionTLS12,
		//InsecureSkipVerify: true,
	}, nil
}

// loadCertificates загружает один или несколько сертификатов
func loadCertificates(paths ...string) (*x509.CertPool, error) {
	certPool := x509.NewCertPool()

	for _, path := range paths {
		if path == "" {
			continue
		}

		// Читаем файл сертификата
		certPEM, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read certificate file %s: %w", path, err)
		}

		// Добавляем сертификат в пул
		if !certPool.AppendCertsFromPEM(certPEM) {
			return nil, fmt.Errorf("failed to parse certificate from %s", path)
		}
	}

	return certPool, nil
}
