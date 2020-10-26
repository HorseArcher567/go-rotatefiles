package rotatefiles

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

type patternConversion struct {
	regexp     *regexp.Regexp
	replRegexp string
}

// RotateFiles represents file that gets automatically
// rotated as you write to it.
type RotateFiles struct {
	filename   string
	timeLayout string
	dir        string

	mutex sync.Mutex
	// Current using file name
	curFilename    string
	generation     int
	lastRotateTime time.Time

	file *os.File
	size int

	cleanMutex sync.Mutex

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
func New(filename string, options ...Option) (*RotateFiles, error) {
	if len(filename) == 0 {
		return nil, errors.New("File pattern must not be empty")
	}

	rf := &RotateFiles{
		filename:   filename,
		timeLayout: `2006-01-02`,
		dir:        "./log/",
		// default, use local time
		clock: Local,
		// default, max age is 7 day
		reserveThreshold: 7 * 24 * time.Hour,
		// default, rotate peroid is a day
		rotatePeroid: 1 * 24 * time.Hour,
		// default, rotate size is 2 GB
		rotateSize: 2 * 1024 * 1024 * 1024,
	}

	rf.withOptions(options...)

	return rf, nil
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
	var truncTime time.Time

	if now.Location() != time.UTC {
		truncTime = time.Date(now.Year(), now.Month(), now.Day(),
			now.Hour(), now.Minute(), now.Second(),
			now.Nanosecond(), time.UTC)
		truncTime = truncTime.Truncate(rf.rotatePeroid)
		truncTime = time.Date(truncTime.Year(), truncTime.Month(), truncTime.Day(),
			truncTime.Hour(), truncTime.Minute(), truncTime.Second(),
			truncTime.Nanosecond(), truncTime.Location())
	} else {
		truncTime = now.Truncate(rf.rotatePeroid)
	}

	if truncTime == rf.lastRotateTime {
		return false
	}

	rf.lastRotateTime = truncTime
	oldGeneration := rf.generation
	rf.generation = 0
	err := rf.genFile()
	if err != nil {
		rf.generation = oldGeneration
		fmt.Println(err)
		return false
	}

	return true
}

func (rf *RotateFiles) rotateBySize() bool {
	if rf.size < rf.rotateSize {
		return false
	}

	err := rf.genFile()
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (rf *RotateFiles) genFile() error {
	err := os.MkdirAll(rf.dir, 0755)
	if err != nil {
		panic(err)
	}

	fileTime := rf.lastRotateTime.Format(rf.timeLayout)
	filesize := 0
	for {
		rf.curFilename =
			rf.filename + fileTime + "_" + strconv.Itoa(rf.generation) + ".log"
		fileinfo, err := os.Stat(rf.dir + rf.curFilename)
		if os.IsNotExist(err) {
			filesize = 0
			break
		} else if err == nil && fileinfo.Size() < int64(rf.rotateSize) {
			filesize = int(fileinfo.Size())
			break
		}

		rf.generation++
	}

	file, err := os.OpenFile(rf.dir+rf.curFilename,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}

	if rf.file != nil {
		rf.file.Close()
	}

	go rf.cleanFile()

	rf.file = file
	rf.size = filesize
	return nil
}

type GeneratedFile struct {
	filePath string
	modTime  time.Time
}
type GeneratedFileSlice []*GeneratedFile

func (s GeneratedFileSlice) Len() int {
	return len(s)
}

func (s GeneratedFileSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s GeneratedFileSlice) Less(i, j int) bool {
	return s[i].modTime.After(s[j].modTime)
}

func (rf *RotateFiles) cleanFile() {
	rf.cleanMutex.Lock()
	defer rf.cleanMutex.Unlock()

	matches, err := filepath.Glob(rf.dir + rf.filename + "*.log")
	if err != nil {
		fmt.Println(err)
		return
	}

	generatedFiles := make([]*GeneratedFile, 0)

	for _, match := range matches {
		_, file := filepath.Split(match)
		index := strings.LastIndexByte(file, '_')
		if index != -1 {
			_, err := time.Parse(rf.timeLayout, file[len(rf.filename):index])
			if err != nil {
				fmt.Println(err)
				continue
			}

			fileInfo, err := os.Stat(match)
			if err != nil {
				continue
			}

			generatedFiles = append(generatedFiles, &GeneratedFile{
				filePath: match,
				modTime:  fileInfo.ModTime(),
			})
		}
	}

	switch rf.reserveThreshold.(type) {
	case int:
		sort.Sort(GeneratedFileSlice(generatedFiles))
		for index, file := range generatedFiles {
			if index >= rf.reserveThreshold.(int) {
				os.Remove(file.filePath)
			}
		}
	case time.Duration:
		criticalTime := time.Now().Add(-rf.reserveThreshold.(time.Duration))
		for _, file := range generatedFiles {
			if criticalTime.After(file.modTime) {
				os.Remove(file.filePath)
			}
		}
	default:
		panic(errors.New("Unknown clean rule."))
	}
}

func (rf *RotateFiles) Write(p []byte) (n int, err error) {
	rf.mutex.Lock()
	defer rf.mutex.Unlock()

	rf.rotate()

	len, err := rf.file.Write(p)
	rf.size += len

	return len, err
}
