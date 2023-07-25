package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"log"
	"os"
	"strconv"
)

type Hero struct {
	ID int `json:"id"`
	Name string `json:"localized_name"`
}

type HeroSummary struct {
	Hero
	WinLoseStat
	TotalMatch int
}

type HeroStats map[string]WinLoseStat

func (hs *HeroStats) addWin(heroID string) {
	// Check if hero stats is already exist
	existingStat, exists := (*hs)[heroID]
	if !exists {
		winLoseStat := new(WinLoseStat)
		winLoseStat.addWin()
		(*hs)[heroID] = *winLoseStat
	} else {
		(&existingStat).addWin()
		(*hs)[heroID] = existingStat
	}

	return
}

func (hs *HeroStats) addLose(heroID string) {
	// Check if hero stats is already exist
	existingStat, exists := (*hs)[heroID]
	if !exists {
		winLoseStat := new(WinLoseStat)
		winLoseStat.addLose()
		(*hs)[heroID] = *winLoseStat
	} else {
		(&existingStat).addLose()
		(*hs)[heroID] = existingStat
	}
}



func loadHeroesFromJSON() map[string]Hero {
	heroes := map[string]Hero{}
	var listHeroes []Hero

	cwd, err := os.Getwd()
	basePath := filepath.Dir(cwd)

	body, err := ioutil.ReadFile(filepath.Join(basePath, ASSEST_PATH, "heroes.json"))
	if err != nil {
		log.Fatal("Cannot read heroes data from JSON")
	}

	err = json.Unmarshal(body, &listHeroes)

	for _, each := range(listHeroes) {
		heroIDStr := strconv.Itoa(each.ID)
		heroes[heroIDStr] = each
	}

	return heroes
}