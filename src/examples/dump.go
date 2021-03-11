package main 

import (
	"net/http"
	"flag"
	"fmt"
	"github.com/sigmaxue/streamdump/src/flv"
)


var url = flag.String("r", "http://6721.liveplay.myqcloud.com/live/6721_d71956d9cc93e4a467b11e06fdaf039a.flv", "http or rtmp url[http://www.xxx.com/live/x.flv]")
var level = flag.Int("level", 1, "log level")

func main() {
	flag.Parse()

	if len(*url) == 0 {
		return
	}

	fmt.Println(*url)
	client := &http.Client{}
	reqest, err := http.NewRequest("GET", *url, nil)

	if err != nil {
	} 

	reqest.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36")
	resp, err := client.Do(reqest)
	if err != nil {
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
	}
	
	var flv flv.Flv
	buf := make([]byte, 102400)          
	for {
		n, err := resp.Body.Read(buf)    
		if err != nil {
			break
		}
		flv.Decode(buf[:n])
		fmt.Println("recv size: ", n)
	}
	
}
