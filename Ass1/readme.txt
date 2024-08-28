# Go Application for Authentication and Access Control System

## Description
This Go program is designed to act as a authentication and access control system. 

## Features
Authentication: The auth package handles user authentication, including verifying login credentials and 
managing user sessions. It retrieves user information such as usernames and hashed passwords from files
(like salt.txt and shadow.txt) and ensures that users are authenticated correctly before granting access 
to the system.

Hashing: The hash package is responsible for generating secure password hashes using the MD5 algorithm. 
It creates a unique salt for each user, combines it with the user's password, and generates a hash that 
is stored securely. This ensures that passwords are not stored in plain text and provides additional 
security through the use of salted hashes.

Models: The models package contains the core data structures and functions for managing the file system 
within the application. It includes operations for creating, reading, writing, appending, and listing 
files, with access control based on the user's security clearance level. The package also manages 
in-memory storage of files during a session and handles the saving of these files back to disk when 
necessary.

## Installation

## Steps 
1. Download the zip file containing the binary executive file named main.

2. Navigate to the directory containing the source code files.

3. To compile the code and build an executable file, navigate to the directory containing the source code
files and run the command:
    go build main.go
It is requried that the user has go version 1.22.2.
Go can be installed by running the following command:
    sudo apt install golang-go

**NOTE TO TUTOR**: CAPA may have blocked permissions to download and install go. For testing and marking,
the executable file has been built on local system and located in the same directory as the source code, 
named main. The marker is able to proceed with steps 4-5 for testing of the authentication system.

4. Input the command ./main -i to run the filesystem in sign-up mode. 

5. input the command ./main to run the filesystem in login mode.

## Description of the Reduction

### User Creation

During initialization (./FileSystem -i), the program prompts for a username, 
password, and security clearance level. If the user does not enter a minimum 
of 8 characters for password, the terminal will display an error and terminate
the program.

Once completed, the program calls the function generateSalt() to make a 8 digit
salt, then concatenates it with the user's password and passes it into the 
HashMD5 function in the hash package. The returned hashed password is then 
storred with the user's secruity clearance in shadow.txt. The username and salt 
value is stored in salt.txt.

The program terminates when the process is completed.

### User Login
During a standard run ./main , the program prompts for the username and password.The system retrieves 
the salt associated with the username from salt.txt and uses it to generate the MD5 hash of the entered 
password. This generated hash is compared with the stored hash in shadow.txt. If the values are equal, 
the user is authenticated and logged in.

If the user is not authenticated, the program will terminate with an error displayed.

### Main Menu Initialization
After a successful login, the mainMenu() function is called. This function does a test by 
calling the HashMD5 function in the hash package and passing in a test string.

The program then calls the LoadFileStore function in the models package. This function 
checks if a Files.store exists, if it does not exists(err != nil) then a new Files.store 
is created. If there is an error in creation of the store file, an error messaged will be 
displayed.

File information in the Files.store takes the following format:
[filename]:[username]:[clearance]:[content]

If the file already exists, then it will extract each file from the store, create a
File struct with the four fields, then append it to an array called files.

### Main Menu Functions

"C" - Create a File:
The program calls models.CreateFile(username, clearance).
This prompts the user for a filename.
It then checks if a file with the same name already exists in memory.
If it doesn’t exist, a new file is created with the user's username as the owner and the user’s clearance level as the file’s classification level.
The created file's clearance level will be the same as the user who created it.

"A" - Append to a File:
If the user chooses "A", the program calls models.AppendFile(username, clearance).
This prompts the user for a filename.
It then checks if the file exists in memory.
If the file exists, it checks if the user’s clearance level allows them to append to the file.
(user clearance level <= file classification level)
If allowed, it prompts the user for content and appends this content to the file.

"R" - Read a File:
If the user chooses "R", the program calls models.ReadFile(username, clearance).
This prompts the user for a filename.
It then checks if the file exists in memory.
If the file exists, it checks if the user’s clearance level allows them to read the file.
(user clearance level >= file classification level)
If allowed, it displays the file’s content.

"W" - Write to a File:
If the user chooses "W", the program calls models.WriteFile(username, clearance).
This prompts the user for a filename.
It then checks if the file exists in memory.
If the file exists, it checks if the user’s clearance level allows them to write to the file.
(user clearance level == file classification level)
If allowed, it prompts the user for new content and replaces the existing content with the new content.

"L" - List All Files:
If the user chooses "L", the program calls models.ListFiles().
This displays the names of all files currently stored in memory. If no files are present, it informs the 
user that the memory is empty.

"S" - Save Files to Disk:
If the user chooses "S", the program calls models.SaveFiles().
This writes the in-memory list of files back to the Files.store file, effectively saving any changes made
during the session.

"E" - Exit the System:
If the user chooses "E", the program calls models.ExitSystem().
This prompts the user to confirm if they really want to exit.
If the user confirms, the program exits; otherwise, it returns to the main menu.

