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
	//Chequeo para saber si esta el host, puerto y username
	if len(os.Args) != 3 {
		fmt.Printf("Error, debe especificar host, puerto o username \n")
		os.Exit(1)
	}
	
	addr := os.Args[1]
	username := os.Args[2]
	fmt.Printf("Estableciendo conexion a la direccion: %s con el nombre de usuario: %s \n", addr, username)
	//conexion con el servidor
	conTCP, err := net.Dial("tcp", addr)
	checkError(err)
	
	helloiam := "helloiam "+username

	_, err = conTCP.Write([]byte(strings.TrimRight(helloiam, "\n")))
	checkError(err)

	respHello := make([]byte, 4096)

	length, err := conTCP.Read(respHello)

	checkError(err)

	if length > 0 {
		fmt.Println(string(respHello))
	}

	_, err = conTCP.Write([]byte("msglen"))
	checkError(err)

	resplen := make([]byte, 4096)

	lengthmsglen, err := conTCP.Read(resplen)

	if lengthmsglen > 0 {
		fmt.Println(string(resplen))
	}
	//creando socket udp para recibir el mensaje
	listenerUDP, err := net.ListenPacket("udp", "192.168.24.42:15006")
	checkError(err)

	_, err = conTCP.Write([]byte("givememsg 15006"))
	checkError(err)

	isOk := make([]byte, 1024)

	lenIsOk, err := conTCP.Read(isOk)

	if lenIsOk > 0{
		
		fmt.Fprintf(os.Stdout, "Mensaje recibido: %s",string(isOk))
		hiddenMsg := make([]byte, 1024)
		//cerrando el socket UDP luego de recibir el mensaje
		defer listenerUDP.Close()

		for {
			//leer el mensaje enviado por UDP desde el servidor
			bytesSent, _ , err := listenerUDP.ReadFrom(hiddenMsg)
			checkError(err)
			if bytesSent > 0 {
				//decodificacion del mensaje enviado en base64
				decodedMsg, _ := b64.StdEncoding.DecodeString(string(hiddenMsg))
				fmt.Fprintf(os.Stdout, "El mensaje es: %s \n", string(decodedMsg))
				//creando hash md5 para confirmar totalidad del mensaje
				hash := md5.New()
				hash.Write(decodedMsg)
				hashedchk := hex.EncodeToString(hash.Sum(nil))
				
				checksum := "chkmsg "+ hashedchk
				//comprobacion con el servidor del checksum del mensaje
				_, err = conTCP.Write([]byte(checksum))
				checkError(err)
				checksumCheck := make([]byte, 1024)

				lenChecksum, _ := conTCP.Read(checksumCheck)
				checkError(err)

				if lenChecksum > 0 {
					fmt.Fprintf(os.Stdout, "Checksum : %s \n", string(checksumCheck))
				}
				break
			}
		}		
		//cerrar conexion con el servidor
		_, err = conTCP.Write([]byte("bye"))
		checkError(err)

		bye := make([]byte, 16)

		_, err := conTCP.Read(bye)
		checkError(err)
		//cerrando el socket TCP
		conTCP.Close()
		byeMsg := strings.Split(string(bye), " ")[1]

		fmt.Printf("Gracias por usar el cliente, %s \n", byeMsg)
	}

}

func checkError(err error){
	//comprobacion de si existe un error
	if err != nil {
		fmt.Printf("Error: %s \n", err.Error())
		os.Exit(1)
	}
	
}