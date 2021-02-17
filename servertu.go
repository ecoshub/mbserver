package mbserver

import (
	"fmt"
	"io"

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
		buffer := make([]byte, 512)

		bytesRead, err := port.Read(buffer)
		if err != nil {
			if err != io.EOF {
				if s.onErrorhandler != nil {
					(*s.onErrorhandler)(fmt.Errorf("[mbserver] serial read error %v", err))
				}
			}
			return
		}

		if bytesRead != 0 {

			// Set the length of the packet to the number of read bytes.
			packet := buffer[:bytesRead]

			frame, err := NewRTUFrame(packet)
			if err != nil {
				if s.onErrorhandler != nil {
					(*s.onErrorhandler)(fmt.Errorf("[mbserver] serial read error %v", err))
				}
				return
			}
			if frame.GetAddress() == slaveID {
				request := &Request{port, frame}
				s.requestChan <- request
			} else {
				if s.onErrorhandler != nil {
					(*s.onErrorhandler)(fmt.Errorf("[mbserver] serial read error %v", err))
				}
			}

		}
	}
}
