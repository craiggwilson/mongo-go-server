package mongo

import (
	"context"
	"crypto/tls"
	"net"
)

// TLSConnectionDecorator decorates a connection with TLS.
func TLSConnectionDecorator(cfg *tls.Config) ConnectionDecorator {
	return ConnectionDecoratorFunc(func(ctx context.Context, c net.Conn) (context.Context, net.Conn) {
		return ctx, tls.Server(c, cfg)
	})
}
