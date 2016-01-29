package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/xlog"
	"github.com/rs/xmux"
	"golang.org/x/net/context"
)

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
	io.Copy(tempfile, file)
	tempfile.Close()

	//append the file extension to the temporary file's name
	filename := tempfile.Name() + filepath.Ext(handler.Filename)
	os.Rename(tempfile.Name(), filename)
	defer os.Remove(filename)

	//Run unoconv to convert the file
	//unoconv's stdout is plugged directly to the httpResponseWriter
	cmd := exec.Command("unoconv", "-f", xmux.Param(ctx, "filetype"), "--stdout", filename)
	cmd.Stdout = w
	err = cmd.Run()
	if err != nil {
		l.Error(err)
		return
	}
}
