package game

type Health = float32

type Level struct {
	Name             string   `json:"name"`
	BulletDamage     Health   `json:"bullet_damage"`
	EnemyHealth      Health   `json:"enemy_health"`
	EnemySteerSpeed  float64  `json:"enemy_steer_speed"`
	EnemyThrustForce float64  `json:"enemy_thrust_force"`
	PlayerHealth     Health   `json:"player_health"`
	Teachers         []string `json:"teachers"`
	Index            int      `json:"-"`
}
