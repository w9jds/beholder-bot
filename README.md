# Beholder Discord Bot

#### Created for Discord Hack Week 2019

The goal of this discord bot is to create a tool that helps DMs run their RPG games! There are loads of different tools online that will help you run a virtual table top RPG. However, Discord is one of the best solutions for setting up a voice and text channel for your group. So, unlike a normal Beholder, this bot will help you keep track and setup what you need.

---
## Commands
- `!creategame [name] [@mention players]` (Done)
    - This command automatically creates a category with a text and voice channel based on the name you pass the command. The person running the command is setup as the `DM` of the group, and gets admin for the category. All mentioned players are added to the permissions list as players, only getting permissions for reading, writing, and using the voice channel.
- `!rollstats` (Planned) 
    - This command will role dice and output values for your stats that you can assign as you wish for a brand new DnD character.
- `!setnextsession [datetime]` [DM Only] (Planned)
    - Setting the next session with a datetime formatted like `mm/dd/yyyy hh:mm` will store when you are planing to play your next game.
- `!nextsession` (Planned)
    - Print out when the next session is scheduled for.
- `!bestdaypoll` [DM Only] (In Progress)
    - Create a reaction based poll to find out which day of the week is best suited for all of your players.
- `!addmap [name]` [DM Only] (In Progress)
    - Register a new map for the game channel you are in (only works in game text channels and requires a file to be attached).
- `!getmap [name]` (In Progress)
    - Pull the map with the name provided. This will pull the file from the message that is saved, and post it again automatically for you.

## Environment Variables

|Variables||
|---|---|
|`BOT_TOKEN`|Discord bot token|
|`POSTGRES_HOST`|Postgres host, in this case Instance connection name `project-name:region:instance-name`|
|`POSTGRES_DB`|database name|
|`POSTGRES_PORT`|connection port|
|`POSTGRES_USER`|bot user|
|`POSTGRES_PASSWORD`|user password|

## Components

This bot is currently powered by Postgres, and configured in the code to connect directly to a Google Cloud Platform Cloud SQL instance. However, it can easily be updated to just connect to a normal Postgres hosted instance.