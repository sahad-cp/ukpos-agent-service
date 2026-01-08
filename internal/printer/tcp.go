package printer

import (
	"fmt"
	"net"
	"time"
)

func Send(ip string, port int, timeout int, data []byte) error {
	address := net.JoinHostPort(ip, fmt.Sprintf("%d", port))

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

// package printer

// import (
// 	"fmt"
// 	"net"
// 	"time"
// )

// func SendTCP(ip string, port int, timeout time.Duration, data []byte) error {
// 	addr := fmt.Sprintf("%s:%d", ip, port)

// 	// 1️⃣ Connect with timeout
// 	conn, err := net.DialTimeout("tcp", addr, timeout)
// 	if err != nil {
// 		return err
// 	}
// 	tcpConn := conn.(*net.TCPConn)
// 	tcpConn.SetNoDelay(true) // IMPORTANT
// 	tcpConn.SetWriteDeadline(time.Now().Add(timeout))
// 	defer conn.Close()

// 	// 2️⃣ Disable Nagle (CRITICAL for printers)
// 	if tcpConn, ok := conn.(*net.TCPConn); ok {
// 		_ = tcpConn.SetNoDelay(true)
// 	}

// 	// 3️⃣ Write deadline (printers hate slow writes)
// 	_ = conn.SetWriteDeadline(time.Now().Add(timeout))

// 	// 4️⃣ Write ALL bytes
// 	total := 0
// 	for total < len(data) {
// 		n, err := conn.Write(data[total:])
// 		if err != nil {
// 			return err
// 		}
// 		total += n
// 		time.Sleep(5 * time.Millisecond)
// 	}

// 	// 5️⃣ Let printer consume buffer (VERY IMPORTANT)
// 	time.Sleep(200 * time.Millisecond)

// 	return nil
// }
