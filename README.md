# MegaBot

MegaBot not a complete piece of software!

MegaBot wants to be a multi service bot with some services. MegaBot never stores your messages. He's cool like that.

## Providers MegaBot supports
* Discord

## Services MegaBot Has
### Responder
Can respond to messages sent by users. Matching is performed by RegEx. Can be configured to always respond only when at'd or DM'd. 

## Services MegaBot Wants
### Kudos
Every user in a group has 5 kudos a day to give to other users when they're awesome. Group leader can customize emoji for kudos.

## Configuring MegaBot
MegaBot uses environment variables for config.
### Required
#### COOKIE_SECRET
Used to encrypt cookie data. should be a big long *uwu* random string. Once you set it, it can never be changed without your users REALLY disliking you. 

#### DB_ENGINE
Postgres connection string.

#### EXT_HOSTNAME
The external hostname the bot will live at. this helps with login callbacks.

#### REDIS_ADDRESS
Address of redis server.

### Optional 
#### DISCORD_KEY
Discord OAUTH key for "login with discord".

#### DISCORD_SECRET
Discord OAUTH secret for "login with discord".

#### DISCORD_TOKEN
Discord bot token. Setting this will enable the discord provider.

#### LOG_LEVEL
changes the logging level of MegaBot. Default is info. Optinos are error, warn, info, debug, and trace. 

Warning: trace level may expose messages in log.

#### REDIS_PASSWORD
If your redis instance has a password set it here.

#### RESPONDER_WORKERS
The number of workers that can send responses. Default is 4. 