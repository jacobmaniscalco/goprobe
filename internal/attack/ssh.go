package attack

import (
	"golang.org/x/crypto/ssh"
)

type sshAttackOptions struct {

	address     string
	addressList []string
	port        string
}


func test() {

	client := ssh.NewClient
}

