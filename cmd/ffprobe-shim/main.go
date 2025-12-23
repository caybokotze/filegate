// ffprobe-shim is a pure Go replacement for ffprobe that outputs compatible JSON
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/filegate/filegate/internal/probe"
)

// FFProbeOutput mimics ffprobe's JSON output format
type FFProbeOutput struct {
	Streams []Stream `json:"streams"`
	Format  Format   `json:"format"`
}

type Stream struct {
	Index       int    `json:"index"`
	CodecName   string `json:"codec_name,omitempty"`
	CodecType   string `json:"codec_type"`
	Width       int    `json:"width,omitempty"`
	Height      int    `json:"height,omitempty"`
	Duration    string `json:"duration,omitempty"`
	BitRate     string `json:"bit_rate,omitempty"`
	SampleRate  string `json:"sample_rate,omitempty"`
	Channels    int    `json:"channels,omitempty"`
}

type Format struct {
	Filename   string `json:"filename"`
	FormatName string `json:"format_name"`
	Duration   string `json:"duration"`
	Size       string `json:"size"`
	BitRate    string `json:"bit_rate"`
}

func main() {
	// Parse arguments - we only care about the filename
	// ffprobe is typically called like: ffprobe -v quiet -print_format json -show_format -show_streams <file>
	var filename string
	for i, arg := range os.Args[1:] {
		// Skip flags and their values
		if arg[0] == '-' {
			continue
		}
		// Check if previous arg was a flag that takes a value
		if i > 0 && os.Args[i][0] == '-' && os.Args[i] != "-v" && os.Args[i] != "-show_format" && os.Args[i] != "-show_streams" {
			continue
		}
		filename = arg
		break
	}

	// Find filename (last non-flag argument)
	for i := len(os.Args) - 1; i >= 1; i-- {
		if os.Args[i][0] != '-' {
			filename = os.Args[i]
			break
		}
	}

	if filename == "" {
		fmt.Fprintln(os.Stderr, "No input file specified")
		os.Exit(1)
	}

	// Probe the file
	info, err := probe.ProbeFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error probing file: %v\n", err)
		os.Exit(1)
	}

	// Get file size
	stat, _ := os.Stat(filename)
	fileSize := int64(0)
	if stat != nil {
		fileSize = stat.Size()
	}

	// Build ffprobe-compatible output
	output := FFProbeOutput{
		Format: Format{
			Filename:   filename,
			FormatName: info.ContentType,
			Duration:   fmt.Sprintf("%.6f", info.Duration.Seconds()),
			Size:       fmt.Sprintf("%d", fileSize),
			BitRate:    fmt.Sprintf("%d", info.Bitrate),
		},
	}

	// Add video stream if we have dimensions
	if info.Width > 0 && info.Height > 0 {
		output.Streams = append(output.Streams, Stream{
			Index:     0,
			CodecType: "video",
			CodecName: info.VideoCodec,
			Width:     info.Width,
			Height:    info.Height,
			Duration:  fmt.Sprintf("%.6f", info.Duration.Seconds()),
			BitRate:   fmt.Sprintf("%d", info.Bitrate),
		})
	}

	// Add audio stream
	if info.AudioCodec != "" {
		output.Streams = append(output.Streams, Stream{
			Index:     len(output.Streams),
			CodecType: "audio",
			CodecName: info.AudioCodec,
			Duration:  fmt.Sprintf("%.6f", info.Duration.Seconds()),
		})
	}

	// Output JSON
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(output)
}
