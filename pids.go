package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/shirou/gopsutil/process"
)

var version = "dev"

func usage() {
	fmt.Println()
	fmt.Println("https://github.com/Beej126/pids")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println()
	fmt.Printf("  -version             # %s\n", version)
	fmt.Println()
	fmt.Println("  -vars                # Output environment variable assignments for each parent process")
	fmt.Println("                       # batch file example: for /f \"delims=\" %%v in ('pids.exe -vars') do %%v")
	fmt.Println("                       # then use %PID0%, %ProcessName2%, etc. in your batch file")
	fmt.Println()
	fmt.Println("  -name [-level N]     # Output process name for the Nth parent (default 0 = current process)")
	fmt.Println("                       # batch example: for /f %%v in ('pids.exe -name -level 3') do if \"%%v\"==\"explorer.exe\" timeout /t 10")
	fmt.Println()
	fmt.Println("  -pid [-level N]      # Output PID for the Nth parent (default 0 = current process)")
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
		showVars     = flag.Bool("vars", false, "Output environment variable assignments for each parent process")
		showPID      = flag.Bool("pid", false, "Output PID for the specified level")
		showName     = flag.Bool("name", false, "Output process name for the specified level")
		levelFlag    = flag.String("level", "0", "Specify which parent level to output (0 = current process)")
		help1        = flag.Bool("h", false, "Show help")
		help2        = flag.Bool("help", false, "Show help")
		help3        = flag.Bool("?", false, "Show help")
		versionFlag1 = flag.Bool("version", false, "Show the program version from the file version resource")
		versionFlag2 = flag.Bool("v", false, "Show the program version from the file version resource")
	)

	flag.Parse()
	helpMerged := *help1 || *help2 || *help3 || len(os.Args) == 1
	// Show usage if no arguments are provided
	if helpMerged {
		usage()
		fmt.Println()
		fmt.Println("-vars output:")
		fmt.Println("-------------")
	}

	chain, err := getProcessChain()
	if err != nil {
		log.Fatalf("Error getting process chain: %v\n", err)
	}

	if *showVars || helpMerged {
		for i, proc := range chain {
			pid := proc.Pid
			name, err := proc.Name()
			if err != nil {
				name = "Unknown"
			}
			fmt.Printf("set PID%d=%d\n", i, pid)
			fmt.Printf("set ProcessName%d=%s\n", i, name)
		}
		if helpMerged {
			fmt.Println()
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

	if *versionFlag1 || *versionFlag2 {
		fmt.Println(version)
		return
	}
}
