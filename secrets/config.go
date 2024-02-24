package config

import (
	"fmt"
	"os"
)

type SshTunnelConfig struct {
	SshHost     string // SSH server host
	SshPort     string // SSH server port
	SshUser     string // SSH username
	SshPassword string // SSH password
	DbHost      string // Remote database host (from the perspective of the SSH server)
	DbPort      string // Remote database port
	LocalPort   string // Local port to forward to dbPort over the SSH tunnel
}

func GetSSHConfig() SshTunnelConfig {
	hostName, hnExists := os.LookupEnv("SSHHN")
	sshPort, sshPoExists := os.LookupEnv("SSHPO")
	sshUser, sshUExists := os.LookupEnv("SSHU")
	sshPassword, sshPExists := os.LookupEnv("SSHP")
	dbHost, dbHExists := os.LookupEnv("DBH")
	dbPort, dbPExists := os.LookupEnv("DBP")
	localPort, localExists := os.LookupEnv("LCP")
	if hnExists && sshPoExists && sshUExists && sshPExists && dbHExists && dbPExists && localExists {
		return SshTunnelConfig{
			SshHost:     hostName,
			SshPort:     sshPort,
			SshUser:     sshUser,
			SshPassword: sshPassword,
			DbHost:      dbHost,
			DbPort:      dbPort,
			LocalPort:   localPort,
		}
	}
	fmt.Println("WARNING! COULD NOT GET ENV VARIABLES. TRYING ALT METHOD")
	hostName = os.Getenv("SSHHN")
	sshPort = os.Getenv("SSHPO")
	sshUser = os.Getenv("SSHU")
	sshPassword = os.Getenv("SSHP")
	dbHost = os.Getenv("DBH")
	dbPort = os.Getenv("DBP")
	localPort = os.Getenv("LCP")
	return SshTunnelConfig{
		SshHost:     hostName,
		SshPort:     sshPort,
		SshUser:     sshUser,
		SshPassword: sshPassword,
		DbHost:      dbHost,
		DbPort:      dbPort,
		LocalPort:   localPort,
	}
}

func Config(local string) map[string]string {
	dbName, exists := os.LookupEnv("DB")
	csUserName, csUNExists := os.LookupEnv("CSUSERNAME")
	csPW, csPWExists := os.LookupEnv("CSPW")
	hostDB, dbExists := os.LookupEnv("DBNAME")
	lcp, lcpExists := os.LookupEnv("LCP")

	if exists && csPWExists && csUNExists && dbExists && lcpExists {
		connstring := csUserName + ":" + csPW + "@tcp(localhost:" + lcp + ")/" + hostDB + "?parseTime=true"
		return map[string]string{
			"db": dbName,
			"cs": connstring,
		}
	}
	dbName = os.Getenv("DB")
	csUserName = os.Getenv("CSUSERNAME")
	csPW = os.Getenv("CSPW")
	hostDB = os.Getenv("DBNAME")
	lcp = os.Getenv("LCP")
	connstring := csUserName + ":" + csPW + "@tcp(localhost:" + lcp + ")/" + hostDB + "?parseTime=true"
	return map[string]string{
		"db": dbName,
		"cs": connstring,
	}
}

func SFAConfig() map[string]string {
	sfaKey, exists := os.LookupEnv("SFAKEY")
	sfaUser, userExists := os.LookupEnv("SFAUSER")
	if exists && userExists {
		return map[string]string{
			"sfaKey":  sfaKey,
			"sfaUser": sfaUser,
		}
	}
	sfaKey = os.Getenv("SFAKEY")
	sfaUser = os.Getenv("SFAUSER")
	return map[string]string{
		"sfaKey":  sfaKey,
		"sfaUser": sfaUser,
	}
}

func AttributeMeans() map[string]map[string]map[string]float32 {
	return map[string]map[string]map[string]float32{
		"Speed": {
			"C":   {"mean": 21, "stddev": 3},
			"CB":  {"mean": 50, "stddev": 7.43},
			"DE":  {"mean": 42, "stddev": 7.12},
			"DT":  {"mean": 29, "stddev": 5.89},
			"FB":  {"mean": 40, "stddev": 8.95},
			"FS":  {"mean": 50, "stddev": 7.32},
			"ILB": {"mean": 47, "stddev": 8.17},
			"K":   {"mean": 13, "stddev": 5.13},
			"OG":  {"mean": 21, "stddev": 3},
			"OLB": {"mean": 47, "stddev": 10.68},
			"OT":  {"mean": 21, "stddev": 3},
			"P":   {"mean": 13, "stddev": 5.43},
			"QB":  {"mean": 47, "stddev": 16.98},
			"RB":  {"mean": 61, "stddev": 7.54},
			"SS":  {"mean": 50, "stddev": 7.23},
			"TE":  {"mean": 47, "stddev": 9.49},
			"WR":  {"mean": 55, "stddev": 10.04},
		},
		"FootballIQ": {
			"C":   {"mean": 27, "stddev": 10.93},
			"CB":  {"mean": 24, "stddev": 6.65},
			"DE":  {"mean": 25, "stddev": 6.25},
			"DT":  {"mean": 24, "stddev": 6.63},
			"FB":  {"mean": 24, "stddev": 6.11},
			"FS":  {"mean": 26, "stddev": 7.45},
			"ILB": {"mean": 29, "stddev": 9.94},
			"K":   {"mean": 23, "stddev": 6.22},
			"OG":  {"mean": 24, "stddev": 6.37},
			"OLB": {"mean": 24, "stddev": 6.53},
			"OT":  {"mean": 25, "stddev": 6.34},
			"P":   {"mean": 23, "stddev": 6.38},
			"QB":  {"mean": 30, "stddev": 9.77},
			"RB":  {"mean": 25, "stddev": 6.01},
			"SS":  {"mean": 27, "stddev": 8.0},
			"TE":  {"mean": 24, "stddev": 6.48},
			"WR":  {"mean": 24, "stddev": 7.08},
		},
		"Agility": {
			"C":   {"mean": 18, "stddev": 6.83},
			"CB":  {"mean": 39, "stddev": 7.17},
			"DE":  {"mean": 34, "stddev": 7.87},
			"DT":  {"mean": 29, "stddev": 7.86},
			"FB":  {"mean": 27, "stddev": 8.48},
			"FS":  {"mean": 39, "stddev": 7.02},
			"ILB": {"mean": 34, "stddev": 7.82},
			"K":   {"mean": 14, "stddev": 5.19},
			"OG":  {"mean": 18, "stddev": 6.54},
			"OLB": {"mean": 35, "stddev": 7.91},
			"OT":  {"mean": 19, "stddev": 6.89},
			"P":   {"mean": 14, "stddev": 5.19},
			"QB":  {"mean": 23, "stddev": 9.84},
			"RB":  {"mean": 33, "stddev": 8.14},
			"SS":  {"mean": 33, "stddev": 6.44},
			"TE":  {"mean": 33, "stddev": 7.31},
			"WR":  {"mean": 35, "stddev": 7.87},
		},
		"Carrying": {
			"C":   {"mean": 13, "stddev": 5.19},
			"CB":  {"mean": 14, "stddev": 5.02},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 19, "stddev": 6.82},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 20, "stddev": 9.55},
			"RB":  {"mean": 26, "stddev": 7.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 23, "stddev": 6.38},
			"WR":  {"mean": 23, "stddev": 7.14},
		},
		"Catching": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 18, "stddev": 10.07},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 22, "stddev": 8.06},
			"FS":  {"mean": 19, "stddev": 10.39},
			"ILB": {"mean": 16, "stddev": 6.08},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 15, "stddev": 5.47},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 14, "stddev": 4.77},
			"RB":  {"mean": 26, "stddev": 9.62},
			"SS":  {"mean": 18, "stddev": 10.2},
			"TE":  {"mean": 33, "stddev": 7.45},
			"WR":  {"mean": 38, "stddev": 9.06},
		},
		"RouteRunning": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 10, "stddev": 5.91},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 8, "stddev": 5.6},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 12, "stddev": 4.89},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 21, "stddev": 7.7},
			"WR":  {"mean": 31, "stddev": 8},
		},
		"ZoneCoverage": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 39, "stddev": 10.33},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 38, "stddev": 10.27},
			"ILB": {"mean": 36, "stddev": 8.5},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 30, "stddev": 9.41},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 37, "stddev": 9.99},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"ManCoverage": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 39, "stddev": 10.86},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 35, "stddev": 9.07},
			"ILB": {"mean": 33, "stddev": 8.32},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 30, "stddev": 9.28},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 35, "stddev": 9.35},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"Strength": {
			"C":   {"mean": 49, "stddev": 7.54},
			"CB":  {"mean": 14, "stddev": 5.2},
			"DE":  {"mean": 36, "stddev": 7.28},
			"DT":  {"mean": 39, "stddev": 9.09},
			"FB":  {"mean": 44, "stddev": 7.04},
			"FS":  {"mean": 17, "stddev": 8.05},
			"ILB": {"mean": 36, "stddev": 8.59},
			"K":   {"mean": 6, "stddev": 3.44},
			"OG":  {"mean": 50, "stddev": 7.07},
			"OLB": {"mean": 34, "stddev": 9.05},
			"OT":  {"mean": 50, "stddev": 7.06},
			"P":   {"mean": 6, "stddev": 3.45},
			"QB":  {"mean": 23, "stddev": 7.59},
			"RB":  {"mean": 24, "stddev": 6.85},
			"SS":  {"mean": 18, "stddev": 7.9},
			"TE":  {"mean": 44, "stddev": 7.92},
			"WR":  {"mean": 23, "stddev": 9.46},
		},
		"Tackle": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 17, "stddev": 5.32},
			"DE":  {"mean": 37, "stddev": 7.09},
			"DT":  {"mean": 39, "stddev": 7.2},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 25, "stddev": 8.36},
			"ILB": {"mean": 36, "stddev": 7.78},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 39, "stddev": 6.99},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 25, "stddev": 8.03},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"PassBlock": {
			"C":   {"mean": 34, "stddev": 9.58},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 42, "stddev": 7.1},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 36, "stddev": 9.73},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 37, "stddev": 9.66},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 26, "stddev": 5.9},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 33, "stddev": 7.67},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"RunBlock": {
			"C":   {"mean": 34, "stddev": 10.14},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 42, "stddev": 7.33},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 34, "stddev": 9.55},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 35, "stddev": 9.74},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 35, "stddev": 7.07},
			"WR":  {"mean": 18, "stddev": 8.99},
		},
		"PassRush": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 35, "stddev": 8},
			"DT":  {"mean": 29, "stddev": 7.9},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 6, "stddev": 3.55},
			"ILB": {"mean": 18, "stddev": 5.2},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 28, "stddev": 13.97},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 6, "stddev": 3.52},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"RunDefense": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 11, "stddev": 3.28},
			"DE":  {"mean": 34, "stddev": 8.17},
			"DT":  {"mean": 32, "stddev": 7.9},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 24, "stddev": 7.02},
			"ILB": {"mean": 38, "stddev": 8.39},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 34, "stddev": 9.18},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 30, "stddev": 7.24},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"ThrowPower": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 12, "stddev": 5.21},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 12, "stddev": 5.04},
			"QB":  {"mean": 38, "stddev": 8.46},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"ThrowAccuracy": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 12, "stddev": 4.99},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 12, "stddev": 5.07},
			"QB":  {"mean": 39, "stddev": 8.15},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"KickAccuracy": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 34, "stddev": 12},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.61},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"KickPower": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 32, "stddev": 11.85},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"PuntAccuracy": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 4.97},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 29, "stddev": 10.05},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"PuntPower": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 5.14},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 30, "stddev": 11.06},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"Stamina": {
			"C":   {"mean": 50, "stddev": 15},
			"CB":  {"mean": 50, "stddev": 15},
			"DE":  {"mean": 50, "stddev": 15},
			"DT":  {"mean": 50, "stddev": 15},
			"FB":  {"mean": 50, "stddev": 15},
			"FS":  {"mean": 50, "stddev": 15},
			"ILB": {"mean": 50, "stddev": 15},
			"K":   {"mean": 50, "stddev": 5.15},
			"OG":  {"mean": 50, "stddev": 15},
			"OLB": {"mean": 50, "stddev": 15},
			"OT":  {"mean": 50, "stddev": 15},
			"P":   {"mean": 50, "stddev": 15},
			"QB":  {"mean": 50, "stddev": 15},
			"RB":  {"mean": 50, "stddev": 15},
			"SS":  {"mean": 50, "stddev": 15},
			"TE":  {"mean": 50, "stddev": 15},
			"WR":  {"mean": 50, "stddev": 15},
		},
		"Injury": {
			"C":   {"mean": 50, "stddev": 15},
			"CB":  {"mean": 50, "stddev": 15},
			"DE":  {"mean": 50, "stddev": 15},
			"DT":  {"mean": 50, "stddev": 15},
			"FB":  {"mean": 50, "stddev": 15},
			"FS":  {"mean": 50, "stddev": 15},
			"ILB": {"mean": 50, "stddev": 15},
			"K":   {"mean": 50, "stddev": 15},
			"OG":  {"mean": 50, "stddev": 15},
			"OLB": {"mean": 50, "stddev": 15},
			"OT":  {"mean": 50, "stddev": 15},
			"P":   {"mean": 50, "stddev": 15},
			"QB":  {"mean": 50, "stddev": 15},
			"RB":  {"mean": 50, "stddev": 15},
			"SS":  {"mean": 50, "stddev": 15},
			"TE":  {"mean": 50, "stddev": 15},
			"WR":  {"mean": 50, "stddev": 15},
		},
	}
}

func NFLAttributeMeans() map[string]map[string]map[string]float32 {
	return map[string]map[string]map[string]float32{
		"Speed": {
			"C":   {"mean": 22, "stddev": 3},
			"CB":  {"mean": 57, "stddev": 7.43},
			"DE":  {"mean": 47, "stddev": 7.12},
			"DT":  {"mean": 28, "stddev": 5.89},
			"FB":  {"mean": 46, "stddev": 9.35},
			"FS":  {"mean": 57, "stddev": 7.32},
			"ILB": {"mean": 54, "stddev": 8.17},
			"K":   {"mean": 13, "stddev": 5.63},
			"OG":  {"mean": 21, "stddev": 3},
			"OLB": {"mean": 54, "stddev": 11.16},
			"OT":  {"mean": 21, "stddev": 3},
			"P":   {"mean": 13, "stddev": 5.43},
			"QB":  {"mean": 52, "stddev": 17.99},
			"RB":  {"mean": 67, "stddev": 7.54},
			"SS":  {"mean": 59, "stddev": 7.23},
			"TE":  {"mean": 54, "stddev": 9.49},
			"WR":  {"mean": 62, "stddev": 9.84},
		},
		"FootballIQ": {
			"C":   {"mean": 35, "stddev": 12.79},
			"CB":  {"mean": 29, "stddev": 6.65},
			"DE":  {"mean": 28, "stddev": 6.25},
			"DT":  {"mean": 29, "stddev": 6.63},
			"FB":  {"mean": 29, "stddev": 6.11},
			"FS":  {"mean": 31, "stddev": 7.45},
			"ILB": {"mean": 35, "stddev": 9.94},
			"K":   {"mean": 27, "stddev": 6.22},
			"OG":  {"mean": 29, "stddev": 6.37},
			"OLB": {"mean": 29, "stddev": 6.53},
			"OT":  {"mean": 29, "stddev": 6.34},
			"P":   {"mean": 27, "stddev": 6.38},
			"QB":  {"mean": 35, "stddev": 10.84},
			"RB":  {"mean": 28, "stddev": 6.01},
			"SS":  {"mean": 33, "stddev": 8.0},
			"TE":  {"mean": 29, "stddev": 6.48},
			"WR":  {"mean": 29, "stddev": 7.08},
		},
		"Agility": {
			"C":   {"mean": 24, "stddev": 6.83},
			"CB":  {"mean": 46, "stddev": 7.17},
			"DE":  {"mean": 38, "stddev": 8.21},
			"DT":  {"mean": 35, "stddev": 7.86},
			"FB":  {"mean": 32, "stddev": 8.48},
			"FS":  {"mean": 45, "stddev": 7.56},
			"ILB": {"mean": 40, "stddev": 7.82},
			"K":   {"mean": 14, "stddev": 5.19},
			"OG":  {"mean": 23, "stddev": 6.54},
			"OLB": {"mean": 40, "stddev": 8.79},
			"OT":  {"mean": 23, "stddev": 6.89},
			"P":   {"mean": 13, "stddev": 5.19},
			"QB":  {"mean": 28, "stddev": 10.65},
			"RB":  {"mean": 38, "stddev": 8.64},
			"SS":  {"mean": 40, "stddev": 6.44},
			"TE":  {"mean": 38, "stddev": 8.01},
			"WR":  {"mean": 41, "stddev": 7.87},
		},
		"Carrying": {
			"C":   {"mean": 13, "stddev": 5.19},
			"CB":  {"mean": 13, "stddev": 5.02},
			"DE":  {"mean": 12, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 24, "stddev": 6.82},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 12, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 19, "stddev": 10.09},
			"RB":  {"mean": 31, "stddev": 7.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 29, "stddev": 6.38},
			"WR":  {"mean": 27, "stddev": 7.14},
		},
		"Catching": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 25, "stddev": 10.07},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 26, "stddev": 8.06},
			"FS":  {"mean": 25, "stddev": 10.39},
			"ILB": {"mean": 16, "stddev": 6.08},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 15, "stddev": 5.47},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 14, "stddev": 4.77},
			"RB":  {"mean": 30, "stddev": 10.20},
			"SS":  {"mean": 25, "stddev": 9.9},
			"TE":  {"mean": 39, "stddev": 7.45},
			"WR":  {"mean": 43, "stddev": 9.3},
		},
		"RouteRunning": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 10, "stddev": 5.91},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 8, "stddev": 5.6},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 14, "stddev": 4.89},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 28, "stddev": 7.7},
			"WR":  {"mean": 37, "stddev": 8},
		},
		"ZoneCoverage": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 45, "stddev": 10.33},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 43, "stddev": 10.27},
			"ILB": {"mean": 43, "stddev": 8.5},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 36, "stddev": 9.41},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 43, "stddev": 9.99},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"ManCoverage": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 44, "stddev": 10.86},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 41, "stddev": 9.07},
			"ILB": {"mean": 39, "stddev": 8.32},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 35, "stddev": 9.28},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 43, "stddev": 9.35},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"Strength": {
			"C":   {"mean": 57, "stddev": 7.54},
			"CB":  {"mean": 20, "stddev": 5.2},
			"DE":  {"mean": 41, "stddev": 7.28},
			"DT":  {"mean": 46, "stddev": 9.49},
			"FB":  {"mean": 50, "stddev": 7.04},
			"FS":  {"mean": 23, "stddev": 8.05},
			"ILB": {"mean": 43, "stddev": 8.59},
			"K":   {"mean": 6, "stddev": 3.44},
			"OG":  {"mean": 57, "stddev": 7.47},
			"OLB": {"mean": 40, "stddev": 9.05},
			"OT":  {"mean": 57, "stddev": 7.06},
			"P":   {"mean": 8, "stddev": 3.45},
			"QB":  {"mean": 27, "stddev": 7.59},
			"RB":  {"mean": 28, "stddev": 6.85},
			"SS":  {"mean": 24, "stddev": 7.9},
			"TE":  {"mean": 51, "stddev": 7.92},
			"WR":  {"mean": 26, "stddev": 9.46},
		},
		"Tackle": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 22, "stddev": 5.32},
			"DE":  {"mean": 42, "stddev": 7.49},
			"DT":  {"mean": 46, "stddev": 7.3},
			"FB":  {"mean": 14, "stddev": 5.5},
			"FS":  {"mean": 30, "stddev": 9.76},
			"ILB": {"mean": 43, "stddev": 8.17},
			"K":   {"mean": 7, "stddev": 5.5},
			"OG":  {"mean": 45, "stddev": 5.5},
			"OLB": {"mean": 39, "stddev": 7.29},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 8, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 31, "stddev": 9.03},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 7, "stddev": 5.5},
		},
		"PassBlock": {
			"C":   {"mean": 40, "stddev": 10.66},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 48, "stddev": 7.35},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 43, "stddev": 10.32},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 42, "stddev": 10.56},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 32, "stddev": 6.1},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 39, "stddev": 8.57},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"RunBlock": {
			"C":   {"mean": 41, "stddev": 10.54},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 49, "stddev": 7.03},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 39, "stddev": 10.85},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 41, "stddev": 10.49},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 41, "stddev": 7.77},
			"WR":  {"mean": 18, "stddev": 9.17},
		},
		"PassRush": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 7, "stddev": 5.5},
			"DE":  {"mean": 39, "stddev": 8.5},
			"DT":  {"mean": 35, "stddev": 8.3},
			"FB":  {"mean": 8, "stddev": 5.5},
			"FS":  {"mean": 7, "stddev": 3.55},
			"ILB": {"mean": 24, "stddev": 5.2},
			"K":   {"mean": 13, "stddev": 5.5},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 33, "stddev": 15.16},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 8, "stddev": 5.5},
			"QB":  {"mean": 8, "stddev": 5.5},
			"RB":  {"mean": 7, "stddev": 5.5},
			"SS":  {"mean": 8, "stddev": 3.52},
			"TE":  {"mean": 18, "stddev": 5.5},
			"WR":  {"mean": 7, "stddev": 5.5},
		},
		"RunDefense": {
			"C":   {"mean": 9, "stddev": 5.5},
			"CB":  {"mean": 11, "stddev": 3.28},
			"DE":  {"mean": 39, "stddev": 9},
			"DT":  {"mean": 38, "stddev": 8.6},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 29, "stddev": 7.82},
			"ILB": {"mean": 44, "stddev": 9.09},
			"K":   {"mean": 7, "stddev": 5.5},
			"OG":  {"mean": 10, "stddev": 5.5},
			"OLB": {"mean": 40, "stddev": 5.58},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 36, "stddev": 5.04},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 8, "stddev": 5.5},
		},
		"ThrowPower": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 12, "stddev": 5.21},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 12, "stddev": 5.04},
			"QB":  {"mean": 44, "stddev": 9.46},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"ThrowAccuracy": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 12, "stddev": 4.99},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 12, "stddev": 5.07},
			"QB":  {"mean": 44, "stddev": 9.15},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"KickAccuracy": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 40, "stddev": 12},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 14, "stddev": 5.61},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"KickPower": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 40, "stddev": 11.85},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 13, "stddev": 5.5},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"PuntAccuracy": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 4.97},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 33, "stddev": 10.05},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"PuntPower": {
			"C":   {"mean": 13, "stddev": 5.5},
			"CB":  {"mean": 13, "stddev": 5.5},
			"DE":  {"mean": 13, "stddev": 5.5},
			"DT":  {"mean": 13, "stddev": 5.5},
			"FB":  {"mean": 13, "stddev": 5.5},
			"FS":  {"mean": 13, "stddev": 5.5},
			"ILB": {"mean": 13, "stddev": 5.5},
			"K":   {"mean": 13, "stddev": 5.14},
			"OG":  {"mean": 13, "stddev": 5.5},
			"OLB": {"mean": 13, "stddev": 5.5},
			"OT":  {"mean": 13, "stddev": 5.5},
			"P":   {"mean": 32, "stddev": 11.06},
			"QB":  {"mean": 13, "stddev": 5.5},
			"RB":  {"mean": 13, "stddev": 5.5},
			"SS":  {"mean": 13, "stddev": 5.5},
			"TE":  {"mean": 13, "stddev": 5.5},
			"WR":  {"mean": 13, "stddev": 5.5},
		},
		"Stamina": {
			"C":   {"mean": 50, "stddev": 15},
			"CB":  {"mean": 50, "stddev": 15},
			"DE":  {"mean": 50, "stddev": 15},
			"DT":  {"mean": 50, "stddev": 15},
			"FB":  {"mean": 50, "stddev": 15},
			"FS":  {"mean": 50, "stddev": 15},
			"ILB": {"mean": 50, "stddev": 15},
			"K":   {"mean": 50, "stddev": 15.15},
			"OG":  {"mean": 50, "stddev": 15},
			"OLB": {"mean": 50, "stddev": 15},
			"OT":  {"mean": 50, "stddev": 15},
			"P":   {"mean": 50, "stddev": 15},
			"QB":  {"mean": 50, "stddev": 15},
			"RB":  {"mean": 50, "stddev": 15},
			"SS":  {"mean": 50, "stddev": 15},
			"TE":  {"mean": 50, "stddev": 15},
			"WR":  {"mean": 50, "stddev": 15},
		},
		"Injury": {
			"C":   {"mean": 50, "stddev": 15},
			"CB":  {"mean": 50, "stddev": 15},
			"DE":  {"mean": 50, "stddev": 15},
			"DT":  {"mean": 50, "stddev": 15},
			"FB":  {"mean": 50, "stddev": 15},
			"FS":  {"mean": 50, "stddev": 15},
			"ILB": {"mean": 50, "stddev": 15},
			"K":   {"mean": 50, "stddev": 15},
			"OG":  {"mean": 50, "stddev": 15},
			"OLB": {"mean": 50, "stddev": 15},
			"OT":  {"mean": 50, "stddev": 15},
			"P":   {"mean": 50, "stddev": 15},
			"QB":  {"mean": 50, "stddev": 15},
			"RB":  {"mean": 50, "stddev": 15},
			"SS":  {"mean": 50, "stddev": 15},
			"TE":  {"mean": 50, "stddev": 15},
			"WR":  {"mean": 50, "stddev": 15},
		},
	}
}

func ESPNModifiers() map[string]map[string]float64 {
	return map[string]map[string]float64{
		"C": {
			"Height": 75.4939759,
			"Weight": 289.0361446,
		},
		"CB": {
			"Height": 71.02755906,
			"Weight": 187.2952756,
		},
		"DE": {
			"Height": 75.84615385,
			"Weight": 268.3076923,
		},
		"DT": {
			"Height": 74.47761194,
			"Weight": 292.7960199,
		},
		"FB": {
			"Height": 71.86567164,
			"Weight": 236.6716418,
		},
		"FS": {
			"Height": 71.46,
			"Weight": 199.4466667,
		},
		"ILB": {
			"Height": 73.39160839,
			"Weight": 230.3426573,
		},
		"K": {
			"Height": 71.74576271,
			"Weight": 190.1525424,
		},
		"OG": {
			"Height": 75.52513966,
			"Weight": 293.4748603,
		},
		"OLB": {
			"Height": 74.44270833,
			"Weight": 237.8958333,
		},
		"OT": {
			"Height": 75.54054054,
			"Weight": 293.3552124,
		},
		"P": {
			"Height": 72.45070423,
			"Weight": 189.1971831,
		},
		"QB": {
			"Height": 73.88826816,
			"Weight": 218.9497207,
		},
		"RB": {
			"Height": 70.3628692,
			"Weight": 209.1561181,
		},
		"SS": {
			"Height": 71.53271028,
			"Weight": 200.3364486,
		},
		"TE": {
			"Height": 75.64393939,
			"Weight": 238.3636364,
		},
		"WR": {
			"Height": 72.40863787,
			"Weight": 193.833887,
		},
		"ATH": {
			"Height": 73.39160839,
			"Weight": 220,
		},
	}
}

func RivalsModifiers() []int {
	return []int{100, 83, 82, 81, 80,
		76, 75, 74, 73, 72,
		69, 68, 67, 66, 65,
		64, 63, 62, 61, 60,
		59, 58, 57, 56, 55,
		53, 53, 53, 53, 53,
		51, 51, 51, 51, 51,
		49, 49, 49, 49, 49,
		47, 47, 47, 47, 47,
		45, 45, 45, 45, 45,
		43, 43, 43, 43, 43,
		41, 41, 41, 41, 41,
		40, 40, 40, 40, 40,
		39, 39, 39, 39, 39,
		38, 38, 38, 38, 38,
		37, 37, 37, 37, 37,
		36, 36, 36, 36, 36,
		35, 35, 35, 35, 35,
		34, 34, 34, 34, 34,
		33, 33, 33, 33, 33,
		32, 32, 32, 32, 32,
		31, 31, 31, 31, 31,
		30, 30, 30, 30, 30,
		29, 29, 29, 29, 29,
		28, 28, 28, 28, 28,
		27, 27, 27, 27, 27,
		26, 26, 26, 26, 26,
		25, 25, 25, 25, 25,
		24, 24, 24, 24, 24,
		23, 23, 23, 23, 23,
		22, 22, 22, 22, 22,
		21, 21, 21, 21, 21,
		20, 20, 20, 20, 20,
		19, 19, 19, 19, 19,
		18, 18, 18, 18, 18,
		17, 17, 17, 17, 17,
		16, 16, 16, 16, 16,
		15, 15, 15, 15, 15,
		14, 14, 14, 14, 14,
		13, 13, 13, 13, 13,
		12, 12, 12, 12, 12,
		11, 11, 11, 11, 11,
		10, 10, 10, 10, 10,
		9, 9, 9, 9, 9,
		8, 8, 8, 8, 8,
		7, 7, 7, 7, 7,
		6, 6, 6, 6, 6,
		5, 5, 5, 5, 5,
		4, 4, 4, 4, 4,
		3, 3, 3, 3, 3,
	}
}
