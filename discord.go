package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

//STATUS is the Server Message
const STATUS = ""

type requestInfo struct {
	authorID  string
	queryType string
	message   string
}

func check(err error) (gotError bool) {
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		return true
	}
	return false
}

func setupPrivate(s *discordgo.Session, author *discordgo.User) (ret *discordgo.Channel) {
	channel, _ := s.UserChannelCreate(author.ID)
	s.State.ChannelAdd(channel)
	return channel
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	// Set the playing status.
	s.UpdateStatus(0, STATUS)
	fmt.Println("Logged in")
}

func sendMessage(s *discordgo.Session, m *discordgo.MessageCreate, message string, private bool, sendTo bool) {
	if private {
		if sendTo {
			for _, name := range m.Mentions {
				if s.State.User.ID == name.ID {
					continue
				}

				channel := setupPrivate(s, name)
				s.ChannelMessageSend(channel.ID, message)
			}
		} else {
			channel := setupPrivate(s, m.Author)
			s.ChannelMessageSend(channel.ID, message)
		}
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		check(err)

	} else {
		s.ChannelMessageSend(m.ChannelID, message)
	}
}

func sendEmbedMessage(s *discordgo.Session, m *discordgo.MessageCreate, embed *Embed, private bool, sendTo bool) {
	if private {
		if sendTo {
			for _, name := range m.Mentions {
				if s.State.User.ID == name.ID {
					continue
				}
				channel := setupPrivate(s, name)
				s.ChannelMessageSendEmbed(channel.ID, embed.MessageEmbed)
			}
		} else {
			channel := setupPrivate(s, m.Author)
			s.ChannelMessageSendEmbed(channel.ID, embed.MessageEmbed)
		}
		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		check(err)
	} else {
		s.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	mentioned := false
	if m.Author.ID == s.State.User.ID {
		return
	}

	if !mentioned {
		match, _ := regexp.MatchString("^!", m.Content)
		if match {
			mentioned = true
		}
	}
	var info requestInfo

	if mentioned {
		username := fmt.Sprintf("%v#%v", m.Author.Username, m.Author.Discriminator)
		info.authorID = username
		info.message = cleanUpMessage(m)
		reType := regexp.MustCompile(`help|calculate`)
		matchType := string(reType.Find([]byte(strings.ToLower(info.message))))

		if matchType == "help" {
			info.queryType = "help"
		} else if matchType == "calculate" {
			info.queryType = "calculate"
		}

		if info.queryType == "help" {
			menu := helpMenu()
			sendMessage(s, m, menu, true, false)
			sendMessage(s, m, "Help sent", false, false)
		} else if info.queryType == "calculate" {
			message, private := calculateDistance(info)
			sendEmbedMessage(s, m, message, private, false)
		}
	}
}

func checkPrivate(message string) (priv bool) {
	matched, _ := regexp.MatchString(`private`, strings.ToLower(message))
	if matched {
		return true
	}
	return false
}

func checkSendTo(message string) (send bool) {
	matched, _ := regexp.MatchString(`send`, strings.ToLower(message))
	if matched {
		return true
	}
	return false
}

func cleanUpMessage(m *discordgo.MessageCreate) (retmessage string) {
	retmessage = m.Content
	for _, name := range m.Mentions {
		mentionid := fmt.Sprintf("<?!?@%v..|<?@?!%v..", name.ID, name.ID)
		reName := regexp.MustCompile(mentionid)
		retmessage = string(reName.ReplaceAll([]byte(retmessage), []byte("")))
	}
	return retmessage
}

func parseEmbedMessage(info requestInfo) (ret *Embed, parsed bool, private bool) {
	embed := NewEmbed()
	matched, _ := regexp.MatchString(`hyperspeed|stats`, strings.ToLower(info.message))
	if matched {
		private := checkPrivate(info.message)
		//mes, parsed := doEntityQuery(info)
		mes := NewEmbed()
		return mes, parsed, private
	}
	return embed, false, false
}
