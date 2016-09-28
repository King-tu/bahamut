// Author: Antoine Mercadal
// See LICENSE file for full LICENSE
// Copyright 2016 Aporeto.

package bahamut

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/aporeto-inc/elemental"
	"github.com/go-zoo/bone"
)

func corsHandler(w http.ResponseWriter, r *http.Request) {
	setCommonHeader(w, r.Referer())
	w.WriteHeader(http.StatusOK)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	WriteHTTPError(w, http.StatusNotFound, elemental.NewError("Not Found", "Unable to find the requested resource", "http", http.StatusNotFound))
}

// an apiServer is the structure serving the api routes.
type apiServer struct {
	config      APIServerConfig
	address     string
	multiplexer *bone.Mux
}

// newAPIServer returns a new apiServer.
func newAPIServer(config APIServerConfig, multiplexer *bone.Mux) *apiServer {

	return &apiServer{
		config:      config,
		multiplexer: multiplexer,
	}
}

// isTLSEnabled checks if the current configuration contains sufficient information
// to estabish of TLS connection.
//
// It basically checks that the TLSCAPath, TLSCertificatePath and TLSKeyPath all correctly defined.
func (a *apiServer) isTLSEnabled() bool {

	return a.config.TLSCAPath != "" && a.config.TLSCertificatePath != "" && a.config.TLSKeyPath != ""
}

// createSecureHTTPServer returns a secure HTTP Server.
//
// It will return an error if any.
func (a *apiServer) createSecureHTTPServer(address string) (*http.Server, error) {

	caCert, err := ioutil.ReadFile(a.config.TLSCAPath)
	if err != nil {
		return nil, err
	}

	CAPool := x509.NewCertPool()
	if !CAPool.AppendCertsFromPEM(caCert) {
		return nil, errors.New("Unable to import CA certificate")
	}

	cert, err := tls.LoadX509KeyPair(a.config.TLSCertificatePath, a.config.TLSKeyPath)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:           []tls.Certificate{cert},
		ClientAuth:             a.config.TLSAuthType,
		ClientCAs:              CAPool,
		RootCAs:                CAPool,
		SessionTicketsDisabled: true,
		MinVersion:             tls.VersionSSL30,
	}

	tlsConfig.BuildNameToCertificate()

	return &http.Server{
		Addr:      address,
		TLSConfig: tlsConfig,
	}, nil
}

// createSecureHTTPServer returns a insecure HTTP Server.
//
// It will return an error if any.
func (a *apiServer) createUnsecureHTTPServer(address string) (*http.Server, error) {

	return &http.Server{
		Addr: address,
	}, nil
}

// installRoutes installs all the routes declared in the APIServerConfig.
func (a *apiServer) installRoutes() {

	for _, route := range a.config.Routes {

		if route.Method == http.MethodHead {
			a.multiplexer.Head(route.Pattern, http.HandlerFunc(route.Handler))
		} else if route.Method == http.MethodGet {
			a.multiplexer.Get(route.Pattern, http.HandlerFunc(route.Handler))
		} else if route.Method == http.MethodPost {
			a.multiplexer.Post(route.Pattern, http.HandlerFunc(route.Handler))
		} else if route.Method == http.MethodPut {
			a.multiplexer.Put(route.Pattern, http.HandlerFunc(route.Handler))
		} else if route.Method == http.MethodDelete {
			a.multiplexer.Delete(route.Pattern, http.HandlerFunc(route.Handler))
		} else if route.Method == http.MethodPatch {
			a.multiplexer.Patch(route.Pattern, http.HandlerFunc(route.Handler))
		}

		log.WithFields(log.Fields{
			"pattern": route.Pattern,
			"method":  route.Method,
			"package": "bahamut",
		}).Debug("API route installed.")
	}

	a.multiplexer.Options("*", http.HandlerFunc(corsHandler))
	a.multiplexer.Get("/", http.HandlerFunc(corsHandler))
	a.multiplexer.NotFound(http.HandlerFunc(notFoundHandler))
}

func (a *apiServer) startProfilingServer() {

	log.WithFields(log.Fields{
		"address": a.config.ProfilingListenAddress,
		"package": "bahamut",
	}).Info("Starting profiling server.")

	srv, err := a.createUnsecureHTTPServer(a.config.ProfilingListenAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"package": "bahamut",
		}).Fatal("Unable to create profiling server.")
	}

	mux := bone.New()
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))

	srv.Handler = mux
	if err := srv.ListenAndServe(); err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"package": "bahamut",
		}).Fatal("Unable to start profiling http server.")
	}
}

func (a *apiServer) startHealthServer() {

	log.WithFields(log.Fields{
		"address":  a.config.HealthListenAddress,
		"endpoint": a.config.HealthEndpoint,
		"package":  "bahamut",
	}).Info("Starting health server.")

	srv, err := a.createUnsecureHTTPServer(a.config.HealthListenAddress)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"package": "bahamut",
		}).Fatal("Unable to create health server.")
	}

	mux := bone.New()
	mux.Get(a.config.HealthEndpoint, a.config.HealthHandler)

	srv.Handler = mux
	if err := srv.ListenAndServe(); err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"package": "bahamut",
		}).Fatal("Unable to start health server.")
	}
}

// start starts the apiServer.
func (a *apiServer) start() {

	if a.config.EnableProfiling {
		go a.startProfilingServer()
	}

	a.installRoutes()

	if a.isTLSEnabled() {

		// if TLS is enabled, we start the dedicated health server.
		if a.config.HealthHandler != nil {
			go a.startHealthServer()
		}

		log.WithFields(log.Fields{
			"address": a.config.ListenAddress,
			"package": "bahamut",
			"routes":  len(a.config.Routes),
		}).Info("Starting secure api server.")

		server, err := a.createSecureHTTPServer(a.config.ListenAddress)
		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"package": "bahamut",
			}).Fatal("Unable to create secure api server.")
		}

		server.Handler = a.multiplexer
		server.SetKeepAlivesEnabled(true)

		if err = server.ListenAndServeTLS("", ""); err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"package": "bahamut",
			}).Fatal("Unable to start secure api server.")
		}

	} else {

		// if TLS is not enabled, we simply add the health handler to the existing server.
		if a.config.HealthHandler != nil {
			a.multiplexer.Get(a.config.HealthEndpoint, a.config.HealthHandler)
		}

		log.WithFields(log.Fields{
			"address": a.config.ListenAddress,
			"package": "bahamut",
			"routes":  len(a.config.Routes),
		}).Info("Starting unsecure api server.")

		server, err := a.createUnsecureHTTPServer(a.config.ListenAddress)
		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"package": "bahamut",
			}).Fatal("Unable to create unsecure api server.")
		}

		server.Handler = a.multiplexer
		if err := server.ListenAndServe(); err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"package": "bahamut",
			}).Fatal("Unable to start unsecure api server.")
		}
	}
}

// stop stops the apiServer.
//
// In reality right now, it does nothing :).
func (a *apiServer) stop() {

}
