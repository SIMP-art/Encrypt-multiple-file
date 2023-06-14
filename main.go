package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"
)

func EnC(filename string, key []byte, wg *sync.WaitGroup) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("Error creating cipher block: %v", err)
		return
	}
	paddingLen := aes.BlockSize - (len(data) % aes.BlockSize)
	padding := make([]byte, paddingLen)
	for i := range padding {
		padding[i] = byte(paddingLen)
	}
	data = append(data, padding...)
	stream := cipher.NewCTR(block, make([]byte, aes.BlockSize))
	stream.XORKeyStream(data, data)
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v", err)
		return
	}

	fmt.Printf("Encryption complete: %s\n", filename)
	defer wg.Done()
}

func DeC(filename string, key []byte, wg *sync.WaitGroup) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file: %v", err)
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Printf("Error creating cipher block: %v", err)
		return
	}
	stream := cipher.NewCTR(block, make([]byte, aes.BlockSize))
	stream.XORKeyStream(data, data)
	paddingLen := int(data[len(data)-1])
	data = data[:len(data)-paddingLen]
	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v", err)
		return
	}

	fmt.Printf("Decryption complete: %s\n", filename)
	defer wg.Done()
}
func rep(s string) string { //make string length to 16 for enc and dec. shitty algorithm btw
	first := string(s[0])
	last := string(s[len(s)-1])
	for len(s) < 16 {
		s = first + s + last
	}
	return s[:16]
}
func main() {
	if len(os.Args) < 4 {
		log.Fatalf("Usage:%s <encrypt[1]/decrypt[2]> <key> <file1> [file2] [file3] ...", os.Args[0])
	}
	mode := os.Args[1]
	key := os.Args[2]
	files := os.Args[3:]
	numFiles := len(files)
	var wg sync.WaitGroup
	wg.Add(numFiles)
	for _, file := range files {
		switch mode {
		case "encrypt":
			go EnC(file, []byte(rep(key)), &wg)
		case "decrypt":
			go DeC(file, []byte(rep(key)), &wg)
		case "1":
			go EnC(file, []byte(rep(key)), &wg)
		case "2":
			go DeC(file, []byte(rep(key)), &wg)
		default:
			log.Fatal("Invalid mode.")
		}
	}
	wg.Wait()
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
