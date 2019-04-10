package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/crypto/ssh"
	"log"
	"os"
	"time"
)

func readConfig(filename string) (data []string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return data
}

func executeCmd(command, remote_host string, username string, password string) string {
	config := sshConfig(username, password)
	client, err := ssh.Dial("tcp", remote_host+":22", config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(command); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println("\x1b[31;1m" + remote_host + "\x1b[0m")
	return b.String()
}

func sshConfig(username string, password string) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config
}

func main() {

	username := flag.String("username", "", "Username for connection to host")
	password := flag.String("password", "", "Password for connection to host")
	command := flag.String("command", "", "Command which will be execute")
	flag.Parse()
	results := make(chan string)
	timeout := time.After(10 * time.Second)
	hosts := readConfig("./hosts")

	for _, hostname := range hosts {
		go func(hostname string) {
			results <- executeCmd(*command, hostname, *username, *password)
		}(hostname)

	}

	for i := 0; i < len(hosts); i++ {
		select {
		case res := <-results:
			fmt.Print(res)
		case <-timeout:
			fmt.Println("Timed out!")
			return
		}
	}
	close(results)

}
