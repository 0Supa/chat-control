package main

import (
	"flag"
	"fmt"
	"regexp"
	"runtime"
	"time"

	"github.com/gempir/go-twitch-irc/v4"
	"github.com/micmonay/keybd_event"
)

var controls = map[*regexp.Regexp]int{
	regexp.MustCompile(`(?i)^(up|jump|w)\b`):     keybd_event.VK_UP,
	regexp.MustCompile(`(?i)^(left|a)\b`):        keybd_event.VK_LEFT,
	regexp.MustCompile(`(?i)^(down|crouch|s)\b`): keybd_event.VK_DOWN,
	regexp.MustCompile(`(?i)^(right|d)\b`):       keybd_event.VK_RIGHT,
}

var pressed = map[int]int{}

func updateKeyState(kb keybd_event.KeyBonding) {
	kb.Clear()
	for key, t := range pressed {
		if t > 0 {
			kb.AddKey(key)
		}
	}
	// kb.Release()
	kb.Press()
	// time.Sleep(2 * time.Second)
	kb.Release()
}

var channelFlag = flag.String("channel", "daedulas_", "Twitch channel name")

func main() {
	fmt.Println("Initializing...")

	flag.Parse()

	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	if runtime.GOOS == "linux" {
		time.Sleep(2 * time.Second)
	}

	client := twitch.NewAnonymousClient()

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		for key, val := range controls {
			if key.MatchString(message.Message) {
				fmt.Println(message.User.DisplayName + ": " + message.Message)

				go func(key int) {
					pressed[key]++
					updateKeyState(kb)
					pressed[key]--
				}(val)
			}
		}
	})

	client.OnSelfJoinMessage(func(message twitch.UserJoinMessage) {
		fmt.Println("Joined #" + message.Channel)
	})

	client.Join(*channelFlag)

	err = client.Connect()
	if err != nil {
		panic(err)
	}
}
