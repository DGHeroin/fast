package fast

import (
    "crypto/tls"
    "crypto/x509"
    "io"
    "io/ioutil"
    "net"
    "runtime"
    "time"
)

func writeFull(conn io.Writer, data []byte) error {
    left := len(data)
    for left > 0 {
        n, err := conn.Write(data)
        if n == left && err == nil { // handle most common case first
            return nil
        }

        if n > 0 {
            data = data[n:]
            left -= n
        }

        if err != nil {
            if !IsTemporary(err) {
                return err
            } else {
                runtime.Gosched()
            }
        }
    }
    return nil
}

func readFull(conn io.Reader, data []byte) error {
    left := len(data)
    for left > 0 {
        n, err := conn.Read(data)
        if n == left && err == nil { // handle most common case first
            return nil
        }

        if n > 0 {
            data = data[n:]
            left -= n
        }

        if err != nil {
            if !IsTemporary(err) {
                return err
            } else {
                runtime.Gosched()
            }
        }
    }
    return nil
}

func IsTemporary(err error) bool {
    if err == nil {
        return false
    }

    e, ok := err.(net.Error)
    return ok && e.Temporary()
}

func tryFlush(conn net.Conn) error {
    if f, ok := conn.(flushable); ok {
        for {
            err := f.Flush()
            if err != nil && !IsTemporary(err) {
                return err
            } else {
                return nil
            }
        }
    }
    return nil
}

func Go(fn func()) {
    go fn()
}

func RunForeverUntilPanic(AutoRestartServerDuration time.Duration, fn func()) {
    isRunning := true
    invoke := func() {
        defer func() {
            //if e := recover(); e != nil {
            //    isRunning = false
            //    log.Println(e)
            //}
        }()
        fn()
    }
    for isRunning {
        invoke()
        time.Sleep(AutoRestartServerDuration)
    }
}
func RunForever(AutoRestartServerDuration time.Duration, fn func()) {
    for {
        fn()
        time.Sleep(AutoRestartServerDuration)
    }
}

func GoN(n int, fn func()) {
    for i := 0; i < n; i++ {
        Go(fn)
    }
}

func LoadTLSConfig(cert, key, ca string) (*tls.Config, error) {
    certConfig, err := tls.LoadX509KeyPair(cert, key)
    if err != nil {
        return nil, err
    }
    var clientCertPool *x509.CertPool
    if ca != "" {
        certBytes, err := ioutil.ReadFile(ca)
        if err != nil {
            return nil, err
        }
        clientCertPool = x509.NewCertPool()
        clientCertPool.AppendCertsFromPEM(certBytes)
        config := &tls.Config{
            Certificates: []tls.Certificate{certConfig},
            RootCAs:      clientCertPool,
        }
        return config, nil
    } else {
        config := &tls.Config{
            Certificates: []tls.Certificate{certConfig},
            ClientAuth:   tls.RequireAndVerifyClientCert,
        }
        return config, nil
    }
}
