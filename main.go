package main

import (
	"GoogleTranslate/cli"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
)

var wg sync.WaitGroup

var sourceLang string
var targetLang string
var sourceText string

func init() {
	flag.StringVar(&sourceLang, "s", "en", "source language [default: en]")
	flag.StringVar(&targetLang, "t", "fr", "target language [default: fr]")
	flag.StringVar(&sourceText, "st", "", "Text to translate")
}

func main() {
	flag.Parse()

	if sourceText == "" {
		fmt.Println("Option:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	strChan := make(chan string)
	wg.Add(1)

	reqBody := &cli.RequestBody{
		SourceLang: sourceLang,
		TargetLang: targetLang,
		SourceText: sourceText,
	}

	go cli.RequestTranslate(reqBody, strChan, &wg)
	processedStr := strings.ReplaceAll(<-strChan, "+", " ")
	fmt.Printf("%s\n", processedStr)
	close(strChan)
	wg.Wait()
}
