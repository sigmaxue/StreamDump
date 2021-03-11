package flv

import (
	//"encoding/binary"
	"github.com/pkg/errors"

)

type TagData interface {
	Decode(raw []byte) (count int , err error) 
	ToString() string
}

type Position struct {
	Offset uint32
	Length uint32
}

type FlvTag struct {
	Position

	Header      TagHeader
	Data        TagData
	Raw         []byte
	PrevTagSize uint32
}

func (f *FlvTag)HeaderSize() uint32 {
	return kTagHeaderLen
}

func (f *FlvTag)TagSize() uint32 {
	return kPrevTagSizeLen + f.Header.DataSize() + f.HeaderSize();
}

func (t *FlvTag)Decode(buffer []byte) (count int, err error){
	n, err := t.Header.Decode(buffer)
	if err != nil {
		return n, err
	}

	if len(buffer) < int(t.TagSize()) {
		return n, errors.Errorf("raw only %d bytes, %d is the minimum length for flv tag", len(buffer), t.TagSize()) 
	}

	if t.Header.IsVideo() {
		t.Data = &VideoData{}
	} else if t.Header.IsAudio() {
		t.Data = &AudioData{}
	} else if t.Header.IsScript() {
		t.Data = &ScriptData{}
	} else {
		return -1, errors.Errorf("unknow tag type: %d", t.Header.tagtype) 
	}
	return t.Data.Decode(buffer[n:t.TagSize()])
}

func (t *FlvTag)ToString() string{
	return t.Header.ToString() + t.Data.ToString()
}

