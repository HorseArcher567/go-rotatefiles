package rotatefiles

import "time"

// WithTimeLayout, default time layout "2006-01-02", you can use
// layout to format your file name, same to time.Format()
func WithTimeLayout(layout string) Option {
	return optionFunc(func(rf *RotateFiles) {
		rf.timeLayout = layout
	})
}

// WithMaxAge, files that have survived for more than age
// will be deleted.
// This option conflicts with WithMaxCount, the latter option will
// override the previous option.
func WithMaxAge(age time.Duration) Option {
	return optionFunc(func(rf *RotateFiles) {
		rf.reserveThreshold = age
	})
}

// WithMaxCount, if files number more than maxCount, the older files
// will be deleted.
// This option conflicts with WithMaxCount, the latter option will
// override the previous option.
func WithMaxCount(count int) Option {
	return optionFunc(func(rf *RotateFiles) {
		rf.reserveThreshold = count
	})
}

// WithRotatePeroid, sets the time between rotation.
func WithRotatePeroid(peroid time.Duration) Option {
	return optionFunc(func(rf *RotateFiles) {
		rf.rotatePeroid = peroid
	})
}

// WithRotateSize, if file size greater than rotate size, rotate will
// occur.
func WithRotateSize(size int) Option {
	return optionFunc(func(rf *RotateFiles) {
		rf.rotateSize = size
	})
}

// WithLocalClock, sets a clock that the RotateFiles object will
// use to determine the current time.
func WithLocalClock() Option {
	return optionFunc(func(rf *RotateFiles) {
		rf.clock = Local
	})
}
