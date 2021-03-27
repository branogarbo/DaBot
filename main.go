package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating discord session:", err)
		return
	}

	dg.Identify.Intents = discordgo.IntentsGuildVoiceStates

	// dg.AddHandler(handleEvent)
	dg.AddHandler(handleEvent)

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}

	fmt.Println("DABABY IS ONLINE, LES GOOOO. PRESS CTRL+C TO PUT HIM TO SLEEP")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

// func handleEvent(s *discordgo.Session, v *discordgo.VoiceConnection) {
// 	// if v.UserID == s.State.User.ID {
// 	// 	return
// 	// }

// 	fmt.Println(v.UserID)

// }

func handleEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		msgHead       string
		targetChannel string
		err           error
	)

	if m.Author.ID == s.State.User.ID || len(m.Content) < 4 {
		return
	}

	msgHead = m.Content[:4]

	if msgHead == "!db " {
		// if m.Content == msgHead {
		// 	err = errors.New("mood not provided")
		// } else {
		// 	moodArg = strings.ToLower(strings.TrimSpace(m.Content[4:]))

		// }

		// if err != nil {
		// 	errMsg = fmt.Sprintf("```Error: %v```", err)

		// 	fmt.Println(errMsg)
		// 	moodString = errMsg
		// }

		// s.ChannelMessageSend(m.ChannelID, moodString)

		targetChannel = strings.TrimSpace(m.Content[4:])

		vc, err := s.ChannelVoiceJoin(m.GuildID, targetChannel, false, false)
		if err != nil {
			fmt.Println("failed to join voice channel:", err)
			return
		}

		dgvoice.PlayAudioFile(vc, "./lesGooo.mp3", make(chan bool))

		vc.Close()
	}
}
