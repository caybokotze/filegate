package probe

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/abema/go-mp4"
	"github.com/dhowden/tag"
)

// MediaInfo contains probed media information
type MediaInfo struct {
	Duration    time.Duration
	Width       int
	Height      int
	VideoCodec  string
	AudioCodec  string
	Bitrate     int64
	ContentType string
}

// ProbeFile probes a media file and returns its information
func ProbeFile(path string) (*MediaInfo, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".mp4", ".m4v", ".mov":
		return probeMP4(path)
	case ".mp3", ".flac", ".ogg", ".m4a":
		return probeAudio(path)
	case ".mkv", ".webm":
		return probeMKV(path)
	default:
		return nil, fmt.Errorf("unsupported format: %s", ext)
	}
}

func probeMP4(path string) (*MediaInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info := &MediaInfo{
		ContentType: "video/mp4",
	}

	// Get file size for bitrate calculation
	stat, _ := f.Stat()
	fileSize := stat.Size()

	// Parse MP4 boxes
	boxes, err := mp4.ExtractBoxWithPayload(f, nil, mp4.BoxPath{mp4.BoxTypeMoov(), mp4.BoxTypeMvhd()})
	if err == nil && len(boxes) > 0 {
		if mvhd, ok := boxes[0].Payload.(*mp4.Mvhd); ok {
			// Duration is in timescale units
			if mvhd.Timescale > 0 {
				durationSec := float64(mvhd.DurationV0) / float64(mvhd.Timescale)
				if mvhd.Version == 1 {
					durationSec = float64(mvhd.DurationV1) / float64(mvhd.Timescale)
				}
				info.Duration = time.Duration(durationSec * float64(time.Second))

				// Calculate bitrate
				if durationSec > 0 {
					info.Bitrate = int64(float64(fileSize*8) / durationSec)
				}
			}
		}
	}

	// Try to get video track info
	f.Seek(0, 0)
	boxes, err = mp4.ExtractBoxWithPayload(f, nil, mp4.BoxPath{mp4.BoxTypeMoov(), mp4.BoxTypeTrak(), mp4.BoxTypeTkhd()})
	if err == nil {
		for _, box := range boxes {
			if tkhd, ok := box.Payload.(*mp4.Tkhd); ok {
				// Use helper methods for width/height
				w := int(tkhd.GetWidthInt())
				h := int(tkhd.GetHeightInt())
				if w > 0 && h > 0 {
					info.Width = w
					info.Height = h
					break
				}
			}
		}
	}

	return info, nil
}

func probeAudio(path string) (*MediaInfo, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	info := &MediaInfo{
		ContentType: "audio/mpeg",
	}

	// Use tag library for audio metadata
	m, err := tag.ReadFrom(f)
	if err == nil {
		if m.FileType() == tag.MP3 {
			info.ContentType = "audio/mpeg"
			info.AudioCodec = "mp3"
		} else if m.FileType() == tag.M4A || m.FileType() == tag.M4B || m.FileType() == tag.M4P {
			info.ContentType = "audio/mp4"
			info.AudioCodec = "aac"
		} else if m.FileType() == tag.FLAC {
			info.ContentType = "audio/flac"
			info.AudioCodec = "flac"
		} else if m.FileType() == tag.OGG {
			info.ContentType = "audio/ogg"
			info.AudioCodec = "vorbis"
		}
	}

	// For MP3, try to get duration from file size and bitrate estimation
	stat, _ := f.Stat()
	fileSize := stat.Size()

	// Rough estimation: assume 128kbps for MP3 if we can't determine
	if strings.HasSuffix(strings.ToLower(path), ".mp3") {
		// Very rough estimate
		estimatedBitrate := int64(128000) // 128 kbps default
		info.Bitrate = estimatedBitrate
		info.Duration = time.Duration(float64(fileSize*8) / float64(estimatedBitrate) * float64(time.Second))
	}

	return info, nil
}

func probeMKV(path string) (*MediaInfo, error) {
	// Basic MKV support - just return minimal info
	// Full MKV parsing would require ebml library
	stat, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	info := &MediaInfo{
		ContentType: "video/x-matroska",
		VideoCodec:  "unknown",
	}

	// Very rough duration estimate based on typical video bitrate
	// Assume ~5 Mbps for video
	estimatedBitrate := int64(5000000)
	info.Bitrate = estimatedBitrate
	info.Duration = time.Duration(float64(stat.Size()*8) / float64(estimatedBitrate) * float64(time.Second))

	return info, nil
}
