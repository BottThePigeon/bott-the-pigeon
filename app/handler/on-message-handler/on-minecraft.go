package onmessagehandlers

import (
	e "bott-the-pigeon/app/error"
	lambdautils "bott-the-pigeon/lib/aws/service/lambda"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/bwmarrin/discordgo"
)

func OnMinecraft(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	functionName := os.Getenv("MC_SERVICE_LAUNCHER_LAMBDA")
	lambdaInvokeIn := &lambda.InvokeInput{
		FunctionName: &functionName,
	}
	if _, err := lambdautils.InvokeLambda(lambdaInvokeIn); err != nil {
		err = fmt.Errorf("failed to run Lambda: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}
	res := genTaskRunSuccessMessage()
	bot.ChannelMessageSendEmbed(msg.ChannelID, res)
}

func genTaskRunSuccessMessage() *discordgo.MessageEmbed {
	msg := &discordgo.MessageEmbed{
		Title:       "Support Group Minecraft server starting...",
		Description: fmt.Sprintf("The server is starting at %s, but it'll take a few minutes. The other pigeon will let you know when it's ready.", os.Getenv("MC_SERVER_DOMAIN_NAME")),
		Color:       0x44DD00,
	}
	return msg
}
