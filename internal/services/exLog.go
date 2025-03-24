package services

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cheezecakee/FitLogr/internal/models"
)

const (
	EntryDelimiter    = "!"
	ExerciseLogPrefix = "["
	ExerciseLogSuffix = "]"
	WorkoutLogPrefix  = "L:"
	EndLogMarker      = "!"
)

func EncodeWorkoutStart(numberOfExercises string) string {
	var sb strings.Builder
	// L:0![]
	sb.WriteString(fmt.Sprintf("%s%s%s%s%s", WorkoutLogPrefix, numberOfExercises, EndLogMarker, ExerciseLogPrefix, ExerciseLogSuffix))
	return sb.String()
}

func EncodeExerciseLog(log models.ExerciseLog) string {
	var sb strings.Builder
	sb.WriteString(ExerciseLogPrefix)

	sb.WriteString(fmt.Sprintf("s=%d%s", log.Sets, EntryDelimiter))
	sb.WriteString(fmt.Sprintf("r=%d%s", log.Reps, EntryDelimiter))
	sb.WriteString(fmt.Sprintf("w=%f%s", log.Weight, EntryDelimiter))
	sb.WriteString(fmt.Sprintf("i=%04d%s", log.Interval, EntryDelimiter))
	sb.WriteString(fmt.Sprintf("r=%04d%s", log.Rest, EntryDelimiter))

	completed := 0
	if log.Completed {
		completed = 1
	}
	sb.WriteString(fmt.Sprintf("c=%d%s", completed, EntryDelimiter))

	// Wrap notes with delimiters
	sb.WriteString(fmt.Sprintf("n=%s%s%s%s", NoteValueDelimiter, log.Notes, NoteValueDelimiter, EntryDelimiter))
	sb.WriteString(fmt.Sprintf("id=%03d%s", log.ExerciseID, EntryDelimiter))

	sb.WriteString(EndLogMarker)
	return sb.String()
}

func CombineExerciseLogs(logs []string) string {
	return strings.Join(logs, "")
}

func parseField(field string) (string, string, error) {
	parts := strings.SplitN(field, "=", 2)
	if len(parts) != 2 {
		return "", "", errors.New("invalid field format")
	}
	return parts[0], parts[1], nil
}

func DecodeExerciseLog(logStr string) (models.ExerciseLog, error) {
	log := models.ExerciseLog{}

	if !strings.HasPrefix(logStr, ExerciseLogPrefix) {
		return log, errors.New("invalid exercise log format")
	}

	logStr = strings.TrimPrefix(logStr, ExerciseLogPrefix)
	logStr = strings.TrimSuffix(logStr, EndLogMarker)

	fields := strings.Split(logStr, EntryDelimiter)

	for _, field := range fields {
		if field == "" {
			continue
		}

		key, value, err := parseField(field)
		if err != nil {
			continue
		}

		switch key {
		case "s":
			log.Sets, _ = strconv.Atoi(value)
		case "r":
			if strings.HasPrefix(field, "r=") && !strings.HasPrefix(field, "r=0") {
				log.Reps, _ = strconv.Atoi(value)
			}
		case "w":
			log.Weight, _ = strconv.ParseFloat(value, 32)
		case "i":
			log.Interval, _ = strconv.Atoi(value)
		case "c":
			completed, _ := strconv.Atoi(value)
			log.Completed = completed == 1
		case "n":
			if strings.HasPrefix(value, NoteValueDelimiter) && strings.HasSuffix(value, NoteValueDelimiter) {
				log.Notes = value[1 : len(value)-1]
			}
		case "id":
			log.ExerciseID, _ = strconv.Atoi(value)
		case "t":
			log.Timestamp, _ = strconv.ParseInt(value, 10, 64)
		}
	}

	return log, nil
}
