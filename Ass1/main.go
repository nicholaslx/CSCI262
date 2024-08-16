package main

import (
	"Ass1/auth"
	"Ass1/hash"
	"Ass1/models"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "-i" {
		// Account creation process
		var username, password, clearanceStr string

		fmt.Print("Enter username: ")
		fmt.Scanln(&username)
		fmt.Print("Enter password (minimum 8 characters): ")
		fmt.Scanln(&password)
		fmt.Print("User clearance (0 or 1 or 2 or 3): ")
		fmt.Scanln(&clearanceStr)

		clearance, err := strconv.Atoi(clearanceStr)
		if err != nil || (clearance != 0 && clearance != 1 && clearance != 2 && clearance != 3) {
			fmt.Println("Error: Invalid clearance level. Must be 0, 1, 2, or 3.")
			return
		}

		err = auth.CreateUser(username, password, clearance)
		if err != nil {
			fmt.Println(err)
		}
	} else {
		// Login process
		var username, password string
		fmt.Print("Enter username: ")
		fmt.Scanln(&username)
		fmt.Print("Enter password: ")
		fmt.Scanln(&password)

		clearance, err := auth.AuthenticateUser(username, password)
		if err != nil {
			fmt.Println(err)
			return
		}

		// After successful authentication, enter the main menu
		mainMenu(username, clearance)
	}
}

func mainMenu(username string, clearance int){
	fmt.Printf("MD5 (\"This is a test\") = %v\n",hash.HashMD5("This is a test"))
	err := models.LoadFileStore()
	if err != nil{
		fmt.Printf("Error loading File.store %v",err)
		return
	}
	for{
		fmt.Println("Options: (C)reate, (A)ppend, (R)ead, (W)rite, (L)ist, (S)ave or (E)xit.")
		fmt.Print("Enter your choice:")
		var choice string
		fmt.Scanln(&choice)
		switch strings.ToUpper(choice) {
		case "C":
			models.CreateFile(username, clearance)
		case "A":
			models.AppendFile(username, clearance)
		case "R":
			models.ReadFile(username, clearance)
		case "W":
			models.WriteFile(username, clearance)
		case "L":
			models.ListFiles()
		case "S":
			models.SaveFiles()
		case "E":
			models.ExitSystem()
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}