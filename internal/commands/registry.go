package commands

import (
	"discord-go-bot/internal/bot/types"
)

var commands []types.Command

// RegisterCommand adiciona um comando ao registro e retorna ele (para uso em init)
func RegisterCommand(cmd types.Command) types.Command {
	commands = append(commands, cmd)
	// Não usar logger aqui para evitar nil pointer durante init()
	return cmd
}

// GetCommands retorna todos os comandos registrados
func GetCommands() []types.Command {
	return commands
}
