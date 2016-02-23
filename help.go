package main

import "os"

func help() {
	var cmd string
	if len(os.Args) > 1 {
		cmd = os.Args[1]
		if len(os.Args) > 2 && cmd == "help" {
			cmd = os.Args[2]
		}
	}
	cli.Commands.loadUsages()
	topic, command := cli.ParseCmd(cmd)
	switch {
	case topic == nil:
		Printf("**THIS IS A PROOF-OF-CONCEPT TOOL NOT AUTHORED BY PARTICLE.IO**\n\n")
		Printf("Usage: particle <command_name> <arguments>\n\n")
		Printf("Help topics, type \"particle help <command_name>\" for more details:\n\n")
		for _, topic := range nonHiddenTopics(cli.Topics) {
			Printf("  particle %-30s# %s\n", topic.Name, topic.Description)
		}
	case command == nil:
		Printf("Usage: particle %s:COMMAND [--app APP] [command-specific-options]\n\n", topic.Name)
		printTopicCommandsHelp(topic)
	case command.Command == "":
		printCommandHelp(command)
		// This is a root command so show the other commands in the topic
		// if there are any
		if len(topic.Commands()) > 1 {
			printTopicCommandsHelp(topic)
		}
	default:
		printCommandHelp(command)
	}
	os.Exit(0)
}

func printTopicCommandsHelp(topic *Topic) {
	commands := topic.Commands()
	if len(commands) > 0 {
		Printf("\nCommands for %s, type \"particle help %s:COMMAND\" for more details:\n\n", topic.Name, topic.Name)
		for _, command := range nonHiddenCommands(commands) {
			Printf(" particle %-30s # %s\n", command.Usage, command.Description)
		}
	}
}

func printCommandHelp(command *Command) {
	Printf("Usage: particle %s\n\n", command.Usage)
	Println(command.buildFullHelp())
}

func nonHiddenTopics(from TopicSet) TopicSet {
	to := make(TopicSet, 0, len(from))
	for _, topic := range from {
		if !topic.Hidden {
			to = append(to, topic)
		}
	}
	return to
}

func nonHiddenCommands(from []*Command) []*Command {
	to := make([]*Command, 0, len(from))
	for _, command := range from {
		if !command.Hidden {
			to = append(to, command)
		}
	}
	return to
}
