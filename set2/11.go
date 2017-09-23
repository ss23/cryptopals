package main

import (
	"encoding/base64"
	"encoding/hex"
	//"io/ioutil"
	//"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	//"crypto/rand"
	"fmt"
	"log"
	"math"
	"math/rand"
	//"os"
	"strings"
	"time"
)

func main() {
	// Set up rng
	rand.Seed(time.Now().UTC().UnixNano())

	// We know the block size is 16 bytes because I wrote it
	bs := 16 // 16 bytes
	plaintext := strings.Repeat("A", bs * 3)
	ciphertext := randomAESEncryption([]byte(plaintext))

	// Split into bs chunks
	chunkedText := make([][]byte, 0)
	for len(ciphertext) > 0 {
		chunkedText = append(chunkedText, ciphertext[:bs])
		ciphertext = ciphertext[bs:]
	}

	// check if any chunks are similar
	if chunkedText[1] == chunkedText[2] {
		fmt.Println("It's ECB")
	} else {
		fmt.Println("It's CBC")
	}

	fmt.Println(chunkedText)
	/*

		// Decode from base64
		var ciphertext []byte
		if ciphertext, err = decodeBase64(string(rawBytes)); err != nil {
			log.Fatal("Failed to decode base64 content: ", err)
		}

		// Create the AES object
		block, err := aes.NewCipher([]byte(key))
		if err != nil {
			log.Fatal("Failed to initialize AES: ", err)
		}

		plaintext := decryptCBC(ciphertext, block, bytes.Repeat([]byte("\x00"), block.BlockSize()))

		fmt.Println(string(plaintext))
		//func decryptCBC(ciphertext []byte, block cipher.Block, iv []byte) []byte { */
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

func padPkcs7(bytesRaw []byte, blocksize int) []byte {
	padding := (blocksize - len(bytesRaw)) % blocksize
	//fmt.Println("padding with bytes: ", padding)
	if padding == 0 {
		return bytesRaw
	}
	//fmt.Println("appending")
	return append(bytesRaw, bytes.Repeat([]byte(string(padding)), padding)...)
}

func decryptCBC(ciphertext []byte, block cipher.Block, iv []byte) []byte {
	bs := block.BlockSize()
	lastBlock := iv

	plaintext := make([]byte, 0)
	for len(ciphertext) > 0 {
		blocktext := make([]byte, bs)
		// xor the last block
		block.Decrypt(blocktext, ciphertext)
		plaintext = append(plaintext, xor(blocktext, lastBlock)...)
		lastBlock = ciphertext[:bs]
		ciphertext = ciphertext[bs:]
	}
	return plaintext
}

func encryptCBC(plaintext []byte, block cipher.Block, iv []byte) []byte {
	bs := block.BlockSize()
	lastBlock := iv

	ciphertext := make([]byte, 0)
	for len(plaintext) > 0 {
		if len(plaintext) < bs {
			plaintext = padPkcs7(plaintext, bs)
		}
		blocktext := make([]byte, bs)
		block.Encrypt(blocktext, xor(plaintext, lastBlock))
		ciphertext = append(ciphertext, blocktext...)
		lastBlock = ciphertext[:bs]
		plaintext = plaintext[bs:]
	}
	return ciphertext
}

func encryptECB(plaintext []byte, block cipher.Block) []byte {
	bs := block.BlockSize()
	ciphertext := make([]byte, 0)

	for len(plaintext) > 0 {
		if len(plaintext) < bs {
			plaintext = padPkcs7(plaintext, bs)
		}
		blocktext := make([]byte, bs)
		block.Encrypt(blocktext, plaintext)
		ciphertext = append(ciphertext, blocktext...)
		plaintext = plaintext[bs:]
	}
	return ciphertext
}


func randomAESEncryption(plaintext []byte) []byte {
	// Generate a random key
	key := make([]byte, 16)
	_, err := rand.Read(key)
	if err != nil {
		log.Fatal("Failed to generat a random AES key: ", err)
	}

	// Construct required AES block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		log.Fatal("Failed to initialize AES: ", err)
	}

	// Append and prepend a random string to the plaintext
	random := []byte(randomString(rand.Intn(6) + 5))
	plaintext = append(random, plaintext...) // random + plaintext
	plaintext = append(plaintext, random...) // plaintext + random

	// Decide whether to use CBC or ECB
	if rand.Intn(2) == 0 {
		// Generate the ciphertext using ECB
		// Create a random IV
		iv := make([]byte, 16)
		_, err := rand.Read(iv)
		if err != nil {
			log.Fatal("Failed to generate random IV for CBC mode: ", err)
		}
		return encryptCBC(plaintext, block, iv)
	} else {
		return encryptECB(plaintext, block)
	}
}

func randomString(length int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
