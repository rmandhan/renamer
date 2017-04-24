package main

import (
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

func logRed(format string, v ...interface{}) {
	c := color.New(color.FgRed)
	logMessageWithColor(fmt.Sprintf(format, v...), c)
}

func logGreen(format string, v ...interface{}) {
	c := color.New(color.FgGreen)
	logMessageWithColor(fmt.Sprintf(format, v...), c)
}

func logYellow(format string, v ...interface{}) {
	c := color.New(color.FgYellow)
	logMessageWithColor(fmt.Sprintf(format, v...), c)
}

func logBlue(format string, v ...interface{}) {
	c := color.New(color.FgBlue)
	logMessageWithColor(fmt.Sprintf(format, v...), c)
}

func logMagenta(format string, v ...interface{}) {
	c := color.New(color.FgMagenta)
	logMessageWithColor(fmt.Sprintf(format, v...), c)
}

func logCyan(format string, v ...interface{}) {
	c := color.New(color.FgCyan)
	logMessageWithColor(fmt.Sprintf(format, v...), c)
}

func logWhite(format string, v ...interface{}) {
	c := color.New(color.FgWhite)
	logMessageWithColor(fmt.Sprintf(format, v...), c)
}

func logBlack(format string, v ...interface{}) {
	c := color.New(color.FgBlack)
	logMessageWithColor(fmt.Sprintf(format, v...), c)
}

func logMessageWithColor(m string, c *color.Color) {
	log.Println(c.Sprintf(m))
}

func logDebug(flag bool, format string, v ...interface{}) {
	if flag {
		logBlue(format, v...)
	}
}

func logDefault(format string, v ...interface{}) {
	log.Println(fmt.Sprintf(format, v...))
}

func logFatal(format string, v ...interface{}) {
	logRed(format, v...)
	os.Exit(1)
}
