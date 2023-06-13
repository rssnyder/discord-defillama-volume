package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

var token, protocol, header, activity, metrics *string
var nickname *bool
var refresh *int
var updates prometheus.Counter

func init() {
	token = flag.String("token", "", "discord bot token")
	protocol = flag.String("protocol", "", "protocol to get tvl for")
	nickname = flag.Bool("nickname", true, "set data in nickname")
	header = flag.String("header", "", "text before data in nickname")
	activity = flag.String("activity", "", "bot activity")
	refresh = flag.Int("refresh", 300, "seconds between refresh")
	metrics = flag.String("metrics", ":8080", "address for prometheus metric serving")
	flag.Parse()

	updates = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "updates",
			Help: "Number of times discord has been updated",
		},
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(updates)
	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	go func() {
		log.Fatal(http.ListenAndServe(*metrics, nil))
	}()
}

func main() {
	dg, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = dg.Open()
	if err != nil {
		log.Fatal(err)
		return
	}

	guilds, err := dg.UserGuilds(100, "", "")
	if err != nil {
		log.Println(err)
		*nickname = false
	}
	if len(guilds) == 0 {
		*nickname = false
	}

	ticker := time.NewTicker(time.Duration(*refresh) * time.Second)
	activityCycles := 1

	for {
		select {
		case <-ticker.C:

			volume, err := GetVolume(*protocol)
			if err != nil {
				log.Println(err)
				continue
			}

			p := message.NewPrinter(language.English)
			var fmtVolume string
			switch {
			case volume.Total24H < 1000000:
				fmtVolume = p.Sprintf("%s$%.2fk", *header, volume.Total24H/1000)
			case volume.Total24H < 1000000000:
				fmtVolume = p.Sprintf("%s$%.2fM", *header, volume.Total24H/1000000)
			case volume.Total24H < 1000000000000:
				fmtVolume = p.Sprintf("%s$%.2fB", *header, volume.Total24H/1000000000)
			case volume.Total24H < 1000000000000000:
				fmtVolume = p.Sprintf("%s$%.2fT", *header, volume.Total24H/1000000000000)
			}

			if *nickname {
				for _, g := range guilds {
					err = dg.GuildMemberNickname(g.ID, "@me", fmtVolume)
					if err != nil {
						log.Println(err)
						continue
					} else {
						log.Printf("Set nickname in %s: %s\n", g.Name, fmtVolume)
						updates.Inc()
					}
				}
			} else {
				err = dg.UpdateWatchStatus(0, fmtVolume)
				if err != nil {
					log.Printf("Unable to set activity: %s\n", err)
				} else {
					log.Printf("Set activity: %s\n", fmtVolume)
					updates.Inc()
				}
			}

			switch {
			case activityCycles == 1:
				if volume.Change1D > 0 {
					*activity = fmt.Sprintf("$%.2fM ↗️", volume.Change1D)
				} else {
					*activity = fmt.Sprintf("$%.2fM ↘️", volume.Change1D)
				}
				activityCycles -= 1
			case activityCycles == 0:
				*activity = "24hr Change"
				activityCycles += 1
			}

			err = dg.UpdateWatchStatus(0, *activity)
			if err != nil {
				log.Printf("Unable to set activity: %s\n", err)
			} else {
				log.Printf("Set activity: %s\n", *activity)
			}
		}
	}
}
