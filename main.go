package main

import (
	"context"
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
}

func main() {

	r := gin.Default()
	r.Static("/public", "./assets")

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

		message := "建立成功"
		success := true
		fileName, err := capturePdf(makePdfInput.Url)

		if err != nil {
			success = false
			message = "建立失敗"
		}
		if fileName != "" {
			fileName = c.Request.Host + fileName
		}

		c.JSON(http.StatusOK, gin.H{
			"success": success,
			"url":     fileName,
			"message": message,
		})
	})
	r.Run()

}

func capturePdf(url string) (string, error) {
	id := uuid.New()

	fileName := id.String() + ".pdf"
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	var buf []byte
	var schbytes []byte
	var err error

	if err = chromedp.Run(ctx, printToPDF(url, &schbytes, &buf)); err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile("./assets/"+fileName, buf, 0o644); err != nil {
		log.Fatal(err)
	}

	return "/public/" + fileName, err
}

func printToPDF(urlstr string, schbytes, res *[]byte) chromedp.Tasks {
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
			*res, _, err = page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}
			return nil

		}),
	}
}
