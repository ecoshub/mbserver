package mbserver

import (
	"io"
	"log"

	"github.com/goburrow/serial"
)

var (
	errStringOpenFailed string = "failed to open %s: %v\n"
	errStringSerialRead string = "serial read error %v\n"
	errStringBadFrame   string = "bad serial frame error %v\n"
	errStringUnknown    string = "unknown error %v\n"
	logPacket           string = "packet content: %v\n"
)

// ListenRTU starts the Modbus server listening to a serial device.
// For example:  err := s.ListenRTU(&serial.Config{Address: "/dev/ttyUSB0"})
func (s *Server) ListenRTU(serialConfig *serial.Config) (err error) {
	port, err := serial.Open(serialConfig)
	if err != nil {
		log.Printf(errStringOpenFailed, serialConfig.Address, err)
		return err
	}
	s.ports = append(s.ports, port)
	go s.acceptSerialRequests(port)
	return nil
}

func (s *Server) acceptSerialRequests(port serial.Port) {
	for {
		buffer := make([]byte, 512)

		bytesRead, err := port.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Printf(errStringSerialRead, err)
				log.Printf(logPacket, string(buffer))
				continue
			}
			log.Printf(errStringUnknown, err)
			continue
		}

		if bytesRead != 0 {
			// Set the length of the packet to the number of read bytes.
			packet := buffer[:bytesRead]
			frame, err := NewRTUFrame(packet)
			if err != nil {
				log.Printf(errStringBadFrame, err)
				continue
			}
			request := &Request{port, frame}
			s.requestChan <- request
		}
	}
}
