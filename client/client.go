package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type HttpRequest struct {
	Method          string
	Uri             string
	Version         string
	Host            string
	Accept          string
	AcceptLanguange string
}

type HttpResponse struct {
	Version         string
	StatusCode      string
	ContentType     string
	ContentLanguage string
	Data            string
}

type Student struct {
	Nama string
	Npm  string
}

const (
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 1024
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input the URL: ")
	urlInput, _ := reader.ReadString('\n')
	urlInput = strings.TrimSpace(urlInput)

	fmt.Print("Input the data type: ")
	mimeType, _ := reader.ReadString('\n')
	mimeType = strings.TrimSpace(mimeType)

	fmt.Print("Input the language: ")
	language, _ := reader.ReadString('\n')
	language = strings.TrimSpace(language)

	urlParts := strings.Split(urlInput, "/")
	hostPort := urlParts[2]
	uri := "/" + strings.Join(urlParts[3:], "/")

	request := HttpRequest{
		Method:         "GET",
		Uri:            uri,
		Version:        "HTTP/1.1",
		Host:           hostPort,
		Accept:         mimeType,
		AcceptLanguage: language,
	}
}

func Fetch(req HttpRequest, connection net.Conn) (HttpResponse, []Student, HttpRequest) {
	//This program handles the request-making to the server
	var res HttpResponse
	var Student []Student

	return res, Student, req

}

func ResponseDecoder(bytestream []byte) HttpResponse {
	var res HttpResponse

	return res

}

func RequestEncoder(req HttpRequest) []byte {
	var result string

	return []byte(result)

}
