package main

import (
	"encoding/base64"
	"encoding/hex"
	//"io/ioutil"
	"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	var file *os.File
	var err error
	if file, err = os.Open("4.txt"); err != nil {
		log.Fatal("Failed to open 4.txt: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		bytesRaw, _ := decodeHex([]byte(line))
		for i := 0; i < 255; i++ {
			byteArray := make([]byte, 1)
			byteArray[0] = byte(i)
			bytesXord := xor(bytesRaw, byteArray)
			if scoreByLetterFreq(bytesXord) < .15 {
				fmt.Printf("Possible key: %v %x with score %v\r\n", string(i), i, scoreByLetterFreq(bytesXord))
				fmt.Println(string(bytesXord))
			}
		}
	}
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

	//fmt.Printf("We are about to xor %v and %v\r\n", input, key)

	dst = make([]byte, len(input))

	for i := 0; i < len(input); i++ {
		dst[i] = input[i] ^ key[i]
	}
	//fmt.Printf("The result is %x\r\n", dst)
	return
}

func scoreByLetterFreq(input []byte) float64 {
	// Start by filtering out non-ASCII character strings
	for i := range input {
		if input[i] > 127 {
			return 100
		}
	}
	// We can probably cheat by just scoring the first few letters
	var letterFreqs = map[string]float64{
		"E": .1202,
		"T": .0910,
		"A": .0812,
		"O": .0768,
		"I": .0731,
		"N": .0695,
		"S": .0628,
		"R": .0602,
		"H": .0592,
		"D": .0423,
		"L": .0398,
		"U": .0288,
		"C": .0271,
		"M": .0261,
		"F": .0230,
		"Y": .0211,
		"W": .0209,
		"G": .0203,
		"P": .0182,
		"B": .0149,
		"V": .0111,
		"K": .0069,
		"X": .0017,
		"Q": .0011,
		"J": .0010,
		"Z": .0007,
	}

	// Loop over the letters we're scoring and calculate a running total
	// A lower score is better (0 is best for example)
	score := float64(0)
	for k, v := range letterFreqs {
		count := bytes.Count(input, []byte(strings.ToUpper(k))) + bytes.Count(input, []byte(strings.ToLower(k)))
		freq := float64(count) / float64(len(input))
		score += math.Pow(freq-v, 2)
	}
	return math.Sqrt(score)
}
