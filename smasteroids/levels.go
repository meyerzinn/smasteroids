package smasteroids

import "time"

var lowerSchoolEnemy = Ship{
	Health:       10,
	Thrust:       50,
	Turn:         2,
	Fire:         time.Second / 2,
	BulletDamage: 2,
}

var middleSchoolEnemy = Ship{
	Health:       15,
	Thrust:       60,
	Turn:         2.5,
	Fire:         time.Second / 2,
	BulletDamage: 4,
}

var biologyEnemy = Ship{
	Health:       20,
	Thrust:       75,
	Turn:         3,
	Fire:         time.Millisecond * 417,
	BulletDamage: 6,
}

var chemistryEnemy = Ship{
	Health:       30,
	Thrust:       80,
	Turn:         3.25,
	Fire:         time.Second / 2,
	BulletDamage: 8,
}

var loraxEnemy = Ship{
	Health:       75,
	Thrust:       50,
	Turn:         3.5,
	Fire:         time.Second / 3 * 2,
	BulletDamage: 15,
}

var physicsEnemy = Ship{
	Health:       50,
	Thrust:       100,
	Turn:         3.75,
	Fire:         time.Second / 3,
	BulletDamage: 15,
}

var administrationEnemy = Ship{
	Health:       100,
	Thrust:       120,
	Turn:         4.25,
	Fire:         10,
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
		Enemies: []Enemy{
			{
				Name: "Carrio",
				Ship: lowerSchoolEnemy,
			},
			{
				Name: "Materre",
				Ship: lowerSchoolEnemy,
			},
			{
				Name: "Pigg",
				Ship: lowerSchoolEnemy,
			},
			{
				Name: "Wetzel",
				Ship: lowerSchoolEnemy,
			},
			{
				Name: "Dillon",
				Ship: lowerSchoolEnemy,
			},
		},
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
		Enemies: []Enemy{
			{
				Name: "Dillon",
				Ship: middleSchoolEnemy,
			},
			{
				Name: "Hoehn",
				Ship: middleSchoolEnemy,
			},
			{
				Name: "Kiehn",
				Ship: middleSchoolEnemy,
			},
			{
				Name: "Smith",
				Ship: middleSchoolEnemy,
			},
			{
				Name: "Mead",
				Ship: middleSchoolEnemy,
			},
		},
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
		Enemies: append([]Enemy{
			{
				Name: "Flint",
				Ship: biologyEnemy,
			},
			{
				Name: "Israni",
				Ship: biologyEnemy,
			},
			{
				Name: "Adame",
				Ship: biologyEnemy,
			},
		},
			mult(2, Enemy{
				Name: "Skeleton",
				Ship: Ship{
					Health:       10,
					Thrust:       80,
					Turn:         3.75,
					Fire:         time.Second / 2,
					BulletDamage: 7.5,
				},
			})...),
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
		Enemies: append([]Enemy{
			{
				Name: "Macaraeg",
				Ship: chemistryEnemy,
			},
			{
				Name: "Owens",
				Ship: chemistryEnemy,
			},
		},
			mult(8, Enemy{
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
		Enemies: append([]Enemy{
			{
				Name: "Northcut",
				Ship: loraxEnemy,
			},
		},
			mult(15, Enemy{
				Name: "Tree",
				Ship: Ship{
					Health:       1,
					Thrust:       20,
					Turn:         5,
					Fire:         time.Second / 3,
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
			Turn:         3,
			Fire:         time.Second / 6,
			BulletDamage: 5,
		},
		Enemies: []Enemy{
			{
				Name: "Houpt",
				Ship: physicsEnemy,
			},
			{
				Name: "Hoehn",
				Ship: physicsEnemy,
			},
			{
				Name: "Balog",
				Ship: physicsEnemy,
			},
			{
				Name: "Carron",
				Ship: physicsEnemy,
			},
		},
	},
	{
		Name: "The Administration",
		Player: Ship{
			Health:       100,
			Thrust:       150,
			Turn:         3,
			Fire:         10,
			BulletDamage: 7.5,
		},
		Enemies: []Enemy{
			{
				Name: "Dini",
				Ship: administrationEnemy,
			},
			{
				Name: "Igoe",
				Ship: administrationEnemy,
			},
		},
	},
}
