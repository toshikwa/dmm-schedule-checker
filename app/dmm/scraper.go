package dmm

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	dmmUrl               = "https://eikaiwa.dmm.com"
	nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
)

func GetDocument(teacherId string) (*goquery.Document, error) {
	// get
	res, err := http.Get(fmt.Sprintf("%s/teacher/index/%s/", dmmUrl, teacherId))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// parse
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %s", err)
	}
	return doc, nil
}

func ParseTeacherName(doc *goquery.Document) (string, error) {
	name := doc.Find("div.area-detail > h1").Text()
	if name == "" {
		return "", errors.New("failded to parse teacher name")
	}
	name = nonAlphanumericRegex.ReplaceAllString(name, "")
	return name, nil
}

func ParseAvailableSlots(doc *goquery.Document, teacherId string) ([]Slot, error) {
	// iterate over slots
	slots := []Slot{}
	nodes := doc.Find("li.date").Nodes
	if len(nodes) == 0 {
		return nil, errors.New("failed to parse available slot")
	}
	for _, node := range nodes {
		date := node.FirstChild.Data
		date = strings.ReplaceAll(date, "月", "-")
		date = strings.ReplaceAll(date, "日", "")
		for elem := node.NextSibling; elem != nil; elem = elem.NextSibling {
			if elem.Data != "li" {
				continue
			}
			time := elem.Attr[0].Val[2:]
			if elem.FirstChild != nil && elem.FirstChild.Data == "a" {
				dateTime := fmt.Sprintf("%s_%s", date, time)
				slots = append(slots, Slot{TeacherId: teacherId, DateTime: dateTime})
			}
		}
	}
	return slots, nil
}
