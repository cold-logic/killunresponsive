/*******************************************
 * Kill Unresponsive Processes v1.0
 * Written By: Christopher Davison
 * Updated: 2015-02-10
 *******************************************
 * Flow
 * 1. Generates a spin dump
 * 2. Parses said dump, extracting unresponsive processes (name and pid)
 * 3. Kills unresponsive processes with extreme prejudice ;)
 */

// A package clause starts every source file.
// Main is a special name declaring an executable rather than a library.
package main

// Import declarations
import (
	"bufio"     // Buffered file reading
	"bytes"     // Bytes of bits
	"log"       // For logging
	"os"        // Platform-independent interface
	"os/exec"   // For executing shell commands
	"os/user"   // For identifying the user
	"regexp"    // Regular expressions
	"strconv"   // String conversion
	s "strings" // String manipulation
)

// Process struct
type Process struct {
	Name string
	PID  int
}

// Error checker
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Parses a spindump file and returns unresponsive processes
func parse(path string) (processes []Process, e error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var p Process
	var insideProcess bool

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`^Process:\s{9}(.+?)\s\[([0-9]+?)\]`)

	for scanner.Scan() {
		if s.Contains(scanner.Text(), "Process: ") {
			insideProcess = true                             // Set a flag to determine when we've found a process
			matches := re.FindStringSubmatch(scanner.Text()) // Search for the process name and PID
			if len(matches) == 3 {
				pname := matches[1]                           // Process name is the first capture group
				pid, e := strconv.ParseInt(matches[2], 10, 0) // PID is the second capture group
				check(e)
				p = Process{pname, int(pid)} // Set p to a new Process type
			}
		} else if insideProcess && s.Contains(scanner.Text(), "Unresponsive") {
			processes = append(processes, p) // Append p to the result
		} else if insideProcess && scanner.Text() == "" {
			insideProcess = false // Clear flag when an empty line is detected
		}
	}
	return processes, scanner.Err()
}

func kill(processes []Process) {
	l := len(processes)
	for i := 0; i < l; i++ {
		log.Printf("Attempting to kill %v [%v]\n", processes[i].Name, processes[i].PID)
		p, _ := os.FindProcess(processes[i].PID)
		err := p.Kill()
		if err != nil {
			log.Println("Process couldn't be killed: ", err)
		}
	}
	return
}

func initiateSpin(path string) {
	cmd := exec.Command("spindump", "-notarget", "1", "-noBulkSymbolication", "-file", path)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// A function definition. Main is special. It is the entry point for the
// executable program. Love it or hate it, Go uses brace brackets.
func main() {
	log.SetFlags(0)          // Disable date prefix for log
	usr, _ := user.Current() // Get the current user
	if usr.Uid == "0" {      // Make sure it's the superuser
		spindump := "/tmp/spindump.txt"   // Location to save the spindump
		initiateSpin(spindump)            // Generate the spindump
		processes, err := parse(spindump) // Process the spindump
		os.Remove(spindump)               // Remove the spindump file
		if err == nil {                   // If there were no errors
			kill(processes) // Kill the unresponsive processes
		} else {
			log.Println("Encountered an error while parsing the spindump file.")
		}
	} else {
		log.Fatal("This app requires superior authority to terminate unresponsive processes.\nTry again with sudo... ;)")
	}
}
