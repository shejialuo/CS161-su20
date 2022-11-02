package client_test

// You MUST NOT change these default imports.  ANY additional imports may
// break the autograder and everyone will be sad.

import (
	// Some imports use an underscore to prevent the compiler from complaining
	// about unused imports.
	"encoding/hex"
	_ "encoding/hex"
	_ "errors"
	"math/rand"
	_ "strconv"
	_ "strings"
	"testing"
	"time"

	// A "dot" import is used here so that the functions in the ginko and gomega
	// modules can be used without an identifier. For example, Describe() and
	// Expect() instead of ginko.Describe() and gomega.Expect().
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	userlib "github.com/cs161-staff/project2-userlib"

	"github.com/cs161-staff/project2-starter-code/client"
)

func TestSetupAndExecution(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "User Registration and User Log Tests")
	RunSpecs(t, "Client Tests")
}

// ================================================
// Global Variables (feel free to add more!)
// ================================================
const defaultPassword = "password"
const emptyString = ""
const contentOne = "Bitcoin is Nick's favorite "
const contentTwo = "digital "
const contentThree = "cryptocurrency!"

// ================================================
// Describe(...) blocks help you organize your tests
// into functional categories. They can be nested into
// a tree-like structure.
// ================================================

var _ = Describe("User Registration Tests", func() {

	BeforeEach(func() {
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic User Registration Test", func() {
		Specify("Register two same username", func() {
			userlib.DebugMsg("Register a new user alice")
			alice, err := client.InitUser("alice", "Qq142536!")
			Expect(err).To(BeNil())
			Expect(alice).ToNot(BeNil())
			userlib.DebugMsg("Register a new user also named alice")
			aliceAnother, err := client.InitUser("alice", "Qq142536!")
			Expect(err).ToNot(BeNil())
			Expect(aliceAnother).To(BeNil())
		})
		Specify("Register two case-sensitive username", func() {
			userlib.DebugMsg("Register a new user alice")
			alice, err := client.InitUser("alice", "Qq142536!")
			Expect(err).To(BeNil())
			Expect(alice).ToNot(BeNil())
			userlib.DebugMsg("Register a new user Alice")
			caseSensitiveAlice, err := client.InitUser("Alice", "Qq142536!")
			Expect(err).To(BeNil())
			Expect(caseSensitiveAlice).ToNot(BeNil())
		})
		Specify("Register a new user with empty username", func() {
			userlib.DebugMsg("Register a new user with empty name")
			empty, err := client.InitUser("", "Qq142536!!!!!")
			Expect(err).ToNot(BeNil())
			Expect(empty).To(BeNil())
		})
		Specify("Register a new user with empty password", func() {
			userlib.DebugMsg("Register a new user alice with empty password")
			alice, err := client.InitUser("alice", "")
			Expect(err).To(BeNil())
			Expect(alice).ToNot(BeNil())
		})
	})
	Describe("Random User Registration Test", func() {
		Specify("Random Test", func() {
			rand.Seed(time.Now().UnixNano())
			isVisited := make(map[string]int)
			for i := 0; i < 10; i++ {
				b := make([]byte, 1)
				rand.Read(b)
				randUsername := hex.EncodeToString(b)
				userlib.DebugMsg("Register a new user %v", randUsername)
				user, err := client.InitUser(randUsername, randUsername)
				if _, ok := isVisited[randUsername]; !ok {
					Expect(err).To(BeNil())
					Expect(user).ToNot(BeNil())
					isVisited[randUsername] = 1
				} else {
					Expect(err).ToNot(BeNil())
					Expect(user).To(BeNil())
				}
			}
		})
	})
})

var _ = Describe("User Login Tests", func() {
	BeforeEach(func() {
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})
	Describe("Basic Login Tests", func() {
		Specify("Login with correct password", func() {
			client.InitUser("alice", "Qq142536!")
			aliceLaptop, err := client.GetUser("alice", "Qq142536!")
			Expect(err).To(BeNil())
			Expect(aliceLaptop).ToNot(BeNil())
		})
		Specify("Login with incorrect password", func() {
			client.InitUser("alice", "Qq142536!")
			aliceLaptop, err := client.GetUser("alice", "qq142536!")
			Expect(err).ToNot(BeNil())
			Expect(aliceLaptop).To(BeNil())
			aliceLaptop, err = client.GetUser("alice", "Qq142536")
			Expect(err).ToNot(BeNil())
			Expect(aliceLaptop).To(BeNil())
			aliceLaptop, err = client.GetUser("alice", "Qq142536!!")
			Expect(err).ToNot(BeNil())
			Expect(aliceLaptop).To(BeNil())
			aliceLaptop, err = client.GetUser("alice", "Qq142536!  ")
			Expect(err).ToNot(BeNil())
			Expect(aliceLaptop).To(BeNil())
		})
		Specify("Login with empty password", func() {
			client.InitUser("alice", "")
			aliceLaptop, err := client.GetUser("alice", "")
			Expect(err).To(BeNil())
			Expect(aliceLaptop).ToNot(BeNil())
		})
	})
	Describe("Attacker Tests", func() {
		Specify("Delete the DataStore entry", func() {
			client.InitUser("alice", "Qq142536!")
			hash := userlib.Hash([]byte("alice_login"))
			id, _ := uuid.FromBytes(hash[:16])
			userlib.DatastoreDelete(id)
			aliceLaptop, err := client.GetUser("alice", "Qq142536!")
			Expect(err).ToNot(BeNil())
			Expect(aliceLaptop).To(BeNil())
		})
		Specify("Simply change the content of the DataStore", func() {
			client.InitUser("alice", "Qq142536!")
			hash := userlib.Hash([]byte("alice_login"))
			id, _ := uuid.FromBytes(hash[:16])
			userlib.DatastoreSet(id, []byte("dasdaf123das"))
			aliceLaptop, err := client.GetUser("alice", "Qq142536!")
			Expect(err).ToNot(BeNil())
			Expect(aliceLaptop).To(BeNil())
		})
		Specify("Retry and replace attack to test integrity", func() {
			client.InitUser("alice", "Qq142536!")
			client.InitUser("attacker", "attacker")
			hashAttacker := userlib.Hash([]byte("attacker_login"))
			idAttacker, _ := uuid.FromBytes(hashAttacker[:16])
			attackerContent, _ := userlib.DatastoreGet(idAttacker)
			hash := userlib.Hash([]byte("alice_login"))
			id, _ := uuid.FromBytes(hash[:16])
			// If there is no integrity support, we could replace the
			// alice password and hack it.
			userlib.DatastoreSet(id, attackerContent)
			attacker, err := client.GetUser("alice", "attacker")
			Expect(err).ToNot(BeNil())
			Expect(attacker).To(BeNil())
		})
	})
})

var _ = Describe("Client Tests", func() {

	// A few user declarations that may be used for testing. Remember to initialize these before you
	// attempt to use them!
	var alice *client.User
	var bob *client.User
	var charles *client.User
	// var doris *client.User
	// var eve *client.User
	// var frank *client.User
	// var grace *client.User
	// var horace *client.User
	// var ira *client.User

	// These declarations may be useful for multi-session testing.
	var alicePhone *client.User
	var aliceLaptop *client.User
	var aliceDesktop *client.User

	var err error

	// A bunch of filenames that may be useful.
	aliceFile := "aliceFile.txt"
	bobFile := "bobFile.txt"
	charlesFile := "charlesFile.txt"
	// dorisFile := "dorisFile.txt"
	// eveFile := "eveFile.txt"
	// frankFile := "frankFile.txt"
	// graceFile := "graceFile.txt"
	// horaceFile := "horaceFile.txt"
	// iraFile := "iraFile.txt"

	BeforeEach(func() {
		// This runs before each test within this Describe block (including nested tests).
		// Here, we reset the state of Datastore and Keystore so that tests do not interfere with each other.
		// We also initialize
		userlib.DatastoreClear()
		userlib.KeystoreClear()
	})

	Describe("Basic Tests", func() {

		Specify("Basic Test: Testing InitUser/GetUser on a single user.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting user Alice.")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())
		})

		Specify("Basic Test: Testing Single User Store/Load/Append.", func() {
			userlib.DebugMsg("Initializing user Alice.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Storing file data: %s", contentOne)
			err = alice.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentTwo)
			err = alice.AppendToFile(aliceFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Appending file data: %s", contentThree)
			err = alice.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Loading file...")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Create/Accept Invite Functionality with multiple users and multiple instances.", func() {
			userlib.DebugMsg("Initializing users Alice (aliceDesktop) and Bob.")
			aliceDesktop, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Getting second instance of Alice - aliceLaptop")
			aliceLaptop, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop storing file %s with content: %s", aliceFile, contentOne)
			err = aliceDesktop.StoreFile(aliceFile, []byte(contentOne))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceLaptop creating invite for Bob.")
			invite, err := aliceLaptop.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob accepting invite from Alice under filename %s.", bobFile)
			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Bob appending to file %s, content: %s", bobFile, contentTwo)
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).To(BeNil())

			userlib.DebugMsg("aliceDesktop appending to file %s, content: %s", aliceFile, contentThree)
			err = aliceDesktop.AppendToFile(aliceFile, []byte(contentThree))
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that aliceDesktop sees expected file data.")
			data, err := aliceDesktop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that aliceLaptop sees expected file data.")
			data, err = aliceLaptop.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Checking that Bob sees expected file data.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))

			userlib.DebugMsg("Getting third instance of Alice - alicePhone.")
			alicePhone, err = client.GetUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that alicePhone sees Alice's changes.")
			data, err = alicePhone.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne + contentTwo + contentThree)))
		})

		Specify("Basic Test: Testing Revoke Functionality", func() {
			userlib.DebugMsg("Initializing users Alice, Bob, and Charlie.")
			alice, err = client.InitUser("alice", defaultPassword)
			Expect(err).To(BeNil())

			bob, err = client.InitUser("bob", defaultPassword)
			Expect(err).To(BeNil())

			charles, err = client.InitUser("charles", defaultPassword)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Alice storing file %s with content: %s", aliceFile, contentOne)
			alice.StoreFile(aliceFile, []byte(contentOne))

			userlib.DebugMsg("Alice creating invite for Bob for file %s, and Bob accepting invite under name %s.", aliceFile, bobFile)

			invite, err := alice.CreateInvitation(aliceFile, "bob")
			Expect(err).To(BeNil())

			err = bob.AcceptInvitation("alice", invite, bobFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err := alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob can load the file.")
			data, err = bob.LoadFile(bobFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Bob creating invite for Charles for file %s, and Charlie accepting invite under name %s.", bobFile, charlesFile)
			invite, err = bob.CreateInvitation(bobFile, "charles")
			Expect(err).To(BeNil())

			err = charles.AcceptInvitation("bob", invite, charlesFile)
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Charles can load the file.")
			data, err = charles.LoadFile(charlesFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Alice revoking Bob's access from %s.", aliceFile)
			err = alice.RevokeAccess(aliceFile, "bob")
			Expect(err).To(BeNil())

			userlib.DebugMsg("Checking that Alice can still load the file.")
			data, err = alice.LoadFile(aliceFile)
			Expect(err).To(BeNil())
			Expect(data).To(Equal([]byte(contentOne)))

			userlib.DebugMsg("Checking that Bob/Charles lost access to the file.")
			_, err = bob.LoadFile(bobFile)
			Expect(err).ToNot(BeNil())

			_, err = charles.LoadFile(charlesFile)
			Expect(err).ToNot(BeNil())

			userlib.DebugMsg("Checking that the revoked users cannot append to the file.")
			err = bob.AppendToFile(bobFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())

			err = charles.AppendToFile(charlesFile, []byte(contentTwo))
			Expect(err).ToNot(BeNil())
		})

	})
})
