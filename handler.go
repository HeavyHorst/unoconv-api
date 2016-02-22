package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/rs/xlog"
	"github.com/rs/xmux"
	"golang.org/x/net/context"
)

func HealthHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	// just return http.StatusOK
	l := xlog.FromContext(ctx)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("200 -OK"))
}

func unoconvHandler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	l := xlog.FromContext(ctx)

	//The whole request body is parsed and up to a total of 34MB bytes of its file parts are stored in memory,
	//with the remainder stored on disk in temporary files.
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		l.Error(err)
		return
	}
	defer file.Close()

	//add the filename to access log
	l.SetField("filename", handler.Filename)

	//create a temporary file and copy the file from the form to it
	tempfile, err := ioutil.TempFile(os.TempDir(), "unoconv-api")
	if err != nil {
		l.Error(err)
		return
	}

	switch filepath.Ext(handler.Filename) {
	case ".txt":
		//read the files content
		data, err := ioutil.ReadAll(file)
		if err != nil {
			l.Error(err)
			return
		}

		//try to convert the textfile (data) to utf-8 and write it to tempfile
		charset, err := toUTF8(data, tempfile)
		l.SetField("charset", charset)
		l.SetField("convertedToUTF8", true)
		if err != nil {
			//Could not convert to utf-8, write the original data to tempfile
			l.Error(err)
			l.SetField("convertedToUTF8", false)
			io.Copy(tempfile, bytes.NewBuffer(data))
		}
	default:
		io.Copy(tempfile, file)
	}

	tempfile.Close()

	//append the file extension to the temporary file's name
	filename := tempfile.Name() + filepath.Ext(handler.Filename)
	os.Rename(tempfile.Name(), filename)
	defer os.Remove(filename)

	//Run unoconv to convert the file
	//unoconv's stdout is plugged directly to the httpResponseWriter
	err = uno.convert(filename, xmux.Param(ctx, "filetype"), w)
	if err != nil {
		l.Error(err)
		return
	}
}
