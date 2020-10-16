package rotatefiles

import (
	"fmt"
	"testing"
	"time"
)

func TestHelloRotateFiles(t *testing.T) {
	fmt.Println("Hello, rotatefiles!")

	str1 := "hello"
	str2 := str1
	fmt.Printf("%s, %s\n", str1, str2)

	str2 = "wolrd"
	fmt.Printf("%s, %s\n", str1, str2)

	fmt.Println("===========================================")
}

func TestTimeTruncate(tt *testing.T) {
	t, _ := time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 12:15:30.918273645")
	fmt.Println(t)
	trunc := []time.Duration{
		time.Nanosecond,
		time.Microsecond,
		time.Millisecond,
		time.Second,
		2 * time.Second,
		time.Minute,
		10 * time.Minute,
		24 * time.Hour,
	}

	for _, d := range trunc {
		fmt.Printf("t.Truncate(%5s) = %s\n",
			d, t.Truncate(d).Format("2006-01-02 15:04:05.999999999"))
	}
	// To round to the last midnight in the local timezone, create a new Date.
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.Local)
	_ = midnight

	fmt.Println("===========================================")
}

func TestNewRotateFiles(t *testing.T) {
	rf, err := New("demo-info",
		WithTimeLayout("2006-01-02_15"),
		WithMaxAge(time.Minute))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rf)
		str := "adfasdfasdfasdfasdfadsfadsfa\n"
		str1 := []byte(str)
		for i := 0; i < 1000; i++ {
			rf.Write(str1)
		}
	}

	fmt.Println("===========================================")
}
