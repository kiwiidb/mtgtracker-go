package mtgtracker

import (
	"mtgtracker/internal/repository"
	"testing"

	"gorm.io/gorm"
)

func TestValidateAndReorderRankings(t *testing.T) {
	// Helper function to create string pointers
	strPtr := func(s string) *string { return &s }

	tests := []struct {
		name              string
		requestRankings   []UpdateRanking
		existingRankings  []repository.Ranking
		expectError       bool
		errorMessage      string
		expectedPositions []int
	}{
		{
			name: "valid reordering",
			requestRankings: []UpdateRanking{
				{RankingID: 2, Position: 0}, // Position ignored
				{RankingID: 1, Position: 0}, // Position ignored
				{RankingID: 3, Position: 0}, // Position ignored
			},
			existingRankings: []repository.Ranking{
				{Model: gorm.Model{ID: 1}, PlayerID: strPtr("player1"), Position: 0},
				{Model: gorm.Model{ID: 2}, PlayerID: strPtr("player2"), Position: 0},
				{Model: gorm.Model{ID: 3}, PlayerID: strPtr("player3"), Position: 0},
			},
			expectError:       false,
			expectedPositions: []int{1, 2, 3}, // Sequential based on request order
		},
		{
			name: "mismatched count - too few in request",
			requestRankings: []UpdateRanking{
				{RankingID: 1, Position: 0},
			},
			existingRankings: []repository.Ranking{
				{Model: gorm.Model{ID: 1}, PlayerID: strPtr("player1"), Position: 0},
				{Model: gorm.Model{ID: 2}, PlayerID: strPtr("player2"), Position: 0},
			},
			expectError:  true,
			errorMessage: "rankings count must match existing rankings",
		},
		{
			name: "mismatched count - too many in request",
			requestRankings: []UpdateRanking{
				{RankingID: 1, Position: 0},
				{RankingID: 2, Position: 0},
				{RankingID: 3, Position: 0},
			},
			existingRankings: []repository.Ranking{
				{Model: gorm.Model{ID: 1}, PlayerID: strPtr("player1"), Position: 0},
				{Model: gorm.Model{ID: 2}, PlayerID: strPtr("player2"), Position: 0},
			},
			expectError:  true,
			errorMessage: "rankings count must match existing rankings",
		},
		{
			name: "invalid ranking ID in request",
			requestRankings: []UpdateRanking{
				{RankingID: 1, Position: 0},
				{RankingID: 999, Position: 0}, // Invalid ranking ID
			},
			existingRankings: []repository.Ranking{
				{Model: gorm.Model{ID: 1}, PlayerID: strPtr("player1"), Position: 0},
				{Model: gorm.Model{ID: 2}, PlayerID: strPtr("player2"), Position: 0},
			},
			expectError:  true,
			errorMessage: "invalid ranking ID in rankings",
		},
		{
			name:              "empty rankings",
			requestRankings:   []UpdateRanking{},
			existingRankings:  []repository.Ranking{},
			expectError:       false,
			expectedPositions: []int{},
		},
		{
			name: "single ranking",
			requestRankings: []UpdateRanking{
				{RankingID: 1, Position: 0},
			},
			existingRankings: []repository.Ranking{
				{Model: gorm.Model{ID: 1}, PlayerID: strPtr("player1"), Position: 0},
			},
			expectError:       false,
			expectedPositions: []int{1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := validateAndReorderRankings(tt.requestRankings, tt.existingRankings)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if err.Error() != tt.errorMessage {
					t.Errorf("expected error message '%s', got '%s'", tt.errorMessage, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expectedPositions) {
				t.Errorf("expected %d rankings, got %d", len(tt.expectedPositions), len(result))
				return
			}

			for i, expectedPos := range tt.expectedPositions {
				if result[i].Position != expectedPos {
					t.Errorf("expected position %d at index %d, got %d", expectedPos, i, result[i].Position)
				}
			}

			// Verify the result maintains the request order for RankingIDs
			if !tt.expectError && len(tt.requestRankings) > 0 {
				for i, reqRanking := range tt.requestRankings {
					if result[i].ID != reqRanking.RankingID {
						t.Errorf("expected RankingID %d at index %d, got %d", reqRanking.RankingID, i, result[i].ID)
					}
				}
			}
		})
	}
}
