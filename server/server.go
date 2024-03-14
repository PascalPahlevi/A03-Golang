package main

import (
	"net"
)

const (
	SERVER_HOST = ""
	SERVER_PORT = ""
	SERVER_TYPE = "tcp"
	BUFFER_SIZE = 1024
	GROUP_NAME  = "CN01"
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

func main() {
	serverAddress err:= net.resolveTCPAddr(serverType, net.joinHotPort(serverIP, serverPort))
	if err!= nil {
		log.Fatalln(err)
	}
	socket, err := net.ListenTCP(serverType, serverAddress)
	if err != nil {
		log.Fatalln(err)
	}

	defer socket.Close()

	fmt.printf("TCP Sever Socket Program in Go\n")
	fmt.printf("[%s] Preparing TCP listening socket on %s\n", serverType, socket.LocalAddr())

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
		var greeting string
		if strings.Contains(req.AcceptLanguage, "id") {
			greeting = "Halo, Kami dari CN01!"
		} else if strings.Contains(req.AcceptLanguage, "en") {
			greeting = "Hello, We are from CN01!"
		}
		res.Data = fmt.Sprintf("<html><body>%s</body></html>", greeting)
	} else {
		res.StatusCode = "404 Not Found"
	}

	return res
}

func RequestDecoder(bytestream []byte) HttpRequest {
	//Put the decoding program for HTTP Request Packet here
	var req HttpRequest

	return req

}

func ResponseEncoder(res HttpResponse) []byte {
	//Put the encoding program for HTTP Response Struct here
	var result string

	return []byte(result)

}
