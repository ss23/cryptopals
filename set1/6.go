package main

import (
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	//"bufio"
	"bytes"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Fatal error reading from stdin")
	}

	// Decode from base64
	var bytesRaw []byte
	if bytesRaw, err = decodeBase64(string(bytes)); err != nil {
		log.Fatal("Failed to decode base64 content: ", err)
	}

	// Attempt to calculate the appropriate keysize
	// We could cheat and do a quick version of this, but if we get the keysize wrong, the next steps will be hard to determine the error for
	keysize := 0
	bestScore := float64(100)
	for i := 6; i <= 40; i++ {
		totalDistance := 0
		numChecked := 0 // we've checked 0 pairs

		// Not very DRY, but slightly easier to write at this point
		if (len(bytesRaw) / i) >= 8 {
			// check the 4th set
			totalDistance += hammingDistance(bytesRaw[i*6:i*7], bytesRaw[i*7:i*8])
			numChecked++
		}
		if (len(bytesRaw) / i) >= 6 {
			// check the 4th set
			totalDistance += hammingDistance(bytesRaw[i*4:i*5], bytesRaw[i*5:i*6])
			numChecked++
		}
		if (len(bytesRaw) / i) >= 4 {
			// check the 4th set
			totalDistance += hammingDistance(bytesRaw[i*2:i*3], bytesRaw[i*3:i*4])
			numChecked++
		}
		if (len(bytesRaw) / i) >= 2 {
			// check the 4th set
			totalDistance += hammingDistance(bytesRaw[i*0:i*1], bytesRaw[i*1:i*2])
			numChecked++
		} else {
			log.Fatal("Error attempting to calculate key length: key is longer than len(input) * 2")
		}

		score := (float64(totalDistance) / float64(numChecked)) / float64(i)
		fmt.Printf("Key Length: %v - Average Score: %v\r\n", i, score)

		// Check if it's the best keysize, and if so, note it down
		if bestScore > score {
			keysize = i
			bestScore = score
		}
	}
	//fmt.Println(encodeHex(bytesRaw))

	//keysize = 29
	fmt.Printf("Attempting decryption of repeating xor using keysize %v\r\n", keysize)

	//fmt.Println(hammingDistance([]byte("this is a test"), []byte("wokka wokka!!!")))

	// Chunk the data up into appropriate blocks
	splitBytes := make([][]byte, keysize)
	//for i := range bytesRaw {
	//	e := make(byte, 0, 0)
	//	splitBytes.append(e)
	//}
	for k, v := range bytesRaw {
		splitBytes[k%keysize] = append(splitBytes[k%keysize], v)
	}

	// Calculate the most likely key candidate for each
	secretKey := ""
	for k, v := range splitBytes {
		fmt.Println("Cracking key ", k, " with value ", v)
		secretKey += crackXor(v)
	}

	fmt.Printf("Key was detected as: %v\r\n", secretKey)

	fmt.Printf("Decrypted message:\r\n%v\r\n", string(xor(bytesRaw, []byte(secretKey))))
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

func decodeBase64(src string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(src)
}

// Will do repeating xor, with one being the input, second being the key
func xor(input []byte, key []byte) (dst []byte) {
	// We need to calculate how long the key should be
	if len(input) > len(key) {
		key = bytes.Repeat(key, int(math.Ceil(float64(len(input))/float64(len(key)))))
	}

	dst = make([]byte, len(input))

	for i := 0; i < len(input); i++ {
		dst[i] = input[i] ^ key[i]
	}
	return
}

func scoreByLetterFreq(input []byte) float64 {
	// Start by filtering out non-ASCII character strings
	for i := range input {
		if input[i] > 127 {
			// ascii too high
			return 100
		} else if (input[i] > 14) && (input[i] < 32) {
			// ascii in the baddie range of kind of valid but not really values
			return 100
		} else if input[i] < 10 {
			// a bit too low for us IMO
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

func hammingDistance(input1 []byte, input2 []byte) int {
	if len(input1) != len(input2) {
		log.Fatal("Cannot calculate hamming distance on unequal length arrays")
	}

	distance := 0
	for i := range input1 {
		xor := input1[i] ^ input2[i]
		for x := xor; x > 0; x >>= 1 {
			if int(x&1) == 1 {
				distance++
			}
		}
	}

	return distance
}

func crackXor(bytesRaw []byte) string {
	candidateKey := ""
	topScore := float64(100)
	for i := 0; i < 255; i++ {
		byteArray := make([]byte, 1)
		byteArray[0] = byte(i)
		bytesXord := xor(bytesRaw, byteArray)
		//if scoreByLetterFreq(bytesXord) < 100 {
		//	fmt.Printf("Possible key: %v %x with score %v\r\n", string(i), i, scoreByLetterFreq(bytesXord))
		//}
		//fmt.Println(string(bytesXord))
		if scoreByLetterFreq(bytesXord) < topScore {
			topScore = scoreByLetterFreq(bytesXord)
			candidateKey = string(i)
		}
	}
	return candidateKey
}
