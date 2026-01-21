package pprof

import (
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
)

type Server struct {
	addr string
}

func New(addr string) *Server {
	return &Server{addr: addr}
}

func (s *Server) Run() {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Printf("pprof error: %v", err)
		return
	}

	log.Printf("pprof server started on %s", s.addr)
	go func() {
		if err := http.Serve(ln, nil); err != nil {
			log.Printf("pprof server stopped: %v", err)
		}
	}()
}
