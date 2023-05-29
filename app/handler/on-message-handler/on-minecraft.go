package onmessagehandlers

import (
	c "bott-the-pigeon/app/common"
	e "bott-the-pigeon/app/error"
	lambdautils "bott-the-pigeon/lib/aws/service/lambda"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/bwmarrin/discordgo"
)

func OnMinecraft(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	serverType := strings.TrimSpace(strings.Replace(msg.Content, ">mc", "", 1))
	var clusterNameOrArn string
	var functionNameOrArn string
	var mcDomain string
	var isCrossAccount bool

	if serverType == "--tekkit" {
		clusterNameOrArn = os.Getenv("MC_CLUSTER_NAME")
		functionNameOrArn = os.Getenv("MC_SERVICE_LAUNCHER_LAMBDA")
		mcDomain = os.Getenv("MINECRAFT_DOMAIN")
		isCrossAccount = false
	} else {
		clusterNameOrArn = os.Getenv("MC_VANILLA_CLUSTER_ARN")
		functionNameOrArn = os.Getenv("MC_VANILLA_SERVICE_LAUNCHER_LAMBDA")
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
		res := genAlreadyActiveSuccessMessage(mcDomain, isCrossAccount)
		bot.ChannelMessageSendEmbed(msg.ChannelID, res)
		return
	}

	lambdaInvokeIn := &lambda.InvokeInput{
		FunctionName: &functionNameOrArn,
	}
	if _, err := lambdautils.InvokeLambda(lambdaInvokeIn); err != nil {
		err = fmt.Errorf("Lambda invocation failed: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}

	res := genStartupSuccessMessage(mcDomain, isCrossAccount)
	bot.ChannelMessageSendEmbed(msg.ChannelID, res)
}

func genStartupSuccessMessage(mcDomain string, isVanilla bool) *discordgo.MessageEmbed {
	var title string
	var description string

	if isVanilla {
		title = "Support Group Minecraft server starting..."
		description = fmt.Sprintf("The Support Group Minecraft server is starting at %s, but it'll take a few minutes. The other pigeon will let you know when it's ready.", mcDomain)
	} else {
		title = "Support Group Tekkit server starting..."
		description = fmt.Sprintf("The Support Group Tekkit server is starting at %s, but it'll take a few minutes. The other pigeon will let you know when it's ready.", mcDomain)
	}

	msg := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0x44DD00,
	}

	return msg
}

func genAlreadyActiveSuccessMessage(mcDomain string, isVanilla bool) *discordgo.MessageEmbed {
	var title string
	var description string

	if isVanilla {
		title = "Support Group Minecraft server already running!"
	} else {
		title = "Support Group Tekkit server already running!"
	}
	description = fmt.Sprintf("Steady on you bird brain, the server's already running at %s!", mcDomain)

	msg := &discordgo.MessageEmbed{
		Title:       title,
		Description: description,
		Color:       0x44DD00,
	}

	return msg
}
