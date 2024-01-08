package main

import (
	"fmt"
	"github.com/DavidGamba/go-getoptions"
	"os"
	"path/filepath"
)

var (
	TokenUrl = "https://services.cancerimagingarchive.net/nbia-api/oauth/token"
	ImageUrl = "https://services.cancerimagingarchive.net/nbia-api/services/v2/getImage"
	MetaUrl  = "https://services.cancerimagingarchive.net/nbia-api/services/v2/getSeriesMetaData"
)

// Options command line parameters
type Options struct {
	Input      string
	Output     string
	Proxy      string
	Concurrent int
	Meta       bool
	Username   string
	Password   string
	Version    bool
	Debug      bool
	Help       bool
	MetaUrl    string
	TokenUrl   string
	ImageUrl   string
	SaveLog    bool

	opt *getoptions.GetOpt
}

func InitOptions() *Options {
	opt := &Options{opt: getoptions.New()}

	setLogger(false, "")

	opt.opt.BoolVar(&opt.Help, "help", false, opt.opt.Alias("h"),
		opt.opt.Description("show help information"))
	opt.opt.BoolVar(&opt.Debug, "debug", false,
		opt.opt.Description("show more info"))
	opt.opt.BoolVar(&opt.SaveLog, "save-log", false,
		opt.opt.Description("save debug log info to file"))
	opt.opt.BoolVar(&opt.Version, "version", false, opt.opt.Alias("v"),
		opt.opt.Description("show version information"))
	opt.opt.StringVar(&opt.Input, "input", "", opt.opt.Alias("i"),
		opt.opt.Description("path to input tcia file"))
	opt.opt.StringVar(&opt.Output, "output", "./", opt.opt.Alias("s"),
		opt.opt.Description("Output directory, or output file when --meta enabled"))
	opt.opt.StringVar(&opt.Proxy, "proxy", "", opt.opt.Alias("x"),
		opt.opt.Description("the proxy to use [http, socks5://user:passwd@host:port]"))
	opt.opt.IntVar(&opt.Concurrent, "processes", 1, opt.opt.Alias("p"),
		opt.opt.Description("start how many download at same time"))
	opt.opt.BoolVar(&opt.Meta, "meta", false, opt.opt.Alias("m"),
		opt.opt.Description("get Meta info of all files"))
	opt.opt.StringVar(&opt.Username, "user", "nbia_guest", opt.opt.Alias("u"),
		opt.opt.Description("username for control data"))
	opt.opt.StringVar(&opt.Password, "passwd", "", opt.opt.Alias("w"),
		opt.opt.Description("password for control data"))
	opt.opt.StringVar(&opt.TokenUrl, "token-url", TokenUrl,
		opt.opt.Description("the api url of login token"))
	opt.opt.StringVar(&opt.MetaUrl, "meta-url", MetaUrl,
		opt.opt.Description("the api url get meta data"))
	opt.opt.StringVar(&opt.ImageUrl, "image-url", ImageUrl,
		opt.opt.Description("the api url to download image data"))

	_, err := opt.opt.Parse(os.Args[1:])
	if err != nil {
		logger.Fatal(err)
	}

	if opt.Debug || opt.SaveLog {
		setLogger(opt.Debug, filepath.Join(opt.Output, "progress.log"))
	}

	if opt.opt.Called("help") || len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, opt.opt.Help())
		os.Exit(1)
	}

	if opt.TokenUrl != "" && opt.TokenUrl != TokenUrl {
		TokenUrl = opt.TokenUrl
		logger.Infof("Using custom token url: %s", TokenUrl)
	}

	if opt.MetaUrl != "" && opt.MetaUrl != MetaUrl {
		MetaUrl = opt.MetaUrl
		logger.Infof("Using custom meta url: %s", MetaUrl)
	}
	if opt.ImageUrl != "" && opt.ImageUrl != ImageUrl {
		ImageUrl = opt.ImageUrl
		logger.Infof("Using custom image url: %s", ImageUrl)
	}
	return opt
}
