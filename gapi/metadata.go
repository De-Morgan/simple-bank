package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardHeader             = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIp  string
}

func (server *Server) extractMetadata(cxt context.Context) *Metadata {
	mtdt := &Metadata{}
	if md, ok := metadata.FromIncomingContext(cxt); ok {
		if useragent := md.Get(grpcGatewayUserAgentHeader); len(useragent) > 0 {
			mtdt.UserAgent = useragent[0]
		}
		if useragent := md.Get(userAgentHeader); len(useragent) > 0 {
			mtdt.UserAgent = useragent[0]
		}
		if clientIps := md.Get(xForwardHeader); len(clientIps) > 0 {
			mtdt.ClientIp = clientIps[0]
		}
	}
	if p, ok := peer.FromContext(cxt); ok {
		mtdt.ClientIp = p.Addr.String()
	}
	return mtdt
}
