package mtgtracker

import (
	"mtgtracker/internal/repository"
	"testing"
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
				{PlayerID: strPtr("player2"), Position: 0}, // Position ignored
				{PlayerID: strPtr("player1"), Position: 0}, // Position ignored
				{PlayerID: strPtr("player3"), Position: 0}, // Position ignored
			},
			existingRankings: []repository.Ranking{
				{PlayerID: strPtr("player1"), Position: 0},
				{PlayerID: strPtr("player2"), Position: 0},
				{PlayerID: strPtr("player3"), Position: 0},
			},
			expectError:       false,
			expectedPositions: []int{1, 2, 3}, // Sequential based on request order
		},
		{
			name: "mismatched count - too few in request",
			requestRankings: []UpdateRanking{
				{PlayerID: strPtr("player1"), Position: 0},
			},
			existingRankings: []repository.Ranking{
				{PlayerID: strPtr("player1"), Position: 0},
				{PlayerID: strPtr("player2"), Position: 0},
			},
			expectError:  true,
			errorMessage: "rankings count must match existing rankings",
		},
		{
			name: "mismatched count - too many in request",
			requestRankings: []UpdateRanking{
				{PlayerID: strPtr("player1"), Position: 0},
				{PlayerID: strPtr("player2"), Position: 0},
				{PlayerID: strPtr("player3"), Position: 0},
			},
			existingRankings: []repository.Ranking{
				{PlayerID: strPtr("player1"), Position: 0},
				{PlayerID: strPtr("player2"), Position: 0},
			},
			expectError:  true,
			errorMessage: "rankings count must match existing rankings",
		},
		{
			name: "invalid player ID in request",
			requestRankings: []UpdateRanking{
				{PlayerID: strPtr("player1"), Position: 0},
				{PlayerID: strPtr("invalid_player"), Position: 0},
			},
			existingRankings: []repository.Ranking{
				{PlayerID: strPtr("player1"), Position: 0},
				{PlayerID: strPtr("player2"), Position: 0},
			},
			expectError:  true,
			errorMessage: "invalid player ID in rankings",
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
				{PlayerID: strPtr("player1"), Position: 0},
			},
			existingRankings: []repository.Ranking{
				{PlayerID: strPtr("player1"), Position: 0},
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

			// Verify the result maintains the request order for PlayerIDs
			if !tt.expectError && len(tt.requestRankings) > 0 {
				for i, reqRanking := range tt.requestRankings {
					if reqRanking.PlayerID == nil || result[i].PlayerID == nil {
						t.Errorf("nil PlayerID found at index %d", i)
						continue
					}
					if *result[i].PlayerID != *reqRanking.PlayerID {
						t.Errorf("expected PlayerID %s at index %d, got %s", *reqRanking.PlayerID, i, *result[i].PlayerID)
					}
				}
			}
		})
	}
}
