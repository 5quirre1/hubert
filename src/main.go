package main
import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"github.com/bwmarrin/discordgo"
)
var commandDefinitions = []*discordgo.ApplicationCommand{
	{
		Name:        "ping",
		Description: "greg",
	},
	{
		Name:        "repeat",
		Description: "repeats message",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "message",
				Description: "e",
				Required:    true,
			},
		},
	},
}
var registeredCommands []*discordgo.ApplicationCommand
func main() {
	token := os.Getenv("DISCORD_BOT_TOKEN")
	if token == "" {
		fmt.Println("DISCORD_BOT_TOKEN is not set")
		return
	}
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}
	dg.AddHandlerOnce(ready)
	dg.AddHandler(interactionCreate)
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	}
	defer dg.Close()
	appID := dg.State.User.ID
	for _, cmd := range commandDefinitions {
		createdCmd, err := dg.ApplicationCommandCreate(appID, "", cmd)
		if err != nil {
			fmt.Printf("cannot create '%s' command: %v\n", cmd.Name, err)
		} else {
			registeredCommands = append(registeredCommands, createdCmd)
		}
	}
	fmt.Println("bot is now running, press ctrl+c to exit")
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop
	for _, cmd := range registeredCommands {
		err := dg.ApplicationCommandDelete(appID, "", cmd.ID)
		if err != nil {
			fmt.Printf("failed to delete command '%s': %v\n", cmd.Name, err)
		}
	}
	fmt.Println("shit down")
}
func ready(s *discordgo.Session, r *discordgo.Ready) {
	fmt.Printf("logged in as: %v#%v\n", s.State.User.Username, s.State.User.Discriminator)
}
func interactionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.ApplicationCommandData().Name {
	case "ping":
		userID := i.Member.User.ID
		response := fmt.Sprintf("hai, <@%s>!", userID)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
	case "repeat":
		msg := ""
		for _, opt := range i.ApplicationCommandData().Options {
			if opt.Name == "message" {
				msg = opt.StringValue()
				break
			}
		}
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content:         msg,
				AllowedMentions: &discordgo.MessageAllowedMentions{},
			},
		})
	}
}
