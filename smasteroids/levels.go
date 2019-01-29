package smasteroids

import "time"

var lowerSchoolShip = Ship{
	Health:       10,
	Thrust:       50,
	Turn:         2,
	Fire:         time.Second / 2,
	BulletDamage: 2,
}

var middleSchoolShip = Ship{
	Health:       15,
	Thrust:       60,
	Turn:         2.5,
	Fire:         time.Second / 2,
	BulletDamage: 4,
}

var biologyShip = Ship{
	Health:       20,
	Thrust:       75,
	Turn:         3,
	Fire:         time.Millisecond * 417,
	BulletDamage: 6,
}

var chemistryShip = Ship{
	Health:       30,
	Thrust:       80,
	Turn:         3.25,
	Fire:         time.Second / 2,
	BulletDamage: 8,
}

var loraxShip = Ship{
	Health:       75,
	Thrust:       50,
	Turn:         3.5,
	Fire:         time.Second / 3 * 2,
	BulletDamage: 15,
}

var physicsShip = Ship{
	Health:       50,
	Thrust:       100,
	Turn:         3.75,
	Fire:         time.Second / 3,
	BulletDamage: 15,
}

var administrationShip = Ship{
	Health:       100,
	Thrust:       120,
	Turn:         4.25,
	Fire:         time.Second / 200,
	BulletDamage: 10,
}

var Levels = []Level{
	{
		Name: "Lower School",
		Player: Ship{
			Health:       20,
			Thrust:       100,
			Turn:         3,
			Fire:         time.Second / 3,
			BulletDamage: 5,
		},
		Enemies: multiple(lowerSchoolShip, "Carrio", "Materre", "Pigg", "Wetzel", "Dillon"),
	},
	{
		Name: "Middle School",
		Player: Ship{
			Health:       30,
			Thrust:       120,
			Turn:         3,
			Fire:         time.Second / 3,
			BulletDamage: 7.5,
		},
		Enemies: multiple(middleSchoolShip, "Dillon", "Hoehn", "Kiehn", "Smith", "Mead"),
	},
	{
		Name: "The Biologists",
		Player: Ship{
			Health:       50,
			Thrust:       130,
			Turn:         3,
			Fire:         time.Second / 4,
			BulletDamage: 10,
		},
		Enemies: append(
			multiple(biologyShip, "Adame", "Flint", "Israni"),
			duplicate(2, Enemy{
				Name: "Skeleton",
				Ship: Ship{
					Health:       10,
					Thrust:       80,
					Turn:         3.75,
					Fire:         time.Second / 2,
					BulletDamage: 7.5,
				},
			})...,
		),
	},
	{
		Name: "The Chemists",
		Player: Ship{
			Health:       60,
			Thrust:       140,
			Turn:         3,
			Fire:         time.Second / 3,
			BulletDamage: 10,
		},
		Enemies: append(
			multiple(chemistryShip, "Macaraeg", "Owens"),
			duplicate(8, Enemy{
				Name: "VisorGogs",
				Ship: Ship{
					Health:       10,
					Thrust:       100,
					Turn:         3.5,
					Fire:         time.Second,
					BulletDamage: 2,
				},
			})...,
		),
	},
	{
		Name: "The Lorax",
		Player: Ship{
			Health:       80,
			Thrust:       150,
			Turn:         3,
			Fire:         time.Second / 3,
			BulletDamage: 15,
		},
		Enemies: append(
			duplicate(15, Enemy{
				Name: "Tree",
				Ship: Ship{
					Health:       1,
					Thrust:       20,
					Turn:         5,
					Fire:         time.Second / 3,
					BulletDamage: 1,
				},
			}),
			Enemy{
				Name: "Northcut",
				Ship: loraxShip,
			},
		),
	},
	{
		Name: "The Physicists",
		Player: Ship{
			Health:       100,
			Thrust:       150,
			Turn:         3,
			Fire:         time.Second / 5,
			BulletDamage: 10,
		},
		Enemies: multiple(physicsShip, "Houpt", "Hoehn", "Balog", "Carron"),
	},
	{
		Name: "Dwarf Fortress",
		Player: Ship{
			Health:       100,
			Thrust:       160,
			Turn:         3,
			Fire:         time.Second,
			BulletDamage: 7.5,
		},
		Enemies: []Enemy{
			{
				Name: "Dwarf King",
				Ship: Ship{
					Health:       100,
					Thrust:       10,
					Turn:         5,
					Fire:         time.Second * 5,
					BulletDamage: 20,
				},
			},
		},
	},
	{
		Name: "The Administration",
		Player: Ship{
			Health:       100,
			Thrust:       150,
			Turn:         3,
			Fire:         time.Second / 100,
			BulletDamage: 7.5,
		},
		Enemies: multiple(administrationShip, "Dini", "Igoe"),
	},
}
