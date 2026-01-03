package natsjs

import (
	"context"

	"github.com/nats-io/nats.go/jetstream"
)

type StreamInfo struct {
	Name     string
	Subjects []string
	State    jetstream.StreamState
}

type HealthStats struct {
	ConnStatus string
	Msg        string

	InMsgsTotal     uint64
	OutMsgsTotal    uint64
	InBytesTotal    uint64
	OutBytesTotal   uint64
	ReconnectsTotal uint64

	Streams []StreamInfo
}

func (c *Client) Health(ctx context.Context) HealthStats {
	if c.conn == nil {
		return HealthStats{ConnStatus: "unknown", Msg: "connection is nil"}
	}

	stats := c.conn.Statistics

	streams := make([]StreamInfo, 0, 8)
	lst := c.js.ListStreams(ctx)
	for si := range lst.Info() {
		if si == nil {
			continue
		}
		streams = append(streams, StreamInfo{
			Name:     si.Config.Name,
			Subjects: si.Config.Subjects,
			State:    si.State,
		})
	}
	msg := "OK"
	if err := lst.Err(); err != nil {
		msg = err.Error()
	}

	return HealthStats{
		ConnStatus: c.conn.Status().String(),
		Msg:        msg,

		InMsgsTotal:     stats.InMsgs,
		OutMsgsTotal:    stats.OutMsgs,
		InBytesTotal:    stats.InBytes,
		OutBytesTotal:   stats.OutBytes,
		ReconnectsTotal: stats.Reconnects,

		Streams: streams,
	}
}
