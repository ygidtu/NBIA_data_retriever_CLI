package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/voxelbrain/goptions"
)

var (
	timeout = time.Duration(120000 * time.Second)
	proxy   = ""
	baseURL = ""
	output  = ""
	meta    = false
)

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean up procedure and exiting the program.
func SetupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func SetLogger() zerolog.Logger {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	output := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339Nano,
		NoColor:    false,
	}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s |", i))
	}

	log.Logger = log.Output(output)

	return zerolog.New(output).With().Timestamp().Logger()
}

func main() {
	SetupCloseHandler()
	SetLogger()

	options := struct {
		InputFile string `goptions:"-i, --input, description='Path to tcia file'"`
		Output    string `goptions:"-o, --output, description='Output directory, or output file when --meta enabled'"`
		Proxy     string `goptions:"-x, --proxy, description='Proxy'"`
		Timeout   int64  `goptions:"-t, --timeout, description='Due to limitation of target server, please set this timeout value as big as possible'"`
		Num       int    `goptions:"-p, --process, description='Start how many download at same time'"`
		Meta      bool   `goptions:"-m, --meta, description='Get Meta info of all files'"`
		Version   bool   `goptions:"-v, --version, description='Show version'"`
		Debug     bool   `goptions:"--debug, description='Show debug info'"`

		Help goptions.Help `goptions:"--help, description='Show this help'"`
	}{
		Output:  "downloads",
		Proxy:   "",
		Timeout: 1200000,
		Num:     1,
		Meta:    false,
	}
	goptions.ParseAndFail(&options)

	if len(os.Args) <= 1 {
		goptions.PrintHelp()
		os.Exit(0)
	}

	if options.Version {
		println("Current version is 0.2.1")
	} else {
		proxy = options.Proxy
		timeout = time.Duration(options.Timeout) * time.Second
		output = options.Output
		meta = options.Meta
		var wg sync.WaitGroup

		files := DecodeTCIA(options.InputFile)

		wg.Add(options.Num)
		inputChan := make(chan *FileInfo, 5)
		for i := 0; i < options.Num; i++ {

			go func(input chan *FileInfo) {
				defer wg.Done()
				for i := range input {
					i.Get()
					if _, err := os.Stat(fmt.Sprintf("%s.json", i.GetOutput(output))); os.IsNotExist(err) {
						if !meta {
							i.Download(output)
							i.ToJson(output)
						}
					} else {
						log.Info().Msgf("Skip %s", i.SeriesUID)
					}
				}
			}(inputChan)
		}

		for _, f := range files {
			inputChan <- f
		}
		close(inputChan)
		wg.Wait()

		if meta {
			ToJson(files, fmt.Sprintf("%s.json", output))
		}
	}
}
