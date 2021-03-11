package flv


import (
	"fmt"
	flvutils "github.com/sigmaxue/streamdump/src/utils"
)

// FLV VIDEO DATA 

// |-------------+-----------------+----------|
// | frame-codec | avc packet type |  cts     |
// |-------------+-----------------+----------|
// |     1B      |       1B        |   3B     |
// |-------------+-----------------+----------|

//  frame-codec in bit
// |---+---+---+---+---+---+---+---+
// | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 |
// |---+---+---+---+---+---+---+---+
// |   frame type  | codec id      | 
// |---+---+---+---+---+---+---+---+

const (
	kFlvAvcPacketTypeSequenceHeader = 0
	kFlvAvcPacketTypeNalu = 1
	kFlvAvcPacketTypeSequenceEnd = 2

	kFlvVideoTagCodecJPEG =             1
	kFlvVideoTagCodecSORENSENH263 =    2
	kFlvVideoTagCodecSCREENVIDEO =     3
	kFlvVideoTagCodecON2VP6 =          4
	kFlvVideoTagCodecON2VP6ALPHA =    5
	kFlvVideoTagCodecSCREENVIDEOV2 =  6
	kFlvVideoTagCodecAVC =              7
	kFlvVideoTagCodecH265 =             12
	kFlvVideoTagCodecAV1 =              13
	kFlvVideoTagCodecVP9 =              14

	kFlvVideoTagFrameTypeKEYFRAME =                1
	kFlvVideoTagFrameTypeINTERFRAME =              2
	kFlvVideoTagFrameTypeDISPOSABLEINTERFRAME =   3
	kFlvVideoTagFrameTypeGENERATEDKEYFRAME =      4
	kFlvVideoTagFrameTypeCOMMANDFRAME =           5

	kFlvNaluTypeSlice = 1     //Coded slice of a non-IDR picture
	kFlvNaluTypeDPA = 2       //Coded slice data partition A
	kFlvNaluTypeDPB = 3       //Coded slice data partition B
	kFlvNaluTypeDPC = 4       //Coded slice data partition C
	kFlvNaluTypeIDR = 5       //Coded slice of an IDR picture
	kFlvNaluTypeSEI = 6       //Supplemental enhancement information
	kFlvNaluTypeSPS = 7       //Sequence parameter set
	kFlvNaluTypePPS = 8       //Picture parameter set
	kFlvNaluTypeAUD = 9       //Access unit delimiter
	kFlvNaluTypeEOSEQ = 10    //End of sequence
	kFlvNaluTypeEOSTREAM = 11 //End of stream
	kFlvNaluTypeFD = 12       //Filler data
	kFlvNaluTypeSPSE = 13     //Sequence parameter set extension
	kFlvNaluTypeUNRESERVEDSTART = 24     //unreserved nalu type
	kFlvNaluTypeUNRESERVEDEND = 31     //unreserved nalu type
	//265
	kFlvNaluTypeVPS265 = 32     // VPS
	kFlvNaluTypeSPS265 = 33     // SPS
	kFlvNaluTypeSei265 = 39     //Sequence parameter set extension
)

type Nalu struct {
	Position

	Type uint8
	Idc  uint8
	Size uint32
}

type VideoData struct {
	Position

	frameType     int 
	codecIdc      int
	avcPacketType uint8
	cts           uint32
	nalus         []Nalu
}

func (v *VideoData)Decode(buffer []byte) (count int, err error) {
	v.frameType = flvutils.FrameType(buffer[0])
	v.codecIdc = flvutils.CodecId(buffer[0])
	v.avcPacketType = buffer[1]
	v.cts = flvutils.BigEndianU32(buffer[2:5])

	return len(buffer), nil 
}

func (t *VideoData)ToString() string{
	out := fmt.Sprintf(", FrameType: %d", t.frameType)
	out += fmt.Sprintf(", CodecId: %d", t.codecIdc)
	out += fmt.Sprintf(", AvcPacketType: %d", t.avcPacketType)
	out += fmt.Sprintf(", Cts: %d", t.cts)
	return out 
}



// H264 avc config
// |---------+---------+---------+-------+----------+-----------+-----------+
// | version | profile |reserve  | level | nalu len | sps entry | pps entry |
// |---------+---------+---------+-------+----------+-----------+-----------+
// | 1B      | 1B      |   1B    | 1B    | 1B       |     X     |     X     |  
// |---------+---------+---------+-------+----------+-----------+-----------+

// for each entry
// +-----+-----------+-------------+--------------+
// | num | nalu_size | nalu header | nalu payload |
// +-----+-----------+-------------+--------------+
// | 1B  | 2B        | 2B          | n            |
// +-----+-----------+-------------+--------------+


// H264 NAL units in byte
// +--------+--------+--------+--------+
// |       0|       1|       2|       3|
// +--------+--------+--------+--------+
// |                           NAL size|
// +--------+--------+--------+--------+
//
// 264 NAL unit header in bit
// +--------+--------+--------+--------+--------+--------+--------+--------+
// |       0|       1|       2|       3|       4|       5|       6|       7|
// +--------+--------+--------+--------+--------+--------+--------+--------+
// |        |        NalRefIdc|                                 NalUnitType|
// +--------+--------+--------+--------+--------+--------+--------+--------+
//


// HEVC AVC
// |-----+------------+----------------+-----------+-----------+-----------+
// | 21B | fixed (1B) |entry num (1B)  | vps entry | sps entry | pps entry | 
// |-----+------------+----------------+-----------+-----------+-----------+    
// for each entry
// |------+-----+-----------+-------------+--------------+
// | type | num | nalu_size | nalu header | nalu payload |
// |------+-----+-----------+-------------+--------------+
// | 1B   | 2B  | 2B        | 2B          | n            |
// |------+-----+-----------+-------------+--------------+

// H265 NAL units in byte 
// +--------+--------+--------+--------+ 
// |       0|       1|       2|       3|
// +--------+--------+--------+--------+
// |                           NAL size|
// +--------+--------+--------+--------+
//                                     
// NAL unit header in bit //2B        
// |---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---|  
// | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 
// |---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---| 
// | F |          TYPE             |        LayerId    |   TID     | 
// |---+---+---+---+---+---+---+---+---+---+---+---+---+---+---+---| 
