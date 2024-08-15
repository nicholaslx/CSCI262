package main

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
	"strconv"
)

// User structure
// type User struct {
// 	Username string
// 	Password string // This will be the hashed password
// 	IsAdmin  bool
// 	Salt     string
// }

// File structure
type File struct {
	Filename       string
	Owner          string
	Classification int
	Content        string
}

var files []File // List of files in the system

// Helper function to generate a random 8-digit salt
func generateSalt(length int) string {
	rand.Seed(time.Now().UnixNano())
	salt := make([]byte, length)
	for i := range salt {
		salt[i] = '0' + byte(rand.Intn(10)) // Generates a digit between '0' and '9'
	}
	return string(salt)
}

// Helper function to hash a password using MD5
func hashMD5(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// Function to create a new user
func createUser(username, password string, clearance int) {
	if userExists(username) {
		fmt.Println("Error: User already exists.")
		return
	}
	if len(password) < 8 {
		fmt.Println("Error: Invalid username, password should be minimum of 8 characters")
		return
	}

	salt := generateSalt(8) // Generates an 8-digit salt
	passSalt := password + salt
	hashedPassword := hashMD5(passSalt)

	saltentry := fmt.Sprintf("%s:%s\n", username, salt)
	err := appendtofile("salt.txt", saltentry)
	if err != nil {
		fmt.Println("Error: Could not write to salt.txt")
		return
	}

	shadowentry := fmt.Sprintf("%s:%s:%v\n", username, hashedPassword, clearance)
	err = appendtofile("shadow.txt", shadowentry)
	if err != nil {
		fmt.Println("Error: Could not write to shadow.txt")
		return
	}

	fmt.Printf("User %s created successfully.", username)
}
func appendtofile(filename, text string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(text)
	return err
}

// Function to check if a user exists
func userExists(username string) bool {
	data, err := ioutil.ReadFile("salt.txt")
	if err != nil {
		return false // File doesn't exist or other error, assume user doesn't exist
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, username+":") {
			return true
		}
	}
	return false
}

// Function to authenticate a user
func authenticateUser(username, password string) {
	// Check if the user exists
	if !userExists(username) {
		fmt.Printf("User %s does not exist. Terminating program.\n", username)
		os.Exit(1) // Exit the program
	}

	// Retrieve the salt and hashed password for the user
	salt, err := getSalt(username)
	if err != nil {
		fmt.Printf("User %s does not exist. Terminating program.\n", username)
		return
	}

	fmt.Printf("%s found in salt.txt\n", username)
	fmt.Printf("Salt retrieved: %s\n", salt)

	passSalt:= password + salt
	hashedPassword:= hashMD5(passSalt)

	fmt.Printf("hash value:%s\n",hashedPassword)

	if authenticateShadow(username, hashedPassword) {
		fmt.Printf("Authentication for user %s complete.\n",username)
		clearance, err := getClearance(username)
		if err != nil{
			fmt.Println("Authentication failed. Incorrect password.")
		}
		fmt.Printf("The clearance for %s is %v\n",username, clearance)
	}

}

func getSalt(username string) (string, error) {
	data, err := ioutil.ReadFile("salt.txt")
	if err != nil {
		return "", err // File doesn't exist or other error, assume user doesn't exist
	}
	//scans salt.txt file to find salt value
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		parts := strings.Split(string(line), ":")
		if len(parts) == 2 && parts[0] == username {
			return parts[1], nil
		}
	}
	return "", errors.New("username not found")
}

func authenticateShadow(username, hashedpwd string)bool{
	data,err := ioutil.ReadFile("shadow.txt")
	if err!= nil{
		fmt.Println("shadow.txt not found")
		return false
	}
	lines:= strings.Split(string(data),"\n")
	for _, line := range lines {
		parts := strings.Split(line,":")
		if len(parts)==3 && parts[0] == username && parts[1] == hashedpwd{
			return true
		}
	}
	return false
}

func getClearance(username string) (int, error) {
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
	return 0, fmt.Errorf("clearance not found")
}

func mainMenu(username string, clearance int){
	for{
		fmt.Println("Options: (C)reate, (A)ppend, (R)ead, (W)rite, (L)ist, (S)ave or (E)xit.")
		fmt.Print("Enter your choice:")
		var choice string
		fmt.Scanln(&choice)
		switch strings.ToUpper(choice) {
		case "C":
			createFile(username, clearance)
		case "A":
			appendToFile(username, clearance)
		case "R":
			readFile(username, clearance)
		case "W":
			writeFile(username, clearance)
		case "L":
			listFiles()
		case "S":
			saveFiles()
		case "E":
			exitSystem()
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
// type File struct {
// 	Filename       string
// 	Owner          string
// 	Classification int
// 	Content        string
// }

// var files []File
func createFile(username string, clearance int){
	fmt.Println("Enter file name:")
	var filename string
	fmt.Scanln(&filename)

	for _,file:= range files{
		if file.Filename  == filename{
			fmt.Println("The file already exists.")
			return 
		}
	}
	files = append(files, File{Filename: filename, Owner:username, Classification:clearance, Content:""})
	fmt.Println("File was created successfully.")
}

func appendToFile(username string, clearance int){
	fmt.Println("Enter file name:")
	var filename string
	fmt.Scanln(&filename)

	for _,file:=range files{
		if file.Filename == filename{
			if clearance<=file.Classification{
				fmt.Println("You do not have enough clearance to append on the file.")
				return 
			}
			fmt.Println("Enter the text to append to file:")
			var content string
			fmt.Scanln(&content)

			file.Content += content
			fmt.Println("Append action was successful.")
			return
		}
	}
	fmt.Printf("Error: file %s cannot be found.",filename)
	return
}

// Main function to handle command-line arguments and run the appropriate process
func main() {
	if len(os.Args) > 1 && os.Args[1] == "-i" {
		// Account creation process
		var username, password, clearanceStr string

		fmt.Print("Enter username: ")
		fmt.Scanln(&username)
		fmt.Print("Enter password (minimum 8 characters): ")
		fmt.Scanln(&password)
		fmt.Print("User clearance (0 or 1 or 2 or 3):")
		fmt.Scanln(&clearanceStr)

		clearance, err := strconv.Atoi(clearanceStr)
		if err != nil && clearance != 0 && clearance != 1 && clearance != 2 && clearance != 3 {
			fmt.Println("Error: Invalid clearance level. Must be 0, 1, 2, or 3.")
			return
		}

		createUser(username, password, clearance)
	} else {
		// Login process
		var username, password string
		fmt.Print("Enter username: ")
		fmt.Scanln(&username)
		fmt.Print("Enter password: ")
		fmt.Scanln(&password)
		authenticateUser(username, password)

	}
}
