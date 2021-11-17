package gof

import (
    "crypto/tls"
    "time"
)

type Option struct {
    Address       string
    TLS           *tls.Config
    RetryDuration time.Duration
}
