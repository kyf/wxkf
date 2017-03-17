package protocol

type Protocol interface {
	New()
}

type ProtocolMgr struct{}

func initProtocl(proName string, conn net.Conn) (Protocol, error) {

}
