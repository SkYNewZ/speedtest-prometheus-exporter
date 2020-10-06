package main

import (
	"encoding/json"
	"os/exec"
	"strings"
	"time"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

type speedtestResult struct {
	Type      string    `json:"type"`
	Timestamp time.Time `json:"timestamp"`
	Ping      struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
	} `json:"ping"`
	Download struct {
		Bandwidth int `json:"bandwidth"`
		Bytes     int `json:"bytes"`
		Elapsed   int `json:"elapsed"`
	} `json:"download"`
	Upload struct {
		Bandwidth int `json:"bandwidth"`
		Bytes     int `json:"bytes"`
		Elapsed   int `json:"elapsed"`
	} `json:"upload"`
	PacketLoss float64 `json:"packetLoss"`
	Isp        string  `json:"isp"`
	Interface  struct {
		InternalIP string `json:"internalIp"`
		Name       string `json:"name"`
		MacAddr    string `json:"macAddr"`
		IsVpn      bool   `json:"isVpn"`
		ExternalIP string `json:"externalIp"`
	} `json:"interface"`
	Server struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		IP       string `json:"ip"`
	} `json:"server"`
	Result struct {
		ID  string `json:"id"`
		URL string `json:"url"`
	} `json:"result"`
}

var (
	c      *cache.Cache
	ticker *time.Ticker
)

func init() {
	// Check if speedtest path is correct
	out, err := exec.Command(speedtestPath, "--version").Output()
	if err != nil {
		log.Fatalf("Invalid speedtest path %s\nSpecify a good one with -speedtest-path=THE_PATH\n", speedtestPath)
	}

	o := string(out)
	o = strings.Split(strings.TrimSuffix(o, "\n"), "\n")[0]
	log.Printf("Executing %s\n", o)

	// Set up cache
	log.Debugln("Set up cache")
	c = cache.New(30*time.Second, 2*time.Minute)
}

// Run the speedtest and parse output
func execute() *speedtestResult {
	log.Debugln("Searching in cache")
	r, found := c.Get("result")
	if found {
		log.Debugln("Found in cache")
		return r.(*speedtestResult)
	}

	log.Debugln("Not found in cache, executing speedtest")
	out, err := exec.Command(speedtestPath, "--accept-license", "--accept-gdpr", "--progress=no", "-f", "json").Output()
	if err != nil {
		log.Fatalln(err)
	}

	var speedtestResult speedtestResult
	err = json.Unmarshal(out, &speedtestResult)
	if err != nil {
		log.Fatalln(err)
	}

	log.Debugln("Write in cache")
	c.Set("result", &speedtestResult, cache.DefaultExpiration)
	return &speedtestResult
}

// Launch Go routine wich execute the speedtest
func startSpeedTestTicker() {
	log.Debugln("Starting ticker for speedtest execution")
	ticker = time.NewTicker(time.Minute)
	go func() {
		for ; true; <-ticker.C {
			execute()
		}
	}()
}
