package rotatefiles

import "time"

type RotateFiles struct {
	filePattern   string
	regexpPattern string
	globPattern   string

	currentFile string

	maxAge   time.Duration
	maxCount int
}

func New(filePattern string, options ...Option) (*RotateFiles, error) {
	rf := &RotateFiles{
		filePattern: filePattern,
	}

	rf = rf.WithOptions(options...)

	return rf, nil
}

func (rf *RotateFiles) clone() *RotateFiles {
	copy := rf
	return copy
}

func (rf *RotateFiles) WithOptions(options ...Option) *RotateFiles {
	copy := rf.clone()
	for _, option := range options {
		option.apply(copy)
	}

	return copy
}
