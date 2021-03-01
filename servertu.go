package mbserver

import (
	"fmt"
	"io"

	"github.com/ecoshub/breakx"
	"github.com/goburrow/serial"
)

// ListenRTU starts the Modbus server listening to a serial device.
// For example:  err := s.ListenRTU(&serial.Config{Address: "/dev/ttyUSB0"})
func (s *Server) ListenRTU(serialConfig *serial.Config, slaveID uint8) (err error) {
	port, err := serial.Open(serialConfig)
	if err != nil {
		return err

	}
	s.ports = append(s.ports, port)
	go s.acceptSerialRequests(port, slaveID)
	return err
}

func (s *Server) acceptSerialRequests(port serial.Port, slaveID uint8) {
	for {
		breakx.Point()
		buffer := make([]byte, 512)

		bytesRead, err := port.Read(buffer)
		breakx.Point()
		if err != nil {
			breakx.Point()
			if err != io.EOF {
				breakx.Point()
				if s.onErrorhandler != nil {
					breakx.Point()
					(*s.onErrorhandler)(fmt.Errorf("[mbserver] serial read error %v", err))
				}
			}
			return
		}

		breakx.Point()
		if bytesRead != 0 {

			// Set the length of the packet to the number of read bytes.
			packet := buffer[:bytesRead]

			breakx.Point()
			frame, err := NewRTUFrame(packet)
			if err != nil {
				breakx.Point()
				if s.onErrorhandler != nil {
					breakx.Point()
					(*s.onErrorhandler)(fmt.Errorf("[mbserver] serial read error %v", err))
				}
				return
			}
			breakx.Point()
			if frame.GetAddress() == slaveID {
				breakx.Point()
				request := &Request{port, frame}
				s.requestChan <- request
			} else {
				breakx.Point()
				if s.onErrorhandler != nil {
					breakx.Point()
					(*s.onErrorhandler)(fmt.Errorf("[mbserver] serial read error %v", err))
				}
			}

		}
	}
}
