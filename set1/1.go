package main

import (
	"encoding/hex"
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"fmt"
)

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Fatal error reading from stdin")
	}
	bytesRaw := decodeHex(bytes)
	base64 := encodeBase64(bytesRaw)

	fmt.Println(base64)
}

func decodeHex(src []byte) ([]byte) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		log.Fatal("Failed to decode hex", err)
	}
	return dst
}

func encodeBase64(src []byte) (string) {
	return base64.StdEncoding.EncodeToString(src)
}
