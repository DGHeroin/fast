package gof

func Serve(network string, delegate Delegate, opt Option) {
    RunForeverUntilPanic(AutoRestartServerDuration, func() {
        switch network {
        case "tcp":
            serveTCP(delegate, opt)
        case "kcp":
            serveKCP(delegate, opt)
        default:
            panic("unsupported protocol")
        }
    })
}
func NewClient(network string, option Option) (*Client, error)  {
    switch network {
    case "tcp":
        return newTCPClient(option)
    case "kcp":
        return newKCPClient(option)
    default:
        panic("unsupported protocol")
    }
}