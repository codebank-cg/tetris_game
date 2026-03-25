package model

import "fmt"

// presets maps preset name to heuristic weights.
var presets = map[string]map[string]float64{
	"aggressive": {
		"aggregateHeight": -0.05,
		"holes":           -0.10,
		"bumpiness":       -0.03,
		"wells":           -0.02,
	},
	"conservative": {
		"aggregateHeight": -0.90,
		"holes":           -1.50,
		"bumpiness":       -0.60,
		"wells":           -0.45,
	},
	"balanced": {
		"aggregateHeight": -0.54,
		"holes":           -0.90,
		"bumpiness":       -0.36,
		"wells":           -0.24,
	},
	"speedrun": {
		"aggregateHeight": -0.24,
		"holes":           -0.45,
		"bumpiness":       -0.15,
		"wells":           -0.09,
	},
	"chaos": {
		"aggregateHeight": -0.01,
		"holes":           -0.02,
		"bumpiness":       -0.01,
		"wells":           -0.01,
	},
}

// GetPreset returns a copy of the named preset's weights.
// Returns an error if the name is not recognized.
func GetPreset(name string) (map[string]float64, error) {
	p, ok := presets[name]
	if !ok {
		return nil, fmt.Errorf("unknown preset %q: valid presets are aggressive, conservative, balanced, speedrun, chaos", name)
	}
	copy := make(map[string]float64, len(p))
	for k, v := range p {
		copy[k] = v
	}
	return copy, nil
}
