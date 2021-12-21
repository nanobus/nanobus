package stream

import (
	"context"
	"fmt"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/go-functions/frames"
	"github.com/nanobus/go-functions/metadata"
)

type Stream interface {
	Metadata() metadata.MD
	RecvData(dst interface{}) error
	SendHeaders(md metadata.MD, end ...bool) error
	SendData(v interface{}, end ...bool) error
	SendUnary(md metadata.MD, v interface{}) error
	SendRequest(path string, v interface{}) error
	SendReply(v interface{}) error
	SendError(err error) error
}

type Streamer struct {
	s     *frames.Stream
	codec functions.Codec
}

func New(stream *frames.Stream, codec functions.Codec) Streamer {
	return Streamer{stream, codec}
}

func (s *Streamer) Metadata() metadata.MD {
	return s.s.Metadata()
}

func (s *Streamer) RecvData(dst interface{}) error {
	msg, err := s.s.RecvData()
	if err != nil {
		return fmt.Errorf("could not receive data: %w", err)
	}

	return s.codec.Decode(msg, dst)
}

func (s *Streamer) SendHeaders(md metadata.MD, end ...bool) error {
	var endVal bool
	if len(end) > 0 {
		endVal = end[0]
	}
	return s.s.SendMetadata(md, endVal)
}

func (s *Streamer) SendData(v interface{}, end ...bool) error {
	var endVal bool
	if len(end) > 0 {
		endVal = end[0]
	}
	var valBytes []byte
	switch v := v.(type) {
	case []byte:
		valBytes = v
	default:
		var err error
		valBytes, err = s.codec.Encode(v)
		if err != nil {
			return fmt.Errorf("could not marshal value to send: %w", err)
		}
	}

	return s.s.SendData(valBytes, endVal)
}

func (s *Streamer) SendUnary(md metadata.MD, v interface{}) error {
	var valBytes []byte
	switch v := v.(type) {
	case nil:
		return s.s.SendMetadata(md, true)
	case []byte:
		valBytes = v
	default:
		var err error
		valBytes, err = s.codec.Encode(v)
		if err != nil {
			return fmt.Errorf("could not marshal value to send: %w", err)
		}
	}

	return s.s.SendUnary(md, valBytes)
}

func (s *Streamer) SendRequest(path string, v interface{}) error {
	return s.SendUnary(metadata.MD{
		":path":        []string{path},
		"content-type": []string{s.codec.ContentType()},
	}, v)
}

func (s *Streamer) SendReply(v interface{}) error {
	return s.SendUnary(metadata.MD{
		":status":      []string{"200"},
		"content-type": []string{s.codec.ContentType()},
	}, v)
}

func (s *Streamer) SendError(err error) error {
	msg := err.Error()
	return s.SendUnary(metadata.MD{
		":status":      []string{"500"},        //strconv.Itoa(e.Status)
		"content-type": []string{"text/plain"}, //s.codec.ContentType()
	}, []byte(msg))
}

type streamKey struct{}

// NewContext creates a new context with incoming `s` attached.
func NewContext(ctx context.Context, s Stream) context.Context {
	return context.WithValue(ctx, streamKey{}, s)
}

func FromContext(ctx context.Context) (s Stream, ok bool) {
	s, ok = ctx.Value(streamKey{}).(Stream)
	return
}
