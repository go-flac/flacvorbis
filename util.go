package flacvorbis

import (
	"bytes"
	"encoding/binary"
	"io"
	"strings"
)

func encodeUint32(n uint32) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.LittleEndian, n); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func readUint32(r io.Reader) (res uint32, err error) {
	err = binary.Read(r, binary.LittleEndian, &res)
	return
}

func packStr(w io.Writer, s string) {
	data := []byte(s)
	w.Write(encodeUint32(uint32(len(data))))
	w.Write(data)
}

func packMapValue(m map[string]string, key string, value string, sep string) {
	if xval, exists := m[key]; exists {
		m[key] = strings.Join([]string{xval, value}, sep)
	} else {
		m[key] = value
	}
}

func unpackMapValue(m map[string]string, key string, sep string) []string {
	if xval, exists := m[key]; exists {
		return strings.Split(xval, sep)
	}
	return []string{}
}
