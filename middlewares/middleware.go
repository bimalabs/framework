package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/bimalabs/framework/v4/loggers"

	"github.com/CAFxX/httpcompression"
	"github.com/CAFxX/httpcompression/contrib/andybalholm/brotli"
	"github.com/CAFxX/httpcompression/contrib/compress/zlib"
	"github.com/CAFxX/httpcompression/contrib/klauspost/pgzip"
)

type (
	Middleware interface {
		Attach(request *http.Request, response http.ResponseWriter) bool
		Priority() int
	}

	Factory struct {
		Debug       bool
		Middlewares []Middleware
	}

	responseWrapper struct {
		http.ResponseWriter
		statusCode int
	}
)

func (w *responseWrapper) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *responseWrapper) StatusCode() int {
	if w.statusCode == 0 {
		return http.StatusOK
	}

	return w.statusCode
}

func (m *Factory) Register(middlewares []Middleware) {
	for _, v := range middlewares {
		m.Add(v)
	}
}

func (m *Factory) Add(middlware Middleware) {
	m.Middlewares = append(m.Middlewares, middlware)
}

func (m *Factory) Sort() {
	sort.Slice(m.Middlewares, func(i, j int) bool {
		return m.Middlewares[i].Priority() > m.Middlewares[j].Priority()
	})
}

func (m *Factory) Attach(handler http.Handler) http.Handler {
	ctx := context.WithValue(context.Background(), "scope", "middleware")
	internal := http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		start := time.Now()
		if !m.Debug {
			for _, middleware := range m.Middlewares {
				if stop := middleware.Attach(request, response); stop {
					return
				}
			}

			handler.ServeHTTP(response, request)

			elapsed := time.Since(start)

			var execution strings.Builder
			execution.WriteString("Execution time: ")
			execution.WriteString(elapsed.String())

			fmt.Println(execution.String())

			return
		}

		wrapper := responseWrapper{ResponseWriter: response}
		for _, middleware := range m.Middlewares {
			if stop := middleware.Attach(request, response); stop {
				var stopper strings.Builder
				stopper.WriteString("Middleware stopped by: ")
				stopper.WriteString(reflect.TypeOf(middleware).Elem().Name())

				loggers.Logger.Debug(ctx, stopper.String())

				return
			}
		}

		handler.ServeHTTP(&wrapper, request)

		elapsed := time.Since(start)

		uri, _ := url.QueryUnescape(request.RequestURI)

		var stdLog strings.Builder
		stdLog.WriteString("[")
		stdLog.WriteString(request.Method)
		stdLog.WriteString("]")
		stdLog.WriteString("\t")
		stdLog.WriteString(strconv.Itoa(wrapper.StatusCode()))
		stdLog.WriteString("\t")
		stdLog.WriteString(elapsed.String())
		stdLog.WriteString("\t")
		stdLog.WriteString(uri)

		fmt.Println(stdLog.String())
	})

	deflateEncoder, _ := zlib.New(zlib.Options{})
	brotliEncoder, _ := brotli.New(brotli.Options{})
	gzipEncoder, _ := pgzip.New(pgzip.Options{
		Level:     pgzip.DefaultCompression,
		BlockSize: 1 << 20,
		Blocks:    4,
	})

	compress, _ := httpcompression.Adapter(
		httpcompression.Compressor(brotli.Encoding, 2, brotliEncoder),
		httpcompression.Compressor(pgzip.Encoding, 1, gzipEncoder),
		httpcompression.Compressor(zlib.Encoding, 0, deflateEncoder),
		httpcompression.Prefer(httpcompression.PreferClient),
		httpcompression.MinSize(165),
	)

	return compress(internal)
}
