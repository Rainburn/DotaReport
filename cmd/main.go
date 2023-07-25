package main

import (
	"net/http"
	"encoding/json"
	"log"
	"strconv"
	"runtime"
	"fmt"
)

type Matches struct {
	MatchID int `json:"match_id"`
	StartTime int `json:"start_time"`
}

type Profile struct {
	RankTier int `json:"rank_tier"`
	ProfileDetail struct {
		Name string `json:"personaname"`
	} `json:"profile"`
}

// Will be changed to parameters
// var USER_ID = 300213178
// var TOTAL_MATCH_ANALYZED = 100
var USER_ID int
var TOTAL_MATCH_ANALYZED int

var ASSEST_PATH = "/assets/"

// Global variables
var heroes map[string]Hero

func getMatches(userID int) ([]Matches, error) {
	var matches_data []Matches
	userIDStr := strconv.Itoa(userID)

	api_link := generateMatchHistoryAPI(userIDStr)
	resp, err := http.Get(api_link)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&matches_data)

	if err != nil {
		log.Fatal(err.Error())
	}

	return matches_data, nil
}

func getProfile(userID int) (Profile, error) {
	var profile Profile

	userIDStr := strconv.Itoa(userID)
	profileAPI := getProfileAPI(userIDStr)

	resp, err := http.Get(profileAPI)

	if err != nil {
		return Profile{}, err
	}

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&profile)
	if err != nil {
		return Profile{}, err
	}

	return profile, nil
}

func main() {
	runtime.GOMAXPROCS(8)

	fmt.Print("USER ID: ")
	fmt.Scan(&USER_ID)

	fmt.Print("TOTAL MATCH ANALYZED: ")
	fmt.Scan(&TOTAL_MATCH_ANALYZED)

	BuildReport(USER_ID, TOTAL_MATCH_ANALYZED)
}

func BuildReport(USER_ID int, TOTAL_MATCH_ANALYZED int) {

	var benchmark Benchmark
	var WLStat WinLoseStat
	var peersStat = PeerStat{}
	heroStats := HeroStats{}

	// Load heroes
	heroes = loadHeroesFromJSON()
	
	// Profile API
	profile, err := getProfile(USER_ID)

	matchesData, err := getMatches(USER_ID)
	if err != nil {
		log.Fatal("Error on fetching matches")
	}

	targetMatches := matchesData[:TOTAL_MATCH_ANALYZED]
	chanMatchesData := handleMatch(targetMatches)

	processMatchData(&WLStat, &benchmark, &heroStats, &peersStat, chanMatchesData)

	benchmarkSummary, heroesSummary, peersSummary := generateSummary(WLStat, benchmark, heroStats, peersStat, TOTAL_MATCH_ANALYZED)

	// Generate the report as PDF
	CreateReport(profile, WLStat, benchmarkSummary, heroesSummary, peersSummary)

}