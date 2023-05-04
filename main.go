package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	e := echo.New()
	e.IPExtractor = echo.ExtractIPDirect()
	logger, _ := zap.NewProduction()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.String("IP", c.RealIP()),
				zap.Int("status", v.Status),
			)

			return nil
		},
	}))
	e.GET("/download/:filename", func(c echo.Context) error {
		filename := c.Param("filename")
		filePath := fmt.Sprintf("local/files/%s", filename)
		return c.File(filePath)
	})

	e.POST("/upload", func(c echo.Context) error {
		ip := c.RealIP()
		if ip == "::1" {
			ip = "127.0.0.1"
		}
		file, err := c.FormFile("file")
		if err != nil {
			return err
		}

		ext := filepath.Ext(file.Filename)
		timestamp := strconv.FormatInt(time.Now().Unix(), 10)
		newFilename := fmt.Sprintf("%s_%s%s", timestamp, strings.TrimSuffix(file.Filename, ext), ext)

		src, err := file.Open()
		if err != nil {
			return err
		}
		defer src.Close()

		err = os.MkdirAll(fmt.Sprintf("uploads/%s", ip), 0755)
		if err != nil {
			return err
		}

		dst, err := os.Create(fmt.Sprintf("uploads/%s/%s", ip, newFilename))
		if err != nil {
			return err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, map[string]string{
			"filename": newFilename,
		})
	})

	e.Logger.Fatal(e.Start(":8080"))
}
