// ffmpegthumbnailer wrapper for Windows
// Translates ffmpegthumbnailer calls to ffmpeg
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	// Parse arguments manually to handle -cpng style (no space)
	input := ""
	output := ""
	size := 128
	timePercent := 10
	quality := 8
	format := "png"

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch {
		case arg == "-i" && i+1 < len(args):
			i++
			input = args[i]
		case arg == "-o" && i+1 < len(args):
			i++
			output = args[i]
		case arg == "-s" && i+1 < len(args):
			i++
			size, _ = strconv.Atoi(args[i])
		case arg == "-t" && i+1 < len(args):
			i++
			timePercent, _ = strconv.Atoi(args[i])
		case arg == "-q" && i+1 < len(args):
			i++
			quality, _ = strconv.Atoi(args[i])
		case arg == "-c" && i+1 < len(args):
			i++
			format = args[i]
		case strings.HasPrefix(arg, "-c"):
			// Handle -cpng or -cjpeg style
			format = strings.TrimPrefix(arg, "-c")
		}
	}

	if input == "" {
		fmt.Fprintln(os.Stderr, "Error: No input file specified")
		os.Exit(1)
	}

	// Get video duration using ffprobe
	duration := getVideoDuration(input)

	// Calculate seek position from percentage
	var seekSec float64
	if duration > 0 {
		seekSec = duration * float64(timePercent) / 100.0
		// Don't seek too close to the end
		if seekSec > duration-1 {
			seekSec = duration / 2
		}
	} else {
		// Fallback: use 2 seconds if we can't get duration
		seekSec = 2
	}

	// Try to generate thumbnail, with fallback to start of video
	err := generateThumbnail(input, output, seekSec, size, quality, format)
	if err != nil {
		// Fallback: try from the very beginning
		err = generateThumbnail(input, output, 0, size, quality, format)
		if err != nil {
			os.Exit(1)
		}
	}
}

func getVideoDuration(input string) float64 {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		input,
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return 0
	}

	durationStr := strings.TrimSpace(out.String())
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return 0
	}

	return duration
}

func generateThumbnail(input, output string, seekSec float64, size, quality int, format string) error {
	args := []string{
		"-hide_banner",
		"-loglevel", "error",
	}

	// Seek before input for faster seeking
	if seekSec > 0 {
		args = append(args, "-ss", fmt.Sprintf("%.2f", seekSec))
	}

	args = append(args, "-i", input)

	// Scale filter
	if size > 0 {
		args = append(args, "-vf", fmt.Sprintf("scale=%d:-1", size))
	}

	// Single frame
	args = append(args, "-vframes", "1")

	// Quality for jpeg
	if format == "jpeg" || format == "jpg" {
		q := 31 - (quality * 3) // Convert 1-10 to ffmpeg's 2-31 scale (inverted)
		if q < 2 {
			q = 2
		}
		args = append(args, "-q:v", strconv.Itoa(q))
	}

	// Output
	if output == "/dev/stdout" || output == "-" {
		args = append(args, "-f", "image2pipe")
		if format == "jpeg" || format == "jpg" {
			args = append(args, "-vcodec", "mjpeg")
		} else {
			args = append(args, "-vcodec", "png")
		}
		args = append(args, "-")
	} else {
		args = append(args, "-y", output)
	}

	cmd := exec.Command("ffmpeg", args...)
	cmd.Stdout = os.Stdout
	// Suppress stderr to avoid noise
	cmd.Stderr = nil

	return cmd.Run()
}
