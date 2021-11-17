package gof

import "crypto/tls"

type Option struct {
    Address string
    TLS     *tls.Config
}
