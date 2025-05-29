package specs

import (
	"encoding/json"
	"fmt"
	"os"
)

// TroopSpec represents the specification for a troop
type TroopSpec struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Health      int     `json:"health"`
	Damage      int     `json:"damage"`
	Defence     int     `json:"defence"`
	Range       float64 `json:"range"`
	Speed       float64 `json:"speed"`
	AttackSpeed float64 `json:"attack_speed"`
	Cost        int     `json:"cost"`
	Target      string  `json:"target"` // "ground", "air", "both"
}

// TowerSpec represents the specification for a tower
type TowerSpec struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // "king", "princess", "cannon"
	Health      int     `json:"health"`
	Damage      int     `json:"damage"`
	Defence     int     `json:"defence"`
	Range       float64 `json:"range"`
	AttackSpeed float64 `json:"attack_speed"`
	Target      string  `json:"target"` // "ground", "air", "both"
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
		if troop.Health <= 0 {
			return fmt.Errorf("invalid health for troop %s: %d", name, troop.Health)
		}
		if troop.Damage < 0 {
			return fmt.Errorf("invalid damage for troop %s: %d", name, troop.Damage)
		}
		if troop.Defence < 0 {
			return fmt.Errorf("invalid defence for troop %s: %d", name, troop.Defence)
		}
		if troop.Range <= 0 {
			return fmt.Errorf("invalid range for troop %s: %f", name, troop.Range)
		}
		if troop.Speed <= 0 {
			return fmt.Errorf("invalid speed for troop %s: %f", name, troop.Speed)
		}
		if troop.AttackSpeed <= 0 {
			return fmt.Errorf("invalid attack speed for troop %s: %f", name, troop.AttackSpeed)
		}
		if troop.Cost <= 0 {
			return fmt.Errorf("invalid cost for troop %s: %d", name, troop.Cost)
		}
		if troop.Target != "ground" && troop.Target != "air" && troop.Target != "both" {
			return fmt.Errorf("invalid target type for troop %s: %s", name, troop.Target)
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
		if tower.Range <= 0 {
			return fmt.Errorf("invalid range for tower %s: %f", name, tower.Range)
		}
		if tower.AttackSpeed <= 0 {
			return fmt.Errorf("invalid attack speed for tower %s: %f", name, tower.AttackSpeed)
		}
		if tower.Target != "ground" && tower.Target != "air" && tower.Target != "both" {
			return fmt.Errorf("invalid target type for tower %s: %s", name, tower.Target)
		}
		if tower.Type != "king" && tower.Type != "princess" && tower.Type != "cannon" {
			return fmt.Errorf("invalid tower type for %s: %s", name, tower.Type)
		}
	}

	return nil
}
