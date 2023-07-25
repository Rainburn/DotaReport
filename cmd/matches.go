package main

import (
	"strconv"
	"net/http"
	"log"
	"encoding/json"
	"time"
	"fmt"
)

type MatchData struct {
	MatchID int `json:"match_id"`
	PlayersData []PlayerData `json:"players"`
}

type PlayerData struct {
	AccountID int `json:"account_id"`
	PersonName string `json:"personaname"`
	IsRadiant bool `json:"isRadiant"`
	HeroID int `json:"hero_id"`
	Kills int `json:"kills"`
	Deaths int `json:"deaths"`
	Assists int `json:"assists"`
	Gpm int `json:"gold_per_min"`
	Xpm int `json:"xp_per_min"`
	Win int `json:"win"` // 1 for Win, 0 for Lose
}

type PeerStat map[int]Peer

type Peer struct {
	AccountID int
	Name string
	WinLoseStat
}

type PeerSummary struct {
	Peer
	TotalMatch int
}


func handleMatch(matches []Matches) (<-chan MatchData) {
	chanMatchesData := make(chan MatchData)
	

	go func(chanMatchesData chan MatchData) {
		for idx, each := range(matches) {
			var matchData MatchData

			shouldGo := make(chan bool)

			go func(flag chan bool) {
				time.Sleep(time.Duration(1.1 * float64(time.Second)))
				shouldGo <- true
			}(shouldGo)

			matchIDStr := strconv.Itoa(each.MatchID)
			fmt.Printf("(%d/%d) - Handling matchid: %s\n", idx+1, len(matches), matchIDStr)

			getMatchAPILink := getMatchDataAPI(matchIDStr)
			resp, err := http.Get(getMatchAPILink)
			if err != nil {
				log.Fatal(err.Error())
			}
			defer resp.Body.Close()

			err = json.NewDecoder(resp.Body).Decode(&matchData)
			if err != nil {
				log.Fatal(err.Error())
			}

			// Insert into channel
			chanMatchesData <- matchData

			flag := <-shouldGo
			if (flag) {
				continue
			}
			// Sleep to avoid too many API calls in 1 minute
			// time.Sleep(time.Duration(1.1 * float64(time.Second)))
		}

		close(chanMatchesData)
	}(chanMatchesData)
	
	return chanMatchesData
}

func processMatchData(WLStat *WinLoseStat, benchmark *Benchmark, heroStats *HeroStats, peerStat *PeerStat, chanMatchesData <-chan MatchData) {
	for matchData := range(chanMatchesData) {
		
		var isRadiant bool
		var wonGame bool

		// Find self-user
		for _, playerData := range(matchData.PlayersData) {
			if (playerData.AccountID == USER_ID) {
				win := playerData.Win
				heroIDStr := strconv.Itoa(playerData.HeroID)
				isRadiant = playerData.IsRadiant

				if (win == 1) {
					heroStats.addWin(heroIDStr)
					WLStat.addWin()
					wonGame = true
				} else {
					heroStats.addLose(heroIDStr)
					WLStat.addLose()
					wonGame = false
				}

				// Benchmarks
				kills := playerData.Kills
				deaths := playerData.Deaths
				assists := playerData.Assists
				gpm := playerData.Gpm
				xpm := playerData.Xpm
				benchmark.addStats(kills, deaths, assists, gpm, xpm)

				break
			}
		}

		// Find teammates
		for _, playerData := range(matchData.PlayersData) {
			if (playerData.IsRadiant == isRadiant && playerData.AccountID != USER_ID) {
				// This user is in the same team as target user but not the user itself
				existingPeer, exists := (*peerStat)[playerData.AccountID]

				if !exists {
					newPeer := Peer{Name: playerData.PersonName,}
					if (wonGame) {
						(&newPeer).addWin()
					} else {
						(&newPeer).addLose()
					}

					(*peerStat)[playerData.AccountID] = newPeer

				} else {
					if (wonGame) {
						(&existingPeer).addWin()
					} else {
						(&existingPeer).addLose()
					}
					(*peerStat)[playerData.AccountID] = existingPeer
				}

			} else {
				continue
			}
		}

	}
}