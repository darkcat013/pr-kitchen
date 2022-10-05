package domain

import (
	"encoding/json"
	"io"
	"os"
	"sort"

	"github.com/darkcat013/pr-kitchen/utils"
)

type CookInfo struct {
	Rank        int    `json:"rank"`
	Proficiency int    `json:"proficiency"`
	Name        string `json:"name"`
	CatchPhrase string `json:"catch-phrase"`
}

func InitializeCooks(jsonPath string) {
	file, err := os.Open(jsonPath)
	if err != nil {
		utils.Log.Fatal("Error opening " + jsonPath)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	var cookInfos []CookInfo

	json.Unmarshal(bytes, &cookInfos)

	if cookInfos == nil {
		utils.Log.Fatal("Failed to decode cooks from " + jsonPath)
	}

	sort.Slice(cookInfos, func(i, j int) bool {
		if cookInfos[i].Rank < cookInfos[j].Rank {
			return true
		}
		return cookInfos[i].Proficiency < cookInfos[j].Proficiency
	})

	Cooks = make([]*Cook, 0, len(cookInfos))

	for i := 0; i < len(cookInfos); i++ {
		cook := NewCook(i, &cookInfos[i])
		Cooks = append(Cooks, cook)
	}

	utils.Log.Info("Cooks decoded and set")
}
