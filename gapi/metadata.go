package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// gateway keys
const (
	gatewayContentTypeKey = "grpcgateway-content-type"
	gatewayUserAgentKey   = "grpcgateway-user-agent"
	gatewayClientIpKey    = "x-forwarded-for"
	gatewayAuthorityKey   = "x-forwarded-host"
)

// raw grpc keys
const (
	contentTypeKey = "content-type"
	userAgentKey   = "user-agent"
	authorityKey   = ":authority"
)

type Metadata struct {
	UserAgent   string
	ClientIp    string
	ContentType string
	Authority   string
}

// if its gateway, get the values as gateway.
func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// first check gateway
		if contentTypes := md.Get(gatewayContentTypeKey); len(contentTypes) > 0 {
			mtdt.ContentType = contentTypes[0]
		}
		if userAgents := md.Get(gatewayUserAgentKey); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}
		if clientIps := md.Get(gatewayClientIpKey); len(clientIps) > 0 {
			mtdt.ClientIp = clientIps[0]
		}
		if authorities := md.Get(gatewayAuthorityKey); len(authorities) > 0 {
			mtdt.Authority = authorities[0]
		}

		if len(mtdt.UserAgent) == 0 || len(mtdt.ClientIp) == 0 {
			// try native call
			if contentTypes := md.Get(contentTypeKey); len(contentTypes) > 0 {
				mtdt.ContentType = contentTypes[0]
			}
			if userAgents := md.Get(userAgentKey); len(userAgents) > 0 {
				mtdt.UserAgent = userAgents[0]
			}
			if authorities := md.Get(authorityKey); len(authorities) > 0 {
				mtdt.Authority = authorities[0]
			}
			if p, ok := peer.FromContext(ctx); ok {
				mtdt.ClientIp = p.Addr.String()
			}
		}
	}

	return mtdt
}
