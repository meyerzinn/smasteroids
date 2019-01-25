package smasteroids

import "time"

var mondayEnemy = Ship{
	Health:       10,
	Thrust:       50,
	Turn:         2,
	Fire:         time.Second / 2,
	BulletDamage: 2,
}

var tuesdayEnemy = Ship{
	Health:       15,
	Thrust:       60,
	Turn:         2.25,
	Fire:         time.Second / 2,
	BulletDamage: 4,
}

var wednesdayEnemy = Ship{
	Health:       20,
	Thrust:       75,
	Turn:         3,
	Fire:         time.Millisecond * 417,
	BulletDamage: 6,
}

var matrimonialEnemy = Ship{
	Health:       30,
	Thrust:       80,
	Turn:         3.25,
	Fire:         time.Second / 2,
	BulletDamage: 8,
}

var parentsEnemy = Ship{
	Health:       75,
	Thrust:       75,
	Turn:         3.5,
	Fire:         time.Second / 3 * 2,
	BulletDamage: 15,
}

var minorParentsEnemy = Ship{
	Health:       25,
	Thrust:       50,
	Turn:         3,
	Fire:         time.Second / 10,
	BulletDamage: 4,
}

//var loraxEnemy = Ship{
//	Health:       75,
//	Thrust:       50,
//	Turn:         3.5,
//	Fire:         time.Second / 3 * 2,
//	BulletDamage: 15,
//}
//
//var physicsEnemy = Ship{
//	Health:       50,
//	Thrust:       100,
//	Turn:         3.75,
//	Fire:         time.Second / 3,
//	BulletDamage: 15,
//}
//
var promisedLandEnemy = Ship{
	Health:       100,
	Thrust:       120,
	Turn:         4.25,
	Fire:         time.Second / 50,
	BulletDamage: 10,
}

var Levels []Level = []Level{
	{
		Name: "Monday",
		Player: Ship{
			Health:       20,
			Thrust:       100,
			Turn:         2.5,
			Fire:         time.Second / 3,
			BulletDamage: 5,
		},
		Enemies: []Enemy{
			{
				Name: "Jamie",
				Ship: mondayEnemy,
			},
			{
				Name: "Jackson",
				Ship: mondayEnemy,
			},
			{
				Name: "Ajay",
				Ship: mondayEnemy,
			},
			{
				Name: "Karen",
				Ship: mondayEnemy,
			},
		},
	},
	{
		Name: "Tuesday",
		Player: Ship{
			Health:       30,
			Thrust:       120,
			Turn:         2.75,
			Fire:         time.Second / 3,
			BulletDamage: 7.5,
		},
		Enemies: []Enemy{
			{
				Name: "Caroline",
				Ship: tuesdayEnemy,
			},
			{
				Name: "Mason",
				Ship: tuesdayEnemy,
			},
			{
				Name: "Madison",
				Ship: tuesdayEnemy,
			},
			{
				Name: "Emma the Younger",
				Ship: Ship{
					Health: 5, Thrust: 100, Turn: 3, Fire: time.Second, BulletDamage: 6,
				},
			},
			{
				Name: "Shayle Cruz",
				Ship: tuesdayEnemy,
			},
			{
				Name: "Lahari",
				Ship: tuesdayEnemy,
			},
		},
	},
	{
		Name: "Wednesday",
		Player: Ship{
			Health:       60,
			Thrust:       130,
			Turn:         3,
			Fire:         time.Second / 4,
			BulletDamage: 10,
		},
		Enemies: []Enemy{
			{
				Name: "Lahari but again",
				Ship: wednesdayEnemy,
			},
			{
				Name: "Meyer",
				Ship: wednesdayEnemy,
			},
			{
				Name: "Mia",
				Ship: wednesdayEnemy,
			},
			{
				Name: "Mikah",
				Ship: wednesdayEnemy,
			},
			{
				Name: "Simone",
				Ship: wednesdayEnemy,
			},
			{
				Name: "Alice",
				Ship: wednesdayEnemy,
			},
			{
				Name: "Emma the Elder",
				Ship: wednesdayEnemy,
			},
		},
	},
	{
		Name: "The Realm of Matrimonial Incongruity",
		Player: Ship{
			Health:       60,
			Thrust:       140,
			Turn:         3,
			Fire:         time.Second / 3,
			BulletDamage: 10,
		},
		Enemies: []Enemy{
			{Name: "Fiona", Ship: matrimonialEnemy},
			{Name: "Juhi", Ship: matrimonialEnemy},
			{Name: "Hailey", Ship: matrimonialEnemy},
			{Name: "Helena", Ship: matrimonialEnemy},
			{Name: "Maisy", Ship: matrimonialEnemy},
		},
	},
	{
		Name: "Meet the Parents",
		Player: Ship{
			Health:       80,
			Thrust:       150,
			Turn:         3,
			Fire:         time.Second / 3,
			BulletDamage: 15,
		},
		Enemies: []Enemy{
			{
				Name: "Ello",
				Ship: minorParentsEnemy,
			},
			{
				Name: "Wheezy",
				Ship: minorParentsEnemy,
			},
			{
				Name: "Sinny",
				Ship: minorParentsEnemy,
			},
			{
				Name: "Eloise",
				Ship: parentsEnemy,
			},
			{
				Name: "Annie",
				Ship: parentsEnemy,
			},
			{
				Name: "Gabe",
				Ship: parentsEnemy,
			},
			{
				Name: "Carolina",
				Ship: parentsEnemy,
			},
		},
	},
	{
		Name: "Props",
		Player: Ship{
			Health:       100,
			Thrust:       150,
			Turn:         3,
			Fire:         time.Second / 6,
			BulletDamage: 5,
		},
		Enemies: append(mult(10, Enemy{
			Name: "Prop",
			Ship: Ship{
				Health:       100,
				Thrust:       20,
				Turn:         3,
				Fire:         time.Second * 5,
				BulletDamage: 20,
			},
		}), Enemy{
			Name: "Catherine",
			Ship: Ship{
				Health:       50,
				Thrust:       100,
				Turn:         3.75,
				Fire:         time.Second / 3,
				BulletDamage: 15,
			},
		}),
	},
	{
		"The Promised Land",
		Ship{

		},
		[]Enemy{
			{
				Name: "Kallos",
				Ship: promisedLandEnemy,
			},
			{
				Name: "Gray",
				Ship: promisedLandEnemy,
			},
		},
	},
}
