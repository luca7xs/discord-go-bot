package events

import (
	"discord-go-bot/internal/bot/types"
	"discord-go-bot/internal/logger"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type ReadyEvent struct{}

func NewReadyEvent() types.Event {
	return &ReadyEvent{}
}

func (e *ReadyEvent) Name() string {
	return "ready"
}

func (e *ReadyEvent) Register(s *discordgo.Session) {
	s.AddHandler(e.Handle)
}

func (e *ReadyEvent) Handle(s *discordgo.Session, r *discordgo.Ready) {
	s.UpdateWatchStatus(0, "Usando comandos slash em Go!")
	logger.Log.Info("Bot está pronto",
		zap.String("user", s.State.User.Username),
		zap.String("id", s.State.User.ID),
		zap.Int("servidores", len(s.State.Guilds)))
}

// Registro automático do evento
var _ = RegisterEvent(NewReadyEvent())
