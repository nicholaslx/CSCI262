package auth

import (
	"Ass1/hash"
	"errors"
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"math/rand"
	"time"
)

func generateSalt(length int) string {
	rand.Seed(time.Now().UnixNano())
	salt := make([]byte, length)
	for i := range salt {
		salt[i] = '0' + byte(rand.Intn(10)) // Generates a digit between '0' and '9'
	}
	return string(salt)
}
// Function to authenticate a user
func CreateUser(username, password string, clearance int) error {
	if UserExists(username) {
		return errors.New("user already exists")
	}
	if len(password) < 8 {
		return errors.New("password should be a minimum of 8 characters")
	}

	salt := generateSalt(8) // Generates an 8-digit salt
	passSalt := password + salt
	hashedPassword := hash.HashMD5(passSalt)

	saltEntry := fmt.Sprintf("%s:%s\n", username, salt)
	err := appendToFile("salt.txt", saltEntry)
	if err != nil {
		return fmt.Errorf("could not write to salt.txt: %v", err)
	}

	shadowEntry := fmt.Sprintf("%s:%s:%d\n", username, hashedPassword, clearance)
	err = appendToFile("shadow.txt", shadowEntry)
	if err != nil {
		return fmt.Errorf("could not write to shadow.txt: %v", err)
	}

	fmt.Printf("User %s created successfully.\n", username)
	return nil
}

// UserExists checks if a user with the given username already exists.
func UserExists(username string) bool {
	data, err := ioutil.ReadFile("salt.txt")
	if err != nil {
		return false
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, username+":") {
			return true
		}
	}
	return false
}

// AuthenticateUser authenticates a user with a given username and password.
func AuthenticateUser(username, password string) (int, error) {
	if !UserExists(username) {
		return 0, errors.New("user does not exist")
	}

	salt, err := GetSalt(username)
	if err != nil {
		return 0, err
	}

	passSalt := password + salt
	hashedPassword := hash.HashMD5(passSalt)

	if authenticateShadow(username, hashedPassword) {
		clearance, err := GetClearance(username)
		if err != nil {
			return 0, errors.New("authentication failed: incorrect password")
		}
		return clearance, nil
	}
	return 0, errors.New("authentication failed: incorrect password")
}

// GetSalt retrieves the salt for a given username from the salt file.
func GetSalt(username string) (string, error) {
	data, err := ioutil.ReadFile("salt.txt")
	if err != nil {
		return "", err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 2 && parts[0] == username {
			return parts[1], nil
		}
	}
	return "", errors.New("username not found")
}

func authenticateShadow(username, hashedpwd string) bool {
	data, err := ioutil.ReadFile("shadow.txt")
	if err != nil {
		fmt.Println("shadow.txt not found")
		return false
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 3 && parts[0] == username && parts[1] == hashedpwd {
			return true
		}
	}
	return false
}

// GetClearance retrieves the clearance level for a given username.
func GetClearance(username string) (int, error) {
	data, err := ioutil.ReadFile("shadow.txt")
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(line, ":")
		if len(parts) == 3 && parts[0] == username {
			var clearance int
			fmt.Sscanf(parts[2], "%d", &clearance)
			return clearance, nil
		}
	}
	return 0, errors.New("clearance not found")
}

func appendToFile(filename, text string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.WriteString(text)
	return err
}