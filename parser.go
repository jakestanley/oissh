package main

import (
	"fmt"
	"io"

	"github.com/gliderlabs/ssh"
)

func Help(s ssh.Session) {
	io.WriteString(s, "You can use commands like\n\t'buy' to buy ships\n\t'embark' to get into one of your ships\n\t'help' to show this message again\n")
}

// TODO return BuyAction
func ParseBuyCommand(s ssh.Session, user *User, commands []string) {

	if len(commands) == 1 {
		io.WriteString(s, "Buy what?")
	} else {
		noun := commands[1]
		if noun == "ship" {
			if len(commands) >= 3 {
				name := commands[2]
				ship := createShip(user, name)
				Ships[ship.id] = ship
				user.ships = append(user.ships, ship)
				// TODO is name allowed function, i.e not empty, valid characters, unique amongst all users ships
			} else {
				io.WriteString(s, "You need to give your ship a name")
			}

		}
	}
}

// TODO return EmbarkAction
func ParseEmbarkCommand(s ssh.Session, user *User, commands []string) {
	if len(commands) == 1 {
		io.WriteString(s, "What would you like to board?")
	} else {
		success := false
		noun := commands[1]

		for _, ship := range user.ships {
			if ship.locationId == user.locationId && ship.name == noun {
				// TODO board
				user.locationType = SHIP
				user.shipId = ship.id
				success = true
			}
		}
		if !success {
			io.WriteString(s, fmt.Sprintf("You don't have any ships at this location with the name %s", noun))
		}
	}
}

func ListShips(s ssh.Session, user *User) {
	io.WriteString(s, fmt.Sprintf("You have %d ships", len(user.ships)))
	for i, ship := range user.ships {
		shipstr := fmt.Sprintf("\n\t%d: '%s', docked at '%s'", i+1, ship.name, Stations[ship.locationId].name)
		if user.shipId == ship.id && user.locationType == SHIP {
			shipstr = fmt.Sprintf("%s (you are on this ship)", shipstr)
		}
		io.WriteString(s, shipstr)
	}
}
