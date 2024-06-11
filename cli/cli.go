package cli

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

type RequestBody struct {
	SourceLang string `json:"sourceLang"`
	TargetLang string `json:"targetLang"`
	SourceText string `json:"sourceText"`
}

const translateUrl = "https://translate.googleapis.com/translate_a/single"

func RequestTranslate(body *RequestBody, str chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := &http.Client{}
	req, err := http.NewRequest("GET", translateUrl, nil)
	if err != nil {
		log.Fatal("there was an error creating the request:", err)
		return
	}

	query := req.URL.Query()
	query.Add("client", "gtx")
	query.Add("sl", body.SourceLang)
	query.Add("tl", body.TargetLang)
	query.Add("dt", "t")
	query.Add("q", body.SourceText)
	req.URL.RawQuery = query.Encode()

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("there was an error translating your request:", err)
		return
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusTooManyRequests {
		str <- "You have exceeded the maximum number of retries"
		return
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("error reading response body:", err)
		return
	}

	var responseData [][]interface{}
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		log.Fatal("error parsing JSON response:", err)
		return
	}

	if len(responseData) > 0 && len(responseData[0]) > 0 {
		translatedText := responseData[0][0].([]interface{})[0].(string)
		str <- translatedText
	} else {
		str <- "Translation failed"
	}
}
