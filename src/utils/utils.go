package utils 

import (
	"fmt"
	"bytes"
)
const (
	kFlvFlagVideo = 0x1
	kFlvFlagAudio = 0x4
)


func BigEndianU32(data []byte) uint32 {
	var out uint32 
	for i:=0; i<len(data); i++ {
		out = out << 8 + uint32(data[i])
	}
	return out
}

func HasVideo(b byte) bool {
	return (b&kFlvFlagVideo) != 0
}

func HasAudio(b byte) bool {
	return (b&kFlvFlagAudio) != 0
}

func FrameType(b byte) int {
	return int((b & 0xF0) >> 4)
}

func CodecId(b byte) int {
	return int(b & 0x0F)
}

func CheckStartCode(buffer []byte) bool {
	if (len(buffer) >= 4) &&
		buffer[0] == 0 &&
		buffer[1] == 0 &&
		buffer[2] == 0 &&
		buffer[3] == 1 { 
		return true
	}
	return false
}

func AudioFormat(b byte) int {
	return int((b & 0xF0) >> 4)
}

func AudioSimpleRate(b byte) int {
	return int((b & 0x06) >> 2)
}

func AudioSimpleLen(b byte) int {
	return int((b & 0x02) >> 1)
}

func AudioCodecType(b byte) int {
	return int(b & 0x01)
}

func DumpHex(data []byte, start_pos int, width int) string {
	if width==0 {
		return ""
	}
	tab := "0123456789ABCDEF";

	var hexline bytes.Buffer
	var charline bytes.Buffer
	
	out := ""
	
	len := len(data)
	for i:=0; i<len; i++ {
		hexline.WriteByte(tab[data[i]>>4])
		hexline.WriteByte(tab[data[i]&0xF])
		hexline.WriteByte(' ')
		
		c := '.'
		if data[i] >= 32 && data[i] <= 126 {
			c = rune(data[i])
		}
		charline.WriteByte(byte(c))

		if ((i+1) % width == 0) {
			out += fmt.Sprintf("%08X: %s %s\n", start_pos+i+1-width, hexline.String(), charline.String())
			charline.Reset()
			hexline.Reset()
		}
	}

	for i:=len; i< (len/width+1)*width;i++ {
		hexline.WriteByte(' ')
		charline.WriteByte(' ')
		if ((i+1) % width == 0) {
			out += fmt.Sprintf("%08X: %s %s\n", start_pos+i+1-width, hexline.String(), charline.String())
			charline.Reset()
			hexline.Reset()
			break;
		}
		i++;
	}
	return out
}
