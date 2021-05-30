package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gocolly/colly"
)

const (
	//means: can't find github repo url
	TypeLangUnknown = "unknown"
	FlagPrefix      = "<span itemprop=\"keywords\" aria-label=\""
	TypeSandbox     = "sandbox"
	TypeIncubating  = "incubating"
	TypeGraduated   = "graduated"
	TypeMember      = "member"
)

var (
	// relations to cncf, if no relation is 'FALSE' in json
	FullRelations = []string{TypeMember, TypeSandbox, TypeIncubating, TypeGraduated}
	Total         int
	LangMap       = map[string]int{}
	LangMapMutex  sync.RWMutex
)

type Pair struct {
	Key   string
	Value int
}

type PairList []Pair

func (p PairList) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p PairList) Len() int { return len(p) }

func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

// A function to turn a map into a PairList, then sort and return it.
func sortMapByValue(m map[string]int) PairList {
	p := make(PairList, len(m))

	i := 0

	for k, v := range m {
		p[i] = Pair{k, v}
		i++
	}

	sort.Sort(sort.Reverse(p))

	return p

}

// Relation find the relation with cncf
func Relation(card [][]interface{}) string {
	for _, kv := range card {
		switch k := kv[0].(type) {
		case string:
			if !strings.Contains(k, "Relation") {
				continue
			}
		}
		switch v := kv[1].(type) {
		case string:
			for _, r := range FullRelations {
				if strings.Contains(v, r) {
					return r
				}
			}
			return ""
		case bool:
			// FALSE
			if !v {
				return "nonmember"
			}
		}
	}
	return ""
}

// GetLangAndPercent get language name and percent
func GetLangAndPercent(str string) (string, float64) {
	prePos := strings.Index(str, FlagPrefix)
	if prePos != -1 {
		start := prePos + len(FlagPrefix)
		tailPos := strings.Index(str[start:], "\"")
		if tailPos != -1 {
			end := start + tailPos
			kv := strings.Split(str[start:end], " ")
			f, _ := strconv.ParseFloat(kv[1], 64)
			return kv[0], f
		}
	}
	return TypeLangUnknown, -1
}

// GetName get project name
func GetName(card [][]interface{}) string {
	for _, kv := range card {
		switch k := kv[0].(type) {
		case string:
			if !strings.Contains(k, "Name") {
				continue
			}
		}
		switch kv[1].(type) {
		case string:
			return kv[1].(string)
		}
	}
	return ""
}

// GetStatus get name|lang|relation|percent and write them to console/file
func GetStatus(card [][]interface{}, kv []interface{}, file *os.File, relation string) {
	name := GetName(card)
	switch v := kv[1].(type) {
	case string:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		c := colly.NewCollector(colly.StdlibContext(ctx), colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"))

		c.OnResponse(func(r *colly.Response) {
			body := string(r.Body)
			lang, percent := GetLangAndPercent(body)
			fmt.Printf("%-40s%-20s%-20s%-10v%s\n", name, relation, lang, percent, v)
			file.WriteString(fmt.Sprintf("%s,%v,%v,%v,%v\n", name, relation, lang, percent, v))
			LangMapMutex.Lock()
			LangMap[lang]++
			Total++
			LangMapMutex.Unlock()
		})
		c.OnError(func(r *colly.Response, err error) {
			fmt.Println("OnError=", err)
		})
		err := c.Visit(v)
		if err != nil {
			fmt.Println("Visit err=", err)
		}
	default:
		// ["Github Repo",null]
		// LangMapMutex.Lock()
		// LangMap[TypeLangUnknown]++
		// Total++
		// LangMapMutex.Unlock()

		// fmt.Printf("%-40s%-20s%-20s%-10v%s\n", name, TypeLangUnknown, relation, -1, "null")
		// file.WriteString(fmt.Sprintf("%s,%v,%v,%v,%v\n", name, TypeLangUnknown, relation, -1, "null"))

	}
}

// CalcResult calc the result and write to file/console
func CalcResult(f *os.File) {
	LangMapMutex.RLock()
	defer LangMapMutex.RUnlock()
	list := sortMapByValue(LangMap)
	for _, pair := range list {
		fmt.Printf("%-15s%-10v%-10.2f\n", pair.Key, pair.Value, float64(pair.Value)/float64(Total))
		l := fmt.Sprintf("%v,%v,%.2f\n", pair.Key, pair.Value, float64(pair.Value)/float64(Total))
		f.WriteString(l)
	}
	fmt.Printf("Total\t %v\n", Total)
	l := fmt.Sprintf("Total,%v\n", Total)
	f.WriteString(l)
}

func main() {
	var all bool
	var fileName string
	var cycle int
	flag.BoolVar(&all, "a", false, "calc all cncf project, otherwise only: sandbox incubating graduated")
	flag.StringVar(&fileName, "o", "output.csv", "output file name")
	flag.IntVar(&cycle, "c", 500, "spider cycle")
	flag.Parse()

	url := "https://landscape.cncf.io/data/items-export.json"

	var cards [][][]interface{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	c := colly.NewCollector(colly.StdlibContext(ctx), colly.UserAgent("Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"))
	c.OnResponse(func(r *colly.Response) {
		err := json.Unmarshal(r.Body, &cards)
		if err != nil {
			fmt.Println(err)
		}
	})

	err := c.Visit(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Printf("%-40s%-20s%-20s%-10v%s\n================================================================================================\n", "Name", "Relation", "Language", "Percent", "Git Repo")
	file.WriteString(fmt.Sprintf("%s,%v,%v,%v,%v\n==========================================================================================\n", "Name", "Relation", "Language", "Percent", "Git Repo"))

	for _, card := range cards {
		r := Relation(card)
		if !all && r != TypeSandbox && r != TypeIncubating && r != TypeGraduated {
			continue
		}

		for _, kv := range card {
			switch k := kv[0].(type) {
			case string:
				if !strings.Contains(k, "Github Repo") {
					continue
				}
			}
			go GetStatus(card, kv, file, r)
			time.Sleep(time.Millisecond * time.Duration(cycle))
		}
	}
	fmt.Println("Wait some minutes......")
	time.Sleep(time.Second * 10)
	fmt.Println("\npercentage of languages used in cncf projects")
	if !all {
		fmt.Println("filter:", TypeSandbox, TypeIncubating, TypeGraduated)
	}
	fmt.Println("=========================")
	CalcResult(file)
}
