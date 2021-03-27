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

	dg.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

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

func handleEvent(s *discordgo.Session, m *discordgo.MessageCreate) {
	var (
		msgHead         string
		targetChannel   *discordgo.Channel
		targetChannelID string
		cmdString       string
		err             error
		errMsg          string
		vc              *discordgo.VoiceConnection
	)

	cmdString = "!db"

	if m.Author.ID == s.State.User.ID || len(m.Content) < len(cmdString) {
		return
	}

	msgHead = m.Content[:len(cmdString)]

	if m.Content == cmdString {
		errMsg = "Error: no channel ID provided"

		fmt.Println(errMsg)
		s.ChannelMessageSend(m.ChannelID, "```"+errMsg+"```")
		return
	}

	if msgHead == cmdString+" " {
		targetChannelID = strings.TrimSpace(m.Content[len(cmdString)+1:])

		targetChannel, err = s.Channel(targetChannelID)
		if err != nil {
			err = fmt.Errorf("channel with ID %v does not exist", targetChannelID)
		} else if targetChannel.Type != discordgo.ChannelTypeGuildVoice {
			err = fmt.Errorf("channel with ID %v is not a voice channel", targetChannelID)
		} else {
			vc, err = s.ChannelVoiceJoin(m.GuildID, targetChannelID, false, false)

			dgvoice.PlayAudioFile(vc, "./lesGooo.mp3", make(chan bool))

			vc.Disconnect()
		}
	}

	if err != nil {
		errMsg = fmt.Sprintf("Error: %v", err)

		fmt.Println(errMsg)
		s.ChannelMessageSend(m.ChannelID, "```"+errMsg+"```")
	}
}
