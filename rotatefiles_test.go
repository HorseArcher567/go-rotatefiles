package rotatefiles

import (
	"fmt"
	"testing"
	"time"
)

func TestNewRotateFiles(t *testing.T) {
	rf, err := New("demo-info",
		//WithTimeLayout("2006010215"),
		//WithDir("./"),
		//WithMaxAge(time.Second*9),
		WithMaxCount(6),
		//WithMaxAge(time.Minute),
		WithRotateSize(128),
	)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rf)
		str := fmt.Sprintf("%s abcdefghijklmnopqrstuvwxyz0123456789中文测试\n",
			time.Now().Local())
		str1 := []byte(str)
		for i := 0; i < 200; i++ {
			rf.Write(str1)
			time.Sleep(time.Millisecond * 100)
		}
	}
}
