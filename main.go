package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

func main() {
	// Esperar conexiones
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port 8080")

	for {
		// Aceptar conexiones
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Cada conexión será una goroutine
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	// Guardar tamaño de la ventana, de los paquetes y el nombre del archivo enviados desde el cliente
	var windowSize, packetSize int32
	err := binary.Read(conn, binary.LittleEndian, &windowSize)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	err = binary.Read(conn, binary.LittleEndian, &packetSize)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fileName := make([]byte, 1024)
	bytesRead, err := conn.Read(fileName)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Crear nuevo archivo o abrir existente
	file, err := os.OpenFile(string(fileName[:bytesRead]), os.O_WRONLY|os.O_CREATE, 0666)
	defer file.Close()

	// Leer y procesar archivo enviado desde el cliente
	data := make([]byte, packetSize)
	for {
		bytesRead, err = conn.Read(data)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Guardar datos al archivo creado o abierto anteriormente
		_, err = file.Write(data[:bytesRead])
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		// Mostrar información de bytes recibidos en consola
		fmt.Println("Recibido: ", bytesRead, "bytes")

		// Enviar acuse de recibido después de cada ventana
		if bytesRead%int(windowSize) == 0 {
			_, err = conn.Write([]byte("ACK"))
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}

	}
}
