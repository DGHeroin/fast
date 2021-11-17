package gof

import (
    "crypto/tls"
    "net"
)

func newTCPClient(opt Option) (*Client, error) {
    var (
        conn net.Conn
        err  error
    )

    if opt.TLS == nil {
        addr, err := net.ResolveTCPAddr("tcp", opt.Address)
        if err != nil {
            return nil, err
        }
        laddr, err := net.ResolveTCPAddr("tcp", ":0")
        if err != nil {
            return nil, err
        }
        conn, err = net.DialTCP("tcp", laddr, addr)
        if err != nil {
            return nil, err
        }
    } else {
        conn, err = tls.Dial("tcp", opt.Address, opt.TLS)
    }
    c := &Client{
        address: opt.Address,
    }
    c.Conn = conn
    return c, err
}
