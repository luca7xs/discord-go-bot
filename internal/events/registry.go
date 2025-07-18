package events

import (
	"discord-go-bot/internal/bot/types"
	"discord-go-bot/internal/logger"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var events []types.Event

// RegisterEvent adiciona um evento ao registro e retorna ele (para uso em init)
func RegisterEvent(event types.Event) types.Event {
	events = append(events, event)
	// Não usar logger aqui para evitar nil pointer durante init()
	return event
}

// GetEvents retorna todos os eventos registrados
func GetEvents() []types.Event {
	return events
}

// RegisterAllEvents registra todos os eventos na sessão do Discord
func RegisterAllEvents(s *discordgo.Session) {
	logger.Log.Info("Registrando eventos na sessão do Discord", zap.Int("total", len(events)))
	for _, event := range events {
		logger.Log.Debug("Registrando evento", zap.String("event", event.Name()))
		event.Register(s)
	}
}
