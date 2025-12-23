package relay

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

var adjectives = []string{
	"brave", "calm", "clever", "eager", "fancy",
	"gentle", "happy", "jolly", "kind", "lively",
	"merry", "noble", "proud", "quick", "quiet",
	"rapid", "sharp", "sleek", "smart", "smooth",
	"soft", "solid", "swift", "tidy", "vivid",
	"warm", "wild", "wise", "witty", "young",
	"amber", "azure", "coral", "crisp", "dusk",
	"early", "faint", "frost", "gold", "green",
	"hazy", "ivory", "jade", "keen", "lime",
	"lunar", "maple", "misty", "navy", "olive",
	"pearl", "plum", "polar", "prime", "ruby",
	"sage", "silk", "snow", "solar", "steel",
	"stone", "sunny", "teal", "titan", "topaz",
	"ultra", "vast", "velvet", "vital", "white",
	"zesty", "agile", "bold", "bright", "cool",
	"daring", "elite", "fair", "fresh", "grand",
	"great", "hardy", "ideal", "just", "light",
	"lucky", "major", "mint", "neat", "new",
	"nice", "open", "pale", "plain", "pure",
	"rare", "real", "rich", "royal", "safe",
}

var nouns = []string{
	"tiger", "river", "cloud", "maple", "falcon",
	"ocean", "forest", "mountain", "valley", "meadow",
	"stream", "canyon", "desert", "island", "glacier",
	"thunder", "breeze", "storm", "sunset", "dawn",
	"eagle", "wolf", "bear", "fox", "hawk",
	"lion", "lynx", "otter", "panda", "raven",
	"shark", "whale", "cobra", "crane", "dove",
	"finch", "heron", "horse", "koala", "lemur",
	"moose", "owl", "robin", "salmon", "seal",
	"sparrow", "swan", "trout", "turtle", "zebra",
	"acorn", "aspen", "birch", "cedar", "daisy",
	"elm", "fern", "grove", "holly", "iris",
	"jasper", "kelp", "lotus", "moss", "oak",
	"orchid", "palm", "pine", "poppy", "reed",
	"rose", "sage", "spruce", "tulip", "vine",
	"willow", "amber", "basalt", "cliff", "coral",
	"crystal", "dune", "ember", "flint", "gem",
	"granite", "jade", "lava", "marble", "opal",
	"pearl", "quartz", "ruby", "sand", "slate",
	"stone", "topaz", "bronze", "chrome", "copper",
}

// GenerateSubdomain creates a random word-based subdomain like "brave-tiger"
func GenerateSubdomain() (string, error) {
	adjIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(adjectives))))
	if err != nil {
		return "", err
	}

	nounIdx, err := rand.Int(rand.Reader, big.NewInt(int64(len(nouns))))
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", adjectives[adjIdx.Int64()], nouns[nounIdx.Int64()]), nil
}
