package client

// CS 161 Project 2

// You MUST NOT change these default imports. ANY additional imports
// may break the autograder!

import (
	"encoding/json"

	userlib "github.com/cs161-staff/project2-userlib"
	"github.com/google/uuid"

	// hex.EncodeToString(...) is useful for converting []byte to string

	// Useful for string manipulation
	"strings"

	// Useful for formatting strings (e.g. `fmt.Sprintf`).
	"fmt"

	// Useful for creating new error messages to return using errors.New("...")
	"errors"

	// Optional.
	_ "strconv"
)

// This serves two purposes: it shows you a few useful primitives,
// and suppresses warnings for imports not being used. It can be
// safely deleted!
func someUsefulThings() {

	// Creates a random UUID.
	randomUUID := uuid.New()

	// Prints the UUID as a string. %v prints the value in a default format.
	// See https://pkg.go.dev/fmt#hdr-Printing for all Golang format string flags.
	userlib.DebugMsg("Random UUID: %v", randomUUID.String())

	// Creates a UUID deterministically, from a sequence of bytes.
	hash := userlib.Hash([]byte("user-structs/alice"))
	deterministicUUID, err := uuid.FromBytes(hash[:16])
	if err != nil {
		// Normally, we would `return err` here. But, since this function doesn't return anything,
		// we can just panic to terminate execution. ALWAYS, ALWAYS, ALWAYS check for errors! Your
		// code should have hundreds of "if err != nil { return err }" statements by the end of this
		// project. You probably want to avoid using panic statements in your own code.
		panic(errors.New("An error occurred while generating a UUID: " + err.Error()))
	}
	userlib.DebugMsg("Deterministic UUID: %v", deterministicUUID.String())

	// Declares a Course struct type, creates an instance of it, and marshals it into JSON.
	type Course struct {
		name      string
		professor []byte
	}

	course := Course{"CS 161", []byte("Nicholas Weaver")}
	courseBytes, err := json.Marshal(course)
	if err != nil {
		panic(err)
	}

	userlib.DebugMsg("Struct: %v", course)
	userlib.DebugMsg("JSON Data: %v", courseBytes)

	// Generate a random private/public keypair.
	// The "_" indicates that we don't check for the error case here.
	var pk userlib.PKEEncKey
	var sk userlib.PKEDecKey
	pk, sk, _ = userlib.PKEKeyGen()
	userlib.DebugMsg("PKE Key Pair: (%v, %v)", pk, sk)

	// Here's an example of how to use HBKDF to generate a new key from an input key.
	// Tip: generate a new key everywhere you possibly can! It's easier to generate new keys on the fly
	// instead of trying to think about all of the ways a key reuse attack could be performed. It's also easier to
	// store one key and derive multiple keys from that one key, rather than
	originalKey := userlib.RandomBytes(16)
	derivedKey, err := userlib.HashKDF(originalKey, []byte("mac-key"))
	if err != nil {
		panic(err)
	}
	userlib.DebugMsg("Original Key: %v", originalKey)
	userlib.DebugMsg("Derived Key: %v", derivedKey)

	// A couple of tips on converting between string and []byte:
	// To convert from string to []byte, use []byte("some-string-here")
	// To convert from []byte to string for debugging, use fmt.Sprintf("hello world: %s", some_byte_arr).
	// To convert from []byte to string for use in a hashmap, use hex.EncodeToString(some_byte_arr).
	// When frequently converting between []byte and string, just marshal and unmarshal the data.
	//
	// Read more: https://go.dev/blog/strings

	// Here's an example of string interpolation!
	_ = fmt.Sprintf("%s_%d", "file", 1)
}

// This is the type definition for the User struct.
// A Go struct is like a Python or Java class - it can have attributes
// (e.g. like the Username attribute) and methods (e.g. like the StoreFile method below).
type User struct {
	Username string

	// You can add other attributes here if you want! But note that in order for attributes to
	// be included when this struct is serialized to/from JSON, they must be capitalized.
	// On the flipside, if you have an attribute that you want to be able to access from
	// this struct's methods, but you DON'T want that value to be included in the serialized value
	// of this struct that's stored in datastore, then you can use a "private" variable (e.g. one that
	// begins with a lowercase letter).
	SessionNumber uint
}

type Password struct {
	Salt           []byte
	HashedPassword []byte
	Signature      []byte
}

type GlobalState struct {
	Users map[string]*User
}

var (
	globalState = &GlobalState{
		Users: make(map[string]*User),
	}
)

// NOTE: The following methods have toy (insecure!) implementations.

// A auxiliary function to generate UUID from string
func generateUUIDFromString(str string) (uuid.UUID, error) {
	hash := userlib.Hash([]byte(str))
	id, err := uuid.FromBytes(hash[:16])
	if err != nil {
		return [16]byte{}, err
	}
	return id, nil
}

// Creates a new User struct and returns a pointer to it.
// Returns an error if:
// 1. A user with the same username exists
// 2. An empty username is provided
func InitUser(username string, password string) (*User, error) {
	// An empty username is provided
	if username == "" {
		return nil, fmt.Errorf("empty user name")
	}

	// Here, we register a public key which is named `<username>_login`
	// This key is to ensure the integrity of the hashedPassword.
	// If there exists such key which means the username is used by others
	// we should returns an error.
	usernameLogin := username + "_login"
	if _, ok := userlib.KeystoreGet(usernameLogin); ok {
		return nil, fmt.Errorf("a user with the same username exists")
	}

	// Generate UUID for DataStore
	usernameLoginUUID, err := generateUUIDFromString(usernameLogin)
	if err != nil {
		return nil, err
	}

	// Generate the random salt
	salt := userlib.RandomBytes(24)

	// Generate the hashed password
	saltWithPasswordByteHash := userlib.Hash(append([]byte(password), salt...))

	privateKey, publicKey, err := userlib.DSKeyGen()
	if err != nil {
		return nil, err
	}

	// Here, we should make the salt and hashed password do not change by the attacker
	// So we need use digital signatures
	signature, err := userlib.DSSign(privateKey, append(saltWithPasswordByteHash, salt...))
	if err != nil {
		return nil, err
	}

	// Set the publicKey in the KeyStore
	if err := userlib.KeystoreSet(usernameLogin, publicKey); err != nil {
		return nil, err
	}

	hashedPasswordInfo := &Password{
		Salt:           salt,
		HashedPassword: saltWithPasswordByteHash,
		Signature:      signature,
	}

	hashedPasswordInfoBytes, err := json.Marshal(hashedPasswordInfo)
	if err != nil {
		return nil, err
	}

	userlib.DatastoreSet(usernameLoginUUID, hashedPasswordInfoBytes)

	user := &User{
		Username:      username,
		SessionNumber: 1,
	}
	globalState.Users[username] = user

	return user, nil
}

// Obtains the User struct of a user who has already been initialized and returns
// a pointer to it.
// Returns an error if:
//  1. There is no initialized user for the given username.
//  2. The user credentials are invalid.
//  3. The User struct cannot be obtained due to malicious action, or the integrity
//     of the user struct has been compromised.
func GetUser(username string, password string) (*User, error) {
	usernameLogin := username + "_login"
	// There is no initialized user for the given username
	publicKey, ok := userlib.KeystoreGet(usernameLogin)
	if !ok {
		return nil, fmt.Errorf("there is no initialized user for the given username")
	}

	usernameLoginUUID, err := generateUUIDFromString(usernameLogin)
	if err != nil {
		return nil, err
	}

	hashedPasswordInfoBytes, ok := userlib.DatastoreGet(usernameLoginUUID)
	if !ok {
		return nil, fmt.Errorf("attacker may delete the UUID information")
	}

	var hashedPasswordInfo Password

	err = json.Unmarshal(hashedPasswordInfoBytes, &hashedPasswordInfo)
	if err != nil {
		return nil, err
	}

	saltWithHashedPassword := append(hashedPasswordInfo.HashedPassword, hashedPasswordInfo.Salt...)

	err = userlib.DSVerify(publicKey, saltWithHashedPassword, hashedPasswordInfo.Signature)
	if err != nil {
		return nil, fmt.Errorf("the integrity check failed")
	}

	saltWithPasswordByteHashProvide := userlib.Hash(append([]byte(password), hashedPasswordInfo.Salt...))
	if !userlib.HMACEqual(hashedPasswordInfo.HashedPassword, saltWithPasswordByteHashProvide) {
		return nil, fmt.Errorf("the password isn't correct")
	}

	globalState.Users[username].SessionNumber++

	return globalState.Users[username], nil
}

func (userdata *User) StoreFile(filename string, content []byte) (err error) {
	storageKey, err := uuid.FromBytes(userlib.Hash([]byte(filename + userdata.Username))[:16])
	if err != nil {
		return err
	}
	contentBytes, err := json.Marshal(content)
	if err != nil {
		return err
	}
	userlib.DatastoreSet(storageKey, contentBytes)
	return
}

func (userdata *User) AppendToFile(filename string, content []byte) error {
	return nil
}

func (userdata *User) LoadFile(filename string) (content []byte, err error) {
	storageKey, err := uuid.FromBytes(userlib.Hash([]byte(filename + userdata.Username))[:16])
	if err != nil {
		return nil, err
	}
	dataJSON, ok := userlib.DatastoreGet(storageKey)
	if !ok {
		return nil, errors.New(strings.ToTitle("file not found"))
	}
	err = json.Unmarshal(dataJSON, &content)
	return content, err
}

func (userdata *User) CreateInvitation(filename string, recipientUsername string) (
	invitationPtr uuid.UUID, err error) {
	return
}

func (userdata *User) AcceptInvitation(senderUsername string, invitationPtr uuid.UUID, filename string) error {
	return nil
}

func (userdata *User) RevokeAccess(filename string, recipientUsername string) error {
	return nil
}
