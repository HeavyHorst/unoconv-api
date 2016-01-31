package main

import (
	"io"
	"os/exec"
)

type request struct {
	filename string
	filetype string
	w        io.Writer
	errChan  chan error
}

type unoconv struct {
	requestChan chan request
}

func initUnoconv() *unoconv {
	uno := new(unoconv)
	uno.requestChan = make(chan request)

	//unoconv can only process one file at a time
	go func(uno *unoconv) {
		for {
			select {
			case data := <-uno.requestChan:
				cmd := exec.Command("unoconv", "-f", data.filetype, "--stdout", data.filename)
				cmd.Stdout = data.w
				err := cmd.Run()
				if err != nil {
					data.errChan <- err
				} else {
					data.errChan <- nil
				}
			}
		}
	}(uno)
	return uno
}

func (u *unoconv) convert(filename, filetype string, w io.Writer) error {
	err := make(chan error)
	req := request{
		filename,
		filetype,
		w,
		err,
	}

	u.requestChan <- req
	return <-err
}
