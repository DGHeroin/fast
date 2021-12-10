package fast

import (
    "bufio"
    "net"
    "sync"
)

type (
    flushable interface {
        Flush() error
    }
    Connection interface {
        net.Conn
        flushable
    }
    bufferedConn struct {
        net.Conn
        bufReader  *bufio.Reader
        bufWriter  *bufio.Writer
        flushMutex sync.Mutex
    }
)

func (c *bufferedConn) Read(b []byte) (n int, err error) {
    return c.bufReader.Read(b)
}

func (c *bufferedConn) Write(b []byte) (n int, err error) {
    return c.bufWriter.Write(b)
}

func (c *bufferedConn) Close() error {
    _ = c.Flush()
    return c.Conn.Close()
}

func (c *bufferedConn) Flush() error {
    c.flushMutex.Lock()
    err := c.bufWriter.Flush()
    c.flushMutex.Unlock()
    if err != nil {
        return err
    }
    if f, ok := c.Conn.(flushable); ok {
        return f.Flush()
    }
    return nil
}

func NewConnection(conn net.Conn, readBufferSize, writeBufferSize int) Connection {
    return &bufferedConn{
        Conn:      conn,
        bufReader: bufio.NewReaderSize(conn, readBufferSize),
        bufWriter: bufio.NewWriterSize(conn, writeBufferSize),
    }
}
