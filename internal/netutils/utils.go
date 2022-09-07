package netutils

import (
	"net"
)

var nq NetQueue = NewQueue()

func SendMsg(conn net.Conn, msg string) error {
	_, err := conn.Write([]byte(msg + "\000"))

	return err
}

// RecvMsg reads the connection, stores all the incoming messages in the queue as
// successive items and returns the next item in the queue as a string
func RecvMsg(conn net.Conn) (string, error) {
	msg := make([]byte, 0)
	buf := make([]byte, 1024)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return "", nil
		}

		// beginning of a new message
		start := 0

		for i := 0; i < n; i++ {
			if buf[i] == 0x00 {
				// append the new message to msg
				msg = append(msg, buf[start:i]...)

				// make a copy of the message
				msgCopy := make([]byte, len(msg))
				copy(msgCopy, msg)

				// add the message copy to the queue
				nqi := NewItem(msgCopy)
				nq.AddItem(nqi)

				// set/reset the variables
				msg = make([]byte, 0)
				start = i
			} else if i == n-1 {
				// append all the buffer from start
				msg = append(msg, buf[start:]...)

				// reset start
				start = 0
			}
		}

		// we reached the end of available data
		if n < len(buf) {
			break
		}
	}

	return string(nq.Next().Content()), nil
}
