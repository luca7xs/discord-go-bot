package commands

import (
	"fmt"
	"sort"
	"strings"

	"discord-go-bot/internal/bot/types"
	"discord-go-bot/internal/logger"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type HelpCommand struct{}

func NewHelpCommand() types.Command {
	// Não usar logger aqui para evitar nil pointer durante init()
	return &HelpCommand{}
}

func (c *HelpCommand) Command() *discordgo.ApplicationCommand {
	return &discordgo.ApplicationCommand{
		Name:        "help",
		Description: "Mostra a lista de comandos disponíveis",
	}
}

func (c *HelpCommand) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Obter todos os comandos registrados
	cmds := GetCommands()

	// Ordenar comandos por nome
	sort.Slice(cmds, func(i, j int) bool {
		return cmds[i].Command().Name < cmds[j].Command().Name
	})

	// Construir a mensagem de ajuda
	var sb strings.Builder
	sb.WriteString("**Comandos Disponíveis:**\n\n")

	for _, cmd := range cmds {
		cmdDef := cmd.Command()
		sb.WriteString(fmt.Sprintf("**/%s** - %s\n", cmdDef.Name, cmdDef.Description))
	}

	sb.WriteString("\n*Use / seguido do nome do comando para executá-lo.*")

	// Responder ao usuário
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
			Flags:   discordgo.MessageFlagsEphemeral, // Mensagem visível apenas para o usuário que executou o comando
		},
	})

	if err != nil {
		logger.Log.Error("Erro ao responder comando help", zap.Error(err))
		return err
	}
	return nil
}

// Registro automático do comando
var _ = RegisterCommand(NewHelpCommand())
