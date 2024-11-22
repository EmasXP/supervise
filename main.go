package main

import (
	"embed"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
)

//go:embed static/*
var static embed.FS

//go:embed templates/*
var views embed.FS

func build(isDev bool, accounts gin.Accounts) *gin.Engine {
	r := gin.Default()

	r.Use(gin.BasicAuth(accounts))

	if !isDev {
		templ := template.Must(template.New("").Funcs(r.FuncMap).ParseFS(views, "templates/*"))
		r.SetHTMLTemplate(templ)

		staticServer := http.FileServer(http.FS(static))
		staticHandler := func(c *gin.Context) {
			staticServer.ServeHTTP(c.Writer, c.Request)
		}
		r.GET("/static/*filepath", staticHandler)
		r.HEAD("/static/*filepath", staticHandler)
	} else {
		r.LoadHTMLGlob("templates/*")
		r.Static("/static", "./static")
	}

	r.GET("/", GetIndex)
	r.GET("/stdout/:program", GetStd("", "Stdout"))
	r.GET("/stderr/:program", GetStd("stderr", "Stderr"))
	r.POST("/start", ManageAction("start"))
	r.POST("/restart", ManageAction("restart"))
	r.POST("/stop", ManageAction("stop"))

	return r
}

func getDefaultConfigPath() string {
	if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		if _, err := os.Stat("/etc/supervise/settings.json"); os.IsNotExist(err) {
			log.Println("/etc/supervise/settings.json does not exist")
			os.Exit(2)
		}
		return "/etc/supervise/settings.json"
	}

	log.Println("Unknown default config file path for " + runtime.GOOS + ". Please use the -config flag.")
	os.Exit(2)
	return ""
}

type Config struct {
	Accounts gin.Accounts `json:"accounts"`
}

func main() {
	dev := flag.Bool("dev", false, "Dev mode. Will not use embedded files")
	configPath := flag.String("config", "", "Path to config file")
	addr := flag.String("addr", "127.0.0.1:9988", "Serve on this address")

	flag.Parse()

	configFile := *configPath
	if configFile == "" {
		configFile = getDefaultConfigPath()
	}

	if _, err := os.Stat(configFile); err != nil {
		log.Println(err)
		os.Exit(2)
	}

	configSource, err := os.ReadFile(configFile)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	var config Config
	err = json.Unmarshal(configSource, &config)
	if err != nil {
		log.Println(err)
		os.Exit(2)
	}

	r := build(
		*dev,
		config.Accounts,
	)

	err = r.Run(*addr)
	if err != nil {
		log.Println(err)
	}
}
