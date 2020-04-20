package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"syscall"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var channels map[string]*epoll

func wsHandler(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query()["channel_id"]
	if channelID[0] == "" {
		return
	}

	if _, exists := channels[channelID[0]]; !exists {
		log.Println("Channel Not Active")
		return
	}

	// Upgrade connection
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		return
	}
	if err := channels[channelID[0]].Add(conn); err != nil {
		log.Printf("Failed to add connection %v", err)
		conn.Close()
	}
}

func handleActivateChannel(w http.ResponseWriter, r *http.Request) {
	channelID := r.URL.Query()["channel_id"]
	if channelID[0] == "" {
		return
	}

	if _, exists := channels[channelID[0]]; exists {
		log.Println("Channel Already Active")
		return
	}

	log.Println(channelID[0])

	// Start epoll
	var err error
	channels[channelID[0]], err = MkEpoll()
	if err != nil {
		panic(err)
	}

	go Start(channelID[0])
}

func main() {
	channels = make(map[string]*epoll)

	// Increase resources limitations
	var rLimit syscall.Rlimit
	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}
	rLimit.Cur = rLimit.Max
	if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		panic(err)
	}

	// Enable pprof hooks
	go func() {
		if err := http.ListenAndServe("localhost:6060", nil); err != nil {
			log.Fatalf("pprof failed: %v", err)
		}
	}()

	http.HandleFunc("/ws/groupchat", wsHandler)
	http.HandleFunc("/channel", handleActivateChannel)
	if err := http.ListenAndServe("0.0.0.0:8000", nil); err != nil {
		log.Fatal(err)
	}
}

func Start(channelID string) {
	for {
		connections, err := channels[channelID].Wait()
		if err != nil {
			log.Printf("Failed to epoll wait %v", err)
			continue
		}
		for _, conn := range connections {
			if conn == nil {
				break
			}
			if _, _, err := wsutil.ReadClientData(conn); err != nil {
				if err := channels[channelID].Remove(conn); err != nil {
					log.Printf("Failed to remove %v", err)
				}
				conn.Close()
				if len(channels[channelID].connections)%10 == 0 {
					log.Println(len(channels[channelID].connections))
				}
			} else {
				// This is commented out since in demo usage, stdout is showing messages sent from > 1M connections at very high rate
				//log.Printf("msg: %s", string(msg))
			}
		}
	}
}
