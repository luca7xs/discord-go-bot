package types

import "github.com/bwmarrin/discordgo"

// Command representa um comando slash do bot
type Command interface {
	// Command retorna a definição do comando para registro na API Discord
	Command() *discordgo.ApplicationCommand

	// Handle é chamado quando o comando é invocado
	Handle(s *discordgo.Session, i *discordgo.InteractionCreate) error
}

// Event representa um evento do Discord que o bot pode processar
type Event interface {
	// Name retorna o nome do evento para fins de registro
	Name() string

	// Register registra o handler do evento na sessão do Discord
	Register(s *discordgo.Session)
}
