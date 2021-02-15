package main

import (
	"log"

	"golang.org/x/crypto/bcrypt"
)

func insertUser(fullname, email string, hashedPassword []byte) AuthResponse {
	insert, err := db.Query("INSERT INTO user_tab VALUES (?, ?, ?)", fullname, email, string(hashedPassword))

	var response AuthResponse
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	if err != nil {
		response = AuthResponse{
			HTTPCode: "200",
			Message:  "Successfully registered user ",
		}
	}

	return response
}

func registerUser(userinfo UserInfo) AuthResponse {
	hash, err := bcrypt.GenerateFromPassword([]byte(userinfo.Password), 10)

	if err != nil {
		log.Printf("Failed to generate password hash for user: %s, error: %s", userinfo.Email, err)
		return AuthResponse{
			HTTPCode: "500",
			Message:  "Server error: Unable to register user",
		}
	}

	return insertUser(userinfo.Fullname, userinfo.Email, hash)
}

func verify(userinfo UserInfo) bool {
	var expectedPasswordHash []byte
	err := db.QueryRow("SELECT password_hashed FROM user_tab WHERE email = ?").Scan(&expectedPasswordHash)
	if err != nil {
		log.Fatalf("Failed to retrieve password hash from table. Error: %s\n", err)
	}

	err = bcrypt.CompareHashAndPassword(expectedPasswordHash, []byte(userinfo.Password))
	return err == nil
}

func loginUser(userinfo UserInfo) AuthResponse {
	var response AuthResponse
	if verify(userinfo) {
		response = AuthResponse{
			HTTPCode: "200",
		}
	} else {
		log.Printf("Unsuccessful login attempt from %s\n", userinfo.Email)
		response = AuthResponse{
			HTTPCode: "200",
			Message:  "Authentication failed: Incorrect ID or Password",
		}
	}

	return response
}
