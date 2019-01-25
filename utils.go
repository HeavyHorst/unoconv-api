package main

import (
	"io"

	"gopkg.in/iconv.v1"
	"github.com/saintfish/chardet"
)

func getFileEncoding(data []byte) (string, error) {
	detector := chardet.NewTextDetector()
	r, err := detector.DetectBest(data)
	if err != nil {
		return "", err
	}
	return r.Charset, nil
}

func toUTF8(data []byte, tempfile io.Writer) (string, error) {
	charset, err := getFileEncoding(data)
	if err != nil {
		return "", err
	}

	cd, err := iconv.Open("utf-8", charset)
	if err != nil {
		return charset, err
	}
	defer cd.Close()

	autoSync := false // buffered or not
	bufSize := 0 // default if zero
	w := iconv.NewWriter(cd, tempfile, bufSize, autoSync)
	w.Write(data)
	w.Sync()

	return charset, nil
}
