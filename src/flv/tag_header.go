package flv

import (
	"github.com/pkg/errors"
	"fmt"
	flvutils "github.com/sigmaxue/streamdump/src/utils"
)


// FLV TAG HEADER
// |------+----------+-----------+--------------+----------|
// | Type | DataSize | Timestamp | timestamp_ex | StreamId |
// |------+----------+-----------+--------------+----------|
// | 1B   | 3B       | 3B        | 1B           | 3B       |
// |------+----------+-----------+--------------+----------|

const (
	kTagHeaderLen = 11
	kFlvTagTypeAudio  = 0x08
	kFlvTagTypeVideo  = 0x09
	kFlvTagTypeScript   = 0x12
	kFlvTagTypeReserved  = 0xFF
)

type TagHeader struct {
	tagtype         uint8
	datasize        uint32
	timestamp       uint32
	timestamp_ex    uint32
	streamid        uint32
}

func (t *TagHeader)IsVideo() bool {
	return t.tagtype == kFlvTagTypeVideo
}

func (t *TagHeader)IsAudio() bool {
	return t.tagtype == kFlvTagTypeAudio
}

func (t *TagHeader)IsScript() bool {
	return t.tagtype == kFlvTagTypeScript
}


func (t *TagHeader)Decode(buffer []byte) (count int, err error) {
	if len(buffer) < kTagHeaderLen {
		return 0, errors.Errorf("raw only %d bytes, %d is the minimum length for tag heder", len(buffer), kTagHeaderLen) 
	}

	t.tagtype = buffer[0]
	t.datasize = flvutils.BigEndianU32(buffer[1:4])
	t.timestamp = flvutils.BigEndianU32(buffer[4:7])
	t.timestamp_ex = uint32(buffer[7])
	t.streamid = flvutils.BigEndianU32(buffer[8:11])

	return kTagHeaderLen, nil
}

func (t *TagHeader)DataSize() uint32{
	return t.datasize
}
// (t *TagHeader) ...
func (t *TagHeader)ToString() string{
	out := fmt.Sprintf("Type: %d", t.tagtype)
	out += fmt.Sprintf(", DataSize: %d", t.datasize)
	out += fmt.Sprintf(", TimeStamp: %d", t.timestamp | t.timestamp_ex<<3)
	return out 
}


