package main

import (
	"encoding/base64"
	"encoding/hex"
	//"io/ioutil"
	"fmt"
	"log"
	"os"
)

func main() {
	string1 := os.Args[1]
	string2 := os.Args[2]

	bytes1 := []byte(string1)
	bytes2 := []byte(string2)

	bytesRaw1 := decodeHex(bytes1)
	bytesRaw2 := decodeHex(bytes2)

	bytesXord := xor(bytesRaw1, bytesRaw2)

	fmt.Println(encodeHex(bytesXord))
}

func decodeHex(src []byte) []byte {
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	if err != nil {
		log.Fatal("Failed to decode hex", err)
	}
	return dst
}

func encodeBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

func xor(src1 []byte, src2 []byte) (dst []byte) {
	if len(src1) != len(src2) {
		log.Fatal("The two specified strings are not the same length")
	}

	dst = make([]byte, len(src1))

	for i := 0; i < len(src1); i++ {
		dst[i] = src1[i] ^ src2[i]
	}
	return
}
