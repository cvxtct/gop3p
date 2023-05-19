// C0303 & C0304 eliminator
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

// Project is a struct that contains information about a project.
type Project struct {
	dir   string
	files []string
}

var buildstamp string
var githash string

// ParseArgument parses project path from command line.
func (p *Project) ParseArgument() {
	var pinfo bool
	// TODO Recursive scan option
	// TODO Validate input
	flag.StringVar(&p.dir, "dir", "", "Directory to scan")
	flag.BoolVar(&pinfo, "version", false, "Program info")
	flag.Parse()

	if pinfo {
		fmt.Println("Version: " + "gop3p" + "-" + githash + " | BuildDate: " + buildstamp)
		os.Exit(0)
	}
}

// ParseFiles parses a directory.
func (p *Project) ParseFiles() {
	var files []fs.DirEntry
	var err error

	log.Println("Scanning: ", p.dir)

	files, err = os.ReadDir(p.dir)
	if err != nil {
		log.Println(err)
	}

	for _, f := range files {
		if strings.Contains(f.Name(), ".py") {
			p.files = append(p.files, f.Name())
		} else {
			log.Println("Skipping: ", f.Name(), "is not a Python file.")
		}
	}
}

// FixFiles opens files to read by line and fixes 0303 & 0304 on each line.
func (p *Project) FixFiles(f string) bool {

	log.Println("Start cleaning: ", f)

	var cleaned []string
	var orig []string
	// Open file.
	file, err := os.Open(p.dir + "/" + f)
	if err != nil {
		log.Println(err)
		return false
	}
	fileScanner := bufio.NewScanner(file)
	// Read file line by line.
	for fileScanner.Scan() {
		line := fileScanner.Text()
		orig = append(orig, line)
		// Fixing C0303: Trailing whitespace (trailing-whitespace).
		res := strings.TrimRight(line, " ")
		cleaned = append(cleaned, res)
	}
	// Close file.
	file.Close()

	// Fixing C0304: Final newline missing (missing-final-newline).
	// TODO if multiple final newlines then just remove them.
	for i := range cleaned {
		if len(cleaned)-1 == i+1 && cleaned[i+1] != "" {
			cleaned = append(cleaned, "")
		}
	}

	// Just in case.
	if len(cleaned)-len(orig) > 1 {
		panic("Too huge difference!")
	}

	// Write cleaned results back into file.
	ioutil.WriteFile(p.dir+"/"+f, []byte(strings.Join(cleaned, "\n")), 0644)

	return true
}

// Runner runs the FixFiles for each file as a separate goroutine.
// boolChan is used to send back success of FixFiles.
func (p *Project) Runner(boolChan chan bool) {
	for _, f := range p.files {
		fixed_file := p.FixFiles(f)
		boolChan <- fixed_file
	}
}

func main() {
	var p Project
	// Channel for sending back results of FixFiles.
	boolChan := make(chan bool)
	defer close(boolChan)

	p.ParseArgument()
	p.ParseFiles()
	// Run Runner in a goroutine.
	go p.Runner(boolChan)
	// Wait for Runner to finish.
	for i := 0; i < len(p.files); i++ {
		res := <-boolChan
		if res {
			log.Println("File: ", p.files[i], " is cleaned!")
		}
		if i == len(p.files)-1 {
			log.Println("All files are cleaned!")
			os.Exit(0)
			break
		}
	}
}
