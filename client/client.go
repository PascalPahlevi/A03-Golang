package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
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

const (
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 1024
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input the URL (e.g., http://localhost:8080/data): ")
	urlInput, _ := reader.ReadString('\n')
	urlInput = strings.TrimSpace(urlInput)

	fmt.Print("Input the data type (e.g., application/json): ")
	mimeType, _ := reader.ReadString('\n')
	mimeType = strings.TrimSpace(mimeType)

	fmt.Print("Input the language (e.g., en-US): ")
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

	conn, err := net.Dial(SERVER_TYPE, hostPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to server: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	response, students := Fetch(request, conn)

	fmt.Printf("Response Status: %s\n", response.StatusCode)
	fmt.Printf("Response ContentType: %s\n", response.ContentType)
	for _, student := range students {
		fmt.Printf("Student Name: %s, NPM: %s\n", student.Nama, student.Npm)
	}
}

func Fetch(req HttpRequest, conn net.Conn) (HttpResponse, []Student) {
	requestBytes := RequestEncoder(req)
	_, err := conn.Write(requestBytes)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return HttpResponse{}, nil
	}

	responseBytes := make([]byte, BUFFER_SIZE)
	n, err := conn.Read(responseBytes)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return HttpResponse{}, nil
	}

	response := ResponseDecoder(responseBytes[:n])

	var students []Student
	if response.ContentType == "application/json" && len(response.Data) > 0 {
		err = json.Unmarshal([]byte(response.Data), &students)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			return response, nil
		}
	}

	return response, students
}

func ResponseDecoder(bytestream []byte) HttpResponse {
	var res HttpResponse
	// Simplified decoding assuming the response is plain text
	resStr := string(bytestream)
	lines := strings.Split(resStr, "\r\n")
	for i, line := range lines {
		if i == 0 {
			parts := strings.Split(line, " ")
			if len(parts) > 2 {
				res.Version = parts[0]
				res.StatusCode = parts[1] + " " + parts[2]
			}
		} else if strings.HasPrefix(line, "Content-Type:") {
			res.ContentType = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if strings.HasPrefix(line, "Content-Language:") {
			res.ContentLanguage = strings.TrimSpace(strings.Split(line, ":")[1])
		} else if len(line) == 0 { // Headers and body are separated by an empty line
			if len(lines) > i+1 {
				res.Data = strings.Join(lines[i+1:], "\r\n")
			}
			break
		}
	}
	return res
}

func RequestEncoder(req HttpRequest) []byte {
	headers := fmt.Sprintf("Host: %s\r\nAccept: %s\r\nAccept-Language: %s\r\n", req.Host, req.Accept, req.AcceptLanguage)
	requestLine := fmt.Sprintf("%s %s %s\r\n", req.Method, req.Uri, req.Version)
	// End the request with an extra CRLF to denote the end of headers (and no body)
	return []byte(requestLine + headers + "\r\n")
}
