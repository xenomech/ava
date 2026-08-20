package hub

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"ava/api/internal/dto"
	"ava/api/internal/model"
)

const userCodeAlphabet = "BCDFGHJKLMNPQRSTVWXZ23456789"

func generateSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("generate secret: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func generateUserCode() (string, error) {
	var code strings.Builder

	for i := range 8 {
		if i == 4 {
			code.WriteByte('-')
		}

		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(userCodeAlphabet))))
		if err != nil {
			return "", fmt.Errorf("generate user code: %w", err)
		}

		code.WriteByte(userCodeAlphabet[n.Int64()])
	}

	return code.String(), nil
}

func normalizeUserCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func toHubResponse(hub *model.Hub) *dto.HubResponse {
	return &dto.HubResponse{
		ID:         hub.ID,
		Name:       hub.Name,
		Status:     string(hub.Status),
		Online:     hub.IsOnline(),
		LastSeenAt: hub.LastSeenAt,
		CreatedAt:  hub.CreatedAt,
	}
}
