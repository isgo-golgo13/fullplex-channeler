package svckit

import (
	"context"
	"io"
	"net"
	"time"
)

// FullPlexChannel defines the interface for full-duplex communication with context support.
type FullPlexChannel interface {
	Send(ctx context.Context, data io.Reader, n int64) (int64, error)
	SendTimeout(ctx context.Context, data io.Reader, n int64, timeout time.Duration) (int64, error)
	Recv(ctx context.Context, writer io.Writer, n int64) (int64, error)
	RecvAll(ctx context.Context, writers []io.Writer, n int64) (int64, error)
	SendAll(ctx context.Context, data []io.Reader, n int64) (int64, error)
	Close(ctx context.Context) error
}

// FullPlexChanneler implements the FullPlexChannel interface.
type FullPlexChanneler struct {
	conn   net.Conn
	closed bool
}

// NewFullPlexChanneler creates a new FullPlexChanneler.
func NewFullPlexChanneler(conn net.Conn) *FullPlexChanneler {
	return &FullPlexChanneler{conn: conn}
}

// Send sends data from an io.Reader to the connection.
func (fc *FullPlexChanneler) Send(ctx context.Context, data io.Reader, n int64) (int64, error) {
	if fc.closed {
		return 0, io.ErrClosedPipe
	}
	return io.CopyN(fc.conn, data, n)
}

// SendTimeout sends data with a timeout.
func (fc *FullPlexChanneler) SendTimeout(ctx context.Context, data io.Reader, n int64, timeout time.Duration) (int64, error) {
	done := make(chan struct{})
	var sentBytes int64
	var sendErr error

	go func() {
		sentBytes, sendErr = fc.Send(ctx, data, n)
		close(done)
	}()

	select {
	case <-done:
		return sentBytes, sendErr
	case <-ctx.Done():
		return sentBytes, ctx.Err()
	case <-time.After(timeout):
		return sentBytes, context.DeadlineExceeded
	}
}

// Recv receives data into an io.Writer from the connection.
func (fc *FullPlexChanneler) Recv(ctx context.Context, writer io.Writer, n int64) (int64, error) {
	if fc.closed {
		return 0, io.ErrClosedPipe
	}
	return io.CopyN(writer, fc.conn, n)
}

// RecvAll receives data into multiple io.Writer targets from the connection.
func (fc *FullPlexChanneler) RecvAll(ctx context.Context, writers []io.Writer, n int64) (int64, error) {
	var totalRecv int64
	for _, writer := range writers {
		recvBytes, err := fc.Recv(ctx, writer, n-totalRecv)
		totalRecv += recvBytes
		if err != nil {
			return totalRecv, err
		}
		if totalRecv >= n {
			break
		}
	}
	return totalRecv, nil
}

// SendAll sends all data from a slice of io.Reader objects.
func (fc *FullPlexChanneler) SendAll(ctx context.Context, data []io.Reader, n int64) (int64, error) {
	var totalSent int64
	for _, reader := range data {
		sent, err := fc.Send(ctx, reader, n-totalSent)
		totalSent += sent
		if err != nil {
			return totalSent, err
		}
		if totalSent >= n {
			break
		}
	}
	return totalSent, nil
}

// Close closes the connection.
func (fc *FullPlexChanneler) Close(ctx context.Context) error {
	fc.closed = true
	return fc.conn.Close()
}
