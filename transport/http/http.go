package http

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/MouseHatGames/mice/logger"
	"github.com/MouseHatGames/mice/options"
	"github.com/MouseHatGames/mice/transport"
)

const headerPrefix = "X-Mice-"

type httpTransport struct {
	log logger.Logger
}

func Transport() options.Option {
	return func(o *options.Options) {
		o.Transport = &httpTransport{
			log: o.Logger.GetLogger("http"),
		}
	}
}

func (t *httpTransport) Listen(ctx context.Context, addr string) (transport.Listener, error) {
	t.log.Infof("listening on %s", addr)

	return &httpListener{addr, t.log}, nil
}

func (t *httpTransport) Dial(ctx context.Context, addr string) (transport.Socket, error) {
	t.log.Infof("dialing %s", addr)

	return &httpOutgoingSocket{
		address: addr,
		resp:    make(chan *http.Response, 1),
	}, nil
}

type httpListener struct {
	addr string
	log  logger.Logger
}

func (l *httpListener) Close() error {
	return nil
}

func (l *httpListener) Accept(ctx context.Context, fn func(transport.Socket)) error {
	handler := http.NewServeMux()
	handler.HandleFunc("/request", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Add("Access-Control-Allow-Origin", "*")
		rw.Header().Add("Access-Control-Allow-Methods", "POST")
		rw.Header().Add("Access-Control-Allow-Headers", "*")

		if r.Method != http.MethodPost {
			return
		}

		l.log.Debugf("got request from %s", r.RemoteAddr)

		close := make(chan struct{}, 1)

		fn(&httpIncomingSocket{
			rw:     rw,
			r:      r,
			log:    l.log,
			closer: close,
		})

		<-close
	})

	l.log.Debugf("accepting connections")
	return http.ListenAndServe(l.addr, handler)
}

func getMiceHeaders(h http.Header) (mh map[string]string) {
	mh = make(map[string]string)

	for k, v := range h {
		if strings.HasPrefix(k, headerPrefix) {
			name := strings.ToLower(strings.TrimPrefix(k, headerPrefix))

			mh[name] = v[0]
		}
	}

	return
}

func unmarshalMessage(r io.Reader, msg *transport.Message) error {
	dec := json.NewDecoder(r)

	if err := dec.Decode(msg); err != nil {
		return fmt.Errorf("parse json: %w", err)
	}

	return nil
}

func marshalMessage(w io.Writer, msg *transport.Message) error {
	enc := json.NewEncoder(w)

	return enc.Encode(msg)
}
