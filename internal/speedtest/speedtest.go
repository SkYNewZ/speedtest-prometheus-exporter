package speedtest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type Result struct {
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

type PushConfig struct {
	URL     string
	Results *Result
}

// Run the speedtest
func Run(ctx context.Context, path string, logger *logrus.Logger) (*Result, error) {
	out, _ := exec.Command(path, "--version").Output()
	o := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")[0]
	logger.Printf("Using %s\n", o)

	out, err := exec.CommandContext(ctx, path, "--accept-license", "--accept-gdpr", "--progress=no", "-f", "json").Output()
	if err != nil {
		return nil, err
	}

	var res Result
	return &res, json.Unmarshal(out, &res)
}

func Push(ctx context.Context, logger *logrus.Logger, config *PushConfig) error {
	logger.Debugf("pushing results to %s", config.URL)

	data, err := json.Marshal(config.Results)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, config.URL+"/results", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("cannot create request: %v", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot make request: %v", err)
	}

	if res.StatusCode >= 400 {
		return fmt.Errorf("invalid status code received: %d", res.StatusCode)
	}

	return nil
}
