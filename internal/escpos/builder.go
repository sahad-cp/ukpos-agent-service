package escpos

import "time"

const (
	ESC = "\x1B"
	GS  = "\x1D"
	LF  = "\x0A"
)

func TestPrint(outlet, printerIP string) []byte {
	data := "\x1B@"              // INIT
	data += "\x1Ba\x01"          // Center
	data += "\x1D!\x30"          // Double size
	data += "UKPOS PRINT TEST\r\n"
	data += "\x1D!\x00"          // Normal
	data += "Printer: " + printerIP + "\r\n"
	data += "Outlet: " + outlet + "\r\n"
	data += "Time: " + time.Now().Format(time.RFC1123) + "\r\n"
	data += "\r\n\r\n"
	data += "\x1D\x56\x41\x00"   // SAFE CUT (GS V 65 0)
	return []byte(data)
}