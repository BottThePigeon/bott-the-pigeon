# Bott The Pigeon

Monorepo for the Discord Bot "Bott The Pigeon" (Or Scott the Pigeon). It is written entirely in Golang, using the [Discord API (For Go!!)](https://github.com/bwmarrin/discordgo), and is cloud-native, using the [AWS SDK (For Go!!)](https://github.com/aws/aws-sdk-go) wherever possible. Currently this is a super simple stack running on an ec2 t2.micro on-demand instance, with IAM access to SSM Parameter Store to get the bot credentials. It starts up and shuts down with the instance, and supports full CI/CD using CodePipeline (The very boilerplate-y .yml files will give you an idea). Oh - feature-wise, it just spits out a message when you tag it, but obviously that will change.

## Usage Info:
 - Just do a `go build` or a `go run main.go`, if you have the AWS access keys (or you use your own and configure similar resources), it'll work.
 - Obviously this uses AWS programmatic access (And it uses it to get the Discord token, too), so you'll need access keys locally if you want anything to work. You can always hook this up to your own AWS account, with your own bots - it'll work exactly the same (provided you either change the Parameter path of /btp/, or use it). But if you need access keys to Scott's resources, contact me.
 - The application is run in test mode by default. If you want to run using the production bot, use the `--prod` flag.
    - "Test mode" refers to using the replica application specifically built for feature testing in development.

## Branch Structure:
 'release' serves as the default branch - what is typically 'main' or 'master'. This is arguably an unusual pattern, but the default branch should, generally, be kept in a buildable state. By extension, it makes sense for the branch that is most likely to be buildable to be the branch that is actually built using CI. 
 'dev' is the main branch for development. Feature branches may be made, but the intention is for this project to be as simple as possible, in line with how Go operates. Using just these two branches in general makes understanding the repository much easier.

## Project Hierarchy:
 This is the logical hierarchy, according to how Go sees it. The directories within the repository should imitate this - perhaps with the exception of the main package, in the root. Each of these should be fairly self-explanatory - capitalised are the modules that should be used as entrypoints to the application. Obviously, this is very much subject to change as feature additions and therefore architectural considerations emerge.
 - Module: Bott-The-Pigeon
   - Package: MAIN
   - Package: AWS-Utils
      - Package: Init
      - Package: AWSEnv
   - Package: Bot-Utils
      - Package: Init
      - Package: Handlers
   - Package: TESTS

## TODOS (Future Features):
 - Random set of answers rather than a fixed one when tagging the bot.
    - These could potentially be stored in a database, or Parameter Store, not sure.
 - Random image of a pigeon when sending "!pigeon".
    - These could be fetched from an S3 bucket. It seems to be kinda slow, so maybe use CloudFront if needed.
 - More to be discussed!

## TODOS (Technical):
 - Add tests (This should be test-driven, but that's kind of difficult with the way the Discord API operates.)
    - An odd, but seemingly viable solution would be to use 3 bots:
        - A production bot, live and hosted: Obviously the one that should be seen 24/7 by end users.
        - A development bot, for testing purposes by the developer. 
        - A "test controller" bot that simulates user activity, such that we can gather the response from the development bot.
    - This is still very much a matter of much chin-stroking, since that seems a little OTT.
 - Furthermore, automate those tests. Will involve some consideration with how to access the test bots to run the build.
 - Delete these TODO lists and actually get them up on a Kanban board or something Adam you lazy sod
