package api

import (
	"encoding/json"
	"fmt"
	mfsdk "github.com/mainflux/mainflux/pkg/sdk/go"
	thmocks "github.com/mainflux/mainflux/things/mocks"
	"github.com/ns1labs/orb/pkg/types"
	"github.com/ns1labs/orb/sinks"
	skmocks "github.com/ns1labs/orb/sinks/mocks"
	"github.com/opentracing/opentracing-go/mocktracer"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	contentType = "application/json"
	token		= "token"
	email		= "user@example.com"
)

var (
	name, _ = types.NewIdentifier("teste")
	sink = sinks.Sink{
		Name: name,
		Config: map[string]interface{}{"test": "data"},
	}
)

type testRequest struct {
	client 		*http.Client
	method 		string
	url 		string
	contentType string
	token 		string
	body 		io.Reader
}

func (tr testRequest) make() (*http.Response, error) {
	req, err := http.NewRequest(tr.method, tr.url, tr.body)
	if err != nil {
		return nil, err
	}
	if tr.token != "" {
		req.Header.Set("Authorization", tr.token)
	}
	if tr.contentType != "" {
		req.Header.Set("Content-Type", tr.contentType)
	}
	return tr.client.Do(req)
}

func newService(tokens map[string]string) sinks.Service {
	auth := thmocks.NewAuthService(tokens)
	sinkRepo := skmocks.NewSinkRepository()
	var logger *zap.Logger

	config := mfsdk.Config{
		BaseURL: "localhost",
		ThingsPrefix: "mf",
	}

	mfsdk := mfsdk.NewSDK(config)
	return sinks.NewSinkService(logger, auth, sinkRepo, mfsdk)
}

func newServer(svc sinks.Service) *httptest.Server {
	mux := MakeHandler(mocktracer.New(),"sinks", svc)
	return httptest.NewServer(mux)
}

func toJSON(data interface{}) string {
	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func TestCreateSinks(t *testing.T) {
	service := newService(map[string]string{token: email})
	server := newServer(service)
	defer server.Close()

	sk := sink
	sk.ID = "key"
	data := toJSON(sk)

	req := testRequest{
		client: server.Client(),
		method: http.MethodPost,
		url: fmt.Sprintf("%s/sinks", server.URL),
		contentType: contentType,
		token: token,
		body: strings.NewReader(data),
	}
	res, err := req.make()
	assert.Nil(t, err, fmt.Sprintf("unexpect erro %s", err))
	if res.StatusCode != http.StatusCreated {
		t.Errorf("waited: %d, received: %d", http.StatusCreated, res.StatusCode)
	}

}
