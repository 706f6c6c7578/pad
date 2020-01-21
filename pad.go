package main

import (
	"bytes"
	"encoding/binary"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
)

//  'pad -p 2048 < input file > output file'
//  'pad -u < input file > output file'

var (
	errNoNeedToPad = errors.New("No need to padding")
	domain         = createDomain()
)

const uint64Bytes = 8

func main() {
	toPad := flag.Int("p", 4096, "size in bytes to do padding")
	toUnPad := flag.Bool("u", false, "specify if we are unpadding")

	flag.Parse()

	if *toUnPad == false {
		w, err := pad(os.Stdin, *toPad)
		if err != nil {
			panic(err)
		}

		if _, err := io.Copy(os.Stdout, w); err != nil {
			panic(err)
		}
		return
	}

	w, err := unpad(os.Stdin)
	if err != nil {
		panic(err)
	}
	if _, err := io.Copy(os.Stdout, w); err != nil {
		panic(err)
	}
}

func createDomain() []byte {
	r := []byte{}
	for i := 'A'; i <= 'Z'; i++ {
		r = append(r, byte(i))
	}
	return r
}

func unpad(r io.Reader) (io.Reader, error) {
	var clean bytes.Buffer

	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("cannot unpad: %w", err)
	}

	padSize, err := littleEndianToUint64(content[(len(content) - uint64Bytes):])
	// fmt.Fprintf(os.Stderr, "%d\n% x\n", padSize, content[(len(content)-uint64Bytes):])

	if _, err := clean.Write(content[0:(uint64(len(content)-uint64Bytes) - padSize)]); err != nil {
		return nil, fmt.Errorf("cannot unpad: %w", err)
	}

	return &clean, nil
}

func pad(r io.Reader, pad int) (io.Reader, error) {
	var w bytes.Buffer

	content, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("cannot pad: %w", err)
	}

	if len(content) > pad {
		return nil, fmt.Errorf("input (%d bytes) > padding (%d bytes): %w", len(content), pad, err)
	}

	padSize := pad - len(content) - uint64Bytes
	padding := make([]byte, padSize+uint64Bytes)
	for i := 0; i < padSize; i++ {
		padding[i] = domain[rand.Intn(len(domain))]
	}

	sizeBuf, err := uint64ToLittleEndian(uint64(padSize))
	if err != nil {
		return nil, err
	}
	for i, v := range sizeBuf {
		padding[padSize+i] = v
	}
	// fmt.Fprintf(os.Stderr, "%d\n% x\n", padSize, padding[padSize:])

	if _, err := w.Write(content); err != nil {
		return nil, fmt.Errorf("cannot pad: %w", err)
	}
	if _, err := w.Write(padding); err != nil {
		return nil, fmt.Errorf("cannot pad: %w", err)
	}

	return &w, nil
}

func uint64ToLittleEndian(n uint64) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, n)
	if err != nil {
		return nil, fmt.Errorf("cannot convert to number to []byte: %w", err)
	}
	return buf.Bytes(), nil
}

func littleEndianToUint64(n []byte) (uint64, error) {
	var r uint64
	buf := bytes.NewReader(n)
	err := binary.Read(buf, binary.LittleEndian, &r)
	if err != nil {
		return 0, fmt.Errorf("cannot convert to []byte to number: %w", err)
	}
	return r, nil
}
