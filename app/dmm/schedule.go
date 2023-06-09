package dmm

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/toshikwa/dmm-schedule-checker/app/line"
)

var (
	teacherTableName  = os.Getenv("TEACHER_TABLE_NAME")
	scheduleTableName = os.Getenv("SCHEDULE_TABLE_NAME")
	teacherIdPattern  = regexp.MustCompile(`^\d{5}$`)
)

type Teacher struct {
	Id string `json:"id" dynamodbav:"id"`
}

type Slot struct {
	TeacherId string `json:"teacherId" dynamodbav:"teacherId"`
	DateTime  string `json:"dateTime" dynamodbav:"dateTime"`
}

type SlotWithTTL struct {
	TeacherId string `json:"teacherId" dynamodbav:"teacherId"`
	DateTime  string `json:"dateTime" dynamodbav:"dateTime"`
	Ttl       int64  `json:"ttl" dynamodbav:"ttl"`
}

func AssertTeacherId(teacherId string) bool {
	return teacherIdPattern.MatchString(teacherId)
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

func DiffSlots(news, exists []Slot) ([]Slot, []Slot, error) {
	// sort
	sort.Slice(news, func(i, j int) bool { return news[i].DateTime < news[j].DateTime })
	sort.Slice(exists, func(i, j int) bool { return exists[i].DateTime < exists[j].DateTime })
	// diff
	adds := []Slot{}
	dels := []Slot{}
	i, j := 0, 0
	for {
		if i == len(news) && j == len(exists) {
			break
		} else if j == len(exists) {
			adds = append(adds, news[i])
			i += 1
		} else if i == len(news) {
			dels = append(dels, exists[j])
			j += 1
		} else if news[i].DateTime < exists[j].DateTime {
			adds = append(adds, news[i])
			i += 1
		} else if news[i].DateTime > exists[j].DateTime {
			dels = append(dels, exists[j])
			j += 1
		} else {
			i += 1
			j += 1
		}
	}
	return adds, dels, nil
}

func SendMessage(id, name string, slots []Slot) error {
	if len(slots) == 0 {
		return nil
	}
	msg := fmt.Sprintf("[%s]", name)
	year, month, _ := time.Now().Date()
	for _, s := range slots {
		dt := strings.Split(s.DateTime, "_")
		date := strings.ReplaceAll(dt[0], "-", "/")
		var y string
		if month == 12 && date[0:2] == "01" {
			y = strconv.Itoa(year + 1)
		} else {
			y = strconv.Itoa(year)
		}
		time := strings.ReplaceAll(dt[1], "-", ":")
		msg += fmt.Sprintf("\n%s/%s %s", y, date, time)
	}
	msg += fmt.Sprintf("\n-> %s/teacher/index/%s/", dmmUrl, id)
	return line.SendMessage(msg)
}
