package protocol

type Client interface {
	Close() error
}

// Protocol describes the basic behavior of any protocol in the system
type Protocol interface {
	IOLoop(Client) error
}

/*
// SendResponse is a server side utility function to prefix data with a length header
func SendResponse(w io.Writer, d []byte) (int, error) {
	err := binary.Write(w, binary.BigEndian, int32(len(d)))
	if err != nil {
		return 0, err
	}
	n, err := w.Write(d)
	if err != nil {
		return 0, err
	}
	return n + 4, nil
}

// SendCommandResponse is a server side utility function to prefix data with a length header and command header
// write to the supplied Writer
func SendCommandResponse(w io.Writer, frameType common.Command, d []byte) (int, error) {
	sizeBuf := make([]byte, 4)
	size := uint32(len(d)) + 4
	binary.BigEndian.PutUint32(sizeBuf, size)
	n, err := w.Write(sizeBuf)
	if err != nil {
		return n, err
	}
	binary.BigEndian.PutUint32(sizeBuf, uint32(frameType))
	n, err = w.Write(sizeBuf)
	if err != nil {
		return n + 4, err
	}

	n, err = w.Write(d)
	return n + 8, err

}


*/
