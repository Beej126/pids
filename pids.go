package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shirou/gopsutil/process"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("  pids -vars                # Output environment variable assignments for each parent process")
	fmt.Println("                            # batch file example: for /f \"delims=\" %%v in ('pids.exe -vars') do %%v")
	fmt.Println("                            # then use %PID0%, %ProcessName2%, etc. in your batch file")
	fmt.Println()
	fmt.Println("  pids -name [-level N]     # Output process name for the Nth parent (default 0 = current process)")
	fmt.Println("                            # example: for /f %%v in ('pids.exe -name -level 3') do if \"%%v\"==\"explorer.exe\" timeout /t 10")
	fmt.Println()
	fmt.Println("  pids -pid [-level N]      # Output PID for the Nth parent (default 0 = current process)")
}

func getProcessChain() ([]*process.Process, error) {
	proc, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return nil, err
	}
	var chain []*process.Process
	for proc != nil {
		chain = append(chain, proc)
		parent, err := proc.Parent()
		if err != nil || parent == nil {
			break
		}
		proc = parent
	}
	return chain, nil
}

func main() {
	var (
		showVars  = flag.Bool("vars", false, "Output environment variable assignments for each parent process")
		showPID   = flag.Bool("pid", false, "Output PID for the specified level")
		showName  = flag.Bool("name", false, "Output process name for the specified level")
		levelFlag = flag.String("level", "0", "Specify which parent level to output (0 = current process)")
		help      = flag.Bool("h", false, "Show help")
	)
	flag.BoolVar(help, "help", false, "Show help")
	flag.Parse()

	// Show usage if no arguments are provided
	if *help || len(os.Args) == 1 {
		usage()
		return
	}

	chain, err := getProcessChain()
	if err != nil {
		log.Fatalf("Error getting process chain: %v\n", err)
	}

	// Only execute -vars logic if -vars is specified
	if *showVars {
		for i, proc := range chain {
			pid := proc.Pid
			name, err := proc.Name()
			if err != nil {
				name = "Unknown"
			}
			fmt.Printf("set PID%d=%d\n", i, pid)
			fmt.Printf("set ProcessName%d=%s\n", i, name)
		}
		return
	}

	// Parse level
	level, err := strconv.Atoi(*levelFlag)
	if err != nil || level < 0 || level >= len(chain) {
		fmt.Fprintf(os.Stderr, "Invalid or out-of-range level: %s\n", *levelFlag)
		os.Exit(1)
	}

	// -pid
	if *showPID {
		fmt.Println(chain[level].Pid)
		return
	}

	// -name
	if *showName {
		name, err := chain[level].Name()
		if err != nil {
			name = "Unknown"
		}
		fmt.Println(name)
		return
	}

	usage()
}
