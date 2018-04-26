// Author: Antoine Mercadal
// See LICENSE file for full LICENSE
// Copyright 2016 Aporeto.

package bahamut

import (
	"crypto/tls"
	"crypto/x509"
	"net/http"
	"time"

	"github.com/aporeto-inc/elemental"
)

// HealthServerFunc is the type used by the Health Server to check the health of the server
type HealthServerFunc func() error

// A Config represents the configuration of Bahamut.
type Config struct {

	// General configuration.
	General struct {

		// Set this to false to disable panic recovery.
		PanicRecoveryDisabled bool
	}

	// ReSTServer contains the configuration for the ReST Server.
	ReSTServer struct {

		// ListenAddress is the general listening address for the API server as
		// well as the PushServer.
		ListenAddress string

		// ReadTimeout defines the read http timeout.
		ReadTimeout time.Duration

		// WriteTimeout defines the write http timeout.
		WriteTimeout time.Duration

		// WriteTimeout defines the idle http timeout.
		IdleTimeout time.Duration

		// DisableKeepalive controls if the ReSTServer should have keepalive activated or not.
		// There is a bug in Go <= 1.7 which makes the server eats all available fd, so DisableKeepalive should
		// be set to true if you are using those versions.
		DisableKeepalive bool

		// Disabled controls if the ReSTServer should be disabled.
		Disabled bool

		// CustomRootHandlerFunc defines a custom handler func for / API.
		CustomRootHandlerFunc http.HandlerFunc
	}

	// PushServer contains the configuration for the Push Server.
	PushServer struct {

		// Service defines the pubsub server to use.
		Service PubSubServer

		// Topic defines the default notification topic to use.
		Topic string

		// DispatchHandler defines the handler that will be used to
		// decide if a push event should be dispatch to push sessions.
		DispatchHandler PushDispatchHandler

		// PublishHandler defines the handler that will be used to
		// decide if an event should be published.
		PublishHandler PushPublishHandler

		// Disabled defines if the the entire websocket server should be disabled.
		// If you set this to true, other options has no effect.
		Disabled bool

		// PublishDisabled disables the publication of events.
		// If PushDisabled is false, this has no incidence.
		PublishDisabled bool

		// DispatchDisabled disables the dispatching of events.
		// If PushDisabled is false, this has no incidence.
		DispatchDisabled bool
	}

	MockServer struct {

		// ListenAddress is the general listening address for the mock server.
		ListenAddress string

		// ReadTimeout defines the read http timeout.
		ReadTimeout time.Duration

		// WriteTimeout defines the write http timeout.
		WriteTimeout time.Duration

		// WriteTimeout defines the idle http timeout.
		IdleTimeout time.Duration

		// Enabled defines if the mock server should be enabled.
		Enabled bool
	}

	// HealthServer contains the configuration for the Health Server.
	HealthServer struct {

		// ListenAddress is the general listening address for the health server.
		ListenAddress string

		// HealthHandler is the type of the function to run to determine the health of the server.
		HealthHandler HealthServerFunc

		// ReadTimeout defines the read http timeout.
		ReadTimeout time.Duration

		// WriteTimeout defines the write http timeout.
		WriteTimeout time.Duration

		// WriteTimeout defines the idle http timeout.
		IdleTimeout time.Duration

		// Disabled defines if the health server should be disabled.
		Disabled bool
	}

	// ProfilingServer contains information about profiling server.
	ProfilingServer struct {

		// ListenAddress is the general listening address for the profiling server.
		// Only matters when mode is "gops".
		ListenAddress string

		// Enabled defines if the profiling server should be enabled.
		Enabled bool

		// Mode represents the mode of the profiling server to run.
		// If can be "gops" or "gcp"
		Mode string

		// Name of the project to report when Mode is set to "gcp".
		GCPProjectID string

		// Set this to add a prefix to your service name when reporting
		// profile to GCP. This allows to differentiate multiple instance
		// of an application running in the same project.
		GCPServicePrefix string
	}

	// TLS contains the TLS configuration.
	TLS struct {

		// RootCAPool is the *x509.CertPool to use for the secure bahamut api server.
		RootCAPool *x509.CertPool

		// ClientCAPool is the *x509.CertPool to use for the authentifying client.
		ClientCAPool *x509.CertPool

		// AuthType defines the tls authentication mode to use for a secure server.
		AuthType tls.ClientAuthType

		// ServerCertificates are the TLS certficates to use for the secure api server.
		// If you set ServerCertificatesRetrieverFunc, the value of ServerCertificates will be ignored.
		ServerCertificates []tls.Certificate

		// ServerCertificatesRetrieverFunc is standard tls GetCertifcate function to use to
		// retrieve the server certificates dynamically.
		// - If you set this, the value of ServerCertificates will be ignored.
		// - If EnableLetsEncrypt is set, this will be ignored
		ServerCertificatesRetrieverFunc func(*tls.ClientHelloInfo) (*tls.Certificate, error)

		// EnableLetsEncrypt defines if the server should get a certificate from letsencrypt automagically.
		EnableLetsEncrypt bool

		// LetsEncryptDomainWhiteList contains the list of white listed domain name to use for
		// issuing certificates.
		LetsEncryptDomainWhiteList []string

		// LetsEncryptCertificateCacheFolder gives the path where to store certificate cache.
		// If empty, the default temp folder of the machine will be used.
		LetsEncryptCertificateCacheFolder string
	}

	// Security contains the Authenticator and Authorizer.
	Security struct {

		// RequestAuthenticators defines the list the RequestAuthenticator to use to authenticate the requests.
		// They are executed in order from index 0 to index n. They will return a bahamut.AuthAction to tell if
		// the current request authenticator grants, denies or let the chain continue. If an error is returned, the
		// chain fails immediately.
		RequestAuthenticators []RequestAuthenticator

		// SessionAuthenticators defines the list of SessionAuthenticator that will be used to
		// initially authentify a websocket connection.
		// They are executed in order from index 0 to index n.They will return a bahamut.AuthAction to tell if
		// the current session authenticator grants, denies or let the chain continue. If an error is returned, the
		// chain fails immediately.
		SessionAuthenticators []SessionAuthenticator

		// Authorizers defines the list Authorizers to use to authorize the requests.
		// They are executed in order from index 0 to index n. They will return a bahamut.AuthAction to tell if
		// the current authorizer grants, denies or let the chain continue. If an error is returned, the
		// chain fails immediately.
		Authorizers []Authorizer

		// Auditer is the Auditer to use to audit the requests.
		// The Audit() method will be run in a go routinel so there is no
		// need to deal with it in the implementation.
		Auditer Auditer
	}

	RateLimiting struct {

		// RateLimiter is the RateLimiter to use eventually limit the rate of some calls.
		RateLimiter RateLimiter
	}

	// Model contains the model configuration.
	Model struct {

		// IdentifiablesFactory is a function that returns a instance of a model
		// according to its identity.
		IdentifiablesFactory elemental.IdentifiableFactory

		// RelationshipsRegistry contains each elemental model RelationshipsRegistry for each version.
		RelationshipsRegistry map[int]elemental.RelationshipsRegistry

		// If ReadOnly is set to true, all write operations will return a Locked HTTP Code (423)
		// This is useful during maintenance.
		ReadOnly bool

		// If ReadOnly is aset to true, this will bypass the readonly mode for the set identities.
		ReadOnlyExcludedIdentities []elemental.Identity

		// Unmarshallers contains a list of custom umarshaller per identity.
		// This allows to create custom function to umarshal the payload of a request.
		// If none is provided for a particular identity, the standard unmarshal function
		// is used.
		Unmarshallers map[elemental.Identity]CustomUmarshaller
	}

	// Meta contains information about the meta apis.
	Meta struct {

		// ServiceName contains the name of the service.
		ServiceName string

		// ServiceVersion contains the version of the service itself.
		ServiceVersion string

		// Version should contain information relative to the service version.
		// like all it's libraries and things like that.
		Version map[string]interface{}

		DisableMetaRoute bool
	}
}
