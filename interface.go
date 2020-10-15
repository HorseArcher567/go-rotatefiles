package rotatefiles

type Option interface {
	apply(*RotateFiles)
}
