package utils

import (
	"io"
	"net"
	"regexp"
)

const DEFAULT_HISLIP_PORT = "4880"

func SplitAddress(address string) (string, error) {
	pattern := `(\d{1,3}(?:\.\d{1,3}){3})(?:.*?,(\d+))?`
	re := regexp.MustCompile(pattern)

	var ip string
	var port string

	matches := re.FindStringSubmatch(address)
	if len(matches) > 2 && matches[2] != "" {
		ip = matches[1]
		port = matches[2]
	} else if len(matches) > 1 {
		ip = matches[1]
		port = DEFAULT_HISLIP_PORT
	}

	final_address := ip + ":" + port

	return final_address, nil
}

func WriteAll(conn net.Conn, buf []byte) error {
	for len(buf) > 0 {
		n, err := conn.Write(buf)
		if err != nil {
			return err
		}
		if n == 0 {
			return io.ErrUnexpectedEOF
		}
		buf = buf[n:]
	}
	return nil
}
