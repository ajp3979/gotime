package main

import (
	"fmt"
	"time"

	"github.com/ajp3979/gotime/pager_duty"
)

type Region struct {
	Name     string
	Location string
}

func (r Region) GenerateRegion(currentTime time.Time) {
	l, _ := time.LoadLocation(r.Location)
	fmt.Printf("%s: %s\n", r.Name, currentTime.In(l).Format("Mon Jan 2 15:04:05"))
}

func main() {
	currentTime := time.Now().UTC()
	regions := []Region{
		{Name: "Hyderabad", Location: "Asia/Kolkata"},
		{Name: "UTC", Location: "UTC"},
		{Name: "Dublin", Location: "Europe/Dublin"},
		{Name: "Eastern USA", Location: "America/New_York"},
		{Name: "Central USA", Location: "America/Chicago"},
		{Name: "Pacific USA", Location: "America/Los_Angeles"},
	}
	for _, region := range regions {
		region.GenerateRegion(currentTime)
	}
	pg_res, err := pager_duty.PagerDuty(currentTime)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(pg_res)
	}

}
