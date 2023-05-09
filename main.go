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
	XMLName       xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Body"`
	Placa         Placa
	PlacaResponse PlacaResponse
}

type Placa struct {
	XMLName xml.Name `xml:"http://webservice.credify.com.br/wscredify.php?wsdl Placa"`
	Valor   string   `xml:"placa"`
}

type PlacaResponse struct {
	XMLName xml.Name `xml:"http://webservice.credify.com.br/wscredify.php?wsdl PlacaResponse"`
	Result  int      `xml:"result"`
}

func main() {
	// Valores para autenticação e consulta.
	idConsulta := "371"
	usuario := "WS00000781"
	senha := "mL7hk9fBvc"
	placaOriginal := "EDQ-4711"

	// Criar mensagem SOAP com a operação de consulta de placa.
	soapBody := Body{
		Placa: Placa{
			Valor: placaOriginal,
		},
	}

	envelope := Envelope{
		Body: soapBody,
	}

	// Serializar a mensagem SOAP.
	payload, err := xml.MarshalIndent(envelope, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling XML: %v\n", err)
		return
	}

	// Enviar a mensagem SOAP com parâmetros de autenticação.
	url := "http://webservice.credify.com.br/wscredify.php?wsdl"
	client := &http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(payload))
	if err != nil {
		fmt.Printf("Error creating HTTP request: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.SetBasicAuth(usuario, senha)
	req.Header.Set("id_consulta", idConsulta)
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

	// Extrair e exibir o resultado da operação de consulta de placa.
	placaResp := respEnvelope.Body.PlacaResponse
	fmt.Printf("Placa %s encontrada com resultado %d\n", placaOriginal, placaResp.Result)
	fmt.Println("Renavam: ", placaResp.Result)
}
