package smasteroids

import "time"

type Ship struct {
	// Health is the amount of health the ship starts with and the maximum health they can have.
	Health float64
	// Thrust is the amount of force applied at every tick the ship is thrusting.
	Thrust float64
	// Turn is the maximum angular velocity the ship can have.
	Turn float64
	// Fire is the minimum number of ticks before the ship can fire again.
	Fire time.Duration
	// BulletDamage is the damage for bullets the ship fires.
	BulletDamage float64
}

// Enemy represents an enemy in a level that must be defeated for the player to advance to the next level or win.
type Enemy struct {
	// Name is the name of the enemy to be displayed. If blank, no name is displayed (but there is still a vertical gap
	// between the ship and its health bar).
	Name string
	Ship
}

// Level represents a game level the player must progress through.
type Level struct {
	Name    string
	Player  Ship
	Enemies []Enemy
}

// multiple returns multiple enemies with the same ship, each with a different name.
func multiple(ship Ship, names ...string) (out []Enemy) {
	for _, name := range names {
		out = append(out, Enemy{
			Name: name,
			Ship: ship,
		})
	}
	return
}

// duplicate returns multiple enemies from the given enemy.
func duplicate(n int, enemy Enemy) (out []Enemy) {
	for i := 0; i < n; i++ {
		out = append(out, enemy)
	}
	return
}
