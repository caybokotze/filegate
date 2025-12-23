package ffprobe

import (
	"errors"
	"os/exec"
)

var (
	exePath          string
	outputFormatFlag = "-of"
)

func exeFound() bool {
	return exePath != ""
}

// Available returns true if ffprobe or avprobe was found in PATH
func Available() bool {
	return exePath != ""
}

func init() {
	var err error
	exePath, err = exec.LookPath("ffprobe")
	if err == nil || errors.Is(err, exec.ErrDot) {
		outputFormatFlag = "-print_format"
		return
	}
	// Don't log "not found" errors - they're expected when ffprobe isn't installed
	exePath, err = exec.LookPath("avprobe")
	if err == nil || errors.Is(err, exec.ErrDot) {
		return
	}
	// Silently continue - ffprobe/avprobe not available
}

func isExecErrNotFound(err error) bool {
	if err == exec.ErrNotFound {
		return true
	}
	execErr, ok := err.(*exec.Error)
	if !ok {
		return false
	}
	return execErr.Err == exec.ErrNotFound
}
