package playlist

import (
	"slices"

	"github.com/lib/pq"
)

func isUniqueConstraintError(err error) bool {
	// For PostgreSQL with pq driver
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505" // unique_violation
	}
	return false
}

// isValidBlockType validates if the block type is supported
func isValidBlockType(blockType string) bool {
	validTypes := []BlockType{BlockTypePlaylist, BlockTypeSuperset, BlockTypeCircuit, BlockTypeDropset}
	return slices.Contains(validTypes, BlockType(blockType))
}
