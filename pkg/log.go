package pkg

import (
	"time"

	"github.com/google/uuid"
)

type Set struct {
	Number   int    // s:<setNumber>
	Reps     int    // r:<reps>
	Weight   int    // w:<weight>
	Interval string // i:<interval>
	Rest     string // R:<rest>
	Note     string // n=<length>:<note>
}

type Exercise struct {
	ID           string // e:<exerciseID>
	TotalSets    int    // S:<TotalSets>
	Sets         []Set
	SupersetWith *Exercise // For supersets
}

type WorkoutSession struct {
	UserID    uuid.UUID `json:"user_id"`
	Exercises []Exercise
	EndedAt   *time.Time `json:"ended_at"`
	CreatedAt time.Time  `json:"created_at"`
}

type Logger struct {
	CurrentSession  *WorkoutSession
	CurrentExercise *Exercise
}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) StartSessionWorkout() {
	l.CurrentSession = &WorkoutSession{
		Exercises: []Exercise{},
	}
}

func (l *Logger) StartSessionExercise(exerciseID string, TotalSets int) {
	l.CurrentExercise = &Exercise{
		ID:        exerciseID,
		TotalSets: TotalSets,
		Sets:      []Set{},
	}
}

func (l *Logger) LogSet(setNumber, reps, weight int, interval, rest, note string) {
	if l.CurrentExercise == nil {
		return
	}

	set := Set{
		Number:   setNumber,
		Reps:     reps,
		Weight:   weight,
		Interval: interval,
		Rest:     rest,
		Note:     note,
	}

	l.CurrentExercise.Sets = append(l.CurrentExercise.Sets, set)
}

func (l *Logger) FinishExercise() {
	if l.CurrentSession == nil || l.CurrentExercise == nil {
		return
	}

	l.CurrentSession.Exercises = append(l.CurrentSession.Exercises, *l.CurrentExercise)
	l.CurrentExercise = nil
}

func (l *Logger) AddSuperSet(exerciseID string) {
	if l.CurrentExercise == nil {
		return
	}

	supersetExercise := &Exercise{
		ID:        exerciseID,
		TotalSets: l.CurrentExercise.TotalSets,
		Sets:      []Set{},
	}

	l.CurrentExercise.SupersetWith = supersetExercise
}

func (l *Logger) LogSuperSet(setNumber, reps, weight int, interval string) {
	if l.CurrentExercise == nil || l.CurrentExercise.SupersetWith == nil {
		return
	}

	set := Set{
		Number:   setNumber,
		Reps:     reps,
		Weight:   weight,
		Interval: interval,
	}

	l.CurrentExercise.SupersetWith.Sets = append(l.CurrentExercise.SupersetWith.Sets, set)
}
