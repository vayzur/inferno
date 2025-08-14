package spark

import "github.com/vayzur/inferno/pkg/httputil"

type SparkClient struct {
	httpClient *httputil.Client
}

func NewSparkClient(httpClient *httputil.Client) *SparkClient {
	return &SparkClient{httpClient: httpClient}
}
