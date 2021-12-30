package main

import (
	"fmt"
	"io"
	"log"
	"strings"
	"sync/atomic"

	"github.com/gliderlabs/ssh"
	"github.com/google/uuid"

	ui "github.com/gizak/termui/v3"
)

type LocationCategory int

const (
	STATION LocationCategory = 0
	SHIP    LocationCategory = 1
)

type User struct {
	username     string
	credits      int64
	ships        []*Ship
	locationType LocationCategory
	locationId   int64
	shipId       uuid.UUID
}

type Station struct {
	owner *User
	name  string
	// capacity int
}

type Ship struct {
	id         uuid.UUID
	owner      *User
	name       string
	x          float64
	y          float64
	docked     bool
	locationId int64
}

// type GameSession struct {}

var (
	kill     int64
	Keys     []ssh.PublicKey
	Users    map[ssh.PublicKey]*User
	Stations []*Station
	Ships    map[uuid.UUID]*Ship
)

func loadUser(s ssh.Session) (*User, bool) {

	writeSystemString(fmt.Sprintf("Users: %d", len(Users)))

	for k, v := range Users {
		if ssh.KeysEqual(k, s.PublicKey()) {
			return v, true
		}
	}

	return nil, false
}

func createUser(s ssh.Session) *User {

	parsedKey, _ := ssh.ParsePublicKey(s.PublicKey().Marshal())

	Users[parsedKey] = &User{
		username:     s.User(),
		credits:      500,
		ships:        make([]*Ship, 0),
		locationType: STATION,
		locationId:   0,
	}

	u, _ := Users[parsedKey]

	return u
}

func createShip(owner *User, name string) *Ship {

	// TODO ensure unique names
	// TODO ensure non-empty names and only sensible characters
	return &Ship{
		id:         uuid.New(),
		owner:      owner,
		x:          0,
		y:          0,
		name:       name,
		docked:     true,
		locationId: 0,
	}
}

func WriteUserInfo(s ssh.Session, user *User) {

	if user.locationType == STATION {
		io.WriteString(s, "\nYou are on a station\n")
	} else if user.locationType == SHIP {
		io.WriteString(s, fmt.Sprintf("\nYou are on the good ship %s\n", Ships[user.shipId].name))
	}
}

func main() {

	// initialise variables and state
	atomic.StoreInt64(&kill, 0)
	Ships = make(map[uuid.UUID]*Ship)
	Users = make(map[ssh.PublicKey]*User)
	Stations = make([]*Station, 0)
	Stations = append(Stations, DefaultStation())

	ssh.Handle(func(s ssh.Session) {

		user, exists := loadUser(s)

		if !exists {
			user = createUser(s)
			// TODO save user to file or database
			writeUsersString(fmt.Sprintf("New user '%s' connected!", user.username))
		} else {
			// TODO load user from file or database
			writeUsersString(fmt.Sprintf("User '%s' connected. Welcome back!\n", user.username))
		}

		WriteUserInfo(s, user)

		fmt.Println(user.username)

		username := s.User()

		s.PublicKey()

		if exists {
			io.WriteString(s, fmt.Sprintf("Hello %s. Welcome back to OiSSH. You have %d credits\n", username, user.credits))
		} else {
			io.WriteString(s, fmt.Sprintf("Hello %s. Welcome to OiSSH. You have %d credits\n", username, user.credits))
		}

		connected := true

		inputBuffer := ""

		for connected {

			// not very useful for non-interactive input
			data := make([]byte, 1000)
			s.Read(data)

			if data[0] != 0 {

				dataStr := strings.TrimSuffix(strings.SplitAfter(string(data), "\n")[0], "\n")
				fmt.Println(dataStr)

				// write back what the user inputted so the session is interactive
				s.Write(data)
				input := data[0]

				// if user presses enter, process the input and clear the input buffer
				if len(dataStr) > 1 {
					if dataStr == "shutdown" {
						connected = false
						defer atomic.StoreInt64(&kill, 1)
					}
				} else if input == 13 && len(inputBuffer) > 0 {

					if inputBuffer == "disconnect" {
						connected = false
					} else if inputBuffer == "shutdown" {
						connected = false
						defer atomic.StoreInt64(&kill, 1)
					} else if inputBuffer == "help" {
						Help(s)
					} else if inputBuffer == "ships" {
						ListShips(s, user)
					} else {
						commands := strings.Fields(inputBuffer)
						// for _, v := range commands {
						// 	fmt.Println(v)
						// }
						if len(commands) == 0 {

						} else if commands[0] == "buy" {
							ParseBuyCommand(s, user, commands)

						} else if commands[0] == "embark" {
							ParseEmbarkCommand(s, user, commands)
							WriteUserInfo(s, user)
						}
					}
					io.WriteString(s, "\n")
					inputBuffer = ""
				} else {
					inputBuffer += string(rune(input))
					// // if user inputs 'd', disconnect the user
					// if rune(input) == 'd' {
					// 	connected = false
					// }
					fmt.Println(inputBuffer)
				}
			}
		}

		io.WriteString(s, "Disconnecting you\n")
		writeUsersString(fmt.Sprintf("User '%s' disconnected. Goodbye!", user.username))
		s.Close()
	})

	// TODO implement this? although should only be used in emergencies
	passwordHandler := ssh.PasswordHandler(func(ctx ssh.Context, password string) bool {
		return false
	})

	EnableUi = false
	if EnableUi {
		initUi()
		defer ui.Close()

		go inputUi()
		go renderUI()
	} else {
		// TODO another way to handle the Ctrl+C interrupt
	}
	go game()

	passwordOption := ssh.PasswordAuth(passwordHandler)

	publicKeyOption := ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		// TODO error if key not valid
		return true
	})

	writeSystemString("Starting server...")
	log.Fatal(ssh.ListenAndServe(":2223", nil, passwordOption, publicKeyOption, ssh.HostKeyFile("/Users/jake/.ssh/oissh")))
}
