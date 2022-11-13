package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type player struct {
	place            room
	isBackpackPacked bool
	isReady          bool
	inventory        []string
}

func (p *player) lookAround() string {
	ans := p.place.surroundings
	for _, places := range p.place.items {
		if len(places.things) != 0 {
			ans += places.where + ": "
			for _, item := range places.things {
				ans += item + ", "
			}
		}
	}
	if ans == p.place.surroundings {
		ans += "пустая комната. "
	}
	if p.place.name == "кухня" {
		if p.isReady {
			ans += "надо идти в универ. "
		} else {
			ans += "надо собрать рюкзак и идти в универ. "
		}
	} else {
		ans = strings.TrimRight(ans, "., ")
		ans += ". "
	}
	ans += "можно пройти - "
	for _, r := range p.place.near {
		ans += r.name + ", "
	}
	ans = strings.TrimRight(ans, ", ")
	ans += "."
	return ans
}

func (p *player) moveTo(location string) string {
	for _, r := range p.place.near {
		if r.name == location && r.door == false {
			p.place = *r
			return r.firstExpression
		}
		if r.name == location && r.door == true {
			return "дверь закрыта"
		}
	}
	return fmt.Sprintf("нет пути в %s", location)
}

func (p *player) equip(gear string) string {
	if gear == "рюкзак" {
		for idx1, obj := range p.place.items {
			for idx2, t := range obj.things {
				if t == gear {
					p.isBackpackPacked = true
					gamer.place.items[idx1].deleteObject(idx2)
					return fmt.Sprintf("Вы надели: %s", gear)
				}
			}
		}
	}
	return "нет такого"
}

func (p *player) take(item string) string {
	if !p.isBackpackPacked {
		return "некуда класть"
	}
	for _, i := range worldItems {
		if item == i {
			for idx1, obj := range p.place.items {
				for idx2, thing := range obj.things {
					if thing == item {
						p.inventory = append(p.inventory, item)
						if item == "конспекты" {
							p.isReady = true
						}
						gamer.place.items[idx1].deleteObject(idx2)
						return fmt.Sprintf("предмет добавлен в инвентарь: %s", item)
					}
				}
			}
		}
	}
	return "нет такого"
}

func (p *player) use(item, to string) string {
	flag := false
	for _, inv := range p.inventory {
		if inv == item {
			flag = true
			break
		}
	}
	if flag == false {
		return fmt.Sprintf("нет предмета в инвентаре - %s", item)
	}
	switch item {
	case "ключи":
		switch to {
		case "дверь":
			for i, rms := range p.place.near {
				if rms.door == true {
					p.place.near[i].door = false
					return "дверь открыта"
				}
			}
		}
	}
	return "не к чему применить"
}

type room struct {
	name            string
	firstExpression string
	surroundings    string
	items           []objects
	near            []*room
	door            bool
}

type objects struct {
	where  string
	things []string
}

func (o *objects) deleteObject(idx int) {
	o.things = append(o.things[0:idx], o.things[idx+1:]...)
}

var (
	kitchen    = room{}
	corridor   = room{}
	myRoom     = room{}
	outside    = room{}
	gamer      = player{}
	equipment  = []string{"рюкзак"}
	worldItems = []string{"конспекты", "ключи", "чай"}
)

func main() {
	var command string
	initGame()
	for command != "идти улица" {
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		fmt.Println(handleCommand(scanner.Text()))
	}
}

func initGame() {
	kitchen = room{
		name:            "кухня",
		firstExpression: "кухня, ничего интересного. можно пройти - коридор",
		surroundings:    "ты находишься на кухне, ",
		items: []objects{
			{where: "на столе", things: []string{"чай"}},
		},
		near: []*room{&corridor},
		door: false,
	}
	corridor = room{
		name:            "коридор",
		firstExpression: "ничего интересного. можно пройти - кухня, комната, улица",
		surroundings:    "ничего интересного",
		items:           []objects{},
		near:            []*room{&kitchen, &myRoom, &outside},
		door:            false,
	}
	myRoom = room{
		name:            "комната",
		firstExpression: "ты в своей комнате. можно пройти - коридор",
		surroundings:    fmt.Sprintf(""),
		items: []objects{
			{where: "на столе", things: []string{"ключи", "конспекты"}},
			{where: "на стуле", things: []string{"рюкзак"}},
		},
		near: []*room{&corridor},
		door: false,
	}
	outside = room{
		name:            "улица",
		firstExpression: "на улице весна. можно пройти - домой",
		surroundings:    "красота",
		items:           []objects{},
		near:            []*room{&corridor},
		door:            true,
	}
	gamer = player{
		place:            kitchen,
		isBackpackPacked: !true,
		isReady:          false,
		inventory:        []string{},
	}

}

func handleCommand(command string) string {
	com := strings.Split(command, " ")
	switch com[0] {
	case "осмотреться":
		return gamer.lookAround()
	case "идти":
		return gamer.moveTo(com[1])
	case "надеть":
		return gamer.equip(com[1])
	case "взять":
		return gamer.take(com[1])
	case "применить":
		return gamer.use(com[1], com[2])
	default:
		return "not implemented"
	}
}
