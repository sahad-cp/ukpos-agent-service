package printer

import (
	"fmt"
	"net"
	"time"
)

func Send(ip string, port int, timeout int, data []byte) error {
	address := fmt.Sprintf("%s:%d", ip, port)

	conn, err := net.DialTimeout(
		"tcp",
		address,
		time.Duration(timeout)*time.Millisecond,
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(data)
	return err
}