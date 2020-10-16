package rotatefiles

import "time"

type Option interface {
	apply(*RotateFiles)
}
type optionFunc func(*RotateFiles)

func (f optionFunc) apply(rf *RotateFiles) {
	f(rf)
}

type Clock interface {
	Now() time.Time
}
type clockFunc func() time.Time

func (f clockFunc) Now() time.Time {
	return f()
}

var UTC = clockFunc(func() time.Time { return time.Now().UTC() })
var Local = clockFunc(time.Now)
