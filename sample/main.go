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

var QWorker gworker.WorkerManage

func Collector(c echo.Context) error {
	delay, err := time.ParseDuration(c.FormValue("delay"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Bad delay value: " + err.Error())
	}

	if delay.Seconds() < 1 || delay.Seconds() > 10 {
		return c.JSON(http.StatusBadRequest, "The delay must be between 1 and 10 seconds.")
	}

	name := c.FormValue("name")
	if name == "" {
		return c.JSON(http.StatusBadRequest, "You must specify a name")
	}

	work := gworker.Job{Name: name, Delay: delay}
	status := QWorker.PutWorkQueue(work)
	if !status{
		return c.JSON(http.StatusInternalServerError, "WorkerQueue Error")
	}
	return c.JSON(http.StatusCreated, nil)

}

func main() {
	flag.Parse()

	QWorker = gworker.InitWorker(*NWorkers)
	QWorker.Start()

	e := echo.New()
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.Gzip())

	e.POST("/work", Collector)

	e.GET("/", func(c echo.Context) error {
		QWorker.Stop()
		return c.String(http.StatusOK, "Hello, World!")
	})
	e.GET("/stop", func(c echo.Context) error {
		QWorker.Stop()
		return c.String(http.StatusOK, "Stop!!!")
	})
	e.GET("/start", func(c echo.Context) error {
		QWorker.Start()
		return c.String(http.StatusOK, "Start!!!")
	})

	e.Logger.Fatal(e.Start(*HTTPAddr))
}