package pprof

import (
	"fmt"
	"net/http"
	"net/http/pprof"
)

// EnablePProfEndpoints starts pprof endpoints on mux provided.
func EnablePProfEndpoints(mux *http.ServeMux, base string) error {
	mux.HandleFunc(fmt.Sprintf("%s/pprof/", base), pprof.Index)
	mux.HandleFunc(fmt.Sprintf("%s/pprof/cmdline", base), pprof.Cmdline)
	mux.HandleFunc(fmt.Sprintf("%s/pprof/profile", base), pprof.Profile)
	mux.HandleFunc(fmt.Sprintf("%s/pprof/symbol", base), pprof.Symbol)
	mux.HandleFunc(fmt.Sprintf("%s/pprof/trace", base), pprof.Trace)

	mux.Handle(fmt.Sprintf("%s/pprof/goroutine", base), pprof.Handler("goroutine"))
	mux.Handle(fmt.Sprintf("%s/pprof/heap", base), pprof.Handler("heap"))
	mux.Handle(fmt.Sprintf("%s/pprof/threadcreate", base), pprof.Handler("threadcreate"))
	mux.Handle(fmt.Sprintf("%s/pprof/block", base), pprof.Handler("block"))
	mux.Handle(fmt.Sprintf("%s/vars", base), http.DefaultServeMux)

	return nil
}
