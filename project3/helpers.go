// Helper functions
package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	log "github.com/sirupsen/logrus"
	"os"

	"golang.org/x/crypto/argon2"
)

// Return true if the given path exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	} else if err != nil {
		log.Error(err)
		return false
	}
	return true
}

// Return a random string of the specified size
func randomByteString(numBytes int) (random string, err error) {
	randomBytes := make([]byte, numBytes)
	_, err = rand.Read(randomBytes)
	random = hex.EncodeToString(randomBytes)
	return random, err
}

// Hash the password with the provided salt, using reasonable defaults
func hashPassword(password string, salt string) (hashedPassword string) {
	hashedPasswordBytes := argon2.IDKey([]byte(password), []byte(salt), 1, 64*1024, 1, 32)
	hashedPassword = base64.StdEncoding.EncodeToString(hashedPasswordBytes)
	return
}

// Called to fully reset the state of the application
func resetState() {
	dropTables()
	createTables()

	// Delete temporary files, in case there are some hanging around
	os.RemoveAll(filePath)
	os.Mkdir(filePath, os.ModePerm)
	os.RemoveAll(filePath)
	os.Mkdir(filePath, os.ModePerm)
}
