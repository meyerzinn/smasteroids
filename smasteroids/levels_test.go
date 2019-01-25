package smasteroids_test

import (
	"github.com/20zinnm/smasteroids/smasteroids"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLevels(t *testing.T) {
	for i, level := range smasteroids.Levels {
		assert.NotEmpty(t, level.Name, "level does not have a name")
		if !testShip(t, level.Player) {
			t.Errorf("levelIndex=%d player spec failed ship test", i)
		}
		if assert.NotEmpty(t, level.Enemies) {
			for j, enemy := range level.Enemies {
				if enemy.Name == "" {
					t.Logf("levelIndex=%d enemyIndex=%d enemy has no name?", i, j)
				}
				if !testShip(t, enemy.Ship) {
					t.Errorf("levelIndex=%d enemyIndex=%d enemy spec failed ship test", i, j)
				}
			}
		}
	}
}

func testShip(t *testing.T, ship smasteroids.Ship, args ...interface{}) bool {
	return assert.NotEmpty(t, ship.Thrust, "ship has no thrust", args) &&
		assert.NotEmpty(t, ship.BulletDamage, "ship has no bullet damage", args) &&
		assert.NotEmpty(t, ship.Health, "ship has no health", args) &&
		assert.NotEmpty(t, ship.Fire, "ship has no firing delay", args) &&
		assert.NotEmpty(t, ship.Turn, "ship cannot turn", args)
}
