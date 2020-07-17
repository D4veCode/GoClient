package main

import (
	"fmt"
	"os"
	"net"
	"strings"
	b64 "encoding/base64"
	"crypto/md5"
	"encoding/hex"
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



	listenerUDP, err := net.ListenPacket("udp", "192.168.24.42:15006")


	checkError(err)

	_, err = conTCP.Write([]byte("givememsg 15006"))

	checkError(err)

	isOk := make([]byte, 1024)

	lenIsOk, err := conTCP.Read(isOk)

	

	if lenIsOk > 0{
		fmt.Fprintf(os.Stdout, "Mensaje recibido: %s",string(isOk))
		hiddenMsg := make([]byte, 1024)
		defer listenerUDP.Close()
		for {
			bytesSent, _ , err := listenerUDP.ReadFrom(hiddenMsg)
			checkError(err)
			if bytesSent > 0 {
				decodedMsg, _ := b64.StdEncoding.DecodeString(string(hiddenMsg))
				fmt.Fprintf(os.Stdout, "El mensaje es: %s \n", string(decodedMsg))
				break
			}
		}
		
		//TODO fix md5 
		hash := md5.New()

		hash.Write(hiddenMsg)
		hashedchk := hex.EncodeToString(hash.Sum(nil))
		
		checksum := "chkmsg "+ hashedchk

		_, err = conTCP.Write([]byte(checksum))
		checkError(err)
		checksumCheck := make([]byte, 1024)

		lenChecksum, _ := conTCP.Read(checksumCheck)
		checkError(err)

		if lenChecksum > 0 {
			fmt.Fprintf(os.Stdout, "Checksum : %s \n", string(checksumCheck))
		}

		_, err = conTCP.Write([]byte("bye"))
		checkError(err)

		bye := make([]byte, 16)

		lenBye, err := conTCP.Read(bye)
		checkError(err)

		if lenBye > 0 {
			fmt.Fprintf(os.Stdout, "%s \n", string(bye))
		}

		conTCP.Close()

		fmt.Fprintln(os.Stdout, "Gracias por usar el cliente, hasta luego")
	}

}

func checkError(err error){
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s \n", err.Error())
		os.Exit(1)
	}
	
}