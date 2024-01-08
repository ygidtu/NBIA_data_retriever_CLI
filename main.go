package main

import (
	"fmt"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

var (
	// version and build info
	buildStamp string
	gitHash    string
	goVersion  string
	version    string
	client     *http.Client
	token      *Token
	logger     *zap.SugaredLogger
)

// SetupCloseHandler creates a 'listener' on a new goroutine which will notify the
// program if it receives an interrupt from the OS. We then handle this by calling
// our clean-up procedure and exiting the program.
func setupCloseHandler() {
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}

func main() {
	setupCloseHandler()

	var options = InitOptions()

	if options.Version {
		logger.Infof("Current version: %s", version)
		logger.Infof("Git Commit Hash: %s", gitHash)
		logger.Infof("UTC Build Time : %s", buildStamp)
		logger.Infof("Golang Version : %s", goVersion)
		os.Exit(0)
	} else {
		client = newClient(options.Proxy)

		err := os.MkdirAll(options.Output, os.ModePerm)
		if err != nil {
			logger.Fatalf("failed to create output directory: %v", err)
		}
		token, err = NewToken(
			options.Username, options.Password,
			filepath.Join(options.Output, fmt.Sprintf("%s.json", options.Username)))

		if err != nil {
			logger.Fatal(err)
		}

		var wg sync.WaitGroup
		files := decodeTCIA(options.Input)

		wg.Add(options.Concurrent)
		inputChan := make(chan *FileInfo, 5)
		for i := 0; i < options.Concurrent; i++ {

			go func(input chan *FileInfo) {
				defer wg.Done()
				for i := range input {
					if _, err := os.Stat(i.MetaFile(options.Output)); os.IsNotExist(err) {
						if !options.Meta {
							if err := i.Download(options.Output); err != nil {
								logger.Warnf("Download %s failed - %s", i.SeriesUID, err)
							} else {
								if err := i.GetMeta(options.Output); err != nil {
									logger.Warnf("save meta info %s failed - %s", i.SeriesUID, err)
								}
							}
						}
					} else {
						logger.Infof("Skip %s", i.SeriesUID)
					}
				}
			}(inputChan)
		}

		for _, f := range files {
			inputChan <- f
		}
		close(inputChan)
		wg.Wait()
	}
}
