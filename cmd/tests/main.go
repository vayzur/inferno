package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	v1 "github.com/vayzur/inferno/pkg/api/v1"
	"github.com/vayzur/inferno/pkg/httputil"
)

func main() {
	cli := httputil.New(time.Second * 2)

	inbound := `
        {
            "listen": null,
            "port": 10800,
            "protocol": "vless",
            "settings": {
                "clients": [
                    {
                        "id": "inferno"
                    }
                ],
                "decryption": "none",
                "fallbacks": []
            },
            "streamSettings": {
                "network": "ws",
                "security": "none",
                "wsSettings": {
                    "acceptProxyProtocol": false,
                    "headers": {},
                    "heartbeatPeriod": 0,
                    "host": "",
                    "path": ""
                },
                "sockopt": {
                    "tcpFastOpen": true,
                    "tcpCongestion": "bbr",
                    "tcpMptcp": true,
                    "tcpNoDelay": true
                }
            },
            "tag": "proxy-10800",
            "sniffing": {
                "enabled": false,
                "destOverride": [
                    "http",
                    "tls",
                    "quic",
                    "fakedns"
                ],
                "metadataOnly": false,
                "routeOnly": false
            },
            "allocate": {
                "strategy": "always",
                "refresh": 5,
                "concurrency": 3
            }
        }
	`

	var conf v1.InboundConfig
	if err := json.Unmarshal([]byte(inbound), &conf); err != nil {
		panic(err)
	}

	conf.Tag = "proxy0"
	conf.Port = 10900

	status, resp, err := cli.Do(http.MethodPost, "http://127.0.0.1:10100/api/v1/inbounds", "token", conf)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %s - Status: %d\n", resp, status)

	status, resp, err = cli.Do(http.MethodDelete, "http://127.0.0.1:10100/api/v1/inbounds/proxy0", "token", nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %s - Status: %d\n", resp, status)
}
