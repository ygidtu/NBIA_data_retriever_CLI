package main

import (
	"fmt"
	"github.com/DavidGamba/go-getoptions"
	"os"
)

var (
	TokenUrl = "https://services.cancerimagingarchive.net/nbia-api/oauth/token"
	TagUrl   = "https://services.cancerimagingarchive.net/nbia-api/services/getDicomTags"
	MetaUrl  = "https://services.cancerimagingarchive.net/nbia-api/services/getSeriesMetadata2"
)

// Options command line parameters
type Options struct {
	Input      string
	Output     string
	Proxy      string
	Timeout    int
	Concurrent int
	Meta       bool
	Username   string
	Password   string
	Version    bool
	Debug      bool
	Help       bool

	opt *getoptions.GetOpt
}

func InitOptions() *Options {
	opt := &Options{opt: getoptions.New()}

	setLogger(false)

	opt.opt.BoolVar(&opt.Help, "help", false, opt.opt.Alias("h"),
		opt.opt.Description("show help information"))
	opt.opt.BoolVar(&opt.Debug, "debug", false,
		opt.opt.Description("show more info"))
	opt.opt.BoolVar(&opt.Version, "version", false, opt.opt.Alias("v"),
		opt.opt.Description("show version information"))
	opt.opt.StringVar(&opt.Input, "input", "", opt.opt.Alias("i"),
		opt.opt.Description("path to input tcia file"))
	opt.opt.StringVar(&opt.Output, "output", "./", opt.opt.Alias("s"),
		opt.opt.Description("Output directory, or output file when --meta enabled"))
	opt.opt.StringVar(&opt.Proxy, "proxy", "", opt.opt.Alias("x"),
		opt.opt.Description("the proxy to use [http, socks5://user:passwd@host:port]"))
	opt.opt.IntVar(&opt.Timeout, "timeout", 120000, opt.opt.Alias("t"),
		opt.opt.Description("due to limitation of target server, please set this timeout value as big as possible"))
	opt.opt.IntVar(&opt.Concurrent, "processes", 1, opt.opt.Alias("p"),
		opt.opt.Description("start how many download at same time"))
	opt.opt.BoolVar(&opt.Meta, "meta", false, opt.opt.Alias("m"),
		opt.opt.Description("get Meta info of all files"))
	opt.opt.StringVar(&opt.Username, "user", "nbia_guest", opt.opt.Alias("u"),
		opt.opt.Description("username for control data"))
	opt.opt.StringVar(&opt.Password, "passwd", "", opt.opt.Alias("w"),
		opt.opt.Description("password for control data"))

	_, err := opt.opt.Parse(os.Args[1:])
	if err != nil {
		logger.Fatal(err)
	}

	if opt.Debug {
		setLogger(opt.Debug)
	}

	if opt.opt.Called("help") || len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, opt.opt.Help())
		os.Exit(1)
	}

	return opt
}
