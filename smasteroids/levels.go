package smasteroids

type Ship struct {
	// Health is the amount of health the ship starts with and the maximum health they can have.
	Health float64
	// Thrust is the amount of force applied at every tick the ship is thrusting.
	Thrust float64
	// Turn is the maximum angular velocity the ship can have.
	Turn float64
	// Fire is the minimum number of ticks before the ship can fire again.
	Fire int
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

var lowerSchoolEnemy = Ship{
	Health:       10,
	Thrust:       50,
	Turn:         2,
	Fire:         30,
	BulletDamage: 2,
}

var middleSchoolEnemy = Ship{
	Health:       15,
	Thrust:       60,
	Turn:         2.5,
	Fire:         30,
	BulletDamage: 4,
}

var biologyEnemy = Ship{
	Health:       20,
	Thrust:       75,
	Turn:         3,
	Fire:         25,
	BulletDamage: 6,
}

var chemistryEnemy = Ship{
	Health:       30,
	Thrust:       80,
	Turn:         3.25,
	Fire:         30,
	BulletDamage: 8,
}

var loraxEnemy = Ship{
	Health:       75,
	Thrust:       120,
	Turn:         3.5,
	Fire:         20,
	BulletDamage: 15,
}

var physicsEnemy = Ship{
	Health:       50,
	Thrust:       100,
	Turn:         3.75,
	Fire:         20,
	BulletDamage: 15,
}

var Levels = []Level{
	{
		Name: "Lower School",
		Player: Ship{
			Health:       20,
			Thrust:       100,
			Turn:         4,
			Fire:         20,
			BulletDamage: 5,
		},
		Enemies: []Enemy{
			{
				Name: "Kay Carrio",
				Ship: lowerSchoolEnemy,
			},
			{
				Name: "Debra Materre",
				Ship: lowerSchoolEnemy,
			},
			{
				Name: "Laura Pigg",
				Ship: lowerSchoolEnemy,
			},
			{
				Name: "Catherine Wetzel",
				Ship: lowerSchoolEnemy,
			},
		},
	},
	{
		Name: "Middle School",
		Player: Ship{
			Health:       30,
			Thrust:       120,
			Turn:         4,
			Fire:         20,
			BulletDamage: 7.5,
		},
		Enemies: []Enemy{
			{
				Name: "Matt Dillon",
				Ship: middleSchoolEnemy,
			},
			{
				Name: "Paul Hoehn",
				Ship: middleSchoolEnemy,
			},
			{
				Name: "Donald Kiehn",
				Ship: middleSchoolEnemy,
			},
			{
				Name: "David Smith",
				Ship: middleSchoolEnemy,
			},
		},
	},
	{
		Name: "The Biologists",
		Player: Ship{
			Health:       50,
			Thrust:       130,
			Turn:         4,
			Fire:         15,
			BulletDamage: 10,
		},
		Enemies: append([]Enemy{
			{
				Name: "Bonnie Flint",
				Ship: biologyEnemy,
			},
			{
				Name: "Nupur Israni",
				Ship: biologyEnemy,
			},
			{
				Name: "Mark Adame",
				Ship: biologyEnemy,
			},
		},
			mult(2, Enemy{
				Name: "Skeleton",
				Ship: Ship{
					Health:       10,
					Thrust:       80,
					Turn:         3.75,
					Fire:         30,
					BulletDamage: 7.5,
				},
			})...),
	},
	{
		Name: "The Chemists",
		Player: Ship{
			Health:       60,
			Thrust:       140,
			Turn:         4,
			Fire:         15,
			BulletDamage: 10,
		},
		Enemies: append([]Enemy{
			{
				Name: "Cristina Macaraeg",
				Ship: chemistryEnemy,
			},
			{
				Name: "Ken Owens",
				Ship: chemistryEnemy,
			},
		},
			mult(8, Enemy{
				Name: "VisorGogs",
				Ship: Ship{
					Health:       10,
					Thrust:       100,
					Turn:         3.5,
					Fire:         60,
					BulletDamage: 2,
				},
			})...
		),
	},
	{
		Name: "The Lorax",
		Player: Ship{
			Health:       80,
			Thrust:       150,
			Turn:         4,
			Fire:         20,
			BulletDamage: 15,
		},
		Enemies: append([]Enemy{
			{
				Name: "Daniel Northcut",
				Ship: loraxEnemy,
			},
			{
				Name: "John Mead",
				Ship: loraxEnemy,
			},
		},
			mult(15, Enemy{
				Name: "Tree",
				Ship: Ship{
					Health:       1,
					Thrust:       20,
					Turn:         5,
					Fire:         10,
					BulletDamage: 1,
				},
			})...,
		),
	},
	{
		Name: "The Physicists",
		Player: Ship{
			Health:       100,
			Thrust:       150,
			Turn:         4,
			Fire:         10,
			BulletDamage: 5,
		},
		Enemies: []Enemy{
			{
				Name: "Stephen Houpt",
				Ship: physicsEnemy,
			},
			{
				Name: "Paul Hoehn",
				Ship: physicsEnemy,
			},
			{
				Name: "Stephen Balog",
				Ship: physicsEnemy,
			},
			{
				Name: "Fletcher Carron",
				Ship: physicsEnemy,
			},
		},
	},
}

func mult(n int, enemy Enemy) (out []Enemy) {
	for i := 0; i < n; i++ {
		out = append(out, enemy)
	}
	return
}
