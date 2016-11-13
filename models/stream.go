package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/dankeroni/gotwitch"
	"github.com/garyburd/redigo/redis"
)

// StreamChunk stores chunk-specific data
type StreamChunk struct {
	ID       int64     `json:"ID"`
	Game     string    `json:"Game"`
	Title    string    `json:"Title"`
	Start    time.Time `json:"Start"`
	LastSeen time.Time `json:"LastSeen"`
}

// Stream stores stream-data about 1 or more chunks
type Stream struct {
	Start  time.Time     `json:"Start"`
	End    *time.Time    `json:"End"`
	Chunks []StreamChunk `json:"Chunks"`
}

// GetLastStream xD
func GetLastStream(channel string) (*Stream, error) {
	conn := db.Pool.Get()
	defer conn.Close()

	bytes, err := redis.ByteSlices(conn.Do("LRANGE", channel+":streams", 0, 0))
	if err != nil {
		return nil, err
	}
	if len(bytes) != 1 {
		return nil, errors.New("No stream available")
	}
	var stream Stream
	err = json.Unmarshal(bytes[0], &stream)
	if err != nil {
		return nil, err
	}
	return &stream, nil
}

// GetLastTwoStreams xD
func GetLastTwoStreams(channel string) ([]Stream, error) {
	conn := db.Pool.Get()
	defer conn.Close()

	bytes, err := redis.ByteSlices(conn.Do("LRANGE", channel+":streams", 0, 1))
	if err != nil {
		return nil, err
	}
	var streams []Stream
	streams = make([]Stream, 2, 2)
	for _, jsonData := range bytes {
		var stream Stream
		err = json.Unmarshal(jsonData, &stream)
		if err != nil {
			return nil, err
		}
		streams = append(streams, stream)
	}
	return streams, nil
}

// ChunkFromTwitchStream xD
func ChunkFromTwitchStream(in gotwitch.Stream) StreamChunk {
	return StreamChunk{
		ID:       in.ID,
		Game:     in.Channel.Game,
		Title:    in.Channel.Status,
		Start:    in.CreatedAt,
		LastSeen: time.Now(),
	}
}

// IsNewStreamChunk xD
func IsNewStreamChunk(channel string, streamChunk StreamChunk) (bool, *Stream, error) {
	conn := db.Pool.Get()
	defer conn.Close()

	streams, err := GetLastTwoStreams(channel)
	if err != nil {
		return false, nil, err
	}

	var lastStream *Stream

	for _, prevStream := range streams {
		if lastStream == nil {
			lastStream = &prevStream
		}
		for _, prevChunk := range prevStream.Chunks {
			if prevChunk.ID == streamChunk.ID {
				return false, &prevStream, nil
			}
		}
	}

	return true, lastStream, nil
}

func addNewStream(conn redis.Conn, channel string, streamChunk StreamChunk) error {
	stream := Stream{
		Start: streamChunk.Start,
		End:   nil,
		Chunks: []StreamChunk{
			streamChunk,
		},
	}
	bytes, err := json.Marshal(&stream)
	if err != nil {
		return err
	}

	conn.Send("LPUSH", channel+":streams", bytes)

	return nil
}

// AddNewStream xD
func AddNewStream(channel string, streamChunk StreamChunk) error {
	conn := db.Pool.Get()
	defer conn.Close()

	err := addNewStream(conn, channel, streamChunk)
	if err != nil {
		return err
	}

	conn.Flush()

	return nil
}

// AddNewStreamP xD
func AddNewStreamP(conn redis.Conn, channel string, streamChunk StreamChunk) error {
	return addNewStream(conn, channel, streamChunk)
}

// AddStreamChunk adds a strema chunk to the streams chunks slice
func (s *Stream) AddStreamChunk(streamChunk StreamChunk) {
	s.Chunks = append(s.Chunks, streamChunk)
}

// UpdateStreamChunk returns true if we updated it, otherwise return false
func (s *Stream) UpdateStreamChunk(streamChunk StreamChunk) bool {
	for i, prevChunk := range s.Chunks {
		if prevChunk.ID == streamChunk.ID {
			s.Chunks[i] = streamChunk
			return true
		}
	}

	return false
}

func updateLastStream(conn redis.Conn, channel string, stream *Stream) error {
	bytes, err := json.Marshal(stream)
	if err != nil {
		return err
	}

	conn.Send("LSET", channel+":streams", 0, bytes)

	return nil
}

// UpdateLastStream xD
func UpdateLastStream(channel string, stream *Stream) error {
	conn := db.Pool.Get()
	defer conn.Close()

	err := updateLastStream(conn, channel, stream)
	if err != nil {
		return err
	}

	conn.Flush()

	return nil
}

// UpdateLastStreamP xD
func UpdateLastStreamP(conn redis.Conn, channel string, stream *Stream) error {
	return updateLastStream(conn, channel, stream)
}
