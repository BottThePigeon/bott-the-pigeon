package onmessagehandlers

import (
	c "bott-the-pigeon/app/common"
	e "bott-the-pigeon/app/error"
	"fmt"
	"os"

	"github.com/bwmarrin/discordgo"
)

func OnMinecraftStatus(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	isActive, err := c.CheckMinecraftServerStatus()
	if err != nil {
		err = fmt.Errorf("ECS describe services failed: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}
	if isActive {
		res := genActiveStatusMessage()
		bot.ChannelMessageSendEmbed(msg.ChannelID, res)
		return
	} else {
		res := genInactiveStatusMessage()
		bot.ChannelMessageSendEmbed(msg.ChannelID, res)
		return
	}
}

func genActiveStatusMessage() *discordgo.MessageEmbed {
	msg := &discordgo.MessageEmbed{
		Title:       "The Support Group Minecraft server is online!",
		Description: fmt.Sprintf("The Support Group Minecraft server is currently online at %s!", os.Getenv("MINECRAFT_DOMAIN")),
		Color:       0x44DD00,
	}
	return msg
}

func genInactiveStatusMessage() *discordgo.MessageEmbed {
	msg := &discordgo.MessageEmbed{
		Title:       "The Support Group Minecraft server is currently offline.",
		Description: fmt.Sprintf("The Support Group Minecraft server at %s is currently offline. You can switch it on by using the `>mc` command!", os.Getenv("MINECRAFT_DOMAIN")),
		Color:       0xDD4400,
	}
	return msg
}
