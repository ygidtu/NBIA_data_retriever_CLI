package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/voxelbrain/goptions"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var (
	timeout = time.Duration(120000 * time.Second)
	proxy   = ""
	baseURL = ""
	output = ""
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
	if gin.IsDebugging() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

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


func downloadWrapper(input chan *FileInfo) {
	for i := range input {
		i.Get()
		i.Download(output)
	}
}



func main() {
	SetupCloseHandler()
	SetLogger()

	options := struct {
		InputFile string `goptions:"-i, --input, description='Path to tcia file'"`
		Output      string `goptions:"-o, --output, description='Output directory'"`
		Proxy     string `goptions:"-x, --proxy, description='Proxy'"`
		Timeout   int64  `goptions:"-t, --timeout, description='Due to limitation of target server, please set this timeout value as big as possible'"`
		Num int `goptions:"-p, --process, description='Start how many download at same time'"`
		Version   bool   `goptions:"-v, --version, description='Show version'"`
		Debug     bool   `goptions:"--debug, description='Show debug info'"`

		Help goptions.Help `goptions:"--help, description='Show this help'"`
	}{
		Output:     "downloads",
		Proxy:   "",
		Timeout: 1200000,
		Num:1,
	}
	goptions.ParseAndFail(&options)

	if len(os.Args) <= 1 {
		goptions.PrintHelp()
		os.Exit(0)
	}

	if options.Version {
		println("Current version is 0.0.1")
	} else {
		proxy = options.Proxy
		timeout = time.Duration(options.Timeout) * time.Second
		output = options.Output

		files := DecodeTCIA(options.InputFile)

		inputChan := make(chan *FileInfo, 5)
		for i := 0; i < options.Num; i ++ {
			go downloadWrapper(inputChan)
		}



		for _, f := range files {
			inputChan <- f
		}
	}
}
