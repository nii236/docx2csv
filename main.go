package main

import (
	"encoding/csv"
	"encoding/xml"
	"flag"
	"log"
	"os"
	"regexp"
	"strings"
)

var patterns string
var inputPath string
var outputPath string

type body struct {
	Paragraph []string `xml:"p>r>t"`
}

type document struct {
	XMLName xml.Name `xml:"document"`
	Body    body     `xml:"body"`
}

func init() {
	flag.StringVar(&patterns, "p", "", "Patterns to match")
	flag.StringVar(&inputPath, "i", "", "File to load")
	flag.StringVar(&outputPath, "o", "", "File to output to")
}

func main() {
	flag.Parse()

	d, err := readDocxFile(inputPath)
	if err != nil {
		log.Println(err)
		return
	}
	b := []byte(d.content)
	document := &document{}
	err = xml.Unmarshal(b, document)
	if err != nil {
		log.Println(err)
		return
	}

	makeTable(document)
}

func match(input string, matchers []string) bool {
	for _, matcher := range matchers {
		matched, err := regexp.MatchString(matcher, input)
		if err == nil && matched {
			return true
		}
	}
	return false
}

func makeTable(d *document) {
	f, err := os.Create(outputPath)
	if err != nil {
		log.Fatalln(err)
	}
	matchers := strings.Split(patterns, ",")
	records := [][]string{[]string{"Instruments"}}
	w := csv.NewWriter(f)
	for _, p := range d.Body.Paragraph {
		if match(p, matchers) {
			records = append(records, []string{p})
		}
	}
	w.WriteAll(records)
}
