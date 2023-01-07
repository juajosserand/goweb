package httpserver

import "net"

type Option func(*server)

func Port(p string) Option {
	return func(s *server) {
		s.server.Addr = net.JoinHostPort("", p)
	}
}
