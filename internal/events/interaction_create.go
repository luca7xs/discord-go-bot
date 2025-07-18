package events

import (
	"discord-go-bot/internal/bot/types"
	"discord-go-bot/internal/logger"
	"discord-go-bot/internal/metrics"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type InteractionCreateEvent struct {
	commands map[string]types.Command
	metrics  *metrics.CommandMetrics
}

func NewInteractionCreateEvent(commands map[string]types.Command, metrics *metrics.CommandMetrics) types.Event {
	// Não usar logger aqui para evitar nil pointer durante init()
	return &InteractionCreateEvent{
		commands: commands,
		metrics:  metrics,
	}
}

func (e *InteractionCreateEvent) Name() string {
	return "interaction_create"
}

func (e *InteractionCreateEvent) Register(s *discordgo.Session) {
	s.AddHandler(e.Handle)
}

func (e *InteractionCreateEvent) Handle(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Tratamento de comandos de aplicação
	if i.Type == discordgo.InteractionApplicationCommand {
		cmdName := i.ApplicationCommandData().Name
		cmd, ok := e.commands[cmdName]
		if !ok {
			logger.Log.Warn("Comando não encontrado",
				zap.String("command", cmdName),
				zap.String("guild_id", i.GuildID))
			// Responder ao usuário que o comando não foi encontrado
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Comando não encontrado. Por favor, tente novamente ou use /help para ver os comandos disponíveis.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}

		user := i.Member.User.Username
		userID := i.Member.User.ID
		guildID := i.GuildID
		channelID := i.ChannelID

		logger.Log.Info("Executando comando",
			zap.String("command", cmdName),
			zap.String("user", user),
			zap.String("userID", userID),
			zap.String("guildID", guildID),
			zap.String("channelID", channelID))

		start := time.Now()
		err := cmd.Handle(s, i)
		duration := time.Since(start)

		// Registrar métricas do comando
		if e.metrics != nil {
			e.metrics.RecordCommandExecution(cmdName, userID, guildID, duration, err != nil)
		}

		if err != nil {
			logger.Log.Error("Erro ao executar comando",
				zap.String("command", cmdName),
				zap.Error(err))
			// Tentar responder ao usuário sobre o erro
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Ocorreu um erro ao executar o comando. Por favor, tente novamente mais tarde.",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
		}
		return
	} else {
		logger.Log.Debug("Interação não suportada", zap.Int("tipo", int(i.Type)))
	}
}
