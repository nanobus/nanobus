package customers

// import (
// 	"context"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"strconv"
// 	"strings"
// 	"sync"

// 	functions "github.com/nanobus/go-functions"
// 	"github.com/nanobus/go-functions/metadata"
// 	"github.com/rsocket/rsocket-go"
// 	"github.com/rsocket/rsocket-go/payload"
// 	"github.com/rsocket/rsocket-go/rx/flux"
// 	"github.com/rsocket/rsocket-go/rx/mono"
// )

// type Server struct {
// 	basePath         string
// 	codec            functions.Codec
// 	handlers         map[string]functions.Handler
// 	statefulHandlers map[string]functions.StatefulHandler
// 	client           rsocket.Client
// }

// type Conn struct {
// 	ctx context.Context
// 	s   *Server
// }

// func NewConnection(basePath string, codec functions.Codec) *Server {
// 	return &Server{
// 		basePath:         basePath,
// 		codec:            codec,
// 		handlers:         make(map[string]functions.Handler),
// 		statefulHandlers: make(map[string]functions.StatefulHandler),
// 	}
// }

// func (s *Server) Register(namespace, operation string, handler functions.Handler) {
// 	s.handlers[namespace+"/"+operation] = handler
// }

// func (s *Server) RegisterStateful(namespace, operation string, handler functions.StatefulHandler) {
// 	s.statefulHandlers[namespace+"/"+operation] = handler
// }

// func (s *Server) Connect(ctx context.Context) error {
// 	// Start a client connection
// 	contentType := s.codec.ContentType()
// 	tp := rsocket.TCPClient().SetHostAndPort("127.0.0.1", 7878).Build() //.UnixClient().SetPath("bus.sock").Build()
// 	var wg sync.WaitGroup
// 	wg.Add(1)
// 	client, err := rsocket.Connect().
// 		OnClose(func(error) {
// 			wg.Done()
// 		}).
// 		DataMimeType(contentType).
// 		MetadataMimeType(contentType).
// 		Acceptor(func(ctx context.Context, socket rsocket.RSocket) rsocket.RSocket {
// 			conn := Conn{ctx, s}
// 			return rsocket.NewAbstractSocket(
// 				rsocket.RequestResponse(conn.requestResponseHandler),
// 			)
// 		}).
// 		Transport(tp).
// 		Start(ctx)
// 	if err != nil {
// 		return err
// 	}

// 	s.client = client

// 	wg.Wait()
// 	return nil
// }

// func (s *Server) InvokeStream(ctx context.Context, namespace, operation string) (functions.Streamer, error) {
// 	path := s.basePath + namespace + "/" + operation
// 	md := metadata.MD{
// 		":path": []string{path},
// 	}

// 	mdBytes, err := s.codec.Encode(md)
// 	if err != nil {
// 		return nil, err
// 	}
// 	pl := payload.New([]byte{0, 0, 0}, mdBytes)

// 	// payloads := make(chan payload.Payload, 100)
// 	// errs := make(chan error, 100)
// 	// sink := &MySink{payloads: payloads, errs: errs}
// 	// sink.Next(pl)
// 	// sink.Complete()
// 	// in := flux.CreateFromChannel(payloads, errs)
// 	// f := s.client.RequestChannel(in)

// 	// c := make(chan payload.Payload, 100)
// 	// e := make(chan error, 1)
// 	// f.DoOnNext(func(pl payload.Payload) error {
// 	// 	if pl != nil {
// 	// 		c <- pl
// 	// 	}
// 	// 	return nil
// 	// }).DoOnError(func(err error) {
// 	// 	e <- err
// 	// }).DoFinally(func(s rx.SignalType) {
// 	// 	close(c)
// 	// 	close(e)
// 	// }).Subscribe(ctx)

// 	c := make(chan payload.Payload, 100)
// 	e := make(chan error, 1)
// 	f := s.client.RequestStream(pl)
// 	f.DoOnNext(func(input payload.Payload) error {
// 		c <- payload.Clone(input)
// 		return nil
// 	}).DoOnError(func(err error) {
// 		e <- err
// 	}).DoOnComplete(func() {
// 		close(c)
// 	}).Subscribe(ctx)
// 	stream := Stream{
// 		server: s,
// 		ctx:    ctx,
// 		s:      nil,
// 		c:      c,
// 		e:      e,
// 	}

// 	//stream.SendMetadata(md, true)
// 	return &stream, nil
// }

// type MySink struct {
// 	payloads chan payload.Payload
// 	errs     chan error
// 	close    sync.Once
// }

// func (m *MySink) Next(v payload.Payload) {
// 	m.payloads <- v
// }

// func (m *MySink) Complete() {
// 	m.close.Do(func() {
// 		close(m.payloads)
// 		close(m.errs)
// 	})
// }

// func (m *MySink) Error(e error) {
// 	m.errs <- e
// }

// type Stream struct {
// 	server *Server
// 	ctx    context.Context
// 	s      flux.Sink
// 	// f      flux.Flux
// 	c <-chan payload.Payload
// 	e <-chan error

// 	sendMd metadata.MD
// 	md     metadata.MD
// }

// func (s *Stream) SendMetadata(md metadata.MD, end ...bool) error {
// 	var endVal bool
// 	if len(end) > 0 {
// 		endVal = end[0]
// 	}
// 	if endVal {
// 		mdBytes, err := s.server.codec.Encode(md)
// 		if err != nil {
// 			return err
// 		}
// 		s.s.Next(payload.New([]byte{}, mdBytes))
// 		s.s.Complete()
// 		return nil
// 	}

// 	s.sendMd = md
// 	return nil
// }

// func (s *Stream) SendData(data []byte, end ...bool) error {
// 	var endVal bool
// 	if len(end) > 0 {
// 		endVal = end[0]
// 	}

// 	var mdBytes []byte
// 	if s.sendMd != nil {
// 		var err error
// 		mdBytes, err = s.server.codec.Encode(s.sendMd)
// 		if err != nil {
// 			return err
// 		}
// 		s.sendMd = nil
// 	}
// 	s.s.Next(payload.New(data, mdBytes))
// 	if endVal {
// 		s.s.Complete()
// 	}
// 	return nil
// }

// func (s *Stream) Close() error {
// 	if s.s != nil {
// 		s.s.Complete()
// 	}
// 	return nil
// }

// type RefCounter interface {
// 	IncRef() int32
// 	Release()
// }

// func (s *Stream) Metadata() metadata.MD { return s.md }
// func (s *Stream) RecvData() ([]byte, error) {
// 	select {
// 	case pl, ok := <-s.c:
// 		if pl == nil || !ok {
// 			return nil, io.EOF
// 		}
// 		if r, ok := pl.(RefCounter); ok {
// 			fmt.Println("INC")
// 			r.IncRef()
// 		}
// 		if payload.Equal(pl, payload.Empty()) {
// 			return nil, io.EOF
// 		}
// 		data := pl.Data()
// 		if mdBytes, ok := pl.Metadata(); ok {
// 			var md metadata.MD
// 			if err := s.server.codec.Decode(mdBytes, &md); err != nil {
// 				return nil, err
// 			}
// 			s.md = md
// 		}
// 		return data, nil
// 	case err := <-s.e:
// 		return nil, err
// 	}
// }

// func (s *Server) Invoke(ctx context.Context, namespace, operation string, data []byte) ([]byte, error) {
// 	path := s.basePath + namespace + "/" + operation
// 	md := metadata.MD{
// 		":path": []string{path},
// 	}
// 	mdBytes, err := s.codec.Encode(md)
// 	if err != nil {
// 		return nil, err
// 	}
// 	resp, err := s.client.RequestResponse(payload.New(data, mdBytes)).Block(ctx)
// 	if err != nil {
// 		return nil, err
// 	}

// 	mdBytes, ok := resp.Metadata()
// 	if !ok {
// 		return nil, nil // TODO
// 	}

// 	var respMD metadata.MD
// 	if err = s.codec.Decode(mdBytes, &respMD); err != nil {
// 		return nil, err
// 	}

// 	statusStr, _ := respMD.Scalar(":status")
// 	status, err := strconv.Atoi(statusStr)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if status/100 != 2 {
// 		return nil, errors.New(string(resp.Data()))
// 	}

// 	return resp.Data(), nil
// }

// func (s *Server) Close() error {
// 	if s.client != nil {
// 		return s.client.Close()
// 	}
// 	return nil
// }

// func (c *Conn) requestResponseHandler(request payload.Payload) mono.Mono {
// 	var path string
// 	if mdBytes, ok := request.Metadata(); ok {
// 		var md metadata.MD
// 		if err := c.s.codec.Decode(mdBytes, &md); err == nil {
// 			path, _ = md.Scalar(":path")
// 		}
// 	}
// 	path = strings.TrimPrefix(path, "/")
// 	parts := strings.Split(path, "/")
// 	data := request.Data()

// 	var response []byte
// 	var err error
// 	switch len(parts) {
// 	case 2:
// 		// Stateless
// 		namespace := parts[0]
// 		operation := parts[1]
// 		h, ok := c.s.handlers[namespace+"/"+operation]
// 		if !ok {
// 			return c.response(metadata.MD{
// 				":status": []string{"404"},
// 			}, []byte("not_found"))
// 		}
// 		response, err = h(c.ctx, data)
// 	case 3:
// 		// Stateful
// 		namespace := parts[0]
// 		id := parts[1]
// 		operation := parts[2]
// 		h, ok := c.s.statefulHandlers[namespace+"/"+operation]
// 		if !ok {
// 			return c.response(metadata.MD{
// 				":status": []string{"404"},
// 			}, []byte("not_found"))
// 		}
// 		response, err = h(c.ctx, id, data)
// 	default:
// 		return c.response(metadata.MD{
// 			":status": []string{"400"},
// 		}, []byte("bad_request"))
// 	}
// 	if err != nil {
// 		return c.response(metadata.MD{
// 			":status": []string{"500"},
// 		}, []byte(err.Error()))
// 	}

// 	return c.response(metadata.MD{
// 		":status": []string{"200"},
// 	}, response)
// }

// func (c *Conn) response(md metadata.MD, data []byte) mono.Mono {
// 	mdBytes, _ := c.s.codec.Encode(md)
// 	return mono.JustOneshot(payload.New(data, mdBytes))
// }
