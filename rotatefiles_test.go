package rotatefiles

import (
	"fmt"
	"testing"
	"time"
)

func TestHelloRotateFiles(t *testing.T) {
	fmt.Println("Hello, rotatefiles!")
}

func TestNewRotateFiles(t *testing.T) {
	rf, err := New("demo-info-%yyyy-%MM-%dd.log",
		WithMaxAge(time.Minute))
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(rf)
	}

	copy := rf.WithOptions(WithMaxAge(time.Hour))
	fmt.Println(copy)
	fmt.Println(rf)
}
