package health

import (
	"context"
	"log"
	"time"

	"github.com/gogo/protobuf/types"
	"safeweb.app/rpc/health"
)

const healthOk = "ok"

type (
	Server struct {
		health.UnimplementedHealthServer
		service map[string]string
	}
)

func NewServer() *Server {
	return &Server{service: make(map[string]string)}
}

func (s *Server) Beat(srv string) {
	s.service[srv] = healthOk
}

func (s *Server) SetStatus(srv string, status string) {
	s.service[srv] = status
}

func (s *Server) Add(ctx context.Context, req *health.HealthAddRequest) (*health.HealthCheckResponse, error) {
	res := &health.HealthCheckResponse{
		Status: healthOk,
	}

	srv := req.GetService()
	status := req.GetStatus()
	if srv != "" {
		s.service[srv] = status
	}

	return res, nil
}

func (s *Server) List(ctx context.Context, req *types.StringValue) (*health.HealthListResponse, error) {
	res := &health.HealthListResponse{}

	names := make([]string, 0)

	for name, _ := range s.service {
		names = append(names, name)
	}

	res.List = names
	res.Total = int64(len(names))

	return res, nil
}

func (s *Server) Check(ctx context.Context, req *types.StringValue) (*health.HealthCheckResponse, error) {
	res := &health.HealthCheckResponse{
		Status: healthOk,
	}

	srv := req.Value

	if srv != "" {
		res.Status = s.service[srv]
	}

	return res, nil

}

func (s *Server) Watch(req *health.HealthCheckRequest, stream health.Health_WatchServer) error {
	// Start a ticker that executes each 1000 milliseconds
	timer := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		// Exit on stream context done
		case <-stream.Context().Done():
			return nil
		case <-timer.C:
			srv := req.GetService()
			output := &health.HealthCheckResponse{
				Status: healthOk,
			}

			if srv != "" {
				output.Status = s.service[srv]
			}

			// Send the response of command execution on the stream
			err := stream.Send(output)
			if err != nil {
				log.Println("[ERROR] Health - server Watch", err.Error())
				return err
			}
		}
	}
}
