// Package controllers registers an HTTP handler at "/api/confusable" that
// renders the given string in raw, transcription of confusables and recovered from ocr
// and gives a judgement whether its profan or not.

package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/finnbear/moderation"
	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gogap/go-wkhtmltox/wkhtmltox"
	"github.com/mtibben/confusables"
	"log"
	"os"
)

type PayloadOCR struct {
	// Prepares the OCR request

	Base64 string `json:"base64"`
	Trim   string `json:"trim"`
}

type PayloadWK struct {
	// Prepares the WebKit request toimage

	To      string `json:"to"`
	Fetcher struct {
		Name   string `json:"name"`
		Params struct {
			Data string `json:"data"`
		} `json:"params"`
	} `json:"fetcher"`
	Converter json.RawMessage `json:"converter"`
}

type Step struct {
	// Orders the results of the different steps

	String  string `json:"string"`
	Profan bool `json:"profan"`
}

type ResponseWK struct {
	// Gives the response from WebKit a structure

	Code    int    `json:"code"`
	Message string `json:"message"`
	Result  struct {
		Data string `json:"data"`
	} `json:"result"`
}

type ResponseOCR struct {
	// Gives the response from GOSSERACT a structure
	Result  string `json:"result"`
	Version string `json:"version"`
}

func DoOCR(b string) string {
	// This function creates the ocr request


	// Create a Resty Client
	client := resty.New()
	request := &PayloadOCR{
		Base64: b,
		Trim:   "\n",
	}
	resp, err := client.R().
		EnableTrace().
		SetBody(request).
		Post(os.Getenv("GOSSERACT_SERVICE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	ocrResponse := ResponseOCR{}
	err = json.Unmarshal(resp.Body(), &ocrResponse)

	if err != nil {
		log.Fatal(err)
	}
	return ocrResponse.Result
}

func RenderImage(buf *bytes.Buffer, u string) string {
	// This function sends the string to WebKit

	testString := "<html lang='en'><head><meta http-equiv='Content-Type' content='text/html; charset=utf-8'/></head> <body> <p>" + u + "</p> </body></html>"
	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString([]byte(testString))
	// Create a Resty Client
	client := resty.New()
    var opts wkhtmltox.ConvertOptions
	request := PayloadWK{}
	request.Fetcher.Params.Data = encoded
	request.To = "image"
	request.Fetcher.Name = "data"
	opts = &wkhtmltox.ToImageOptions{
		Quality: 100,
	}
	request.Converter , _ = json.Marshal(opts)

	b, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
	}
	resp, err := client.R().
		EnableTrace().
		SetBody(b).
		Post(os.Getenv("WKHTMLTOX_SERVICE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	wkResponse := ResponseWK{}
	err = json.Unmarshal(resp.Body(), &wkResponse)

	if err != nil {
		log.Fatal(err)
	}
	return wkResponse.Result.Data

}

func CheckConfusable(c *fiber.Ctx) error {
	// Handlerfunction that calls all steps and returns the result
	
	type Request struct {
		ToCheck string `json:"toCheck"`
	}

	var body Request

	err := c.BodyParser(&body)
	// if error
	if err != nil {
		log.Println(err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Cannot parse JSON",
		})
	}
	transcribed := confusables.Skeleton(body.ToCheck)
	buf := new(bytes.Buffer)
	Base64Image := RenderImage(buf, body.ToCheck)
	OcrResult := DoOCR(Base64Image)
	Raw := &Step{
		String: body.ToCheck,
		Profan: moderation.IsInappropriate(body.ToCheck),
	}
	Transcribed := &Step{
		String: transcribed,
		Profan: moderation.IsInappropriate(transcribed),
	}
	Ocr := &Step{
		String: OcrResult,
		Profan: moderation.IsInappropriate(OcrResult),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"raw": Raw,
		"transcribed": Transcribed,
		"ocr": Ocr,
	})
}
