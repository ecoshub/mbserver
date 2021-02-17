package mbserver

import (
	"fmt"
	"io"
	"net"
)

func (s *Server) accept(listen net.Listener) error {
	for {
		conn, err := listen.Accept()
		if err != nil {
			// if strings.Contains(err.Error(), "use of closed network connection") {
			// 	return nil
			// }
			return fmt.Errorf("[mbserver] Unable to accept connections: %v", err)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			for {
				packet := make([]byte, 512)
				bytesRead, err := conn.Read(packet)
				if err != nil {
					if err != io.EOF {
						if s.errorHandler != nil {
							(*s.errorHandler)(fmt.Errorf("[mbserver] read error %v", err))
						}
					}
					return
				}
				// Set the length of the packet to the number of read bytes.
				packet = packet[:bytesRead]

				frame, err := NewTCPFrame(packet)
				if err != nil {
					if s.errorHandler != nil {
						(*s.errorHandler)(fmt.Errorf("[mbserver] bad packet error %v", err))
					}
					return
				}

				request := &Request{conn, frame}

				s.requestChan <- request
			}
		}(conn)
	}
}

// ListenTCP starts the Modbus server listening on "address:port".
func (s *Server) ListenTCP(addressPort string) (err error) {
	listen, err := net.Listen("tcp", addressPort)
	if err != nil {
		return err
	}
	s.listeners = append(s.listeners, listen)
	go s.accept(listen)
	return err
}
