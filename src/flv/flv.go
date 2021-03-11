package flv 

import (
	"fmt"
	"bytes"
	"encoding/binary"
	"github.com/pkg/errors"
	flvutils "github.com/sigmaxue/streamdump/src/utils"
)

type Parser interface {
	Decode(buffer []byte) (count int, err error) 
	ToString() string
	HeaderSize() uint32
	TagSize() uint32
}

const(
	kFlvMagicLen = 3
	kFlvHeaderLen = 9 
	kPrevTagSizeLen = 4
)
// FLV HEADER
// |-----------+---------+-------+-------------|
// | signature | version | flags | Header Len  |
// |-----------+---------+-------+-------------|
// | 3B        | 1B      | 1B    | 4B          |
// |-----------+---------+-------+-------------|

type FlvHeader struct {
	Signature  [3]byte
	Version    uint8
	Flags      uint8
	HeaderLen uint32
	Raw     []byte
	PrevTagSize uint32
}

func (f *FlvHeader)HeaderSize() uint32 {
	return kFlvHeaderLen
}

func (f *FlvHeader)TagSize() uint32 {
	return f.HeaderLen + kPrevTagSizeLen
}

func (f *FlvHeader)Decode(data []byte) (count int, err error){
	if len(data) < kFlvHeaderLen + kPrevTagSizeLen {
		return 0, errors.Errorf("raw only %d bytes, %d is the minimum length for flv header", len(data), kFlvHeaderLen+kPrevTagSizeLen) 
	}
	if data[0]!='F' || 
		data[1]!='L' ||
		data[2]!='V' {
		return -1, errors.Errorf("not flv header") 
	}

	f.Signature[0] = data[0]
	f.Signature[1] = data[1]
	f.Signature[2] = data[2]
	f.Version = data[3]
	f.Flags = data[4]
	f.HeaderLen = binary.BigEndian.Uint32(data[5:9])
	if len(data) < int(f.HeaderSize() + kPrevTagSizeLen) {
		return 9, errors.Errorf("raw only %d bytes, %d is the minimum length for flv header", len(data), f.HeaderLen + kPrevTagSizeLen) 
	}
	f.PrevTagSize = binary.BigEndian.Uint32(data[f.HeaderLen : f.HeaderLen + kPrevTagSizeLen])
	copy(f.Raw, data[f.HeaderSize() + kPrevTagSizeLen:])
	return int(f.TagSize()), nil
}


func (f *FlvHeader)ToString() string{
	out := fmt.Sprintf("Magic: FLV")
	out += fmt.Sprintf(", Version: %d", f.Version)
	out += fmt.Sprintf(", HasVideo: %t", flvutils.HasVideo(f.Flags))
	out += fmt.Sprintf(", HasAudio: %t", flvutils.HasAudio(f.Flags))
	return out 
}

type Flv struct {
	dataBuffer        bytes.Buffer 
	nextTagSize       uint32
	Parser
}

func NewFlv() Flv {
	return Flv{nextTagSize: 0}
}

func (f *Flv)Decode(data []byte) (count int, err error) {
	if f.dataBuffer.Len() > 10*1024*1024 {
		return -1, errors.Errorf("raw only %d bytes, %d is the max length for flv tag", f.dataBuffer.Len(), 10*1024*1024) 
	}

	f.dataBuffer.Write(data)

	if f.dataBuffer.Len() < int(f.nextTagSize) {
		return 0, errors.Errorf("raw only %d bytes, %d is the minimum length for flv header", f.dataBuffer.Len(), f.nextTagSize) 
	}

	f.Parser = &FlvHeader{}

	Header := true
	for i:=0;i<10000;i++ {
		//fmt.Println(flvutils.DumpHex(f.dataBuffer.Bytes(), 0, 16))
		n, err := f.Parser.Decode(f.dataBuffer.Bytes())
		if n < 0  {
			if Header {
				f.Parser = &FlvTag{}
				Header = false
			} else {
				return -1, errors.Errorf("fatal error")
			}
		} else if n == 0 {
			f.nextTagSize = f.Parser.HeaderSize() 
			return n, err
		} else if n > 0 && err != nil {
			f.nextTagSize = f.Parser.TagSize()
			return n, err
		} else {
			// success
			fmt.Println(f.ToString(), ", TagSize: ", f.Parser.TagSize())

			f.dataBuffer.Next(int(f.Parser.TagSize())) 
			f.nextTagSize = 0
		}	
	}
	return -1, errors.Errorf("fatal error")
}
