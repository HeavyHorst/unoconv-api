package main

import (
	"bytes"
	"io"

	"github.com/qiniu/iconv"
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

func toUTF8(file []byte, tempfile io.Writer) (string, error) {
	charset, err := getFileEncoding(file)
	if err != nil {
		return "", err
	}

	cd, err := iconv.Open("utf-8", charset)
	if err != nil {
		return charset, err
	}
	defer cd.Close()

	bufSize := 0
	reader := iconv.NewReader(cd, bytes.NewBuffer(file), bufSize)

	io.Copy(tempfile, reader)
	return charset, nil
}
