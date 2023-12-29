package integration

import (
	"net/http/httptest"
	"testing"

	"github.com/pirosiki197/emoine/pkg/handler"
	"github.com/pirosiki197/emoine/pkg/infra/proto/api/v1/apiv1connect"
	"github.com/pirosiki197/emoine/pkg/infra/repository"
)

type APIServer struct {
	srv    *httptest.Server
	client apiv1connect.APIServiceClient
}

func NewAPIServer(t *testing.T) *APIServer {
	t.Helper()
	db := NewDB(t)
	t.Cleanup(db.Cleanup)
	repo := repository.NewRepository(db.BunDB())
	_, handler := apiv1connect.NewAPIServiceHandler(handler.NewHandler(repo))
	srv := httptest.NewUnstartedServer(handler)
	srv.EnableHTTP2 = true
	srv.StartTLS()
	client := apiv1connect.NewAPIServiceClient(srv.Client(), srv.URL)
	return &APIServer{
		srv:    srv,
		client: client,
	}
}

func (s *APIServer) Close() {
	s.srv.Close()
}

func (s *APIServer) Client() apiv1connect.APIServiceClient {
	return s.client
}
