package protocol

import (
	"bufio"
	"net"
	"regexp"
)

type ProtocolType int

const (
	ProtocolUnknown ProtocolType = iota
	ProtocolThrift
	ProtocolHttp
)

var httpReg = regexp.MustCompile(`^(?:GET |POST|PUT|DELE|HEAD|OPTI|CONN|TRAC|PATC)$`)

// HTTP协议方法以GET/POST等名称开头
func isHttpProtocol(data []byte) bool {
	return httpReg.Match(data)
}

func DectProtocol(conn net.Conn) (ProtocolType, error) {
	reader := bufio.NewReader(conn)
	peek, err := reader.Peek(8)
	if err != nil {
		return ProtocolUnknown, err
	}

	if isHttpProtocol(peek) {
		return ProtocolHttp, nil
	}

	if len(peek) < 4 {
		return ProtocolUnknown, nil
	}

	return ProtocolThrift, nil
}
