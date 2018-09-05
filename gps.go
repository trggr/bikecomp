package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"os/exec"
	"time"
)

const (
	MAXITERATIONS = 10000
	POLLINTERVAL  = 60 * time.Second
	LOCATIONCMD   = "termux-location"
	LOGFILENM     = "gpsdb.txt"
)

type Location struct {
	Latitude  float64
	Longitude float64
	Altitude  float64
	Accuracy  float64
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

// From https://gist.github.com/cdipaolo/d3f8db3848278b49db68
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	const (
		rad = math.Pi / 180 // convert to radians
		r   = 6378100       // Earth radius in meters
	)

	la1, lo1, la2, lo2 := lat1*rad, lon1*rad, lat2*rad, lon2*rad

	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)
	return 2 * r * math.Asin(math.Sqrt(h))
}

func log(text string) {
	fp, err := os.OpenFile(LOGFILENM, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("open error:", err)
	}

	defer fp.Close()

	if _, err = fp.WriteString(text); err != nil {
		fmt.Println("write error:", err)
	}
}

func main() {
	var loc, priorloc Location
	var dist, totaldist, seconds, totalseconds, speed, avgspeed float64
	var priortime time.Time
	var logtext string

	for i := 0; i < MAXITERATIONS; i++ {
		out, err := exec.Command(LOCATIONCMD).Output()
		if err != nil {
			fmt.Println("error obtaining location:", err)
			time.Sleep(POLLINTERVAL)
			continue
		}
		now := time.Now()

		err = json.Unmarshal([]byte(out), &loc)
		if err != nil {
			fmt.Println("error parsing location:", err)
			continue
		}

		if priorloc.Latitude != 0 {
			dist = distance(priorloc.Latitude, priorloc.Longitude,
				loc.Latitude, loc.Longitude)
			seconds = time.Since(priortime).Seconds()

			totaldist += dist
			totalseconds += seconds

			if seconds > 0 {
				speed = dist * 3.6 / seconds
			} else {
				speed = 0
			}

			if totalseconds > 0 {
				avgspeed = totaldist * 3.6 / totalseconds
			} else {
				avgspeed = 0
			}
		}
		logtext = fmt.Sprintf("%s|%f|%f|%f|%f|%f|%f\n",
			now.Format("01/02/2006 15:04:05"),
			loc.Latitude, loc.Longitude,
			totaldist, totalseconds, speed, avgspeed)
		log(logtext)

		fmt.Printf("%s %6.0fm %5.0fs %5.1fkm/h %5.1fkm/h\n",
			now.Format("15:04:05"),
			totaldist, totalseconds, speed, avgspeed)
		priorloc = loc
		priortime = now

		time.Sleep(POLLINTERVAL)
	}
}
