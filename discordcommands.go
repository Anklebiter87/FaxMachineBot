package main

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

func helpMenu() (retstr string) {
	return "Help would go here if I had one"
}

func doEntityQuery(info requestInfo) {

}

func convertMinToTime(mins float64) (timestr string) {
	days := int(math.Floor(mins / 24.0 / 60.0))
	hours := int(mins) / 60 % 24
	minutes := int(mins) % 60
	return fmt.Sprintf("%v Days, %v Hours, %v mins", days, hours, minutes)
}

func calculateDistance(info requestInfo) (mesg *Embed, private bool) {
	mesg = NewEmbed()
	private = checkPrivate(info.message)

	//get coords
	reType := regexp.MustCompile(`[(]?[-]?\d{1,3}[,]?\s?[-]?\d{1,3}[)]?`)
	coords := reType.FindAll([]byte(strings.ToLower(info.message)), 2)
	reType = regexp.MustCompile(`[-]?\d{1,3}`)
	firstpoint := reType.FindAll(coords[0], 2)
	secondpoint := reType.FindAll(coords[1], 2)

	//get hyper
	reType = regexp.MustCompile(`(speed|space)\s?\d`)
	hyper := string(reType.Find([]byte(strings.ToLower((info.message)))))
	reType = regexp.MustCompile(`\d`)
	hyperint, err := strconv.Atoi(string(reType.Find([]byte(hyper))))
	check(err)
	mesg.AddField("Point 1:", string(coords[0]), true)
	mesg.AddField("Point 2:", string(coords[1]), true)
	mesg.AddField("Hyperspace:", hyper, true)
	x1, err := strconv.Atoi(string(firstpoint[0]))
	check(err)
	y1, err := strconv.Atoi(string(firstpoint[1]))
	check(err)
	x2, err := strconv.Atoi(string(secondpoint[0]))
	check(err)
	y2, err := strconv.Atoi(string(secondpoint[1]))
	check(err)
	//Calculate distance
	distance := 0.0
	distanceX := math.Abs(float64(x2 - x1))
	distanceY := math.Abs(float64(y2 - y1))
	if distanceX > distanceY {
		distance = distanceX
	} else {
		distance = distanceY
	}
	minEta := (120.0 * distance) / float64(hyperint) * (1.0)
	maxEta := (120.0 * distance) / float64(hyperint) * (1.0 - (0.04 * 5.0))
	minEtaStr := convertMinToTime(minEta)
	maxEtaStr := convertMinToTime(maxEta)
	mesg.AddField("Minimum Eta:", minEtaStr, true)
	mesg.AddField("Maximum Eta:", maxEtaStr, true)
	return mesg, private
}
