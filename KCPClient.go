package gof

import (
    "github.com/xtaci/kcp-go"
    "net"
)

func newKCPClient(opt Option) (*Client, error) {
    var (
        conn net.Conn
        err  error
    )

    conn, err = kcp.Dial(opt.Address)
    c := &Client{
        address: opt.Address,
    }
    c.Conn = conn
    return c, err
}
