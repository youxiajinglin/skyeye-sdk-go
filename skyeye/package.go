package skyeye

import (
	"bufio"
	"bytes"
	"encoding/binary"
)

func Encode(id []byte, content []byte) []byte {
	buffer := new(bytes.Buffer)
	binary.Write(buffer, binary.BigEndian, id)
	binary.Write(buffer, binary.BigEndian, uint32(len(content)))
	binary.Write(buffer, binary.BigEndian, content)
	return buffer.Bytes()
}

func Decode(reader *bufio.Reader) ([]byte, error) {
	haedByte, err := reader.Peek(6)
	if err != nil {
		return nil, err
	}
	lengthBuff := bytes.NewBuffer(haedByte[2:6])

	var lenght int32
	err = binary.Read(lengthBuff, binary.BigEndian, &lenght)
	if err != nil {
		return nil, err
	}
	if lenght + 6 > int32(reader.Buffered()) {
		return nil, err
	}

	pack := make([]byte, int(6 + lenght))
	_, err = reader.Read(pack)
	if err != nil {
		return nil, err
	}
	return pack[6:], nil
}