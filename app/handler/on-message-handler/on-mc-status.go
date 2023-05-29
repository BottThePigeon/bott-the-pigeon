package onmessagehandlers

import (
	c "bott-the-pigeon/app/common"
	e "bott-the-pigeon/app/error"
	"fmt"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func OnMinecraftStatus(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	serverType := strings.TrimSpace(strings.Replace(msg.Content, ">mc-status", "", 1))
	var clusterNameOrArn string
	var mcDomain string
	var isCrossAccount bool

	if serverType == "--tekkit" {
		clusterNameOrArn = os.Getenv("MC_CLUSTER_NAME")
		mcDomain = os.Getenv("MINECRAFT_DOMAIN")
		isCrossAccount = false
	} else {
		clusterNameOrArn = os.Getenv("MC_VANILLA_CLUSTER_ARN")
		mcDomain = os.Getenv("MINECRAFT_VANILLA_DOMAIN")
		isCrossAccount = true
	}

	isActive, err := c.CheckMinecraftServerStatus(clusterNameOrArn, os.Getenv("MC_SERVICE_NAME"), isCrossAccount)
	if err != nil {
		err = fmt.Errorf("ECS describe services failed: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}
	if isActive {
		res := genActiveStatusMessage(mcDomain, isCrossAccount)
		bot.ChannelMessageSendEmbed(msg.ChannelID, res)
		return
	} else {
		res := genInactiveStatusMessage(mcDomain, isCrossAccount)
		bot.ChannelMessageSendEmbed(msg.ChannelID, res)
		return
	}
}

func genActiveStatusMessage(mcDomain string, isVanilla bool) *discordgo.MessageEmbed {
	var title string
	var description string

	if isVanilla {
		title = "Support Group Minecraft server online!"
		description = fmt.Sprintf("The Support Group Minecraft server is currently online at %s!", mcDomain)
	} else {
		title = "Support Group Tekkit server online!"
		description = fmt.Sprintf("The Support Group Tekkit server is currently online at %s!", mcDomain)
	}

	msg := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0x44DD00,
	}

	return msg
}

func genInactiveStatusMessage(mcDomain string, isVanilla bool) *discordgo.MessageEmbed {
	var title string
	var description string

	if isVanilla {
		title = "Support Group Minecraft server currently offline."
		description = fmt.Sprintf("The Support Group Minecraft server at %s is currently offline. You can switch it on by using the `>mc` command!", mcDomain)
	} else {
		title = "Support Group Tekkit server currently offline."
		description = fmt.Sprintf("The Support Group Tekkit server at %s is currently offline. You can switch it on by using the `>mc --tekkit` command!", mcDomain)
	}

	msg := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0xDD4400,
	}

	return msg
}
