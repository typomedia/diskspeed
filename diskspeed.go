package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/pflag"
	"github.com/typomedia/diskspeed/app"
	"github.com/typomedia/diskspeed/bench"
	"github.com/typomedia/diskspeed/mem"
	"log"
	"os"
	"runtime"
	"time"
)

func main() {
	var (
		jsonout, verbose, version          bool
		numberOfRuns, numberOfSecondsToRun int
		err                                error
	)

	bm := bench.Mark{
		Start: time.Now(),
	}

	currentDir, err := os.Getwd()
	check(err)

	bm.PhysicalMemory, err = mem.Get()
	check(err)

	pflag.BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
	pflag.BoolVarP(&version, "version", "V", false, "Show version information")
	pflag.BoolVarP(&jsonout, "json", "j", false, "Will print JSON-formatted results to stdout")
	pflag.IntVarP(&numberOfRuns, "runs", "r", 1, "The number of test runs")
	pflag.IntVarP(&numberOfSecondsToRun, "seconds", "s", 0, "The time in seconds to run the test")
	pflag.IntVarP(&bm.NumReadersWriters, "threads", "t", runtime.NumCPU(), "The number of concurrent readers/writers. Defaults to the number of CPU cores")
	pflag.Float64VarP(&bm.AggregateTestFilesSizeInGiB, "gb", "g", float64(2*int(bm.PhysicalMemory>>20))/1024,
		"The amount of disk space to use (in GiB), defaults to twice the physical RAM")
	pflag.Float64VarP(&bm.IODuration, "iops-duration", "i", 15.0,
		"The duration in seconds to run the IOPS benchmark, set to 0.5 for quick feedback during development")
	pflag.StringVarP(&currentDir, "dir", "d", currentDir, "The directory to use for the test. Defaults to the current directory")
	pflag.Parse()

	if version {
		fmt.Printf("%s %s\n", app.App.Name, app.App.Version)
		fmt.Println(app.App.Description)
		fmt.Println(app.App.Author)
		os.Exit(0)
	}

	check(bm.SetTempDir(currentDir))
	defer os.RemoveAll(bm.TempDir)

	if !jsonout {
		fmt.Printf("%s %s\n", app.App.Name, app.App.Version)
	}
	if verbose {
		fmt.Printf("runs: %d, seconds: %d, threads: %d, disk space to use: %d MB\n",
			numberOfRuns,
			numberOfSecondsToRun,
			bm.NumReadersWriters,
			int(bm.AggregateTestFilesSizeInGiB*(1<<10)))
		log.Printf("Number of CPU cores: %d", runtime.NumCPU())
		log.Printf("Total system RAM: %d MB", bm.PhysicalMemory>>20)
		log.Printf("Working directory: %s", currentDir)
	}

	check(bm.CreateRandomBlock())
	go bench.ClearBufferCacheEveryThreeSeconds() // flush the Buffer Cache every 3 seconds

	finishTime := bm.Start.Add(time.Duration(numberOfSecondsToRun) * time.Second)
	for i := 0; (i < numberOfRuns) || time.Now().Before(finishTime); i++ {
		check(bm.RunSequentialWriteTest())
		if verbose {
			log.Printf("Written (MiB): %d\n", bm.Results[i].WrittenBytes>>20)
			log.Printf("Written (MB): %f\n", float64(bm.Results[i].WrittenBytes)/1000000)
			log.Printf("Duration (seconds): %f\n", bm.Results[i].WrittenDuration.Seconds())
		}
		if !jsonout {
			fmt.Printf("Sequential Write MB/s: %0.2f\n",
				bench.MegaBytesPerSecond(bm.Results[i].WrittenBytes, bm.Results[i].WrittenDuration))
		}

		check(bm.RunSequentialReadTest())
		if verbose {
			log.Printf("Read (MiB): %d\n", bm.Results[i].ReadBytes>>20)
			log.Printf("Read (MB): %f\n", float64(bm.Results[i].ReadBytes)/1000000)
			log.Printf("Duration (seconds): %f\n", bm.Results[i].ReadDuration.Seconds())
		}
		if !jsonout {
			fmt.Printf("Sequential Read MB/s: %0.2f\n",
				bench.MegaBytesPerSecond(bm.Results[i].ReadBytes, bm.Results[i].ReadDuration))
		}

		check(bm.RunIOPSTest())
		if verbose {
			log.Printf("operations %d\n", bm.Results[i].IOOperations)
			log.Printf("Duration (seconds): %f\n", bm.Results[i].IODuration.Seconds())
		}
		if !jsonout {
			fmt.Printf("IOPS: %0.0f\n",
				bench.IOPS(bm.Results[i].IOOperations, bm.Results[i].IODuration))
		}
	}
	if jsonout {
		json.NewEncoder(os.Stdout).Encode(bm)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
