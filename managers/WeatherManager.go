package managers

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CalebRose/SimFBA/dbprovider"
	"github.com/CalebRose/SimFBA/structs"
	"github.com/CalebRose/SimFBA/util"
	"gorm.io/gorm"
)

func GenerateWeatherForGames() {
	db := dbprovider.GetInstance().GetDB()
	fmt.Println(time.Now().UnixNano())
	rand.Seed(time.Now().UnixNano())
	regions := getRegionalWeather()
	rainForecasts := getRainChart()
	mixForecasts := getMixChart()
	snowForecasts := getSnowChart()
	ts := GetTimestamp()
	stadiums := GetAllStadiums()

	teamRegions := getRegionsForSchools()

	stadiumMap := make(map[uint]structs.Stadium)

	for _, stadium := range stadiums {
		stadiumMap[stadium.ID] = stadium
	}
	seasonID := strconv.Itoa(int(ts.CollegeSeasonID))
	collegeGames := GetCollegeGamesBySeasonID(seasonID)

	for _, game := range collegeGames {
		GenerateWeatherForGame(db, game, teamRegions, regions, rainForecasts, mixForecasts, snowForecasts)
	}
}

func GenerateWeatherForGame(db *gorm.DB, game structs.CollegeGame, teamRegions map[string]string, regions map[string]structs.WeatherRegion, rainForecasts map[float64]map[int]string, mixForecasts map[float64]map[int]string, snowForecasts map[float64]map[int]string) {
	regionName := teamRegions[game.HomeTeam]
	region := regions[regionName]
	chances := []structs.WeatherChance{}

	precip := ""
	lowTemp := 0.0
	highTemp := 0.0
	gameTemp := 0.0
	cloud := ""
	wind := 0.0
	windCategory := ""

	if region.Forecasts[game.Week].DaysOfRain != 0 {
		chances = append(chances, structs.WeatherChance{Weather: "Rain", DaysOfWeather: region.Forecasts[game.Week].DaysOfRain})
	}

	if region.Forecasts[game.Week].DaysOfMix != 0 {
		chances = append(chances, structs.WeatherChance{Weather: "Mix", DaysOfWeather: region.Forecasts[game.Week].DaysOfMix})
	}

	if region.Forecasts[game.Week].DaysOfSnow != 0 {
		chances = append(chances, structs.WeatherChance{Weather: "Snow", DaysOfWeather: region.Forecasts[game.Week].DaysOfSnow})
	}

	var prev float64 = 0.0

	for _, chance := range chances {
		chance.ApplyChances(prev)
		prev = chance.DaysOfWeather
	}

	chances = append(chances, structs.WeatherChance{Weather: "Clear", DaysOfWeather: 30.0})

	weatherRoll := util.GenerateFloatFromRange(0, 30)

	for _, chance := range chances {
		if weatherRoll > chance.DaysOfWeather {
			continue
		}
		precip = chance.Weather
		break
	}

	if precip == "Rain" {
		inchesPerEvent := region.Forecasts[game.Week].InchesPerRain
		for k := range rainForecasts {
			if inchesPerEvent > k {
				continue
			}
			roll := util.GenerateIntFromRange(1, 12)
			precip = rainForecasts[k][roll]
			break
		}
	} else if precip == "Mix" {
		inchesPerEvent := region.Forecasts[game.Week].InchesPerRain
		for k := range mixForecasts {
			if inchesPerEvent > k {
				continue
			}
			roll := util.GenerateIntFromRange(1, 12)
			precip = rainForecasts[k][roll]
			break
		}
	} else if precip == "Snow" {
		inchesPerEvent := region.Forecasts[game.Week].InchesPerSnow
		for k := range snowForecasts {
			if inchesPerEvent > k {
				continue
			}
			roll := util.GenerateIntFromRange(1, 12)
			precip = rainForecasts[k][roll]
			break
		}
	}

	cloudChances := []structs.WeatherChance{}
	if precip == "Clear" {
		cloudChances = append(cloudChances, structs.WeatherChance{Weather: "Clear", DaysOfWeather: region.Forecasts[game.Week].Clear})
		cloudChances = append(cloudChances, structs.WeatherChance{Weather: "Mostly Clear", DaysOfWeather: region.Forecasts[game.Week].MostlyClear})
		cloudChances = append(cloudChances, structs.WeatherChance{Weather: "Partly Cloudy", DaysOfWeather: region.Forecasts[game.Week].PartlyCloudy})
	}
	cloudChances = append(cloudChances, structs.WeatherChance{Weather: "Mostly Cloudy", DaysOfWeather: region.Forecasts[game.Week].MostlyCloudy})
	cloudChances = append(cloudChances, structs.WeatherChance{Weather: "Overcast", DaysOfWeather: region.Forecasts[game.Week].Overcast})

	prev = 0
	for _, chance := range cloudChances {
		chance.ApplyChances(prev)
		prev = chance.DaysOfWeather
	}

	roll := util.GenerateFloatFromRange(0, cloudChances[0].DaysOfWeather)
	for _, chance := range cloudChances {
		if roll > chance.DaysOfWeather {
			continue
		}
		cloud = chance.Weather
	}

	meanWind := region.Forecasts[game.Week].WindSpeedAvg
	maxWind := region.Forecasts[game.Week].WindSpeedMax
	stdDev := (maxWind - meanWind) / 1.28
	speed := rand.NormFloat64()*stdDev + meanWind

	if strings.Contains(precip, "Light") {
		roll = util.GenerateFloatFromRange(0, 3)
		speed = speed + roll
	} else if strings.Contains(precip, "Moderate") {
		roll = util.GenerateFloatFromRange(1, 4)
		speed = speed + roll
	} else if strings.Contains(precip, "Heavy") {
		roll = util.GenerateFloatFromRange(2, 5)
		speed = speed + roll
	}

	if speed < 0 {
		speed = 0
	}
	wind = speed

	if speed < 4 {
		windCategory = "Calm"
	} else if speed < 10 {
		windCategory = "Breezy"
	} else if speed < 15 {
		windCategory = "Slightly Windy"
	} else if speed < 20 {
		windCategory = "Windy"
	} else {
		windCategory = "Very Windy"
	}

	lowMean := float64(region.Forecasts[game.Week].AvgLow)
	lowMin := float64(region.Forecasts[game.Week].MinLow)
	lowStdDev := (lowMean - lowMin) / 1.28
	lowTemp = rand.NormFloat64()*lowStdDev + lowMean

	highMean := float64(region.Forecasts[game.Week].AvgHigh)
	highMin := float64(region.Forecasts[game.Week].MinHigh)
	highStdDev := (highMean - highMin) / 1.28
	highTemp = rand.NormFloat64()*highStdDev + highMean

	if lowTemp > highTemp {
		tempo := lowTemp
		lowTemp = highTemp
		highTemp = tempo
	}

	if game.IsNightGame {
		gameTemp = lowTemp
	} else {
		gameTemp = highTemp
	}

	if strings.Contains(windCategory, "Slight") {
		mod := util.GenerateFloatFromRange(0, 3)
		gameTemp += mod
	} else if strings.Contains(windCategory, "Very") {
		mod := util.GenerateFloatFromRange(0, 3)
		gameTemp -= mod
	} else if strings.Contains(windCategory, "Windy") {
		mod := util.GenerateFloatFromRange(2, 5)
		gameTemp -= mod
	}

	if game.Week < 11 {
		// Summer Weather
		if cloud == "Clear" {
			mod := util.GenerateFloatFromRange(0, 3)
			gameTemp += mod
		} else if cloud == "Mostly Cloudy" {
			mod := util.GenerateFloatFromRange(0, 3)
			gameTemp -= mod
		} else if cloud == "Overcast" {
			mod := util.GenerateFloatFromRange(0, 3)
			gameTemp -= mod
		}
	} else {
		// IT'S FALL, BABY!
		// Summer Weather
		if cloud == "Clear" {
			mod := util.GenerateFloatFromRange(0, 3)
			gameTemp -= mod
		} else if cloud == "Overcast" {
			mod := util.GenerateFloatFromRange(0, 3)
			gameTemp += mod
		}
	}

	if strings.Contains(precip, "Rain") {
		if gameTemp < 33 {
			gameTemp = util.GenerateFloatFromRange(33, 35)
		}
	} else if strings.Contains(precip, "Mix") {
		if gameTemp > 34 {
			gameTemp = util.GenerateFloatFromRange(33, 35)
		} else if gameTemp < 29 {
			gameTemp = util.GenerateFloatFromRange(33, 35)
		}
	} else if strings.Contains(precip, "Snow") {
		if gameTemp > 32 {
			gameTemp = util.GenerateFloatFromRange(29, 32)
		}
	}

	if gameTemp < lowTemp {
		lowTemp = gameTemp
	} else if gameTemp > highTemp {
		highTemp = gameTemp
	}

	game.ApplyWeather(precip, lowTemp, highTemp, gameTemp, cloud, wind, windCategory, regionName)

	db.Save(&game)
}

func GenerateWeatherForNFLGame(db *gorm.DB, game structs.NFLGame, teamRegions map[string]string, regions map[string]structs.WeatherRegion, rainForecasts map[float64]map[int]string, mixForecasts map[float64]map[int]string, snowForecasts map[float64]map[int]string) {
	regionName := teamRegions[game.HomeTeam]
	region := regions[regionName]
	chances := []structs.WeatherChance{}

	precip := ""
	lowTemp := 0.0
	highTemp := 0.0
	gameTemp := 0.0
	cloud := ""
	wind := 0.0
	windCategory := ""

	if region.Forecasts[game.Week].DaysOfRain != 0 {
		chances = append(chances, structs.WeatherChance{Weather: "Rain", DaysOfWeather: region.Forecasts[game.Week].DaysOfRain})
	}

	if region.Forecasts[game.Week].DaysOfMix != 0 {
		chances = append(chances, structs.WeatherChance{Weather: "Mix", DaysOfWeather: region.Forecasts[game.Week].DaysOfMix})
	}

	if region.Forecasts[game.Week].DaysOfSnow != 0 {
		chances = append(chances, structs.WeatherChance{Weather: "Snow", DaysOfWeather: region.Forecasts[game.Week].DaysOfSnow})
	}

	var prev float64 = 0.0

	for _, chance := range chances {
		chance.ApplyChances(prev)
		prev = chance.DaysOfWeather
	}

	chances = append(chances, structs.WeatherChance{Weather: "Clear", DaysOfWeather: 30.0})

	weatherRoll := util.GenerateFloatFromRange(0, 30)

	for _, chance := range chances {
		if weatherRoll > chance.DaysOfWeather {
			continue
		}
		precip = chance.Weather
		break
	}

	if precip == "Rain" {
		inchesPerEvent := region.Forecasts[game.Week].InchesPerRain
		for k := range rainForecasts {
			if inchesPerEvent > k {
				continue
			}
			roll := util.GenerateIntFromRange(1, 12)
			precip = rainForecasts[k][roll]
			break
		}
	} else if precip == "Mix" {
		inchesPerEvent := region.Forecasts[game.Week].InchesPerRain
		for k := range mixForecasts {
			if inchesPerEvent > k {
				continue
			}
			roll := util.GenerateIntFromRange(1, 12)
			precip = rainForecasts[k][roll]
			break
		}
	} else if precip == "Snow" {
		inchesPerEvent := region.Forecasts[game.Week].InchesPerSnow
		for k := range snowForecasts {
			if inchesPerEvent > k {
				continue
			}
			roll := util.GenerateIntFromRange(1, 12)
			precip = rainForecasts[k][roll]
			break
		}
	}

	cloudChances := []structs.WeatherChance{}
	if precip == "Clear" {
		cloudChances = append(cloudChances, structs.WeatherChance{Weather: "Clear", DaysOfWeather: region.Forecasts[game.Week].Clear})
		cloudChances = append(cloudChances, structs.WeatherChance{Weather: "Mostly Clear", DaysOfWeather: region.Forecasts[game.Week].MostlyClear})
		cloudChances = append(cloudChances, structs.WeatherChance{Weather: "Partly Cloudy", DaysOfWeather: region.Forecasts[game.Week].PartlyCloudy})
	}
	cloudChances = append(cloudChances, structs.WeatherChance{Weather: "Mostly Cloudy", DaysOfWeather: region.Forecasts[game.Week].MostlyCloudy})
	cloudChances = append(cloudChances, structs.WeatherChance{Weather: "Overcast", DaysOfWeather: region.Forecasts[game.Week].Overcast})

	prev = 0
	for _, chance := range cloudChances {
		chance.ApplyChances(prev)
		prev = chance.DaysOfWeather
	}

	roll := util.GenerateFloatFromRange(0, cloudChances[0].DaysOfWeather)
	for _, chance := range cloudChances {
		if roll > chance.DaysOfWeather {
			continue
		}
		cloud = chance.Weather
	}

	meanWind := region.Forecasts[game.Week].WindSpeedAvg
	maxWind := region.Forecasts[game.Week].WindSpeedMax
	stdDev := (maxWind - meanWind) / 1.28
	speed := rand.NormFloat64()*stdDev + meanWind

	if strings.Contains(precip, "Light") {
		roll = util.GenerateFloatFromRange(0, 3)
		speed = speed + roll
	} else if strings.Contains(precip, "Moderate") {
		roll = util.GenerateFloatFromRange(1, 4)
		speed = speed + roll
	} else if strings.Contains(precip, "Heavy") {
		roll = util.GenerateFloatFromRange(2, 5)
		speed = speed + roll
	}

	if speed < 0 {
		speed = 0
	}
	wind = speed

	if speed < 4 {
		windCategory = "Calm"
	} else if speed < 10 {
		windCategory = "Breezy"
	} else if speed < 15 {
		windCategory = "Slightly Windy"
	} else if speed < 20 {
		windCategory = "Windy"
	} else {
		windCategory = "Very Windy"
	}

	lowMean := float64(region.Forecasts[game.Week].AvgLow)
	lowMin := float64(region.Forecasts[game.Week].MinLow)
	lowStdDev := (lowMean - lowMin) / 1.28
	lowTemp = rand.NormFloat64()*lowStdDev + lowMean

	highMean := float64(region.Forecasts[game.Week].AvgHigh)
	highMin := float64(region.Forecasts[game.Week].MinHigh)
	highStdDev := (highMean - highMin) / 1.28
	highTemp = rand.NormFloat64()*highStdDev + highMean

	if lowTemp > highTemp {
		tempo := lowTemp
		lowTemp = highTemp
		highTemp = tempo
	}

	if game.IsNightGame {
		gameTemp = lowTemp
	} else {
		gameTemp = highTemp
	}

	if strings.Contains(windCategory, "Slight") {
		mod := util.GenerateFloatFromRange(0, 3)
		gameTemp += mod
	} else if strings.Contains(windCategory, "Very") {
		mod := util.GenerateFloatFromRange(0, 3)
		gameTemp -= mod
	} else if strings.Contains(windCategory, "Windy") {
		mod := util.GenerateFloatFromRange(2, 5)
		gameTemp -= mod
	}

	if game.Week < 11 {
		// Summer Weather
		if cloud == "Clear" {
			mod := util.GenerateFloatFromRange(0, 3)
			gameTemp += mod
		} else if cloud == "Mostly Cloudy" {
			mod := util.GenerateFloatFromRange(0, 3)
			gameTemp -= mod
		} else if cloud == "Overcast" {
			mod := util.GenerateFloatFromRange(0, 3)
			gameTemp -= mod
		}
	} else {
		// IT'S FALL, BABY!
		// Summer Weather
		if cloud == "Clear" {
			mod := util.GenerateFloatFromRange(0, 3)
			gameTemp -= mod
		} else if cloud == "Overcast" {
			mod := util.GenerateFloatFromRange(0, 3)
			gameTemp += mod
		}
	}

	if strings.Contains(precip, "Rain") {
		if gameTemp < 33 {
			gameTemp = util.GenerateFloatFromRange(33, 35)
		}
	} else if strings.Contains(precip, "Mix") {
		if gameTemp > 34 {
			gameTemp = util.GenerateFloatFromRange(33, 35)
		} else if gameTemp < 29 {
			gameTemp = util.GenerateFloatFromRange(33, 35)
		}
	} else if strings.Contains(precip, "Snow") {
		if gameTemp > 32 {
			gameTemp = util.GenerateFloatFromRange(29, 32)
		}
	}

	if gameTemp < lowTemp {
		lowTemp = gameTemp
	} else if gameTemp > highTemp {
		highTemp = gameTemp
	}

	game.ApplyWeather(precip, lowTemp, highTemp, gameTemp, cloud, wind, windCategory, regionName)

	db.Save(&game)
}

func getRegionalWeather() map[string]structs.WeatherRegion {
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\WeatherData"

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	regions := make(map[string]structs.WeatherRegion)

	for _, file := range files {
		filePath := path + "\\" + file.Name()
		f, err := os.Open(filePath)
		if err != nil {
			log.Fatal("Unable to read input file "+filePath, err)
		}
		defer f.Close()

		csvReader := csv.NewReader(f)
		rows, err := csvReader.ReadAll()
		if err != nil {
			log.Fatal("Unable to parse file as CSV for "+filePath, err)
		}

		temp := file.Name()[15:]
		key := temp[:len(temp)-4]

		region := structs.WeatherRegion{
			RegionName: key,
			Forecasts:  make(map[int]structs.WeatherForecast),
		}
		for idx, row := range rows {
			if idx < 4 {
				continue
			}
			week := util.ConvertStringToInt(row[1])
			forecast := structs.WeatherForecast{
				ReferenceWeek:   row[0],
				Week:            uint(week),
				MinLow:          util.ConvertStringToInt(row[2]),
				AvgLow:          util.ConvertStringToInt(row[3]),
				MaxLow:          util.ConvertStringToInt(row[4]),
				MinHigh:         util.ConvertStringToInt(row[5]),
				AvgHigh:         util.ConvertStringToInt(row[6]),
				MaxHigh:         util.ConvertStringToInt(row[7]),
				Clear:           util.ConvertStringToFloat(row[8]),
				MostlyClear:     util.ConvertStringToFloat(row[9]),
				PartlyCloudy:    util.ConvertStringToFloat(row[10]),
				MostlyCloudy:    util.ConvertStringToFloat(row[11]),
				Overcast:        util.ConvertStringToFloat(row[12]),
				DaysOfRain:      util.ConvertStringToFloat(row[13]),
				DaysOfMix:       util.ConvertStringToFloat(row[14]),
				DaysOfSnow:      util.ConvertStringToFloat(row[15]),
				MonthlyRainfall: util.ConvertStringToFloat(row[16]),
				MonthlySnowfall: util.ConvertStringToFloat(row[17]),
				InchesPerRain:   util.ConvertStringToFloat(row[18]),
				InchesPerSnow:   util.ConvertStringToFloat(row[19]),
				WindSpeedMin:    util.ConvertStringToFloat(row[20]),
				WindSpeedAvg:    util.ConvertStringToFloat(row[21]),
				WindSpeedMax:    util.ConvertStringToFloat(row[22]),
			}
			region.AssignForecast(week, forecast)
		}
		regions[region.RegionName] = region
	}
	return regions
}

func getRainChart() map[float64]map[int]string {
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\WeatherSetup\\Weather Data - Rain Chart.csv"

	rainChartCSV := util.ReadCSV(path)

	return getChartMap(rainChartCSV)
}

func getMixChart() map[float64]map[int]string {
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\WeatherSetup\\Weather Data - Mix Chart.csv"

	mixChartCSV := util.ReadCSV(path)

	return getChartMap(mixChartCSV)
}

func getSnowChart() map[float64]map[int]string {
	path := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\WeatherSetup\\Weather Data - Snow Chart.csv"

	snowChartCSV := util.ReadCSV(path)

	return getChartMap(snowChartCSV)
}

func getChartMap(csvRecords [][]string) map[float64]map[int]string {
	chartMap := make(map[float64]map[int]string)

	for idx, row := range csvRecords {
		if idx < 2 {
			continue
		}
		var key float64
		temp := make(map[int]string)
		for i, ele := range row {

			if i < 1 {
				continue
			} else if i < 2 {
				key = util.ConvertStringToFloat(ele)
				continue
			}
			temp[i-1] = ele
		}
		chartMap[key] = temp
	}

	return chartMap
}

func getRegionsForSchools() map[string]string {
	FBSPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\WeatherSetup\\Weather Data - FBS Assigns.csv"
	FCSPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\WeatherSetup\\Weather Data - FCS Assigns.csv"
	NFLPath := "C:\\Users\\ctros\\go\\src\\github.com\\CalebRose\\SimFBA\\data\\WeatherSetup\\Weather Data - NFL Assigns.csv"

	teamMap := make(map[string]string)

	fbsCSV := util.ReadCSV(FBSPath)
	fcsCSV := util.ReadCSV(FCSPath)
	nflCSV := util.ReadCSV(NFLPath)

	for _, fbs := range fbsCSV {
		teamMap[fbs[0]] = fbs[1]
	}
	for _, fcs := range fcsCSV {
		teamMap[fcs[0]] = fcs[1]
	}
	for _, nfl := range nflCSV {
		teamMap[nfl[0]] = nfl[1]
	}

	return teamMap
}
