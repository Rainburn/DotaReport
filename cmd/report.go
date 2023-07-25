package main

import (
	"github.com/jung-kurt/gofpdf"
	"path/filepath"
	"os"
	"time"
	"fmt"
	"log"
	"strconv"
)

func generateMedalIcon(rankTier int) string {
	medalStr := strconv.Itoa(rankTier)
	medalLink := fmt.Sprintf("../assets/icons/ranks/" + medalStr + ".png")
	return medalLink
}

func generateHeroIcon(heroID int) string {
	heroIDStr := strconv.Itoa(heroID)
	heroIconLink := fmt.Sprintf("../assets/icons/heroes/" + heroIDStr + ".png")
	return heroIconLink
}

func generateResultPath(username string) string {
	resultPath := fmt.Sprintf("../results/" + username + ".pdf")
	return resultPath
}

func CreateReport(profile Profile, wlStat WinLoseStat, benchmarkSummary BenchmarkSummary, heroesSummary []HeroSummary, peersSummary []PeerSummary) {
	// A6 Page is approximately 105mm x 150mm
	cwd, _ := os.Getwd()
	parent := filepath.Dir(cwd)

	pdf := gofpdf.New("L", "mm", "A6", filepath.Join(parent, "assets", "font"))
	
	// Add Custom Font
	pdf.AddUTF8Font("GlossAndBloom", "", "GlossAndBloom.ttf")
	
	// Add new page
	pdf.AddPage()

	// Background
	pdf.Image("../assets/bg-scroll.jpeg", 0, 0, 150, 105, false, "", 0, "")

	// Create Title
	pdf.SetFont("GlossAndBloom", "", 22)
	pdf.Text(52, 12, "Dota 2 Journey")

	// Create user data
	pdf.SetFont("GlossAndBloom", "", 16)
	pdf.Text(10, 20, profile.ProfileDetail.Name)

	// Medal Rank
	
	medalLink := generateMedalIcon(profile.RankTier)
	pdf.Image(medalLink, 125, 15, 16, 16, false, "", 0, "")

	// General W/L and winrate
	winrate := float64(wlStat.Win * 100) / float64(wlStat.Win + wlStat.Lose)
	pdf.SetFont("GlossAndBloom", "", 10)
	pdf.Text(10, 25, fmt.Sprintf("Has destroyed %d enemies' ancient but lost %d times", wlStat.Win, wlStat.Lose))
	pdf.Text(10, 30, fmt.Sprintf("Wieldin' %.2f%% victory rate from last %d battles", winrate, wlStat.Win + wlStat.Lose))


	// Header of columns
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(30, 35, "yer favorite soldier")
	pdf.Text(76, 35, "yer average stats")
	pdf.Text(113, 35, "yer favorite matey")
	
	// Hero winrate 1
	hero1 := heroesSummary[0]
	hero1ID := hero1.ID
	hero1Link := generateHeroIcon(hero1ID)
	pdf.Image(hero1Link, 10, 37, 16, 9, false, "", 0, "")
	pdf.SetFont("GlossAndBloom", "", 8)
	pdf.Text(30, 41, hero1.Name)
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(30, 45, fmt.Sprintf("%d Wins n' %d Losses", hero1.Win, hero1.Lose))

	// Hero winrate 2
	hero2 := heroesSummary[1]
	hero2ID := hero2.ID
	hero2Link := generateHeroIcon(hero2ID)
	pdf.Image(hero2Link, 10, 49, 16, 9, false, "", 0, "")
	pdf.SetFont("GlossAndBloom", "", 8)
	pdf.Text(30, 53, hero2.Name)
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(30, 57, fmt.Sprintf("%d Wins n' %d Losses", hero2.Win, hero2.Lose))

	// Hero winrate 3
	hero3 := heroesSummary[2]
	hero3ID := hero3.ID
	hero3Link := generateHeroIcon(hero3ID)
	pdf.Image(hero3Link, 10, 61, 16, 9, false, "", 0, "")
	pdf.SetFont("GlossAndBloom", "", 8)
	pdf.Text(30, 65, hero3.Name)
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(30, 69, fmt.Sprintf("%d Wins n' %d Losses", hero3.Win, hero3.Lose))

	// Hero winrate 4
	hero4 := heroesSummary[3]
	hero4ID := hero4.ID
	hero4Link := generateHeroIcon(hero4ID)
	pdf.Image(hero4Link, 10, 73, 16, 9, false, "", 0, "")
	pdf.SetFont("GlossAndBloom", "", 8)
	pdf.Text(30, 77, hero4.Name)
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(30, 81, fmt.Sprintf("%d Wins n' %d Losses", hero4.Win, hero4.Lose))

	// Hero winrate 5
	hero5 := heroesSummary[4]
	hero5ID := hero5.ID
	hero5Link := generateHeroIcon(hero5ID)
	pdf.Image(hero5Link, 10, 85, 16, 9, false, "", 0, "")
	pdf.SetFont("GlossAndBloom", "", 8)
	pdf.Text(30, 89, hero5.Name)
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(30, 93, fmt.Sprintf("%d Wins n' %d Losses", hero5.Win, hero5.Lose))

	
	// General Stats
	pdf.SetFont("GlossAndBloom", "", 14)
	// Kill
	pdf.Image("../assets/icons/misc/kill.png", 60, 37, 10, 10, false, "", 0, "")
	pdf.Text(75, 45, fmt.Sprintf("%.2f", benchmarkSummary.Kill))

	// Death
	pdf.Image("../assets/icons/misc/death.png", 60, 49, 10, 10, false, "", 0, "")
	pdf.Text(75, 57, fmt.Sprintf("%.2f", benchmarkSummary.Death))

	// Assist
	pdf.Image("../assets/icons/misc/assist.png", 60, 61, 10, 10, false, "", 0, "")
	pdf.Text(75, 69, fmt.Sprintf("%.2f", benchmarkSummary.Assist))

	// GPM
	pdf.Image("../assets/icons/misc/gold.png", 60, 73, 10, 10, false, "", 0, "")
	pdf.Text(75, 81, fmt.Sprintf("%.2f", benchmarkSummary.Gpm))

	// XPM
	pdf.Image("../assets/icons/misc/exp.png", 60, 85, 10, 10, false, "", 0, "")
	pdf.Text(75, 93, fmt.Sprintf("%.2f", benchmarkSummary.Xpm))


	// Teammates
	// Teammate 1
	person1 := peersSummary[0]
	pdf.Image("../assets/icons/misc/person.png", 100, 37, 10, 10, false, "", 0, "")
	pdf.SetFont("GlossAndBloom", "", 8)
	pdf.Text(113, 42, person1.Name)
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(113, 46, fmt.Sprintf("%d Wins n' %d Losses", person1.Win, person1.Lose))

	// Teammate 2
	person2 := peersSummary[1]
	pdf.Image("../assets/icons/misc/person.png", 100, 49, 10, 10, false, "", 0, "")
	pdf.SetFont("GlossAndBloom", "", 8)
	pdf.Text(113, 54, person2.Name)
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(113, 58, fmt.Sprintf("%d Wins n' %d Losses", person2.Win, person2.Lose))


	// Teammate 3
	person3 := peersSummary[2]
	pdf.Image("../assets/icons/misc/person.png", 100, 61, 10, 10, false, "", 0, "")
	pdf.SetFont("GlossAndBloom", "", 8)
	pdf.Text(113, 66, person3.Name)
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(113, 70, fmt.Sprintf("%d Wins n' %d Losses", person3.Win, person3.Lose))

	// Teammate 4
	person4 := peersSummary[3]
	pdf.Image("../assets/icons/misc/person.png", 100, 73, 10, 10, false, "", 0, "")
	pdf.SetFont("GlossAndBloom", "", 8)
	pdf.Text(113, 78, person4.Name)
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(113, 82, fmt.Sprintf("%d Wins n' %d Losses", person4.Win, person4.Lose))

	// Teammate 5
	person5 := peersSummary[4]
	pdf.Image("../assets/icons/misc/person.png", 100, 85, 10, 10, false, "", 0, "")
	pdf.SetFont("GlossAndBloom", "", 8)
	pdf.Text(113, 90, person5.Name)
	pdf.SetFont("GlossAndBloom", "", 6)
	pdf.Text(113, 94, fmt.Sprintf("%d Wins n' %d Losses", person5.Win, person5.Lose))



	// As of date
	pdf.SetFont("GlossAndBloom", "", 8)
	now := time.Now()
	var dateFormat = now.Format("02 January 2006 15:04")
	pdf.Text(110, 102, fmt.Sprintf("As o' %v", dateFormat))

	// Generate the report as output
	resultPath := generateResultPath(profile.ProfileDetail.Name)
	err := pdf.OutputFileAndClose(resultPath)
	if err != nil {
		log.Println("ERROR:", err.Error())
	}



	return
}