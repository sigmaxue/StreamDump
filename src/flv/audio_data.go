package flv

import (
	"fmt"
	flvutils "github.com/sigmaxue/streamdump/src/utils"
)

// FLV AUDIO DATA 
// |-----------------------+---------------------+--------------+------------+
// | 0   |   1 |   2 |   3 |   4      |   5      |   6          |   7        |
// |---------------------------------------------+--------------+------------+
// |  audio format         |     simple rate     | simple length| audio type |   
// |-----------------------+---------------------+--------------+------------+



type AudioData struct {
	format    int 
	simplerate int 
	simplesize int
	audiotype  int
}

func (a *AudioData)Decode(buffer []byte) (count int, err error){
	a.format = flvutils.AudioFormat(buffer[0])
	a.simplerate = flvutils.AudioSimpleRate(buffer[0])
	a.simplesize = flvutils.AudioSimpleLen(buffer[0])
	a.audiotype = flvutils.AudioCodecType(buffer[0])
	return len(buffer), nil 
}

func (a *AudioData)ToString() string{
	out := fmt.Sprintf(", Format: %d", a.format)
	out += fmt.Sprintf(", SimpleRate: %d", a.simplerate)
	out += fmt.Sprintf(", SimpleSize: %d", a.simplesize)
	out += fmt.Sprintf(", AudioType: %d", a.audiotype)
	return out
}



