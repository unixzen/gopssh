package main

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"time"
)

type Config struct {
	SshKeyPath  string   `mapstructure:"ssh_key_path"`
	SshUser     string   `mapstructure:"ssh_user"`
	SshPassword string   `mapstructure:"ssh_password"`
	Command     string   `mapstructure:"command"`
	Passphrase  string   `mapstructure:"passphrase"`
	SshPort     string   `mapstructure:"ssh_port"`
	Method      string   `mapstructure:"method"`
	Hosts       []string `mapstructure:"hosts"`
}

var Conf Config

// Read config file
func readConfig(filename string) {
	viper.SetConfigName(filename)
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		log.Fatal("Unable to decode into struct, %v", err)
	}
}

// Execute command at remote host
func executeCmd(remote_host string) string {
	readConfig("config")
	//	config, err := sshAuthKey(Conf.SshUser, Conf.SshKeyPath, Conf.Passphrase)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	var config *ssh.ClientConfig

	if Conf.Method == "password" {
		config = sshAuthPassword(Conf.SshUser, Conf.SshPassword)
	} else if Conf.Method == "key" {
		config = sshAuthKey(Conf.SshUser, Conf.SshKeyPath, Conf.Passphrase)
		//		if err != nil {
		//			log.Fatal(err)
		//		}
	} else {
		log.Fatal(`Please set method "password" or "key" at configuration file`)
	}

	client, err := ssh.Dial("tcp", remote_host+":"+Conf.SshPort, config)
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
	if err := session.Run(Conf.Command); err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
	fmt.Println("\x1b[31;1m" + remote_host + "\x1b[0m")
	return b.String()
}

// Pass username and password for authenticate by ssh at remote host
func sshAuthPassword(username string, password string) *ssh.ClientConfig {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config
}

// Pass username, path of private ssh key and passphrase for authenticate by ssh at remote host
func sshAuthKey(username string, ssh_key_path string, passphrase string) *ssh.ClientConfig {
	privateKey, err := ioutil.ReadFile(ssh_key_path)
	if err != nil {
		return &ssh.ClientConfig{}
	}

	signer, err := ssh.ParsePrivateKeyWithPassphrase(privateKey, []byte(passphrase))
	if err != nil {
		return &ssh.ClientConfig{}
	}

	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	return config
}

func main() {
	readConfig("config")
	results := make(chan string)
	timeout := time.After(10 * time.Second)

	for _, hostname := range Conf.Hosts {
		go func(hostname string) {
			results <- executeCmd(hostname)
		}(hostname)

	}

	for i := 0; i < len(Conf.Hosts); i++ {
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
