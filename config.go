// Copyright 2019 Aporeto Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//     http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bahamut

import (
	"crypto/tls"
	"crypto/x509"
	"net"
	"net/http"
	"time"

	opentracing "github.com/opentracing/opentracing-go"
	"go.aporeto.io/elemental"
	"golang.org/x/time/rate"
)

// HealthServerFunc is the type used by the Health Server to check the health of the server.
type HealthServerFunc func() error

// HealthStatFunc is the type used by the Health Server to return additional custom health info.
type HealthStatFunc func(http.ResponseWriter, *http.Request)

// TraceCleaner is the type of function that can be used to clean a trace data
// before it is sent to OpenTracing server. You can use this to strip passwords
// or other sensitive data.
type TraceCleaner func(elemental.Identity, []byte) []byte

// A config represents the configuration of Bahamut.
type config struct {
	general struct {
		panicRecoveryDisabled bool
	}

	restServer struct {
		listenAddress         string
		readTimeout           time.Duration
		writeTimeout          time.Duration
		idleTimeout           time.Duration
		disableCompression    bool
		disableKeepalive      bool
		enabled               bool
		customRootHandlerFunc http.HandlerFunc
		customListener        net.Listener
	}

	pushServer struct {
		service         PubSubClient
		topic           string
		endpoint        string
		dispatchHandler PushDispatchHandler
		publishHandler  PushPublishHandler
		enabled         bool
		publishEnabled  bool
		dispatchEnabled bool
	}

	healthServer struct {
		listenAddress  string
		healthHandler  HealthServerFunc
		readTimeout    time.Duration
		writeTimeout   time.Duration
		idleTimeout    time.Duration
		enabled        bool
		customStats    map[string]HealthStatFunc
		metricsManager MetricsManager
	}

	profilingServer struct {
		listenAddress    string
		enabled          bool
		mode             string
		gcpProjectID     string
		gcpServicePrefix string
	}

	tls struct {
		clientCAPool                    *x509.CertPool
		authType                        tls.ClientAuthType
		serverCertificates              []tls.Certificate
		serverCertificatesRetrieverFunc func(*tls.ClientHelloInfo) (*tls.Certificate, error)
	}

	security struct {
		requestAuthenticators []RequestAuthenticator
		sessionAuthenticators []SessionAuthenticator
		authorizers           []Authorizer
		auditer               Auditer
		CORSOrigin            string
	}

	rateLimiting struct {
		rateLimiter *rate.Limiter
	}

	model struct {
		modelManagers              map[int]elemental.ModelManager
		readOnly                   bool
		readOnlyExcludedIdentities []elemental.Identity
		unmarshallers              map[elemental.Identity]CustomUmarshaller
	}

	meta struct {
		serviceName      string
		serviceVersion   string
		version          map[string]interface{}
		disableMetaRoute bool
	}

	opentracing struct {
		tracer             opentracing.Tracer
		excludedIdentities map[string]struct{}
		traceCleaner       TraceCleaner
	}

	hooks struct {
		postStart func(Server) error
		preStop   func(Server) error
	}
}
