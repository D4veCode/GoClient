package main

import (
	"fmt"
	"os"
	"net"
	"strings"
)
func main(){	

	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "Error, debe especificar host, puerto o username \n")
		os.Exit(1)
	}
	
	username := os.Args[2]
	addr := os.Args[1]
	fmt.Fprintf(os.Stdout, "Estableciendo conexion a la direccion: %s con el nombre de usuario: %s \n", addr, username)
	
	conTCP, err := net.Dial("tcp", addr)
	checkError(err)
	
	helloiam := "helloiam "+username

	_, err = conTCP.Write([]byte(strings.TrimRight(helloiam, "\n")))
	checkError(err)

	respHello := make([]byte, 4096)

	length, err := conTCP.Read(respHello)

	checkError(err)

	if length > 0 {
		fmt.Fprintf(os.Stdout, string(respHello))
	}

	_, err = conTCP.Write([]byte("msglen"))
	checkError(err)

	resplen := make([]byte, 4096)

	lengthmsglen, err := conTCP.Read(resplen)

	if lengthmsglen > 0 {
		fmt.Fprintf(os.Stdout, string(resplen))
	}


	//TODO Fix Listener UDP
	/*
	listenerUDP, err := net.Listen("udp", "192.168.24.42:15006")

	checkError(err)

	_, err = conTCP.Write([]byte("givememsg 15006"))

	checkError(err)

	isOk := make([]byte, 1024)

	lenIsOk, err := conTCP.Read(isOk)

	if lenIsOk > 0{
		fmt.Fprintf(os.Stdout, "Mensaje recibido: %s",string(isOk))
		conUDP, err := listenerUDP.Accept()
		checkError(err)

		hiddenMsg := make([]byte, 4096)
		lenHiddenMsg, err := conUDP.Read(hiddenMsg)
		checkError(err)

		if lenHiddenMsg > 0 {
			fmt.Fprintf(os.Stdout, "El mensaje es: %s", string(hiddenMsg))
		}
		
	}
	*/

	




	

}

func checkError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s \n", err.Error())
		os.Exit(1)
	}
	
}