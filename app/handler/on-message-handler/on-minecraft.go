package onmessagehandlers

import (
	e "bott-the-pigeon/app/error"
	ecsutils "bott-the-pigeon/lib/aws/service/ecs"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/bwmarrin/discordgo"
)

func OnMinecraft(bot *discordgo.Session, msg *discordgo.MessageCreate) {
	cluster := os.Getenv("AWS_ECS_MC_CLUSTER_ARN")
	taskDef := os.Getenv("AWS_ECS_MC_TASK_DEF_ARN")
	ecsTags := true
	taskIn := &ecs.RunTaskInput{
		Cluster:              &cluster,
		TaskDefinition:       &taskDef,
		EnableECSManagedTags: &ecsTags,
	}
	if _, err := ecsutils.RunTask(taskIn); err != nil {
		err = fmt.Errorf("failed to start task: %v", err)
		e.ThrowBotError(bot, msg.ChannelID, err)
		return
	}
	res := genTaskRunSuccessMessage()
	bot.ChannelMessageSendEmbed(msg.ChannelID, res)
}

func genTaskRunSuccessMessage() *discordgo.MessageEmbed {
	msg := &discordgo.MessageEmbed{
		Title:       "Minecraft Server Started (Hopefully).",
		Description: fmt.Sprintf("I'll tentatively say it worked, lad, but it'll take a few minutes."),
		Color:       0x44DD00,
	}
	return msg
}
