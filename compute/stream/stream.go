package stream

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/nanobus/go-functions"
	"github.com/nanobus/go-functions/frames"
	transport_stream "github.com/nanobus/go-functions/transports/stream"
	"go.nanomsg.org/mangos/v3/protocol/pair"

	"github.com/nanobus/nanobus/compute"
	"github.com/nanobus/nanobus/config"
	"github.com/nanobus/nanobus/resolve"
	"github.com/nanobus/nanobus/stream"
)

var ErrInvalidURISyntax = errors.New("invalid invocation URI syntax")

type MuxConfig struct {
	BasePath      string `mapstructure:"basePath"`
	SocketAddress string `mapstructure:"socketAddress"`
}

// Stream is the NamedLoader for the mux compute.
func Stream() (string, compute.Loader) {
	return "stream", StreamLoader
}

func StreamLoader(with interface{}, resolver resolve.ResolveAs) (*compute.Compute, error) {
	c := MuxConfig{
		BasePath:      "/",
		SocketAddress: "ipc://bus.sock",
	}
	if err := config.Decode(with, &c); err != nil {
		return nil, err
	}

	var msgpackcodec functions.Codec
	var busInvoker compute.BusInvoker
	var stateInvoker compute.StateInvoker
	if err := resolve.Resolve(resolver,
		"codec:msgpack", &msgpackcodec,
		"bus:invoker", &busInvoker,
		"state:invoker", &stateInvoker); err != nil {
		return nil, err
	}

	sock, err := pair.NewSocket()
	if err != nil {
		return nil, fmt.Errorf("could not create nanomsg pair: %w", err)
	}
	err = sock.Listen(c.SocketAddress)
	if err != nil {
		return nil, fmt.Errorf("could not listen on address %s: %w", c.SocketAddress, err)
	}
	framer := frames.NewFramer(sock)
	conn := frames.NewConnection("server", framer, frames.ServerStartingStreamID)

	handler := func(ctx context.Context, strm *frames.Stream) {
		s := stream.New(strm, msgpackcodec)
		ctx = stream.NewContext(ctx, &s)

		path, _ := s.Metadata().Scalar(":path")
		path = strings.TrimPrefix(path, "/")
		parts := strings.Split(path, "/")
		if len(parts) != 2 {
			s.SendError(ErrInvalidURISyntax)
			return
		}

		namespace := parts[0]
		function := parts[1]

		lastDot := strings.LastIndexAny(namespace, ".:")
		if lastDot < 0 {
			s.SendError(ErrInvalidURISyntax)
			return
		}
		service := namespace[lastDot+1:]
		namespace = namespace[:lastDot]

		if namespace == "nanobus" && service == "state" {
			type Args struct {
				Namespace string `json:"namespace" msgpack:"namespace"`
				ID        string `json:"id" msgpack:"id"`
				Key       string `json:"key" msgpack:"key"`
			}

			var args Args
			err := s.RecvData(&args)
			if err != nil {
				s.SendError(err)
				return
			}
			output, err := stateInvoker(ctx, args.Namespace, args.ID, args.Key)
			if err != nil {
				s.SendError(err)
				return
			}

			// No need to decode bytes
			s.SendReply(output)
			return
		}

		var input interface{}
		err := s.RecvData(&input)
		if err != nil {
			s.SendError(err)
			return
		}

		output, err := busInvoker(ctx, namespace, service, function, input)
		if err != nil {
			s.SendError(err)
			return
		}

		s.SendReply(output)
	}
	conn.SetHandler(handler)

	m := transport_stream.New(conn, c.BasePath, msgpackcodec.ContentType())
	invoker := functions.NewInvoker(m.Invoke, m.InvokeStream, msgpackcodec)

	return &compute.Compute{
		Invoker: invoker,
		Start: func() error {
			return conn.ReceiveLoop()
		},
		WaitUntilShutdown: func() error {
			conn.WaitUntilShutdown()
			return nil
		},
		Close: func() error {
			return conn.Close()
		},
		Environ: func() []string {
			return []string{
				fmt.Sprintf("BUS_SOCKET_ADDR=%s", c.SocketAddress),
			}
		},
	}, nil
}