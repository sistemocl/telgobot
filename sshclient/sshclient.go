package sshclient

import (
	"fmt"
	"io"
	"os"
	"time"

	"golang.org/x/crypto/ssh"
)

func ConexionErlang(user, password, ip, command string) (string, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", ip+":22", config)
	if err != nil {
		return "", fmt.Errorf("Error conectando al servidor %s: %v", ip, err)
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", fmt.Errorf("Error creando sesión SSH: %v", err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return "", fmt.Errorf("Error obteniendo stdin: %v", err)
	}

	stdout, err := session.StdoutPipe()
	if err != nil {
		return "", fmt.Errorf("Error obteniendo stdout: %v", err)
	}

	session.Stderr = os.Stderr

	// Inicia la consola remota de Erlang
	if err := session.Shell(); err != nil {
		return "", fmt.Errorf("Error iniciando shell: %v", err)
	}

	// Envía el comando para iniciar Erlang en modo consola remota
	_, err = fmt.Fprintln(stdin, "sudo /opt/butler_server/bin/butler_server remote_console")
	if err != nil {
		return "", fmt.Errorf("Error enviando comando para iniciar Erlang: %v", err)
	}

	// Agrega un delay para asegurar que el prompt de Erlang aparezca
	time.Sleep(2 * time.Second)

	// Envía el comando de login de Erlang
	_, err = fmt.Fprintln(stdin, `login("")`)
	if err != nil {
		return "", fmt.Errorf("Error enviando comando de login a Erlang: %v", err)
	}

	// Agrega un pequeño delay para que el prompt esté listo para recibir comandos
	time.Sleep(1 * time.Second)

	// Envía el comando especificado por el usuario
	_, err = fmt.Fprintln(stdin, command)
	if err != nil {
		return "", fmt.Errorf("Error enviando el comando a Erlang: %v", err)
	}

	// Leer la respuesta del comando ejecutado en Erlang
	output := make([]byte, 2048) // Tamaño de buffer aumentado para comandos más largos
	n, err := stdout.Read(output)
	if err != nil && err != io.EOF {
		return "", fmt.Errorf("Error leyendo respuesta de Erlang: %v", err)
	}

	// Devolver la salida del comando ejecutado
	return string(output[:n]), nil
}
