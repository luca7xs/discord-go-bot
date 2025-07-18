package bot

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"discord-go-bot/internal/bot/types"
	"discord-go-bot/internal/events"
	"discord-go-bot/internal/logger"
	"discord-go-bot/internal/metrics"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type Bot struct {
	session  *discordgo.Session
	stop     chan os.Signal
	commands map[string]types.Command
	metrics  *metrics.CommandMetrics
}

func NewBot(token string, cmds []types.Command) (*Bot, error) {
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	commandsMap := make(map[string]types.Command)
	for _, cmd := range cmds {
		commandsMap[cmd.Command().Name] = cmd
	}

	logger.Log.Info("Inicializando bot")

	b := &Bot{
		session:  dg,
		stop:     make(chan os.Signal, 1),
		commands: commandsMap,
		metrics:  metrics.NewCommandMetrics(),
	}

	// Registra o evento de interação com métricas
	events.RegisterEvent(events.NewInteractionCreateEvent(commandsMap, b.metrics))

	// Registra todos os eventos na sessão
	events.RegisterAllEvents(dg)

	// Configura as intents necessárias para os eventos
	dg.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentsDirectMessages | discordgo.IntentsMessageContent

	return b, nil
}

func (b *Bot) Start(guildID string) error {
	if err := b.session.Open(); err != nil {
		logger.Log.Error("Falha ao abrir conexão", zap.Error(err))
		return err
	}
	logger.Log.Info("Bot iniciado com sucesso")

	// Obter comandos existentes no Discord
	existingCommands, err := b.session.ApplicationCommands(b.session.State.User.ID, guildID)
	if err != nil {
		logger.Log.Error("Falha ao obter comandos existentes", zap.Error(err))
		return err
	}

	// Mapear comandos existentes por nome para facilitar a busca
	existingCommandsMap := make(map[string]*discordgo.ApplicationCommand)
	for _, cmd := range existingCommands {
		existingCommandsMap[cmd.Name] = cmd
	}

	// Registrar comandos que não existem e atualizar os existentes
	logger.Log.Info("Registrando comandos no Discord", zap.Int("total", len(b.commands)))
	for name, cmd := range b.commands {
		cmdDef := cmd.Command()
		if _, exists := existingCommandsMap[name]; exists {
			logger.Log.Debug("Registrando comando", zap.String("command", name))
			// Remover do mapa para saber quais comandos devem ser removidos depois
			delete(existingCommandsMap, name)
		} else {
			// Criar novo comando
			logger.Log.Info("Criando novo comando", zap.String("command", name))
			_, err := b.session.ApplicationCommandCreate(b.session.State.User.ID, guildID, cmdDef)
			if err != nil {
				logger.Log.Error("Falha ao criar comando", zap.String("command", name), zap.Error(err))
				return err
			}
		}
	}

	// Remover comandos que não existem mais no bot
	for name, cmd := range existingCommandsMap {
		logger.Log.Info("Removendo comando não utilizado", zap.String("command", name))
		err := b.session.ApplicationCommandDelete(b.session.State.User.ID, guildID, cmd.ID)
		if err != nil {
			logger.Log.Error("Falha ao remover comando", zap.String("command", name), zap.Error(err))
			// Não retornar erro aqui para não interromper o bot por falha na remoção
		}
	}

	signal.Notify(b.stop, os.Interrupt, syscall.SIGTERM)

	return nil
}

func (b *Bot) Wait() {
	// Iniciar goroutine para registrar métricas periodicamente
	ticker := time.NewTicker(1 * time.Hour)
	go func() {
		for {
			select {
			case <-ticker.C:
				b.metrics.LogMetrics()
			case <-b.stop:
				return
			}
		}
	}()

	<-b.stop
	logger.Log.Info("Encerrando bot...")

	// Registrar métricas finais antes de encerrar
	b.metrics.LogMetrics()

	ticker.Stop()
	b.session.Close()
}
