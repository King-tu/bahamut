package bahamut

import (
	"errors"
	"net/http"
	"testing"

	"github.com/aporeto-inc/elemental/test/model"

	"github.com/aporeto-inc/elemental"
	"github.com/go-zoo/bone"

	. "github.com/smartystreets/goconvey/convey"
)

type mockPubSubServer struct {
	publications []*Publication
	PublishErr   error
}

func (p *mockPubSubServer) Connect() Waiter   { return nil }
func (p *mockPubSubServer) Disconnect() error { return nil }

func (p *mockPubSubServer) Publish(publication *Publication) error {
	p.publications = append(p.publications, publication)
	return p.PublishErr
}

func (p *mockPubSubServer) Subscribe(pubs chan *Publication, errors chan error, topic string, args ...interface{}) func() {
	return nil
}

type mockSessionAuthenticator struct {
	action AuthAction
	err    error
}

func (a *mockSessionAuthenticator) AuthenticateSession(Session) (AuthAction, error) {
	return a.action, a.err
}

type mockSessionHandler struct {
	onPushSessionInitCalled  int
	onPushSessionInitOK      bool
	onPushSessionInitErr     error
	onPushSessionStartCalled int
	onPushSessionStopCalled  int
	shouldPublishCalled      int
	shouldPublishOK          bool
	shouldPublishErr         error
	shouldPushCalled         int
	shouldPushOK             bool
	shouldPushErr            error
}

func (h *mockSessionHandler) OnPushSessionInit(PushSession) (bool, error) {
	h.onPushSessionInitCalled++
	return h.onPushSessionInitOK, h.onPushSessionInitErr
}

func (h *mockSessionHandler) OnPushSessionStart(PushSession) {
	h.onPushSessionStartCalled++
}

func (h *mockSessionHandler) OnPushSessionStop(PushSession) {
	h.onPushSessionStopCalled++
}

func (h *mockSessionHandler) ShouldPublish(*elemental.Event) (bool, error) {
	h.shouldPublishCalled++
	return h.shouldPublishOK, h.shouldPublishErr
}

func (h *mockSessionHandler) ShouldPush(PushSession, *elemental.Event) (bool, error) {
	h.shouldPushCalled++
	return h.shouldPushOK, h.shouldPushErr
}

func TestWebsocketServer_newWebsocketServer(t *testing.T) {

	Convey("Given I have a processor finder", t, func() {

		pf := func(identity elemental.Identity) (Processor, error) {
			return struct{}{}, nil
		}

		Convey("When I create a new websocket server with push and wsapi", func() {

			mux := bone.New()
			cfg := Config{}
			wss := newWebsocketServer(cfg, mux, pf)

			Convey("Then the websocket sever should be correctly initialized", func() {
				So(wss.sessions, ShouldResemble, map[string]internalWSSession{})
				So(wss.multiplexer, ShouldEqual, mux)
				So(wss.config, ShouldResemble, cfg)
				So(wss.processorFinder, ShouldEqual, pf)
			})

			Convey("Then the handlers should be installed in the mux", func() {
				So(len(mux.Routes), ShouldEqual, 1)
				So(len(mux.Routes["GET"]), ShouldEqual, 2)
				So(mux.Routes["GET"][0].Path, ShouldEqual, "/events")
				So(mux.Routes["GET"][1].Path, ShouldEqual, "/wsapi")
			})
		})

		Convey("When I create a new websocket server with push disabled", func() {

			mux := bone.New()
			cfg := Config{}
			cfg.WebSocketServer.PushDisabled = true
			_ = newWebsocketServer(cfg, mux, pf)

			Convey("Then the handlers should be installed in the mux", func() {
				So(len(mux.Routes), ShouldEqual, 1)
				So(len(mux.Routes["GET"]), ShouldEqual, 1)
				So(mux.Routes["GET"][0].Path, ShouldEqual, "/wsapi")
			})
		})

		Convey("When I create a new websocket server with api disabled", func() {

			mux := bone.New()
			cfg := Config{}
			cfg.WebSocketServer.APIDisabled = true
			_ = newWebsocketServer(cfg, mux, pf)

			Convey("Then the handlers should be installed in the mux", func() {
				So(len(mux.Routes), ShouldEqual, 1)
				So(len(mux.Routes["GET"]), ShouldEqual, 1)
				So(mux.Routes["GET"][0].Path, ShouldEqual, "/events")
			})
		})

		Convey("When I create a new websocket server with everything disabled", func() {

			mux := bone.New()
			cfg := Config{}
			cfg.WebSocketServer.PushDisabled = true
			cfg.WebSocketServer.APIDisabled = true
			_ = newWebsocketServer(cfg, mux, pf)

			Convey("Then the handlers should be installed in the mux", func() {
				So(len(mux.Routes), ShouldEqual, 0)
			})
		})
	})
}

func TestWebsockerServer_SessionRegistration(t *testing.T) {

	Convey("Given I have a websocket server", t, func() {

		pf := func(identity elemental.Identity) (Processor, error) {
			return struct{}{}, nil
		}

		req, _ := http.NewRequest("GET", "bla", nil)
		mux := bone.New()
		cfg := Config{}
		h := &mockSessionHandler{}
		cfg.WebSocketServer.SessionsHandler = h

		wss := newWebsocketServer(cfg, mux, pf)

		Convey("When I register a valid push session", func() {

			s := newWSPushSession(req, cfg, nil)
			wss.registerSession(s)

			Convey("Then the session should correctly registered", func() {
				So(len(wss.sessions), ShouldEqual, 1)
				So(wss.sessions[s.Identifier()], ShouldEqual, s)
			})

			Convey("Then handler.onPushSessionStart should have been called", func() {
				So(h.onPushSessionStartCalled, ShouldEqual, 1)
			})

			Convey("When I unregister it", func() {

				wss.unregisterSession(s)

				Convey("Then the session should correctly unregistered", func() {
					So(len(wss.sessions), ShouldEqual, 0)
				})

				Convey("Then handler.onPushSessionStop should have been called", func() {
					So(h.onPushSessionStopCalled, ShouldEqual, 1)
				})
			})
		})

		Convey("When I register a valid session that is not a push session", func() {

			s := newWSSession(req, cfg, nil)
			wss.registerSession(s)

			Convey("Then the session should correctly registered", func() {
				So(len(wss.sessions), ShouldEqual, 1)
				So(wss.sessions[s.Identifier()], ShouldEqual, s)
			})

			Convey("Then handler.onPushSessionStart should have been called", func() {
				So(h.onPushSessionStartCalled, ShouldEqual, 0)
			})

			Convey("When I unregister it", func() {

				wss.unregisterSession(s)

				Convey("Then the session should correctly unregistered", func() {
					So(len(wss.sessions), ShouldEqual, 0)
				})

				Convey("Then handler.onPushSessionStop should have been called", func() {
					So(h.onPushSessionStopCalled, ShouldEqual, 0)
				})
			})
		})

		Convey("When I register a valid session with no id", func() {

			s := &wsSession{}

			Convey("Then it should panic", func() {
				So(func() { wss.registerSession(s) }, ShouldPanicWith, "cannot register websocket session. empty identifier")
			})
		})

		Convey("When I unregister a valid session with no id", func() {

			s := &wsSession{}

			Convey("Then it should panic", func() {
				So(func() { wss.unregisterSession(s) }, ShouldPanicWith, "cannot unregister websocket session. empty identifier")
			})
		})
	})
}

func TestWebsocketServer_authSession(t *testing.T) {

	Convey("Given I have a websocket server", t, func() {

		pf := func(identity elemental.Identity) (Processor, error) {
			return struct{}{}, nil
		}

		req, _ := http.NewRequest("GET", "bla", nil)
		mux := bone.New()

		Convey("When I call authSession on when there is no authenticator configured", func() {

			cfg := Config{}

			wss := newWebsocketServer(cfg, mux, pf)

			s := newWSSession(req, cfg, nil)
			err := wss.authSession(s)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I call authSession with a configured authenticator that is ok", func() {

			a := &mockSessionAuthenticator{}
			a.action = AuthActionOK

			cfg := Config{}
			cfg.Security.SessionAuthenticators = []SessionAuthenticator{a}

			wss := newWebsocketServer(cfg, mux, pf)

			s := newWSSession(req, cfg, nil)
			err := wss.authSession(s)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I call authSession with a configured authenticator that is not ok", func() {

			a := &mockSessionAuthenticator{}
			a.action = AuthActionKO

			cfg := Config{}
			cfg.Security.SessionAuthenticators = []SessionAuthenticator{a}

			wss := newWebsocketServer(cfg, mux, pf)

			s := newWSSession(req, cfg, nil)
			err := wss.authSession(s)

			Convey("Then err should not be nil", func() {
				So(err.Error(), ShouldEqual, "error 401 (bahamut): Unauthorized: You are not authorized to start a session")
			})
		})

		Convey("When I call authSession with a configured authenticator that returns an error", func() {

			a := &mockSessionAuthenticator{}
			a.action = AuthActionOK // we wan't to check that error takes precedence
			a.err = errors.New("nope")

			cfg := Config{}
			cfg.Security.SessionAuthenticators = []SessionAuthenticator{a}

			wss := newWebsocketServer(cfg, mux, pf)

			s := newWSSession(req, cfg, nil)
			err := wss.authSession(s)

			Convey("Then err should not be nil", func() {
				So(err.Error(), ShouldEqual, "error 401 (bahamut): Unauthorized: nope")
			})
		})
	})
}

func TestWebsocketServer_initPushSession(t *testing.T) {

	Convey("Given I have a websocket server", t, func() {

		pf := func(identity elemental.Identity) (Processor, error) {
			return struct{}{}, nil
		}

		req, _ := http.NewRequest("GET", "bla", nil)
		mux := bone.New()

		Convey("When I call initSession on when there is no session handler configured", func() {

			cfg := Config{}

			wss := newWebsocketServer(cfg, mux, pf)

			s := newWSPushSession(req, cfg, nil)
			err := wss.initPushSession(s)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I call initSession on when there is a session handler that is ok", func() {

			h := &mockSessionHandler{}
			h.onPushSessionInitOK = true

			cfg := Config{}
			cfg.WebSocketServer.SessionsHandler = h

			wss := newWebsocketServer(cfg, mux, pf)

			s := newWSPushSession(req, cfg, nil)
			err := wss.initPushSession(s)

			Convey("Then err should be nil", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When I call initSession on when there is a session handler that is not ok", func() {

			h := &mockSessionHandler{}
			h.onPushSessionInitOK = false

			cfg := Config{}
			cfg.WebSocketServer.SessionsHandler = h

			wss := newWebsocketServer(cfg, mux, pf)

			s := newWSPushSession(req, cfg, nil)
			err := wss.initPushSession(s)

			Convey("Then err should not be nil", func() {
				So(err.Error(), ShouldEqual, "error 403 (bahamut): Forbidden: You are not authorized to initiate a push session")
			})
		})

		Convey("When I call initSession on when there is a session handler that returns an error", func() {

			h := &mockSessionHandler{}
			h.onPushSessionInitOK = true // we wan't to check that error takes precedence
			h.onPushSessionInitErr = errors.New("nope")

			cfg := Config{}
			cfg.WebSocketServer.SessionsHandler = h

			wss := newWebsocketServer(cfg, mux, pf)

			s := newWSPushSession(req, cfg, nil)
			err := wss.initPushSession(s)

			Convey("Then err should not be nil", func() {
				So(err.Error(), ShouldEqual, "error 403 (bahamut): Forbidden: nope")
			})
		})
	})
}

func TestWebsocketServer_pushEvents(t *testing.T) {

	Convey("Given I have a websocket server", t, func() {

		pf := func(identity elemental.Identity) (Processor, error) {
			return struct{}{}, nil
		}

		mux := bone.New()

		Convey("When I call pushEvents when no service is configured", func() {

			cfg := Config{}

			wss := newWebsocketServer(cfg, mux, pf)
			wss.pushEvents(nil)

			Convey("Then nothing special should happen", func() {
			})
		})

		Convey("When I call pushEvents with a service is configured but no sessions handler", func() {

			srv := &mockPubSubServer{}

			cfg := Config{}
			cfg.WebSocketServer.Service = srv

			wss := newWebsocketServer(cfg, mux, pf)
			wss.pushEvents(elemental.NewEvent(elemental.EventCreate, testmodel.NewList()))

			Convey("Then I should find one publication", func() {
				So(len(srv.publications), ShouldEqual, 1)
				So(string(srv.publications[0].Data), ShouldStartWith, `{"entity":{"creationOnly":"","date":"0001-01-01T00:00:00Z","description":"","name":"","readOnly":"","slice":null,"ID":"","parentID":"","parentType":""},"identity":"list","type":"create",`)
			})
		})

		Convey("When I call pushEvents with a service is configured and sessions handler that is ok to push", func() {

			srv := &mockPubSubServer{}
			h := &mockSessionHandler{}
			h.shouldPublishOK = true

			cfg := Config{}
			cfg.WebSocketServer.Service = srv
			cfg.WebSocketServer.SessionsHandler = h

			wss := newWebsocketServer(cfg, mux, pf)
			wss.pushEvents(elemental.NewEvent(elemental.EventCreate, testmodel.NewList()))

			Convey("Then I should find one publication", func() {
				So(len(srv.publications), ShouldEqual, 1)
				So(string(srv.publications[0].Data), ShouldStartWith, `{"entity":{"creationOnly":"","date":"0001-01-01T00:00:00Z","description":"","name":"","readOnly":"","slice":null,"ID":"","parentID":"","parentType":""},"identity":"list","type":"create",`)
			})
		})

		Convey("When I call pushEvents with a service is configured and sessions handler that is not ok to push", func() {

			srv := &mockPubSubServer{}
			h := &mockSessionHandler{}
			h.shouldPublishOK = false

			cfg := Config{}
			cfg.WebSocketServer.Service = srv
			cfg.WebSocketServer.SessionsHandler = h

			wss := newWebsocketServer(cfg, mux, pf)
			wss.pushEvents(elemental.NewEvent(elemental.EventCreate, testmodel.NewList()))

			Convey("Then I should find one publication", func() {
				So(len(srv.publications), ShouldEqual, 0)
			})
		})

		Convey("When I call pushEvents with a service is configured and sessions handler that returns an error", func() {

			srv := &mockPubSubServer{}
			h := &mockSessionHandler{}
			h.shouldPublishOK = true // we want to be sure error takes precedence
			h.shouldPublishErr = errors.New("nop")

			cfg := Config{}
			cfg.WebSocketServer.Service = srv
			cfg.WebSocketServer.SessionsHandler = h

			wss := newWebsocketServer(cfg, mux, pf)
			wss.pushEvents(elemental.NewEvent(elemental.EventCreate, testmodel.NewList()))

			Convey("Then I should find one publication", func() {
				So(len(srv.publications), ShouldEqual, 0)
			})
		})
	})
}