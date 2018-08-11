package models

import (
	"time"
	"strings"
	"regexp"
	"github.com/araddon/dateparse"
)

var (
	formatTypes = []string{"FIRST_FORMAT", "SECOND_FORMAT"}
)

type LogRecord struct {
	LogTime    time.Time `bson:"log_time"`
	LogMsg       string    `bson:"log_msg"`
	FileName       string    `bson:"file_name"`
	LogFormat       string    `bson:"log_format"`
}

func NewRecord(line, filepath string) LogRecord {
	parts := strings.Split(line, " | ")
	if len(parts) != 2 {
		panic("len(parts) != 2 for line '" + line + "'")
	}
	return LogRecord{
		LogTime:   toTime(parts[0]),
		LogMsg:    parts[1],
		FileName:  filepath,
		LogFormat: detectFormatType(parts[0]),
	}
}

func detectFormatType(timeStr string) string {
	var patterns = []*regexp.Regexp{
		//"Feb 1, 2018 at 3:04:05pm (UTC)",
		regexp.MustCompile("[A-Z][a-z]{2} \\d{1,2}, \\d{4} at \\d{1,2}:\\d{2}:\\d{2}(am|pm) \\([A-Z]{3}\\)"),
		//"2018-02-01T15:04:05Z",
		regexp.MustCompile("\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}Z"),
	}
	for i,p := range patterns {
		if p.MatchString(timeStr) {
			return formatTypes[i]
		}
	}
	return "Unrecognized format [" + timeStr + "]"
}

func toTime(timeStr string) time.Time {
	t, err := dateparse.ParseLocal(timeStr)
	if err != nil {
		panic(err)
	}
	return t
}
