package rotatefiles

import "time"

type optionFunc func(*RotateFiles)

func (f optionFunc) apply(rf *RotateFiles) {
	f(rf)
}

func WithMaxAge(maxAge time.Duration) Option {
	return optionFunc(func(rf *RotateFiles) {
		rf.maxAge = maxAge
	})
}
