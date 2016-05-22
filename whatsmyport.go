package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var services_tcp map[int]string
var services_udp map[int]string
var services_ddp map[int]string

func load() {
	fileData, err := ioutil.ReadFile("/etc/services")
	if err != nil {
		log.Fatal("Unable to load file /etc/services")
		os.Exit(1)
	}

	services_tcp = make(map[int]string)
	services_udp = make(map[int]string)
	services_ddp = make(map[int]string)

	fileDataStr := string(fileData)
	lines := strings.Split(fileDataStr, "\n")
	r := regexp.MustCompile("([0-9a-z-]+)\\s*([0-9]+)/([tcp|udp|ddp]+).*#(.*)$")
	r2 := regexp.MustCompile("([0-9a-z-]+)\\s*([0-9]+)/([tcp|udp|ddp]+).*$")
	for _, line := range lines {
		var name string
		var port int
		var protocol string
		var description string
		var err error

		matches := r.FindStringSubmatch(line)
		if len(matches) == 5 {
			name = matches[1]
			port, err = strconv.Atoi(matches[2])
			protocol = matches[3]
			description = strings.TrimSpace(matches[4])
		} else {
			alt_matches := r2.FindStringSubmatch(line)
			if len(alt_matches) == 4 {
				name = alt_matches[1]
				port, err = strconv.Atoi(alt_matches[2])
				protocol = alt_matches[3]
			}
		}

		if port != 0 && err == nil {
			var completeDescription string
			if description != "" {
				completeDescription = fmt.Sprintf("%s (%s)", name, description)
			} else {
				completeDescription = fmt.Sprintf("%s", name)
			}

			switch protocol {
			case "tcp":
				services_tcp[port] = completeDescription
			case "udp":
				services_udp[port] = completeDescription
			case "ddp":
				services_ddp[port] = completeDescription
			}
		}
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	portStr := r.URL.Path[1:]
	if portStr == "" {
		portStr = "0"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 0
	}
	tcp_desc := services_tcp[port]
	udp_desc := services_udp[port]
	ddp_desc := services_ddp[port]

	if tcp_desc != "" {
		fmt.Fprintf(w, "%d/tcp: %s\n", port, tcp_desc)
	}

	if udp_desc != "" {
		fmt.Fprintf(w, "%d/udp: %s\n", port, udp_desc)
	}

	if ddp_desc != "" {
		fmt.Fprintf(w, "%d/ddp: %s\n", port, ddp_desc)
	}
}

func main() {
	load()
	http.HandleFunc("/", root)

	log.Println("Starting, listening on :30004")
	log.Fatal(http.ListenAndServe(":30004", nil))
}
