package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"strings"
)

const (
	SERVER_HOST = "127.0.0.1"
	SERVER_PORT = "8080"
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 1024
	GROUP_NAME  = "CN01"
)

type HttpRequest struct {
	Method         string
	Uri            string
	Version        string
	Host           string
	Accept         string
	AcceptLanguage string
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

func main() {
	serverAddress, err := net.ResolveTCPAddr(SERVER_TYPE, net.JoinHostPort(SERVER_HOST, SERVER_PORT))
	if err != nil {
		log.Fatalln(err)
	}
	socket, err := net.ListenTCP(SERVER_TYPE, serverAddress)
	if err != nil {
		log.Fatalln(err)
	}

	defer socket.Close()

	fmt.Printf("TCP Server Socket Program in Go\n")
	fmt.Printf("[%s] Preparing TCP listening socket on %s\n", SERVER_TYPE, socket.Addr().String())

	for {
		connection, err := socket.AcceptTCP()
		if err != nil {
			log.Fatalln(err)
		}

		go HandleConnection(connection)
	}
}

func HandleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, BUFFER_SIZE)
	length, err := conn.Read(buffer)
	if err != nil {
		log.Println(err)
		return
	}

	rawRequest := buffer[:length]
	request := RequestDecoder(rawRequest)
	response := HandleRequest(request)
	encodedResponse := ResponseEncoder(response)

	_, err = conn.Write(encodedResponse)
	if err != nil {
		log.Println(err)
	}
}

func HandleRequest(req HttpRequest) HttpResponse {
	var res HttpResponse
	res.Version = "HTTP/1.1"

	students := []Student{
		{"Muhamad Pascal Alfin Pahlevi", "2206046752"},
		{"Tiva Adhisti Nafira Putri", "2206046840"},
		{"Raisa Fadilla", "2206822414"},
	}

	if req.Uri == "/" || req.Uri == "/?name="+GROUP_NAME {
		res.StatusCode = "200 OK"
		res.ContentType = "text/html"
		res.Data = fmt.Sprintf("<html><body>Welcome to %s homepage!</body></html>", GROUP_NAME)
	} else if req.Uri == "/data" {
		if strings.Contains(req.Accept, "application/xml") {
			res.ContentType = "application/xml"
			xmlData, err := xml.Marshal(students)
			if err != nil {
				log.Println("XML marshalling error:", err)
				res.StatusCode = "500 Internal Server Error"
				return res
			}
			res.Data = string(xmlData)
		} else { // Default to JSON
			res.ContentType = "application/json"
			jsonData, err := json.Marshal(students)
			if err != nil {
				log.Println("JSON marshalling error:", err)
				res.StatusCode = "500 Internal Server Error"
				return res
			}
			res.Data = string(jsonData)
		}
		res.StatusCode = "200 OK"
	} else if req.Uri == "/greeting" {
        res.StatusCode = "200 OK"
        res.ContentType = "text/html"
        greeting := "Hello, We are from CN01!" // Default greeting
        if strings.Contains(strings.ToLower(req.AcceptLanguage), "id") {
            greeting = "Halo, Kami dari CN01!"
        }
        res.Data = fmt.Sprintf("<html><body>%s</body></html>", greeting)
    } else {
        res.StatusCode = "404 Not Found"
        res.Data = "404 Page Not Found"
    }

	return res
}

func RequestDecoder(bytestream []byte) HttpRequest {
	var req HttpRequest
	lines := strings.Split(string(bytestream), "\r\n")
	for i, line := range lines {
		if i == 0 {
			parts := strings.Split(line, " ")
			if len(parts) >= 3 {
				req.Method = parts[0]
				req.Uri = parts[1]
				req.Version = parts[2]
			}
		} else {
			headerParts := strings.SplitN(line, ": ", 2)
			if len(headerParts) == 2 {
				switch headerParts[0] {
				case "Host":
					req.Host = headerParts[1]
				case "Accept":
					req.Accept = headerParts[1]
				case "Accept-Language":
					req.AcceptLanguage = headerParts[1]
				}
			}
		}
	}
	return req
}

func ResponseEncoder(res HttpResponse) []byte {
	//Put the encoding program for HTTP Response Struct here
	var result strings.Builder
	result.WriteString(fmt.Sprintf("%s %s\r\n", res.Version, res.StatusCode))
	result.WriteString(fmt.Sprintf("Content-Type: %s\r\n", res.ContentType))
	if res.ContentLanguage != "" {
		result.WriteString(fmt.Sprintf("Content-Language: %s\r\n", res.ContentLanguage))
	}
	result.WriteString("\r\n")
	result.WriteString(res.Data)
	return []byte(result.String())
}
