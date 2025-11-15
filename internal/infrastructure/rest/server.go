package rest

import (
	"log/slog"
)

type QNAServer struct {
	ServerInstance *Server
}

func NewApp(
	log *slog.Logger,
	addr *string,
	// tokenSecret *[]byte,
	service QNADispatcher,
) *QNAServer {
	server := NewServer(log, service, addr)
	return &QNAServer{
		ServerInstance: server,
	}
}
