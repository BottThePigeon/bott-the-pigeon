# Bott The Pigeon

Monorepo for the Discord Bot "Bott The Pigeon" (Or Scott the Pigeon). It is written entirely in Golang, using the [Discord API](https://github.com/bwmarrin/discordgo), and is cloud-native, using the [AWS SDK](https://github.com/aws/aws-sdk-go) wherever possible. Currently this is a super simple stack running on an ec2 t2.micro on-demand instance, with IAM access to SSM Parameter Store to get the bot credentials. It starts up and shuts down with the instance, and supports full CI/CD using CodePipeline (The very boilerplate-y .yml files will give you an idea).

## Usage Info:
 - Just do a `go build` or a `go run main.go`, if you have the AWS access keys (or you use your own and configure similar resources), it'll work.
 - Obviously this uses AWS programmatic access (And it uses it to get the Discord token, too), so you'll need access keys locally if you want anything to work. You can always hook this up to your own AWS account, with your own bots - it'll work exactly the same (provided you either change the Parameter path of /btp/, or use it). But if you need access keys to Scott's resources, contact me.
 - The application is run in test mode by default. If you want to run using the production bot, use the `--prod` flag.
    - "Test mode" refers to using the replica application specifically built for feature testing in development.

## Branch Structure:
 'release' serves as the default branch - what is typically 'main' or 'master'. This is arguably an unusual pattern, but the default branch should, generally, be kept in a buildable state. By extension, it makes sense for the branch that is most likely to be buildable to be the branch that is actually built using CI. 
 'dev' is the main branch for development. There are feature branches, but generally the intention is for this project to be as simple as possible, in line with how Go operates. Using just these two branches in general makes understanding the repository much easier.

## Project Hierarchy:
 This is the logical hierarchy, according to how Go sees it. The directories within the repository should imitate this - perhaps with the exception of the main package, in the root. Each of these should be fairly self-explanatory - capitalised are the modules that should be used as entrypoints to the application. Obviously, this is very much subject to change as feature additions and therefore architectural considerations emerge.
 - Module: Bott-The-Pigeon
   - Package: MAIN
   - Package: AWS-Utils (This is completely generic AWS stuff, like session init. Things like getting something from a specific ARN on a particular AWS service should be written within the context of the logic it is used in - probably the Bot-Utils handlers. Otherwise 90% of the application would be inside AWS-Utils, which doesn't make as much sense as it all being within the context of the bot itself.)
      - Package: Session
      - Package: AWSEnv
   - Package: Bot-Utils (Most stuff should go here - any specific features that the bot offers, basically, which is most of the application.)
      - Package: Init
      - Package: Handlers
         - Package: On-Message-Handlers (Since many different kinds of events can occur based on a message being sent.)
   - Package: TESTS

## Naming Conventions:
 Golang naming conventions are pretty interesting with regards to capitalisation and camel casing, but follow those. Additionally, folders and files should be lower-case and hyphenated. We can use _* to indicate some sort of functional distinction, such as _test for test files, so ensure hyphens are instead used for word spacing. For example, bot-utils/on-message. Packages should match the folder structure - packages cannot contain hyphens, so should simply be an un-hyphenated version of the folder name. A bit ugly, but any possible ambiguity is removed by the folder name.

## [TODO](https://github.com/users/adad-mitch/projects/1)