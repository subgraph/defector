package torcontrol

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
)


func Connect(method string, address string) (net.Conn, error) {
	connection, err := net.Dial(method, address)
	if err != nil {
		return connection, err
	}
	return connection, nil
}

func SendCommand(connection net.Conn, command string) (string, error) {
	fmt.Fprintf(connection, "%s\n", command)
	buf := make([]byte, 2048)
	nBytes, err := connection.Read(buf)
	if err != nil {
		fmt.Printf("%d bytes read\n", nBytes)
		panic(err)
	}
	response := string(buf)
	return response, nil
}

func Authenticate(password string) (string) {
	command := "AUTHENTICATE"
	if password != "" {
		command = fmt.Sprintf("AUTHENTICATE %s", password)
	}
	return command
}

func GetInfoStatusBootstrapPhase(connection net.Conn) (string) {
	cmd := fmt.Sprintf("GETINFO status/bootstrap-phase")
	response, err := SendCommand(connection, cmd)
	if err != nil {
		panic(err)
	}
	progressRe, err := regexp.Compile(`PROGRESS=(\d{0,3})`)
	progress := progressRe.FindStringSubmatch(response)[1]
	return progress				
}

func IsTorConnected() (bool) {
	connection, err := Connect("tcp", "127.0.0.1:9051")
	if err != nil {
		panic(err)
	}
	defer connection.Close()
	response, err := SendCommand(connection, Authenticate(""))
        if err != nil {
                panic(err)
        }
	response = GetInfoStatusBootstrapPhase(connection)
	progress, err := strconv.Atoi(response)
	if err != nil {
		panic(err)
	}
	if progress == 100 {
		return true
	}
	return false
}

