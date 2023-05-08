package main

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type Envelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    Body
}

type Body struct {
	XMLName     xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Sum         Sum
	SumResponse SumResponse
}

type Sum struct {
	XMLName xml.Name `xml:"http://www.example.com/soap/calculator/ Sum"`
	A       int      `xml:"a"`
	B       int      `xml:"b"`
}

type SumResponse struct {
	XMLName xml.Name `xml:"http://www.example.com/soap/calculator/ SumResponse"`
	Result  int      `xml:"result"`
}

func main() {
	// Parâmetros para a operação de soma.
	a, b := 5, 3

	// Criar mensagem SOAP com a operação de soma.
	envelope := Envelope{
		Body: Body{
			Sum: Sum{
				A: a,
				B: b,
			},
		},
	}

	// Serializar a mensagem SOAP.
	payload, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling XML: %v\n", err)
		return
	}

	// Enviar a mensagem SOAP.
	url := "http://www.example.com/soap/calculator"
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(payload))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	// Processar a resposta SOAP.
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading HTTP response: %v\n", err)
		return
	}
	var respEnvelope Envelope
	if err := xml.Unmarshal(respData, &respEnvelope); err != nil {
		fmt.Printf("Error unmarshaling XML response: %v\n", err)
		return
	}

	// Extrair e exibir o resultado da operação de soma.
	sumResp := respEnvelope.Body.SumResponse
	fmt.Printf("%d + %d = %d\n", a, b, sumResp.Result)
}
