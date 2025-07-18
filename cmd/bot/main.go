package main

import (
	"log"
	"os"

	"discord-go-bot/internal/bot"
	"discord-go-bot/internal/commands"
	"discord-go-bot/internal/events"
	"discord-go-bot/internal/logger"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Nenhum arquivo .env encontrado: %v", err)
	}

	if err := logger.Init(); err != nil {
		log.Fatalf("Falha ao inicializar logger: %v", err)
	}
	defer logger.Sync()

	token := os.Getenv("BOT_TOKEN")
	if token == "" {
		logger.Log.Fatal("BOT_TOKEN não definido")
	}

	guildID := os.Getenv("GUILD_ID")
	if guildID == "" {
		logger.Log.Fatal("GUILD_ID não definido")
	}

	// Carrega todos os eventos registrados automaticamente
	evts := events.GetEvents()
	logger.Log.Info("Eventos carregados", zap.Int("total", len(evts)))

	// Carrega todos os comandos registrados automaticamente
	cmds := commands.GetCommands()
	logger.Log.Info("Comandos carregados", zap.Int("total", len(cmds)))

	b, err := bot.NewBot(token, cmds)
	if err != nil {
		logger.Log.Fatal("Erro ao criar bot", zap.Error(err))
	}

	if err := b.Start(guildID); err != nil {
		logger.Log.Fatal("Erro ao iniciar bot", zap.Error(err))
	}

	b.Wait()
}
