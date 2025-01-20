package rustmaps

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/maintc/rustmaps-cli/pkg/common"
)

func (g *Generator) IsApiKeySet() bool {
	return g.config.APIKey != ""
}

func (g *Generator) Pending() bool {
	for _, m := range g.maps {
		if m.Status == common.StatusPending {
			return true
		}
	}
	return false
}

func (g *Generator) Generating() bool {
	for _, m := range g.maps {
		if m.Status == common.StatusGenerating {
			return true
		}
	}
	return false
}

func (g *Generator) ContainCustomMaps() bool {
	for _, m := range g.maps {
		if m.SavedConfig != "" {
			return true
		}
	}
	return false
}

func (g *Generator) GetRandomSeed() string {
	// Seed the random number generator using the current time
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Define the maximum seed value
	const maxSeed = 2147483647

	// Generate a random integer between 0 and maxSeed (inclusive)
	seed := rng.Intn(maxSeed + 1)

	// Convert the integer to a string
	return fmt.Sprintf("%d", seed)
}

func parseInt(s string) int {
	var i int
	fmt.Sscanf(s, "%d", &i)
	return i
}

func parseBool(s string) bool {
	return strings.ToLower(s) == "true"
}
