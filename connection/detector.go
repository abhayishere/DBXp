package connection

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/abhayishere/DBXp/db"
	"github.com/docker/docker/api/types" // Add this import
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	Name     string
	Type     string
	Port     string
	Database string
	Status   string
}

func DetectDatabases(onConnect func(db.DatabaseConfig)) ([]db.DatabaseConfig, error) {
	// This function is a placeholder for detecting databases.
	// It should run `docker ps` and parse the output to find running database containers.
	// For each container, it should run `docker inspect` to get environment variables and port information.
	// Finally, it should build a slice of db.DatabaseConfig for each detected database and return it.

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, fmt.Errorf("failed to create Docker client: %w", err)
	}
	defer cli.Close()
	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list containers: %w", err)
	}

	detected := []db.DatabaseConfig{}
	for _, container := range containers {
		// Inspect the container to get detailed information
		if isDatabaseContainer(container) {
			inspect, err := cli.ContainerInspect(context.Background(), container.ID)
			if err != nil {
				return nil, fmt.Errorf("failed to inspect container %s: %w", container.ID, err)
			}
			config := db.DatabaseConfig{
				Host: "localhost",
			}
			if strings.Contains(container.Image, "postgres") {
				config.Type = "PostgreSQL"
				config.Port = extractPort(container.Ports, 5432)
				config.User = extractEnvVar(inspect.Config.Env, "POSTGRES_USER", "postgres")
				config.Password = extractEnvVar(inspect.Config.Env, "POSTGRES_PASSWORD", "")
				config.Database = extractEnvVar(inspect.Config.Env, "POSTGRES_DB", "postgres")
			} else if strings.Contains(container.Image, "mysql") {
				config.Type = "MySQL"
				config.Port = extractPort(container.Ports, 3306)
				config.User = extractEnvVar(inspect.Config.Env, "MYSQL_USER", "root")
				config.Password = extractEnvVar(inspect.Config.Env, "MYSQL_PASSWORD", "")
				config.Database = extractEnvVar(inspect.Config.Env, "MYSQL_DATABASE", "mysql")
			}
			detected = append(detected, config)
		}
	}

	if len(detected) == 0 {
		return nil, errors.New("no databases detected")
	}
	return detected, nil
}

func isDatabaseContainer(container types.Container) bool {
	return strings.Contains(container.Image, "postgres") || strings.Contains(container.Image, "mysql") || strings.Contains(container.Image, "sqlite")
}

func extractPort(ports []types.Port, defaultPort int) string {
	for _, port := range ports {
		if port.PrivatePort == uint16(defaultPort) && port.PublicPort > 0 {
			return fmt.Sprintf("%d", port.PublicPort)
		}
	}
	return fmt.Sprintf("%d", defaultPort)
}

func extractEnvVar(env []string, key, defaultValue string) string {
	for _, envVar := range env {
		if strings.HasPrefix(envVar, key+"=") {
			return strings.TrimPrefix(envVar, key+"=")
		}
	}
	return defaultValue
}
