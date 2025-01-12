package dependencies

import (
	"fmt"
	"os/exec"
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
