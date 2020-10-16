package rotatefiles

import (
	"errors"
	"os"
	"regexp"
	"sync"
	"time"
)

type patternConversion struct {
	regexp     *regexp.Regexp
	replRegexp string
	replGlob   string
}

var pcs = []*patternConversion{
	{
		regexp:     regexp.MustCompile(`%yyyy`),
		replRegexp: `\d{4}`,
		replGlob:   `????`,
	},
	{
		regexp:     regexp.MustCompile(`%MM`),
		replRegexp: `\d{2}`,
		replGlob:   `??`,
	},
	{
		regexp:     regexp.MustCompile(`%dd`),
		replRegexp: `\d{2}`,
		replGlob:   `??`,
	},
	{
		regexp:     regexp.MustCompile(`%HH`),
		replRegexp: `\d{2}`,
		replGlob:   `??`,
	},
	{
		regexp:     regexp.MustCompile(`%mm`),
		replRegexp: `\d{2}`,
		replGlob:   `??`,
	},
	{
		regexp:     regexp.MustCompile(`%ss`),
		replRegexp: `\d{2}`,
		replGlob:   `??`,
	},
}

// RotateFiles represents file that gets automatically
// rotated as you write to it.
type RotateFiles struct {
	// File name template
	filePattern string
	// Regular expression to match file name
	regexpPattern string
	// Glob to match file name
	globPattern string

	mutex sync.Mutex
	// Current using file name
	currentFile string
	generation  int
	file        *os.File
	size        int

	// Use local time or
	clock Clock
	// Files reserve threshold, may be max age or max count
	reserveThreshold interface{}
	// Files rotate period
	rotatePeroid time.Duration
	// File max size
	rotateSize int
}

// New constructs a new RotateFiles from the provided file pattern and options.
// filePattern stipulate file name pattern, for example
// demo-info-%yyyy-%MM-%dd_%HH-%mm-%ss.log will create file like
// demo-info-2020-10-16_12-59-59.log
func New(filePattern string, options ...Option) (*RotateFiles, error) {
	if len(filePattern) == 0 {
		return nil, errors.New("File pattern must not be empty")
	}

	rf := &RotateFiles{
		filePattern: filePattern,
		// default, max age is 7 day
		reserveThreshold: 7 * 24 * time.Hour,
		// default, rotate peroid is a day
		rotatePeroid: 1 * 24 * time.Hour,
		// default, rotate size is 2 GB
		rotateSize: 2 * 1024 * 1024 * 1024,
	}

	rf.convertPattern()
	rf.withOptions(options...)

	return rf, nil
}

func (rf *RotateFiles) convertPattern() {
	for _, pc := range pcs {
		rf.regexpPattern =
			pc.regexp.ReplaceAllString(rf.regexpPattern, pc.replRegexp)
		rf.globPattern =
			pc.regexp.ReplaceAllString(rf.globPattern, pc.replGlob)
	}
}

func (rf *RotateFiles) withOptions(options ...Option) {
	for _, option := range options {
		option.apply(rf)
	}
}

func (rf *RotateFiles) rotate() {
	rf.mutex.Lock()
	defer rf.mutex.Unlock()

	if len(rf.currentFile) == 0 {
	}
}

func (rf *RotateFiles) genFileName() {

}

func (rf *RotateFiles) Write(p []byte) (n int, err error) {
	return rf.file.Write(p)
}
