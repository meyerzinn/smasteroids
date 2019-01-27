package smasteroids

import "time"

var seniorShip = Ship{
	Health:       30,
	Thrust:       80,
	Turn:         3.25,
	Fire:         time.Second / 2,
	BulletDamage: 8,
}

var minorSeniorShip = Ship{
	Health:       15,
	Thrust:       100,
	Turn:         3,
	Fire:         time.Second * 2,
	BulletDamage: 8,
}

var promisedLandEnemy = Ship{
	Health:       100,
	Thrust:       120,
	Turn:         4.25,
	Fire:         time.Second / 50,
	BulletDamage: 2,
}

var Levels = []Level{
	{
		Name: "Freshmen",
		Player: Ship{
			Health:       20,
			Thrust:       110,
			Turn:         2.5,
			Fire:         time.Second / 3,
			BulletDamage: 5,
		},
		Enemies: multiple(Ship{
			Health:       10,
			Thrust:       50,
			Turn:         1.75,
			Fire:         time.Second,
			BulletDamage: 2,
		}, "Jamesithati", "Madison", "Madeline", "Caroline"),
	},
	{
		Name: "Sophomores",
		Player: Ship{
			Health:       30,
			Thrust:       120,
			Turn:         2.75,
			Fire:         time.Second / 3,
			BulletDamage: 7.5,
		},
		Enemies: append(multiple(Ship{
			Health:       15,
			Thrust:       60,
			Turn:         2.25,
			Fire:         time.Second / 2,
			BulletDamage: 4,
		}, "Lahari", "Jamie", "Ajay", "Joy", "Mia", "Mikah"),
			Enemy{
				"Emma the Younger",
				Ship{
					Health: 5, Thrust: 100, Turn: 3, Fire: time.Second, BulletDamage: 6,
				}}),
	},
	{
		Name: "Juniors",
		Player: Ship{
			Health:       150,
			Thrust:       130,
			Turn:         3,
			Fire:         time.Second / 4,
			BulletDamage: 8,
		},
		Enemies: multiple(Ship{
			Health:       20,
			Thrust:       75,
			Turn:         3,
			Fire:         time.Millisecond * 417,
			BulletDamage: 6,
		}, "Abby", "Helena", "Alice", "Catherine", "Emma the Older", "Fiona", "Hailey", "Judson", "Faraz", "Jackson", "Juhi", "Karen", "Maisy", "Mason", "Meyer", "Simone"),
	},
	{
		Name: "Senior(s)",
		Player: Ship{
			Health:       60,
			Thrust:       140,
			Turn:         3,
			Fire:         time.Second / 3,
			BulletDamage: 10,
		},
		Enemies: append(
			multiple(seniorShip, "Eloise", "Meghna", "Gabe"),
			multiple(minorSeniorShip, "Ello", "Sinny", "Wheezy")...,
		),
	},
	{
		"God",
		Ship{
			Health:       100,
			Thrust:       140,
			Turn:         3,
			Fire:         time.Second / 10,
			BulletDamage: 5,
		},
		multiple(promisedLandEnemy, "Kallos", "Annie", "Gray"),
	},
}
