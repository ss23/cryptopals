package main

import (
	"encoding/base64"
	"encoding/hex"
	//"io/ioutil"
	"bytes"
	"fmt"
	//"log"
	"math"
	"os"
)

func main() {
	string1 := os.Args[1]

	bytes1 := []byte(string1)

	bytesRaw1, _ := decodeHex(bytes1)

	fmt.Println(bytesRaw1)
}

func decodeHex(src []byte) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(src)))
	_, err := hex.Decode(dst, src)
	return dst, err
}

func encodeHex(src []byte) string {
	return hex.EncodeToString(src)
}

func encodeBase64(src []byte) string {
	return base64.StdEncoding.EncodeToString(src)
}

// Will do repeating xor, with one being the input, second being the key
func xor(input []byte, key []byte) (dst []byte) {
	// We need to calculate how long the key should be
	if len(input) > len(key) {
		key = bytes.Repeat(key, len(input)/len(key))
	}

	dst = make([]byte, len(input))

	for i := 0; i < len(input); i++ {
		dst[i] = input[i] ^ key[i]
	}
	return
}

func scoreByLetterFreq(input []byte) float64 {
	freq := float64(bytes.Count(input, []byte("e"))) / float64(len(input))
	return math.Abs(freq - float64(.1202))
}
