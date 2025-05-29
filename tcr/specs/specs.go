package specs

import (
	"encoding/json"
	"fmt"
	"os"
)

// TroopSpec represents the specification for a troop
type TroopSpec struct {
	Name    string `json:"name"`
	Health  int    `json:"health"`
	Damage  int    `json:"damage"`
	Defence int    `json:"defence"`
	Cost    int    `json:"cost"`
}

// TowerSpec represents the specification for a tower
type TowerSpec struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Health  int    `json:"health"`
	Damage  int    `json:"damage"`
	Defence int    `json:"defence"`
}

// Specs holds all game specifications
type Specs struct {
	Troops map[string]TroopSpec `json:"troops"`
	Towers map[string]TowerSpec `json:"towers"`
}

// LoadSpecs reads and parses the specifications file
func LoadSpecs(filename string) (*Specs, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var specs Specs
	if err := json.Unmarshal(data, &specs); err != nil {
		return nil, err
	}

	// Validate specifications
	if err := validateSpecs(&specs); err != nil {
		return nil, err
	}

	return &specs, nil
}

// validateSpecs checks if the specifications are valid
func validateSpecs(specs *Specs) error {
	// Validate troops
	if len(specs.Troops) == 0 {
		return fmt.Errorf("no troops defined")
	}

	for name, troop := range specs.Troops {
		if name == "" {
			return fmt.Errorf("troop name cannot be empty")
		}
		if troop.Health < 0 {
			return fmt.Errorf("invalid health for troop %s: %d", name, troop.Health)
		}
		if troop.Damage < 0 {
			return fmt.Errorf("invalid damage for troop %s: %d", name, troop.Damage)
		}
		if troop.Defence < 0 {
			return fmt.Errorf("invalid defence for troop %s: %d", name, troop.Defence)
		}
		if troop.Cost < 0 {
			return fmt.Errorf("invalid cost for troop %s: %d", name, troop.Cost)
		}
	}

	// Validate towers
	if len(specs.Towers) == 0 {
		return fmt.Errorf("no towers defined")
	}

	for name, tower := range specs.Towers {
		if name == "" {
			return fmt.Errorf("tower name cannot be empty")
		}
		if tower.Health <= 0 {
			return fmt.Errorf("invalid health for tower %s: %d", name, tower.Health)
		}
		if tower.Damage < 0 {
			return fmt.Errorf("invalid damage for tower %s: %d", name, tower.Damage)
		}
		if tower.Defence < 0 {
			return fmt.Errorf("invalid damage for tower %s: %d", name, tower.Defence)
		}
		if tower.Type != "king" && tower.Type != "guard" {
			return fmt.Errorf("invalid tower type for %s: %s", name, tower.Type)
		}
	}

	return nil
}
