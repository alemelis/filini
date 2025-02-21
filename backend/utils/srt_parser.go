package utils

import (
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/alemelis/filini/models"
)

var timestampRegex = regexp.MustCompile(`(\d{2}):(\d{2}):(\d{2}),(\d{3})`)

// parseTimestamp converts "HH:MM:SS,MS" to float seconds.
func parseTimestamp(s string) (float64, error) {
	matches := timestampRegex.FindStringSubmatch(s)
	if len(matches) != 5 {
		return 0, nil
	}
	hours, _ := strconv.Atoi(matches[1])
	minutes, _ := strconv.Atoi(matches[2])
	seconds, _ := strconv.Atoi(matches[3])
	milliseconds, _ := strconv.Atoi(matches[4])

	totalSeconds := float64(hours*3600+minutes*60+seconds) + float64(milliseconds)/1000
	return totalSeconds, nil
}

// ParseSRT reads an .srt file and extracts subtitles.
func ParseSRT(reader io.Reader, videoID int) ([]models.Subtitle, error) {
	var subtitles []models.Subtitle
	scanner := bufio.NewScanner(reader)

	var id int
	var startTime, endTime float64
	var textLines []string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// If the line is empty, finalize the subtitle block
		if line == "" && len(textLines) > 0 {
			subtitles = append(subtitles, models.Subtitle{
				VideoID:   videoID,
				Text:      strings.Join(textLines, " "),
				StartTime: startTime,
				EndTime:   endTime,
			})
			textLines = nil
			continue
		}

		// Detect subtitle ID
		if id == 0 {
			id, _ = strconv.Atoi(line)
			continue
		}

		// Detect timestamps
		if strings.Contains(line, "-->") {
			parts := strings.Split(line, " --> ")
			startTime, _ = parseTimestamp(parts[0])
			endTime, _ = parseTimestamp(parts[1])
			continue
		}

		// Collect subtitle text
		textLines = append(textLines, line)
	}

	// Handle last subtitle block
	if len(textLines) > 0 {
		subtitles = append(subtitles, models.Subtitle{
			VideoID:   videoID,
			Text:      strings.Join(textLines, " "),
			StartTime: startTime,
			EndTime:   endTime,
		})
	}

	return subtitles, scanner.Err()
}
