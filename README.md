# Bott The Pigeon

Monorepo for the Discord Bot "Bott The Pigeon" (Or Scott the Pigeon). It is written entirely in Go, using the [Discord API](https://github.com/bwmarrin/discordgo), and it is cloud-native, using the [AWS SDK](https://github.com/aws/aws-sdk-go). Currently this is a super simple stack running on an ec2 t2.micro on-demand instance, with IAM access to SSM Parameter Store to get the bot credentials. It starts up and shuts down with the instance, and supports full CI/CD using CodePipeline (The very boilerplate-y .yml files will give you an idea). The plan is for this to become serverless in future, but most current developments are building core functions that will accelerate development later on.

## [TODO LIST](https://github.com/users/adad-mitch/projects/1)