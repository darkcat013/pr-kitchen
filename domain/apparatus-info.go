package domain

import (
	"encoding/json"
	"io"
	"os"

	"github.com/darkcat013/pr-kitchen/utils"
)

type ApparatusInfo struct {
	Name   string `json:"name"`
	Amount int    `json:"amount"`
}

func InitializeApparatuses(jsonPath string) {
	file, err := os.Open(jsonPath)
	if err != nil {
		utils.Log.Fatal("Error opening " + jsonPath)
	}
	defer file.Close()

	bytes, _ := io.ReadAll(file)
	var apparatusInfos []ApparatusInfo

	json.Unmarshal(bytes, &apparatusInfos)

	if apparatusInfos == nil {
		utils.Log.Fatal("Failed to decode apparatuses from " + jsonPath)
	}

	Apparatuses = make(map[string][]*Apparatus)

	for i := 0; i < len(apparatusInfos); i++ {
		specificApparatuses := make([]*Apparatus, 0, apparatusInfos[i].Amount)
		ApparatusesChans[apparatusInfos[i].Name] = make(chan ApparatusFoodInfo, 100)
		for j := 0; j < apparatusInfos[i].Amount; j++ {
			apparatus := NewApparatus(j, apparatusInfos[i].Name, ApparatusesChans[apparatusInfos[i].Name])
			specificApparatuses = append(specificApparatuses, apparatus)
		}
		Apparatuses[apparatusInfos[i].Name] = specificApparatuses
	}
	utils.Log.Info("Apparatuses decoded and set")
}
