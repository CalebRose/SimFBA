package dbprovider

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"

	config "github.com/CalebRose/SimFBA/secrets"
	"github.com/CalebRose/SimFBA/structs"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/ssh"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Provider
type Provider struct {
}

var db *gorm.DB
var once sync.Once
var instance *Provider

func GetInstance() *Provider {
	once.Do(func() {
		instance = &Provider{}
	})
	return instance
}

func (p *Provider) InitDatabase() bool {
	fmt.Println("Database initializing...")

	sshConfig := config.GetSSHConfig()
	localPort, localErr := setupSSHTunnel(&sshConfig)
	if localErr != nil {
		log.Fatalf("Failed to establish SSH tunnel: %v", localErr)
	}

	var err error
	c := config.Config(localPort) // c["cs"]
	db, err = gorm.Open(mysql.Open(c["cs"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
		return false
	}

	sqlDB, err := db.DB()
	if err != nil {
		return false
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// AutoMigrations -- uncomment when needing to update a table
	//
	// General
	// db.AutoMigrate(&structs.Stadium{})

	// College

	db.AutoMigrate(&structs.CollegePlayer{})
	// db.AutoMigrate(&structs.HistoricCollegePlayer{})
	db.AutoMigrate(&structs.TransferPortalProfile{})
	// db.AutoMigrate(&structs.Player{})
	// db.AutoMigrate(&structs.CollegePlayerSeasonStats{})
	// db.AutoMigrate(&structs.CollegePlayerStats{})
	// db.AutoMigrate(&structs.CollegePlayerGameSnaps{})
	// db.AutoMigrate(&structs.CollegePlayerSeasonSnaps{})
	// db.AutoMigrate(&structs.CollegePromise{})
	// db.AutoMigrate(&structs.UnsignedPlayer{})
	// db.AutoMigrate(&structs.CollegeTeam{})
	// db.AutoMigrate(&structs.CollegeRival{})
	// db.AutoMigrate(&structs.League{})
	// db.AutoMigrate(&structs.CollegeConference{})
	// db.AutoMigrate(&structs.CollegeDivision{})
	// db.AutoMigrate(&structs.CollegeTeamStats{})
	// db.AutoMigrate(&structs.CollegeTeamSeasonStats{})
	// db.AutoMigrate(&structs.CollegeTeamDepthChart{})
	// db.AutoMigrate(&structs.CollegeDepthChartPosition{})
	// db.AutoMigrate(&structs.CollegeGameplan{})
	// db.AutoMigrate(&structs.CollegeGame{})
	// db.AutoMigrate(&structs.CollegeCoach{})
	// db.AutoMigrate(&structs.CollegeStandings{})
	// db.AutoMigrate(&structs.CollegeWeek{})
	// db.AutoMigrate(&structs.CollegeSeason{})
	// db.AutoMigrate(&structs.CollegePlayByPlay{})
	// db.AutoMigrate(&structs.CollegePollOfficial{})
	// db.AutoMigrate(&structs.CollegePollSubmission{})
	// TEST
	// db.AutoMigrate(&structs.CollegeTeamDepthChartTEST{})
	// db.AutoMigrate(&structs.CollegeDepthChartPositionTEST{})
	// db.AutoMigrate(&structs.CollegeGameplanTEST{})
	//
	// Recruit
	// db.AutoMigrate(&structs.Recruit{})
	// db.AutoMigrate(&structs.RecruitPlayerProfile{})
	// db.AutoMigrate(&structs.RecruitingTeamProfile{})
	// db.AutoMigrate(&structs.RecruitPointAllocation{})
	// db.AutoMigrate(&structs.RecruitRegion{})
	// db.AutoMigrate(&structs.RecruitState{})
	// db.AutoMigrate(&structs.ProfileAffinity{})

	// NFL
	// db.AutoMigrate(&structs.FreeAgencyOffer{})
	// db.AutoMigrate(&structs.NFLCapsheet{})
	// db.AutoMigrate(&structs.NFLDraftPick{})
	// db.AutoMigrate(&models.NFLDraftee{})
	// db.AutoMigrate(&models.EventResults{})
	// db.AutoMigrate(&models.NFLWarRoom{})
	// db.AutoMigrate(&models.ScoutingProfile{})
	// db.AutoMigrate(&structs.NFLContract{})
	// db.AutoMigrate(&structs.NFLDepthChart{})
	// db.AutoMigrate(&structs.NFLDepthChartPosition{})
	// db.AutoMigrate(&structs.NFLExtensionOffer{})
	// db.AutoMigrate(&structs.NFLGame{})
	// db.AutoMigrate(&structs.NFLGameplan{})
	// db.AutoMigrate(&structs.NFLPlayByPlay{})
	// db.AutoMigrate(&structs.NFLTradePreferences{})
	// db.AutoMigrate(&structs.NFLTradeProposal{})
	// db.AutoMigrate(&structs.NFLTradeOption{})
	// db.AutoMigrate(&structs.NFLPlayer{})
	// db.AutoMigrate(&structs.NFLRetiredPlayer{})
	// db.AutoMigrate(&structs.NFLPlayerSeasonStats{})
	// db.AutoMigrate(&structs.NFLPlayerStats{})
	// db.AutoMigrate(&structs.NFLPlayerGameSnaps{})
	// db.AutoMigrate(&structs.NFLPlayerSeasonSnaps{})
	// db.AutoMigrate(&structs.NFLUser{})
	// db.AutoMigrate(&structs.NFLTeam{})
	// db.AutoMigrate(&structs.NFLTeamStats{})
	// db.AutoMigrate(&structs.NFLTeamSeasonStats{})
	// db.AutoMigrate(&structs.NFLStandings{})
	// db.AutoMigrate(&structs.NFLRequest{})
	// db.AutoMigrate(&structs.NFLWaiverOffer{})

	// All
	// db.AutoMigrate(&structs.AdminRecruitModifier{})
	// db.AutoMigrate(&structs.Affinity{})
	// db.AutoMigrate(&structs.TeamRequest{})
	// db.AutoMigrate(&structs.Timestamp{})
	// db.AutoMigrate(&structs.FaceData{})
	// db.AutoMigrate(&structs.NewsLog{})
	// db.AutoMigrate(&structs.Notification{})
	return true
}

func (p *Provider) GetDB() *gorm.DB {
	return db
}

// setupSSHTunnel establishes an SSH tunnel and forwards a local port to the remote database port.
// Returns the local port and any error encountered.
func setupSSHTunnel(config *config.SshTunnelConfig) (string, error) {
	sshConfig := &ssh.ClientConfig{
		User: config.SshUser,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.SshPassword),
		},
		// CAUTION: In production, you should use a more secure HostKeyCallback.
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// Connect to SSH server
	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(config.SshHost, config.SshPort), sshConfig)
	if err != nil {
		return "", err
	}

	// Setup local port forwarding
	localListener, err := net.Listen("tcp", "localhost:"+config.LocalPort)
	if err != nil {
		return "", err
	}

	go func() {
		defer localListener.Close()
		for {
			localConn, err := localListener.Accept()
			if err != nil {
				log.Printf("Failed to accept local connection: %s", err)
				continue
			}

			// Handle the connection in a new goroutine
			go func() {
				defer localConn.Close()

				// Connect to the remote database server through the SSH tunnel
				remoteConn, err := sshClient.Dial("tcp", net.JoinHostPort(config.DbHost, config.DbPort))
				if err != nil {
					log.Printf("Failed to dial remote server: %s", err)
					return
				}
				defer remoteConn.Close()

				// Copy data between the local connection and the remote connection
				copyConn(localConn, remoteConn)
			}()
		}
	}()

	return localListener.Addr().String(), nil
}

// copyConn copies data between two io.ReadWriteCloser objects (e.g., network connections)
func copyConn(localConn, remoteConn io.ReadWriteCloser) {
	// Start goroutine to copy data from local to remote
	go func() {
		_, err := io.Copy(remoteConn, localConn)
		if err != nil {
			log.Printf("Error copying from local to remote: %v", err)
		}
		localConn.Close()
		remoteConn.Close()
	}()

	// Copy data from remote to local in the main goroutine (or vice versa)
	_, err := io.Copy(localConn, remoteConn)
	if err != nil {
		log.Printf("Error copying from remote to local: %v", err)
	}
	// Ensure connections are closed when copying is done or an error occurs
	localConn.Close()
	remoteConn.Close()
}
