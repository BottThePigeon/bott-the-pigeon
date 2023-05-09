package onmessagehandlers

import (
	c "bott-the-pigeon/app/common"
	e "bott-the-pigeon/app/error"
	lambdautils "bott-the-pigeon/lib/aws/service/lambda"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/bwmarrin/discordgo"
)

func OnMinecraft(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	isActive, err := c.CheckMinecraftServerStatus()
	if err != nil {
		err = fmt.Errorf("ECS describe services failed: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}
	if isActive {
		res := genAlreadyActiveSuccessMessage()
		bot.ChannelMessageSendEmbed(msg.ChannelID, res)
		return
	}

	functionName := os.Getenv("MC_SERVICE_LAUNCHER_LAMBDA")
	lambdaInvokeIn := &lambda.InvokeInput{
		FunctionName: &functionName,
	}
	if _, err := lambdautils.InvokeLambda(lambdaInvokeIn); err != nil {
		err = fmt.Errorf("Lambda invocation failed: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}

	res := genStartupSuccessMessage()
	bot.ChannelMessageSendEmbed(msg.ChannelID, res)
}

func genStartupSuccessMessage() *discordgo.MessageEmbed {
	msg := &discordgo.MessageEmbed{
		Title:       "Support Group Minecraft server starting...",
		Description: fmt.Sprintf("The Support Group Minecraft server is starting at %s, but it'll take a few minutes. The other pigeon will let you know when it's ready.", os.Getenv("MINECRAFT_DOMAIN")),
		Color:       0x44DD00,
	}
	return msg
}

func genAlreadyActiveSuccessMessage() *discordgo.MessageEmbed {
	msg := &discordgo.MessageEmbed{
		Title:       "Support Group Minecraft server already running!",
		Description: fmt.Sprintf("Steady on you bird brain, the server's already running at %s!", os.Getenv("MINECRAFT_DOMAIN")),
		Color:       0x44DD00,
	}
	return msg
}
