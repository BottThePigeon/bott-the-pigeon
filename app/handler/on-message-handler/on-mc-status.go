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
	serverType := strings.TrimSpace(strings.Replace(msg.Content, ">todo", "", 1))
	var clusterNameOrArn string
	var mcDomain string
	if serverType == "--vanilla" || serverType == "" {
		clusterNameOrArn = os.Getenv("MC_VANILLA_CLUSTER_ARN")
		mcDomain = os.Getenv("MINECRAFT_VANILLA_DOMAIN")
	} else if serverType == "--tekkit" {
		clusterNameOrArn = os.Getenv("MC_CLUSTER_NAME")
		mcDomain = os.Getenv("MINECRAFT_DOMAIN")
	}
	isActive, err := c.CheckMinecraftServerStatus(clusterNameOrArn, os.Getenv("MC_SERVICE_NAME"))
	if err != nil {
		err = fmt.Errorf("ECS describe services failed: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}
	if isActive {
		res := genActiveStatusMessage(mcDomain)
		bot.ChannelMessageSendEmbed(msg.ChannelID, res)
		return
	} else {
		res := genInactiveStatusMessage(mcDomain)
		bot.ChannelMessageSendEmbed(msg.ChannelID, res)
		return
	}
}

func genActiveStatusMessage(mcDomain string) *discordgo.MessageEmbed {
	msg := &discordgo.MessageEmbed{
		Title:       "Support Group Minecraft server online!",
		Description: fmt.Sprintf("The Support Group Minecraft server is currently online at %s!", mcDomain),
		Color:       0x44DD00,
	}
	return msg
}

func genInactiveStatusMessage(mcDomain string) *discordgo.MessageEmbed {
	msg := &discordgo.MessageEmbed{
		Title:       "Support Group Minecraft server currently offline.",
		Description: fmt.Sprintf("The Support Group Minecraft server at %s is currently offline. You can switch it on by using the `>mc` command, and specifying which type of server you want with `--vanilla` or `--tekkit`!", mcDomain),
		Color:       0xDD4400,
	}
	return msg
}
