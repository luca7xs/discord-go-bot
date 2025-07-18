package metrics

import (
	"sync"
	"time"

	"discord-go-bot/internal/logger"

	"go.uber.org/zap"
)

// CommandMetrics armazena métricas de uso dos comandos
type CommandMetrics struct {
	mutex           sync.RWMutex
	CommandCounts   map[string]int           // Contagem de uso por comando
	CommandDuration map[string]time.Duration // Duração média por comando
	CommandErrors   map[string]int           // Contagem de erros por comando
	UserCounts      map[string]int           // Contagem de uso por usuário
	GuildCounts     map[string]int           // Contagem de uso por servidor
}

// NewCommandMetrics cria uma nova instância de CommandMetrics
func NewCommandMetrics() *CommandMetrics {
	return &CommandMetrics{
		CommandCounts:   make(map[string]int),
		CommandDuration: make(map[string]time.Duration),
		CommandErrors:   make(map[string]int),
		UserCounts:      make(map[string]int),
		GuildCounts:     make(map[string]int),
	}
}

// RecordCommandExecution registra a execução de um comando
func (m *CommandMetrics) RecordCommandExecution(command, userID, guildID string, duration time.Duration, hasError bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Incrementar contagem do comando
	m.CommandCounts[command]++

	// Atualizar duração média
	currentDuration, exists := m.CommandDuration[command]
	if !exists {
		m.CommandDuration[command] = duration
	} else {
		count := m.CommandCounts[command]
		// Calcular nova média: (média_atual * (count-1) + nova_duração) / count
		m.CommandDuration[command] = (currentDuration*time.Duration(count-1) + duration) / time.Duration(count)
	}

	// Incrementar contagem de erros se aplicável
	if hasError {
		m.CommandErrors[command]++
	}

	// Incrementar contagem por usuário
	m.UserCounts[userID]++

	// Incrementar contagem por servidor
	m.GuildCounts[guildID]++
}

// LogMetrics registra as métricas atuais no logger
func (m *CommandMetrics) LogMetrics() {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	logger.Log.Info("Métricas de comandos",
		zap.Int("total_comandos_executados", m.getTotalCommandCount()),
		zap.Int("total_usuarios_unicos", len(m.UserCounts)),
		zap.Int("total_servidores_unicos", len(m.GuildCounts)),
		zap.Int("total_erros", m.getTotalErrorCount()),
	)

	// Registrar detalhes por comando
	for cmd, count := range m.CommandCounts {
		errors := m.CommandErrors[cmd]
		duration := m.CommandDuration[cmd]
		logger.Log.Debug("Métricas do comando",
			zap.String("comando", cmd),
			zap.Int("execucoes", count),
			zap.Int("erros", errors),
			zap.Duration("duracao_media", duration),
			zap.Float64("taxa_erro", float64(errors)/float64(count)),
		)
	}
}

// getTotalCommandCount retorna o total de comandos executados
func (m *CommandMetrics) getTotalCommandCount() int {
	total := 0
	for _, count := range m.CommandCounts {
		total += count
	}
	return total
}

// getTotalErrorCount retorna o total de erros ocorridos
func (m *CommandMetrics) getTotalErrorCount() int {
	total := 0
	for _, count := range m.CommandErrors {
		total += count
	}
	return total
}
