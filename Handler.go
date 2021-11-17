package gof

import (
    "net"
)

type (
    ClientHandler struct {
        ClientDelegate
    }
    ServerHandler struct {
        ServerDelegate
    }
)
// client
func (c *ClientHandler) OnOpen(client *Client)                       {}
func (c *ClientHandler) OnClose(client *Client)                      {}
func (c *ClientHandler) OnError(client *Client, err error)           {}
func (c *ClientHandler) HandlePacket(client *Client, packet *Packet) {}
// server
func (s *ServerHandler) OnStartServe(addr net.Addr)                  {}
func (s *ServerHandler) OnNew(client *Client)                         {}
func (s *ServerHandler) OnClose(client *Client)                       {}
func (s *ServerHandler) HandlePacket(client *Client, packet *Packet) {}
