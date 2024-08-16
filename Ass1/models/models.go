package models

import (
	"fmt"
	"os"
	"strings"
	"io/ioutil"
	"strconv"
)



type File struct {
	Filename       string
	Owner          string
	Classification int
	Content        string
}

var files []File

//Loading instance of FileStore, if does not exist, create a new one
func LoadFileStore() error{
	_, err := os.Stat("Files.store")
	if err != nil {
		if os.IsNotExist(err) {
			// File doesn't exist, create it
			fmt.Println("Files.store does not exist. Creating a new one.")
			file, err := os.Create("Files.store")
			if err != nil{
				fmt.Printf("Error creating file Files.Store,%v",err)
			}
			defer file.Close()

			fmt.Println("File Files.store created successfully.")
			return nil
		}else{
			return fmt.Errorf("Error checking Files.store: %v",err)
		}
	}
	// If the file exists, load its contents into memory
    data, err := ioutil.ReadFile("Files.store")
    if err != nil {
        return fmt.Errorf("error reading Files.store: %v", err)
    }

    lines := strings.Split(string(data), "\n")
    for _, line := range lines {
        parts := strings.SplitN(line, ":", 4) // Split into at most 4 parts: Filename, Owner, Classification, Content
        if len(parts) != 4 {
            continue // skip invalid lines
        }

        clearance, err := strconv.Atoi(parts[2])
        if err != nil {
            fmt.Printf("Invalid clearance level in file metadata: %v\n", err)
            continue
        }

        file := File{
            Filename:       parts[0],
            Owner:          parts[1],
            Classification: clearance,
            Content:        parts[3],
        }

        files = append(files, file)
    }

    fmt.Println("Files loaded into memory successfully.")
    return nil
}

//O_APPEND opens file in append mode, content is written to end of the file
//O_CREATE creates file if it does not exist


// func checkFileExists(filename string) bool {
// 	file, err := ioutil.ReadFile("Files.store")
// 	if err != nil{
// 		fmt.Printf("Error reading file %v",err)
// 	}

// 	lines:= strings.Split(string(file), "\n")
// 	for _,line := range lines{
// 		parts := strings.Split(line, ":")
// 		if len(parts)==3 && parts[0] == filename{
// 			return true
// 		}
// 	}
// 	return false
// }
func checkFileInMemory(filename string)bool{
	for _, file:= range files{
		if file.Filename == filename{
			return true
		}
	}
	return false
}

// func CreateFile(username string, clearance int){
// 	fmt.Println("Enter file name:")
// 	var filename string
// 	fmt.Scanln(&filename)

// 	if checkFileExists(filename)||checkFileInMemory(filename){
// 		fmt.Println("File already exists.")
// 		return
// 	}
// 	file, err := os.OpenFile("Files.store", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
// 	if err != nil{
// 		fmt.Printf("Error opening file %v",err)
// 	}
// 	defer file.Close()
// 	fileob:= File{Filename:filename, Owner:username, Classification:clearance, Content:""}
// 	files = append(files, fileob)
// 	fmt.Println("File created successfully.")
// }

func CreateFile(username string, clearance int) {
    fmt.Println("Enter file name:")
    var filename string
    fmt.Scanln(&filename)

    if checkFileInMemory(filename) {
        fmt.Println("File already exists.")
        return
    }

    // Create the File object and append it to the in-memory slice
    fileob := File{Filename: filename, Owner: username, Classification: clearance, Content: ""}
    files = append(files, fileob)

    fmt.Println("File created successfully in memory.")
}

func AppendFile(username string, clearance int) {
    fmt.Println("Enter file name:")
    var filename string
    fmt.Scanln(&filename)

    if !checkFileInMemory(filename) {
        fmt.Println("File does not exist.")
        return
    }

    for i, file := range files {
        if file.Filename == filename {
            if clearance > file.Classification {
                fmt.Println("You do not have enough clearance to append to the file.")
                return
            }
            fmt.Println("Enter the text to append to file:")
            var content string
            fmt.Scanln(&content)

            files[i].Content += content
            fmt.Println("Content appended successfully in memory.")
            return
        }
    }
    fmt.Printf("Error: file %s cannot be found.\n", filename)
}

// func ReadFile(username string, clearance int){
// 	fmt.Println("Enter file name to read:")
// 	var filename string
// 	fmt.Scanln(&filename)

// 	if !checkFileExists(filename) && !checkFileInMemory(filename){
// 		fmt.Println("File does not exist.")
// 		return
// 	}
// 	for _, file := range files{
// 		if file.Filename == filename{
// 			if clearance < file.Classification{
// 				fmt.Println("You do not have enough clearance to read the file.")
// 				return
// 			}
// 			fmt.Println(file.Content)
// 			return
// 		}
// 	}
// 	fmt.Printf("Error: file %s cannot be found.\n", filename)
// }

func ReadFile(username string, clearance int) {
    fmt.Println("Enter file name to read:")
    var filename string
    fmt.Scanln(&filename)

    if !checkFileInMemory(filename) {
        fmt.Println("File does not exist in memory.")
        return
    }

    for _, file := range files {
        if file.Filename == filename {
            if clearance < file.Classification {
                fmt.Println("You do not have enough clearance to read the file.")
                return
            }
            fmt.Printf("filename:%v,owner:%v,class:%v\nContent of the file '%s':\n%s\n",file.Filename,
			file.Owner,file.Classification, filename, file.Content,)
            return
        }
    }

    fmt.Printf("Error: file %s cannot be found in memory.\n", filename)
}

func WriteFile(username string, clearance int){
	fmt.Println("Enter file name to write:")
	var filename string
	fmt.Scanln(&filename)

	if !checkFileInMemory(filename){
		fmt.Println("File does not exist.")
		return
	}
	for i, file := range files{
		if file.Filename == filename{
			if clearance > file.Classification{
				fmt.Println("You do not have enough clearance to write to the file.")
				return
			}
			fmt.Println("Enter the text to write to file:")
			var content string
			fmt.Scanln(&content)

			files[i].Content = content
			fmt.Println("Content written successfully in memory.")
			return
		}
	}
	fmt.Printf("Error: file %s cannot be found in memory.\n", filename)
}

func ListFiles(){
	if len(files)==0{
		fmt.Println("No files in memory.")
	}
	for _, file := range files{
		fmt.Println(file.Filename)
	}
}

func SaveFiles() error {
	filestore,err:= os.OpenFile("Files.store", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil{
		return fmt.Errorf("error opening file %v",err)
	}
	defer filestore.Close()
	for _, file := range files{
		_, err = filestore.WriteString(file.Filename + ":" + file.Owner + ":" + strconv.Itoa(file.Classification) + ":" + file.Content + "\n")
		if err != nil{
			return fmt.Errorf("error writing to file %v",err)
		}
	}
	return nil
}

func ExitSystem(){
	fmt.Println("Shut down the FileSystem? (Y)es or (N)o")
	var input string
	fmt.Scanln(&input)
	if strings.ToUpper(input) == "Y"{
		os.Exit(0)
	}
}