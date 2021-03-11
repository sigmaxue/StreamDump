package flv

import (
	"fmt"
	//flvutils "github.com/sigmaxue/streamdump/src/utils"
)


type ScriptData struct {
}

func (s *ScriptData)Decode(buffer []byte) (count int, err error){
	return len(buffer), nil
}

func (s *ScriptData)ToString() string{
	out := fmt.Sprintf("")
	return out
}


