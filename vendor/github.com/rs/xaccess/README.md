# Go HTTP Handler Access Log

[![godoc](http://img.shields.io/badge/godoc-reference-blue.svg?style=flat)](https://godoc.org/github.com/rs/xaccess) [![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/xaccess/master/LICENSE) [![Build Status](https://travis-ci.org/rs/xaccess.svg)](https://travis-ci.org/rs/xaccess) [![Coverage](http://gocover.io/_badge/github.com/rs/xaccess)](http://gocover.io/github.com/rs/xaccess)

Package xaccess is a middleware that logs all access requests performed on the sub handler using [xlog](https://github.com/rs/xlog) and [xstats](https://github.com/rs/xstats) stored in context if any.

## Usage

```go
c := xhandler.Chain{}

c.UseC(xlog.NewHandler(xlog.Config{}))
c.UseC(xstats.NewHandler(dogstatsd.New(statsdWriter, flushInterval), tags))

c.UseC(xaccess.NewHandler())

http.Handle("/", c.Handler(xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello World"))
})))

if err := http.ListenAndServe(":8080", nil); err != nil {
    log.Fatal(err)
}
```

## Licenses

All source code is licensed under the [MIT License](https://raw.github.com/rs/xlog/master/LICENSE).
