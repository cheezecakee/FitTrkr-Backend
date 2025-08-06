package playlist

import (
	"fmt"
	"strings"
)

// Validation methods for custom types
func (v Visibility) IsValid() bool {
	switch v {
	case VisibilityPrivate, VisibilityPublic, VisibilityUnlisted:
		return true
	}
	return false
}

func (bt BlockType) IsValid() bool {
	switch bt {
	case BlockTypeStandard, BlockTypeSuperset, BlockTypeTriset, BlockTypeCircuit,
		BlockTypeDropset, BlockTypeCardio, BlockTypeWarmup, BlockTypeCooldown:
		return true
	}
	return false
}

// Smart block naming based on context
func GenerateBlockName(blockType BlockType, blockOrder int, hasExercises bool) string {
	// If it's the first block and has no exercises yet, use a primary name
	if blockOrder == 1 && !hasExercises {
		switch blockType {
		case BlockTypeWarmup:
			return "Warm-up"
		case BlockTypeCardio:
			return "Cardio"
		case BlockTypeCooldown:
			return "Cool-down"
		default:
			return "Main Block" // Instead of "Workout"
		}
	}

	// For subsequent blocks, use descriptive names
	switch blockType {
	case BlockTypeStandard:
		return fmt.Sprintf("Block %d", blockOrder)
	case BlockTypeSuperset:
		return fmt.Sprintf("Superset %d", blockOrder)
	case BlockTypeTriset:
		return fmt.Sprintf("Triset %d", blockOrder)
	case BlockTypeCircuit:
		return fmt.Sprintf("Circuit %d", blockOrder)
	case BlockTypeDropset:
		return fmt.Sprintf("Dropset %d", blockOrder)
	case BlockTypeCardio:
		return fmt.Sprintf("Cardio %d", blockOrder)
	case BlockTypeWarmup:
		return "Warm-up"
	case BlockTypeCooldown:
		return "Cool-down"
	default:
		return fmt.Sprintf("Block %d", blockOrder)
	}
}

// Smart block type suggestion based on context
func SuggestBlockType(exerciseCount int, isFirstBlock bool, playlistTitle string) BlockType {
	// Check playlist title for hints
	title := strings.ToLower(playlistTitle)
	if strings.Contains(title, "cardio") {
		return BlockTypeCardio
	}
	if strings.Contains(title, "warmup") || strings.Contains(title, "warm-up") {
		return BlockTypeWarmup
	}
	if strings.Contains(title, "cooldown") || strings.Contains(title, "cool-down") {
		return BlockTypeCooldown
	}

	// For first block, default to standard
	if isFirstBlock {
		return BlockTypeStandard
	}

	// Suggest based on exercise count
	switch exerciseCount {
	case 2:
		return BlockTypeSuperset
	case 3:
		return BlockTypeTriset
	case 4, 5, 6:
		return BlockTypeCircuit
	default:
		return BlockTypeStandard
	}
}
