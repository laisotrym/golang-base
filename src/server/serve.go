package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// Serve server listen for HTTP and GRPC
func (s *Server) WebServe(httpAddress string, grpcAddress string, webroot string) error {
	stop := make(chan os.Signal, 1)
	errch := make(chan error)
	signal.Notify(stop, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// HTTP Server serve Restful API
	httpMux := http.NewServeMux()
	// httpMux.Handle("/api/", s.mux)
	httpMux.Handle("/", http.FileServer(http.Dir(webroot)))

	httpMux.HandleFunc("/api/", func(rw http.ResponseWriter, req *http.Request) {
		if origin := req.Header.Get("Origin"); origin != "" {
			rw.Header().Set("Access-Control-Allow-Origin", origin)
			rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			rw.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Stop here if its Preflighted OPTIONS request
		if req.Method == "OPTIONS" {
			return
		}
		s.mux.ServeHTTP(rw, req)
	})

	httpServer := http.Server{
		Addr:    httpAddress,
		Handler: httpMux,
	}

	// GRPC Server serve Grpc-Web & Protobuf requests
	grpcWebServer := grpcweb.WrapServer(s.gRPC)
	http2Server := &http.Server{
		Addr: grpcAddress,
		Handler: h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 {
				grpcWebServer.ServeHTTP(w, r)
			} else {
				// w.Header().Set("Access-Control-Allow-Origin", "*")
				// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
				// w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
				// w.Header().Set("grpc-status", "")
				// w.Header().Set("grpc-message", "")
				if grpcWebServer.IsGrpcWebRequest(r) {
					grpcWebServer.ServeHTTP(w, r)
				}
			}
		}), &http2.Server{}),
	}

	go func() {
		log.Println("================= HTTP Server Starting at", httpAddress, "=================")
		if err := httpServer.ListenAndServe(); err != nil {
			errch <- err
		}
	}()
	go func() {
		log.Println("================= GRPC Server Starting at", grpcAddress, "=================")
		if err := http2Server.ListenAndServe(); err != nil {
			errch <- err
		}

		// listener, err := net.Listen("tcp", s.cfg.GRPC.String())
		// if err != nil {
		// 	errch <- err
		// 	return
		// }
		// if err := s.gRPC.Serve(listener); err != nil {
		// 	errch <- err
		// }
	}()
	for {
		select {
		case <-stop:
			ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
			httpServer.Shutdown(ctx)
			http2Server.Shutdown(ctx)
			s.gRPC.GracefulStop()
			return nil
		case err := <-errch:
			return err
		}
	}
}
