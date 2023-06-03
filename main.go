package main

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	ApiKey = "123456789"
)

type MakePdfInput struct {
	ApiKey string `json:"apiKey" binding:"required"`
	Url    string `json:"url" binding:"required"`
	Size   string `json:"size" binding:"required"`
}

type SizeInput struct {
	Width  float64
	Height float64
}

func main() {
	setSystemLog()
	setGinLog()

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	r := setupRouter()
	r.Run()

}

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.Use(gin.Recovery())
	r.Static("/public", "./assets")
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	r.POST("/pdf", func(c *gin.Context) {
		var makePdfInput MakePdfInput

		if err := c.ShouldBindJSON(&makePdfInput); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"url":     "",
				"message": err.Error()})
			return
		}

		if makePdfInput.ApiKey != ApiKey {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"url":     "",
				"message": "ApiKey無效"})
			return
		}

		fileName, err := capturePdf(makePdfInput.Url, makePdfInput.Size)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"url":     fileName,
				"message": err.Error(),
			})
			return
		}
		if fileName != "" {
			fileName = c.Request.Host + fileName
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"url":     fileName,
			"message": "建立成功",
		})
	})
	return r
}

func capturePdf(url string, size string) (string, error) {
	var sizeDict = map[string]SizeInput{}

	sizeDict["A4"] = SizeInput{Width: 8.3, Height: 11.7}
	sizeDict["A5"] = SizeInput{Width: 5.8, Height: 8.3}
	sizeDict["A6"] = SizeInput{Width: 4.1, Height: 5.8}
	sizeDict["A7"] = SizeInput{Width: 2.9, Height: 4.1}

	sizeInput, ok := sizeDict[size]

	if !ok {
		return "", errors.New("尺吋格式錯誤")
	}

	id := uuid.New()

	fileName := id.String() + ".pdf"
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	var schbytes []byte
	var err error

	if err = chromedp.Run(ctx, printToPDF(url, &schbytes, &buf, sizeInput)); err != nil {
		return "", err
	}

	if err = os.WriteFile("./assets/"+fileName, buf, 0o644); err != nil {
		return "", err
	}

	return "/public/" + fileName, err
}

func printToPDF(urlstr string, schbytes, res *[]byte, size SizeInput) chromedp.Tasks {

	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.ActionFunc(func(ctx context.Context) error {

			err := emulation.SetDefaultBackgroundColorOverride().WithColor(&cdp.RGBA{R: 0, G: 0, B: 0, A: 0}).Do(ctx)
			if err != nil {
				return err
			}
			*schbytes, err = page.CaptureScreenshot().Do(ctx)
			if err != nil {
				return err
			}
			*res, _, err = page.PrintToPDF().WithPrintBackground(true).WithPaperWidth(size.Width).WithPaperHeight(size.Height).Do(ctx)
			if err != nil {
				return err
			}
			return nil

		}),
	}
}

func setSystemLog() {
	file, e1 := openLogFile("./system.log")
	if e1 != nil {
		log.Fatal(e1)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)
}

func setGinLog() {
	file, e2 := openLogFile("gin.log")
	if e2 != nil {
		log.Fatal(e2)
	}

	gin.DefaultWriter = io.MultiWriter(file)
}

func openLogFile(path string) (*os.File, error) {
	logFile, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return logFile, nil
}
