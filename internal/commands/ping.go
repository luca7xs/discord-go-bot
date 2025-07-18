package commands

import (
	"discord-go-bot/internal/bot/types"
	"discord-go-bot/internal/logger"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type PingCommand struct{}

func NewPingCommand() types.Command {
	// Não usar logger aqui para evitar nil pointer durante init()
	return &PingCommand{}
}

func (c *PingCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "ping",
		Description: "Responde com Pong!",
	}
}

func (c *PingCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "Pong!",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})

	if err != nil {
		logger.Log.Error("Erro ao responder comando ping", zap.Error(err))
		return err
	}
	return nil
}

// Registro automático do comando
var _ = RegisterCommand(NewPingCommand())
