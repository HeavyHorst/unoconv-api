package main

import (
	"os"
	"bytes"
	"io/ioutil"
	"log"
)

var loremIpsum = []byte(`Lorem ipsum dolor sit amet, consectetur adipiscing elit,
sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris
nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in
reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla
pariatur. Excepteur sint occaecat cupidatat non proident, sunt in
culpa qui officia deserunt mollit anim id est laborum.`)

var pdfSignature = []byte("%PDF-1.5\n")

// Try to convert a plain text lorem ipsum to pdf,
// and check if the result appears to be a pdf (ie starts with %PDF-1.5)
func checkUnoconv(uno *unoconv) bool {
	tempfile, err := ioutil.TempFile(os.TempDir(), "watchdog*.txt")
	if err != nil {
		log.Print(err)
		return false
	}
	tempfile.Write(loremIpsum)
	tempfile.Close()
	filename := tempfile.Name()
	defer os.Remove(filename)

	var converted bytes.Buffer
	err = uno.convert(filename, "pdf", &converted)
	if err != nil {
		log.Print(err)
		return false
	}

	convertedHeader := converted.Next(len(pdfSignature))
	return bytes.Equal(convertedHeader, pdfSignature)
}
