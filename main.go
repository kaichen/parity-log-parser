package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/araddon/dateparse"
	"github.com/olekukonko/tablewriter"
)

type Reqeust struct {
	Id      int64         `json:"id"`
	Method  string        `json:"method"`
	Version string        `json:"jsonrpc"`
	Params  []interface{} `json:"params"`
}

type LogLine struct {
	Requests  []Reqeust
	Timestamp string
}

func main() {
	var logfile string
	flag.StringVar(&logfile, "logfile", "", "parity log file you wanna analysis")
	flag.Parse()
	fmt.Println("start process log file:", logfile)
	file, err := os.Open(logfile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	stats := make(map[string]uint64)
	var lines []*LogLine
	scanner := bufio.NewScanner(file)
	scanner.Buffer([]byte{}, bufio.MaxScanTokenSize*10)
	for scanner.Scan() {
		lineTxt := scanner.Text()
		line := processLine(lineTxt)
		//fmt.Println(lineTxt)
		if line != nil {
			lines = append(lines, line)
			if line.Requests != nil {
				for _, req := range line.Requests {
					if v, ok := stats[req.Method]; ok {
						stats[req.Method] = v + 1
					} else {
						stats[req.Method] = 1
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	printLines(lines)
	printStats(stats)
}

func printStats(stats map[string]uint64) {
	data := make([][]string, 0)
	for m, c := range stats {
		data = append(data, []string{m, strconv.FormatUint(uint64(c), 10)})
	}

	sort.SliceStable(data, func(i, j int) bool {
		ii, _ := strconv.ParseUint(data[i][1], 10, 32)
		jj, _ := strconv.ParseUint(data[j][1], 10, 32)
		return ii > jj
	})

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Method", "Count"})

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}

func printLines(llines []*LogLine) {
	t1, _ := dateparse.ParseLocal(llines[0].Timestamp)
	t2, _ := dateparse.ParseLocal(llines[len(llines)-1].Timestamp)
	reqCount := len(llines)
	fmt.Println("from:\t", t1)
	fmt.Println("to:\t", t2)
	fmt.Println("requests:\t", reqCount)

	period := t2.Sub(t1)
	if period == 0 {
		period = 1
	}
	fmt.Println("qps:\t", float64(reqCount)/period.Seconds())
}

// example line:
// 2019-06-13 21:58:46 UTC jsonrpc-eventloop-0 TRACE rpc  Request: [].
// ----------------------- ------------------- ----- --- ------------
// timestmap               logger              level com message
func processLine(line string) *LogLine {
	var lline LogLine
	if strings.Contains(line, "rpc") && strings.Contains(line, "Request") {
		words := strings.Fields(line)
		timestamp := strings.Join(words[:3], " ")
		text := strings.Join(words[7:], " ")
		stripJson := text[:len(text)-1]

		lline.Timestamp = timestamp

		var reqs []Reqeust
		if err := json.Unmarshal([]byte(stripJson), &reqs); err != nil {
			var singleRequest Reqeust
			if err := json.Unmarshal([]byte(stripJson), &singleRequest); err != nil {
				//fmt.Println(text)
				//panic(err)
				return nil
			}
			reqs = append(reqs, singleRequest)
		}
		//if len(reqs) == 0 {
		//	panic(line)
		//}

		var stripReqs []Reqeust
		for _, r := range reqs {
			if r.Method != "" {
				stripReqs = append(stripReqs, r)
			}
		}
		lline.Requests = stripReqs
		return &lline
	}
	return nil
}
