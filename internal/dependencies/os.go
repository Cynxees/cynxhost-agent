package dependencies

import (
	"bufio"
	"cynxhostagent/internal/model/entity"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type OSManager struct {
}

func NewOsManager() *OSManager {
	return &OSManager{}
}

func (manager *OSManager) RunBashScript(script string) error {
	// Check if the script exists
	cmd := exec.Command("bash", "-c", script) // Assuming the script is a bash script

	// Capture the output and error
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error executing script: %v, output: %s", err, string(output))
	}

	// You can log the output here if necessary
	fmt.Printf("Script Output: %s\n", string(output))
	return nil
}

func (*OSManager) ReadServerProperties(filePath string) ([]entity.ServerProperty, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open server.properties: %w", err)
	}
	defer file.Close()

	properties := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" { // Skip comments and empty lines
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			properties[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading server.properties: %w", err)
	}

	var serverProperties []entity.ServerProperty
	for key, value := range properties {
		serverProperties = append(serverProperties, entity.ServerProperty{
			Key:   key,
			Value: value,
		})
	}

	return serverProperties, nil
}

func (*OSManager) SetServerProperties(filePath string, serverProperties []entity.ServerProperty) error {

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create server.properties: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, property := range serverProperties {
		_, err := writer.WriteString(fmt.Sprintf("%s=%s\n", property.Key, property.Value))
		if err != nil {
			return fmt.Errorf("failed to write to server.properties: %w", err)
		}
	}
	return writer.Flush()
}
