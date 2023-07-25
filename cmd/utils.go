package main

import (
	"fmt"
	"sort"
	"log"
	"strconv"
)

type Benchmark struct {
	KillTotal int
	DeathTotal int
	AssistTotal int
	GpmTotal int
	XpmTotal int
}

type BenchmarkSummary struct {
	Kill float64
	Death float64
	Assist float64
	Gpm float64
	Xpm float64
}

func (b *Benchmark) addStats(kills, deaths, assists, gpm, xpm int) {
	b.KillTotal += kills
	b.DeathTotal += deaths
	b.AssistTotal += assists
	b.GpmTotal += gpm
	b.XpmTotal += xpm
}

type WinLoseStat struct {
	Win int
	Lose int
}

func (s *WinLoseStat) addWin() {
	s.Win += 1
	return
} 

func (s *WinLoseStat) addLose() {
	s.Lose += 1
	return
}

func generateMatchHistoryAPI(userIDStr string) string {
	return fmt.Sprintf("https://api.opendota.com/api/players/" + userIDStr + "/matches")
}

func getMatchDataAPI(matchIDStr string) string {
	return fmt.Sprintf("https://api.opendota.com/api/matches/" + matchIDStr)
}

func getProfileAPI(userIDStr string) string {
	return fmt.Sprintf("https://api.opendota.com/api/players/" + userIDStr)
}

func generateSummary(wlStat WinLoseStat, benchmark Benchmark, heroStats HeroStats, peersStat PeerStat, totalMatch int) (BenchmarkSummary, []HeroSummary, []PeerSummary) {
	var benchmarkSummary = BenchmarkSummary{}
	
	fmt.Printf("Stats for user: %d\n", USER_ID)
	fmt.Println("------------------------------------")
	// General Benchmark
	avgGpm := calculateAvg(benchmark.GpmTotal, totalMatch)
	avgXpm := calculateAvg(benchmark.XpmTotal, totalMatch)
	avgKills := calculateAvg(benchmark.KillTotal, totalMatch)
	avgDeaths := calculateAvg(benchmark.DeathTotal, totalMatch)
	avgAssists := calculateAvg(benchmark.AssistTotal, totalMatch)
	kda := calculateAvg(benchmark.KillTotal + benchmark.AssistTotal, benchmark.DeathTotal)

	// Insert into benchmark summary
	benchmarkSummary.Kill = avgKills
	benchmarkSummary.Death = avgDeaths
	benchmarkSummary.Assist = avgAssists
	benchmarkSummary.Gpm = avgGpm
	benchmarkSummary.Xpm = avgXpm


	wins := wlStat.Win
	loss := wlStat.Lose
	winrate := calculateAvg(wins, wins+loss)

	fmt.Printf("Winrate: %f | Wins: %d | Loss: %d | Matches: %d\n", winrate, wins, loss, wins+loss)
	fmt.Printf("GPM: %f | XPM: %f\n", avgGpm, avgXpm)
	fmt.Printf("Kills: %f | Deaths: %f | Assists: %f\n", avgKills, avgDeaths, avgAssists)
	fmt.Println("KDA ratio: ", kda)

	fmt.Println("Hero benchmark")
	// Hero Benchmark
	var heroesSummary = []HeroSummary{}
	for key, heroStat := range(heroStats) {
		wins := heroStat.Win
		loss := heroStat.Lose
		totalMatchEachHero := wins + loss

		// winrate := calculateAvg(wins, totalMatchEachHero)
		heroName := heroes[key].Name

		heroSummary := HeroSummary{}
		keyInInt, err := strconv.Atoi(key)
		if err != nil {
			log.Fatal(err.Error())
		}

		heroSummary.ID = keyInInt
		heroSummary.Name = heroName
		heroSummary.Win = wins
		heroSummary.Lose = loss
		heroSummary.TotalMatch = totalMatchEachHero

		heroesSummary = append(heroesSummary, heroSummary)

		
	}

	sort.Slice(heroesSummary, func(i, j int) bool {
		return heroesSummary[i].TotalMatch > heroesSummary[j].TotalMatch
	})

	// Peer benchmark
	var peersSummary = []PeerSummary{}
	for key, peerStat := range(peersStat) {
		peerSummary := PeerSummary{}
		peerSummary.AccountID = key
		peerSummary.Name = peerStat.Name
		peerSummary.Win = peerStat.Win
		peerSummary.Lose = peerStat.Lose
		peerSummary.TotalMatch = peerStat.Win + peerStat.Lose

		peersSummary = append(peersSummary, peerSummary)
	}

	sort.Slice(peersSummary, func(i, j int) bool {
		return peersSummary[i].TotalMatch > peersSummary[j].TotalMatch
	})

	// Remove private profile from peersSummary (ID == 0)
	for idx, each := range(peersSummary) {
		if (each.AccountID == 0) {
			peersSummary = append(peersSummary[:idx], peersSummary[idx+1:]...)
			break
		}
	}


	// Test heroes
	for _, heroSummary := range(heroesSummary) {
		heroID := heroSummary.ID
		heroName := heroSummary.Name
		wins := heroSummary.Win
		loss := heroSummary.Lose
		totalMatch := heroSummary.TotalMatch
		winrate := float64(wins) * float64(100) / float64(wins+loss)

		fmt.Printf("%d | %s | Wins: %d | Loss: %d | Total Match: %d | Winrate: %.2f%%\n", heroID, heroName, wins, loss, totalMatch, winrate)
	}

	// Test peers
	for _, peerSummary := range(peersSummary) {
		peerName := peerSummary.Name
		accountID := peerSummary.AccountID
		wins := peerSummary.Win
		loss := peerSummary.Lose
		totalMatch := peerSummary.TotalMatch
		winrate := float64(wins) * float64(100) / float64(wins+loss)
		fmt.Printf("%s (%d) - Wins: %d | Loss: %d | Total Match: %d | Winrate %.2f%%\n", peerName, accountID, wins, loss, totalMatch, winrate)
	}
	
	return benchmarkSummary, heroesSummary, peersSummary

}

func calculateAvg(target int, divisor int) float64 {
	return float64(target) / float64(divisor)
}