// +build varsity

package smasteroids

var varsityEnemy = Enemy{
	Name: "Quiz",
	Ship: Ship{
		Health:       1,
		Thrust:       20,
		Turn:         5,
		Fire:         20,
		BulletDamage: 1,
	},
}

var perrymanEnemy = Ship{
	Health:       100,
	Thrust:       110,
	Turn:         4,
	Fire:         120,
	BulletDamage: 20,
}

var Levels = []Level{
	{
		Name: "The Varsity",
		Player: Ship{
			Health:       20,
			Thrust:       100,
			Turn:         4,
			Fire:         20,
			BulletDamage: 5,
		},
		Enemies: append([]Enemy{
			{
				Name: "John Perryman",
				Ship: perrymanEnemy,
			},
		},
			mult(20, varsityEnemy)...),
	},
}
