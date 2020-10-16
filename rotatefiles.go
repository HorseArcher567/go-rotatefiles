package rotatefiles

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
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
	fileName string
	// File name template
	filePattern string
	timeLayout  string
	// Regular expression to match file name
	regexpPattern string
	// Glob to match file name
	globPattern string

	mutex sync.Mutex
	// Current using file name
	curFileName    string
	generation     int
	lastRotateTime time.Time

	file *os.File
	size int

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
		fileName:      filePattern,
		filePattern:   filePattern,
		regexpPattern: filePattern,
		globPattern:   filePattern,
		// default, use local time
		clock: Local,
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
	_ = rf.rotateByTime() || rf.rotateBySize()
}

func (rf *RotateFiles) rotateByTime() bool {
	now := rf.clock.Now()
	truncTime := now.Truncate(rf.rotatePeroid)
	if truncTime == rf.lastRotateTime {
		return false
	}

	rf.lastRotateTime = truncTime
	rf.generation = 0
	rf.genFile()

	return true
}

func (rf *RotateFiles) rotateBySize() bool {
	if rf.size < rf.rotateSize {
		return false
	}

	rf.genFile()

	return true
}

func (rf *RotateFiles) genFile() {
	fileTime := rf.lastRotateTime.Format(rf.timeLayout)
	for {
		rf.curFileName =
			rf.fileName + fileTime + "_" + strconv.Itoa(rf.generation) + ".log"
		_, err := os.Stat(rf.curFileName)
		if err != nil {
			break
		}

		rf.generation++
	}

	file, err := os.OpenFile(rf.curFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	if rf.file != nil {
		rf.file.Close()
	}

	rf.file = file
}

func (rf *RotateFiles) Write(p []byte) (n int, err error) {
	rf.mutex.Lock()
	defer rf.mutex.Unlock()

	rf.rotate()

	return rf.file.Write(p)
}
