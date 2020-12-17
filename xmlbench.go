package main

import (
	"bufio"
	"compress/gzip"
	"flag"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"

	xmlparser_dev "github.com/imirkin/xml-stream-parser"
	xmlparser "github.com/tamerh/xml-stream-parser"
)

var profile bool
var devel bool
var sampleSize int
var label string

func main() {

	pid := os.Getpid()
	pidstr := strconv.Itoa(pid)

	log.Printf("Process id  %s", pidstr)

	flag.IntVar(&sampleSize, "sample", 100, "")
	flag.BoolVar(&profile, "profile", false, "")
	flag.BoolVar(&devel, "devel", false, "")
	flag.Parse()

	if devel {
		label = "devel"
	} else {
		label = "prod"
	}

	if profile {
		os.Remove("memprof" + label + ".out")
		os.Remove("cpuprof" + label + ".out")

		ff, err := os.Create("cpuprof" + label + ".out")
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(ff); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	f, err := os.Open("uniprot_sprot.xml.gz")

	if err != nil {
		panic(err)
	}

	gz, err := gzip.NewReader(f)

	if err != nil {
		panic(err)
	}

	br := bufio.NewReaderSize(gz, 65536)

	// collect stats
	cmd := exec.Command("./top.sh", pidstr, label)
	cmd.Stdout = os.Stdout
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	if devel {
		parser := xmlparser_dev.NewXMLParser(br, "entry")
		runDevel(parser)
	} else {
		parser := xmlparser.NewXMLParser(br, "entry")
		runProd(parser)
	}

	elapsed := time.Since(start)
	log.Printf("finished took %s", elapsed)

	if profile {
		f2, err := os.Create("memprof" + label + ".out")
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f2); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
		f2.Close()
	}

}

// same func as runDevel
func runProd(parser *xmlparser.XMLParser) {
	index := 0
	samples := make([]string, sampleSize)
	for xml := range parser.Stream() {

		if xml.Childs["name"] == nil {
			continue
		}

		entryid := xml.Childs["name"][0].InnerText

		if index > 0 {
			rand := rand.Intn(index)

			if index < sampleSize {
				samples[index] = entryid
			} else if rand < 100 {
				samples[rand] = entryid
			}

		} else {
			samples[0] = entryid
		}

		index++

	}
}

// same func as runProd
func runDevel(parser *xmlparser_dev.XMLParser) {
	index := 0
	samples := make([]string, sampleSize)
	for xml := range parser.Stream() {

		if xml.Childs["name"] == nil {
			continue
		}

		entryid := xml.Childs["name"][0].InnerText

		if index > 0 {
			rand := rand.Intn(index)

			if index < sampleSize {
				samples[index] = entryid
			} else if rand < 100 {
				samples[rand] = entryid
			}

		} else {
			samples[0] = entryid
		}

		index++

	}
}
