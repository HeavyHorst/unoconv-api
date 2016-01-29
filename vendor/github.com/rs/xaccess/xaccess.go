// Package xaccess is a middleware that logs all access requests performed on the sub handler
// using github.com/rs/xlog and github.com/rs/xstats stored in context.
package xaccess

import (
	"net/http"
	"strconv"
	"time"

	"github.com/rs/xhandler"
	"github.com/rs/xlog"
	"github.com/rs/xstats"
	"github.com/zenazn/goji/web/mutil"
	"golang.org/x/net/context"
)

// NewHandler returns a handler that log access information about each request performed
// on the provided sub-handlers. Uses context's github.com/rs/xlog and
// github.com/rs/xstats if present for logging.
func NewHandler() func(next xhandler.HandlerC) xhandler.HandlerC {
	return func(next xhandler.HandlerC) xhandler.HandlerC {
		return xhandler.HandlerFuncC(func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			// Time request
			reqStart := time.Now()

			// Sniff the status and content size for logging
			lw := mutil.WrapWriter(w)

			// Call the next handler
			next.ServeHTTPC(ctx, lw, r)

			// Conpute request duration
			reqDur := time.Since(reqStart)

			// Get request status
			status := responseStatus(ctx, lw.Status())

			// Log request stats
			sts := xstats.FromContext(ctx)
			tags := []string{
				"status:" + status,
				"status_code:" + strconv.Itoa(lw.Status()),
			}
			sts.Timing("request_time", reqDur, tags...)
			sts.Histogram("request_size", float64(lw.BytesWritten()), tags...)

			// Log access info
			log := xlog.FromContext(ctx)
			log.Infof("%s %s %03d", r.Method, ellipsize(r.URL.String(), 100), lw.Status(), xlog.F{
				"method":      r.Method,
				"uri":         r.URL.String(),
				"type":        "access",
				"status":      status,
				"status_code": lw.Status(),
				"duration":    reqDur.Seconds(),
				"size":        lw.BytesWritten(),
			})
		})
	}
}

func responseStatus(ctx context.Context, statusCode int) string {
	if ctx.Err() != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "timeout"
		}
		return "canceled"
	} else if statusCode >= 200 && statusCode < 300 {
		return "ok"
	}
	return "error"
}

// ellipsize shorten a string using ellises in the middle if the string
// is longer than max.
func ellipsize(s string, max int) string {
	if max <= 3 {
		s = "..."[:max]
	} else if l := len(s); l > max {
		s = s[:max/2-1] + "..." + s[l-(max/2)+1:]
	}
	return s
}
