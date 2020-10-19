package rotatefiles

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestNewRotateFiles(t *testing.T) {
	var wg1 sync.WaitGroup
	for j := 0; j < 10; j++ {
		wg1.Add(1)
		filename := fmt.Sprintf("demo-info_%d_", j)
		go func(filename string) {
			rf, err := New(filename,
				WithTimeLayout("2006-01-02_15"),
				//WithDir("./"),
				//WithMaxAge(time.Second*9),
				WithMaxCount(6),
				WithRotateSize(1024*1024*2),
			)

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(rf)

				var wg sync.WaitGroup

				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func() {
						str := fmt.Sprintf("%s abcdefghijklmnopqrstuvwxyz0123456789中文测试\n",
							time.Now().Local())
						str1 := []byte(str)
						for i := 0; i < 20000; i++ {
							rf.Write(str1)
							//time.Sleep(time.Millisecond * 10)
						}

						wg.Done()
					}()
				}

				wg.Wait()
			}

			wg1.Done()
		}(filename)
	}

	wg1.Wait()

	time.Sleep(10 * time.Second)
}
