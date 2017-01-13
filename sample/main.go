package main

import (
	"flag"
	"net/http"
	"time"

	// github library
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"

	// Test library
	"github.com/cutedogspark/go-deamon"
)

var (
	NWorkers = flag.Int("n", 4, "The number of workers")
	HTTPAddr = flag.String("http", "127.0.0.1:8000", "Address to listen for HTTP requests on")
)

func Collector(c echo.Context) error {
	x := gworker.NewJob()
	x.SetJobType(gworker.JobTypeCMD)
	x.At(time.Now().UTC().Add(5 * time.Second).Format("2006-01-02 15:04:05"))
	x.SetJobContext("ls")

	status := gworker.Service.PutWorkQueue(x)
	if !status{
		return c.JSON(http.StatusInternalServerError, "WorkerQueue Error")
	}
	return c.JSON(http.StatusCreated, nil)

}

func main() {
	flag.Parse()

	/* GWorker Init & close */
	gworker.InitWorker(*NWorkers)
	gworker.Service.Start()
	defer gworker.Service.Stop()

	e := echo.New()
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	e.POST("/work", Collector)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "")
	})

	e.GET("/add", Collector)

	e.Logger.Fatal(e.Start(*HTTPAddr))
}