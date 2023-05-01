package dmm

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/toshikwa/dmm-schedule-checker/app/line"
)

type Slot struct {
	TeacherId string `json:"teacherId" dynamodbav:"teacherId"`
	DateTime  string `json:"dateTime" dynamodbav:"dateTime"`
}

func CheckSchedule(teacherId string) (string, []Slot, error) {
	// get document
	doc, err := GetDocument(teacherId)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get document: %s", err)
	}
	// teacher name
	name, err := ParseTeacherName(doc)
	if err != nil {
		return "", nil, err
	}
	// find available slots
	slots, err := ParseAvailableSlots(doc, teacherId)
	if err != nil {
		return "", nil, err
	}
	return name, slots, nil
}

func FindSlotDiff(all, sub []Slot) ([]Slot, error) {
	if len(all) == len(sub) {
		return []Slot{}, nil
	}
	// sort
	sort.Slice(all, func(i, j int) bool { return all[i].DateTime < all[j].DateTime })
	sort.Slice(sub, func(i, j int) bool { return sub[i].DateTime < sub[j].DateTime })
	// diff
	diff := []Slot{}
	i, j := 0, 0
	for i < len(all) {
		if j != len(sub) && all[i] == sub[j] {
			j += 1
		} else {
			diff = append(diff, all[i])
		}
		i += 1
	}
	return diff, nil
}

func SendMessage(name string, slots []Slot) error {
	if len(slots) == 0 {
		return nil
	}
	msg := fmt.Sprintf("[%s]", name)
	year, month, _ := time.Now().Date()
	for _, s := range slots {
		dt := strings.Split(s.DateTime, "_")
		date := strings.ReplaceAll(dt[0], "-", "/")
		y := strconv.Itoa(year)
		if month == 12 && date[0:2] == "01" {
			y = strconv.Itoa(year + 1)
		}
		time := strings.ReplaceAll(dt[1], "-", ":")
		msg += fmt.Sprintf("\n%s/%s %s", y, date, time)
	}
	return line.SendMessage(msg)
}
