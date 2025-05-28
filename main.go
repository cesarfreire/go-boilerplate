/*
Copyright © 2025 Cesar Freire iceesar@live.com
*/
package main

import (
	"github.com/cesarfreire/go-boilerplate/cmd"
	"github.com/cesarfreire/go-boilerplate/internal/infra/logger"
	"log"
)

func main() {
	loggerConfig := logger.Config{
		// IsDevelopment: isDevFromConfig, // Se não definir, usará a ENV "DEV_MODE"
	}

	appLogger, err := logger.NewLogger(loggerConfig)
	if err != nil {
		log.Fatalf("Falha ao inicializar o logger: %v", err) // Usando log padrão do Go para este erro crítico
	}
	defer appLogger.Sync() // Importante para garantir que os logs sejam escritos antes de sair

	cmd.Execute()
}
