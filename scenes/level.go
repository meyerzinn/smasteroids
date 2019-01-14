package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/meyerzinn/smastroids/assets"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
)

const Decay = .99994

var shipPoints = []pixel.Vec{{-12, -12}, {0, 24}, {12, -12}}

type ship struct {
	pos        pixel.Vec
	vel        pixel.Vec
	acc        pixel.Vec
	angle      float64
	angularVel float64
	angularAcc float64

	//body     *cp.Body
	health             int
	label              *text.Text
	ticksSinceLastShot int
	name               string // for debug really
}

const ShipMaxSpeed = 1
const ShipMaxForce = .0005
const ShipTurnSpeed = .00001
const BulletForce = .05
const BulletDamage = 10

func (ship *ship) seek(target pixel.Vec) pixel.Vec {
	desired := target.Sub(ship.pos)
	//desired := target.Sub(ship.body.Position())
	steer := desired.Normal().Scaled(ShipMaxSpeed).Sub(ship.vel)
	if steer.Len() > ShipMaxForce {
		steer = steer.Normal().Scaled(ShipMaxForce)
	}
	return steer
}

func (ship *ship) draw(canvas *pixelgl.Canvas, imd *imdraw.IMDraw) {
	imd.Color = colornames.White
	//centroid := pixel.Vec(ship.body.Position())
	centroid := ship.pos
	for _, p := range shipPoints {
		imd.Push(centroid.Add(p.Rotated(ship.angle)))
	}
	imd.Polygon(2)
	ship.label.Draw(canvas, pixel.IM.Moved(ship.pos.Sub(pixel.V(0, 50)).Sub(ship.label.Bounds().Center())))
}
func (ship *ship) update() {
	ship.vel = ship.vel.Add(ship.acc)
	ship.acc = pixel.ZV
	ship.vel = ship.vel.Scaled(Decay)
	ship.pos = ship.pos.Add(ship.vel)
	// wraparound behavior
	if ship.pos.X < -1920/2 {
		ship.pos.X = 1920 / 2
	} else if ship.pos.X > 1920/2 {
		ship.pos.X = -1920 / 2
	}
	if ship.pos.Y < -1080/2 {
		ship.pos.Y = 1080 / 2
	} else if ship.pos.Y > 1080/2 {
		ship.pos.Y = -1080 / 2
	}
	ship.angularVel += ship.angularAcc
	ship.angularAcc = 0
	ship.angularVel *= Decay
	ship.angle += ship.angularVel
}

func (ship *ship) broadphaseTestCircle(b *bullet) bool {
	// tests if a bullet is possibly within the circle inscribing the triangle
	// qwik mafs
	return b.pos.Sub(ship.pos).Len() < 12*math.Sqrt2+4
}

func (ship *ship) nearphaseTestBullet(bullet *bullet) bool {
	// adapted from https://github.com/mattdesl/point-in-triangle/blob/master/index.js
	cx := bullet.pos.X
	cy := bullet.pos.Y
	t0 := shipPoints[0].Rotated(ship.angle).Add(ship.pos)
	t1 := shipPoints[1].Rotated(ship.angle).Add(ship.pos)
	t2 := shipPoints[2].Rotated(ship.angle).Add(ship.pos)
	v0x := t2.X - t0.X
	v0y := t2.Y - t0.Y
	v1x := t1.X - t0.X
	v1y := t1.Y - t0.Y
	v2x := cx - t0.X
	v2y := cy - t0.Y
	dot00 := v0x*v0x + v0y*v0y
	dot01 := v0x*v1x + v0y*v1y
	dot02 := v0x*v2x + v0y*v2y
	dot11 := v1x*v1x + v1y*v1y
	dot12 := v1x*v2x + v1y*v2y
	b := dot00*dot11 - dot01*dot01
	var inv float64
	if b != 0 {
		inv = 1.0 / b
	}
	u := (dot11*dot02 - dot01*dot12) * inv
	v := (dot00*dot12 - dot01*dot02) * inv
	return u >= 0 && v >= 0 && (u+v) < 1
}

func (ship *ship) broadphaseTestShip(o *ship) bool {
	return o.pos.Sub(ship.pos).Len() < 2*12*math.Sqrt2
}

func (ship *ship) cross4(s, o *ship) bool {
	// adapted from https://stackoverflow.com/a/44269990/3152168
	pa := shipPoints[0].Rotated(o.angle).Add(o.pos)
	pb := shipPoints[1].Rotated(o.angle).Add(o.pos)
	pc := shipPoints[2].Rotated(o.angle).Add(o.pos)
	p0 := shipPoints[0].Rotated(s.angle).Add(s.pos)
	p1 := shipPoints[1].Rotated(s.angle).Add(s.pos)
	p2 := shipPoints[2].Rotated(s.angle).Add(s.pos)
	dXa := pa.X - p2.X
	dYa := pa.Y - p2.Y
	dXb := pb.X - p2.X
	dYb := pb.Y - p2.Y
	dXc := pc.X - p2.X
	dYc := pc.Y - p2.Y
	dX21 := p2.X - p1.X
	dY12 := p1.Y - p2.Y
	d := dY12*(p0.X-p2.X) + dX21*(p0.Y-p2.Y)
	sa := dY12*dXa + dX21*dYa
	sb := dY12*dXb + dX21*dYb
	sc := dY12*dXc + dX21*dYc
	ta := (p2.Y-p0.Y)*dXa + (p0.X-p2.X)*dYa
	tb := (p2.Y-p0.Y)*dXb + (p0.X-p2.X)*dYb
	tc := (p2.Y-p0.Y)*dXc + (p0.X-p2.X)*dYc
	if d < 0 {
		return (sa >= 0 && sb >= 0 && sc >= 0) || (ta >= 0 && tb >= 0 && tc >= 0) || (sa+ta <= d && sb+tb <= d && sc+tc <= d)
	}
	return (sa <= 0 && sb <= 0 && sc <= 0) || (ta <= 0 && tb <= 0 && tc <= 0) || (sa+ta >= d && sb+tb >= d && sc+tc >= d)
}

func (ship *ship) nearphaseTestShip(o *ship) bool {
	return !(ship.cross4(ship, o) || ship.cross4(o, ship))
}

func newShip(name string, health int) *ship {
	var ship ship
	ship.health = health
	ship.label = text.New(pixel.ZV, assets.FontLabel)
	_, _ = ship.label.WriteString(name)
	ship.name = name
	return &ship
}

type bullet struct {
	pos pixel.Vec
	vel pixel.Vec
	acc pixel.Vec
	//body     *cp.Body
	collided bool
	alive    int
	player   bool
}

func (b *bullet) draw(imd *imdraw.IMDraw) {
	imd.Color = colornames.White
	imd.Push(pixel.Vec(b.pos))
	imd.Circle(4, 2)
}

func (b *bullet) update() {
	b.vel = b.vel.Add(b.acc).Scaled(Decay)
	b.acc = pixel.ZV
	b.pos = b.pos.Add(b.vel)
}
func newBullet(player bool) *bullet {
	var bullet bullet
	bullet.alive = 600
	bullet.player = player
	return &bullet
}

type LevelScene struct {
	levelIndex   int
	imd          *imdraw.IMDraw
	player       *ship
	enemies      []*ship
	bullets      []*bullet
	playerTarget pixel.Vec
}

func (s *LevelScene) Render(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	s.player.ticksSinceLastShot++
	if win.Pressed(pixelgl.KeySpace) {
		if s.player.ticksSinceLastShot > 60*(.5) {
			bullet := newBullet(true)
			bullet.pos = s.player.pos
			bullet.vel = s.player.vel
			bullet.acc = pixel.V(0, BulletForce).Rotated(s.player.angle)
			s.bullets = append(s.bullets, bullet)
			s.player.ticksSinceLastShot = 0
		}
	}
	for _, e := range s.enemies {
		e.ticksSinceLastShot++
		if e.ticksSinceLastShot > 60*(.5) {
			b := newBullet(false)
			b.pos = e.pos
			b.vel = e.vel
			b.acc = b.acc.Add(pixel.Vec{Y: BulletForce}.Rotated(e.angle))
			s.bullets = append(s.bullets, b)
			e.ticksSinceLastShot = 0
		}
	}
	// run the physics loop 300 times
	for step := 0; step < 300; step++ {
		if win.Pressed(pixelgl.KeyW) {
			s.player.acc = s.player.acc.Add(pixel.V(0, ShipMaxForce/300).Rotated(s.player.angle))
		}
		if win.Pressed(pixelgl.KeyA) {
			if !win.Pressed(pixelgl.KeyD) {
				s.player.angularAcc = ShipTurnSpeed / 2 / 300
			}
		} else if win.Pressed(pixelgl.KeyD) {
			s.player.angularAcc = -ShipTurnSpeed / 2 / 300
		}
		if len(s.enemies) == 0 {
			panic("next level")
		}
		if s.player.health <= 0 {
			panic("death message")
		}

		// update bullets
		for _, b := range s.bullets {
			b.update()
		} // do bullet-bullet collision testing here (simple circle radius distance check)
		for i := len(s.bullets) - 1; i >= 0; i-- {
			// clear collided bullets
			bullet := s.bullets[i]
			if bullet.alive < 0 {
				goto deleteBullet // skip expensive collision detection
			}
			for j := i - 1; j > 0; j-- {
				if bullet.player == s.bullets[j].player && bullet.pos.Sub(s.bullets[j].pos).Len() < 8 {
					bullet.collided = true
					s.bullets[j].collided = true
					break
				}
			}
			continue
		deleteBullet:
			if bullet.collided || bullet.alive < 0 {
				// delete without memory leak
				copy(s.bullets[i:], s.bullets[i+1:])
				s.bullets[len(s.bullets)-1] = nil
				s.bullets = s.bullets[:len(s.bullets)-1]
			}
		}
		// update ship physics
		s.player.update()
		for _, e := range s.enemies {
			e.update()
			if s.player.nearphaseTestShip(e) {
				b := e.health - s.player.health
				s.player.health -= e.health
				e.health = b
			}
		}
		for _, b := range s.bullets {
			if !b.player && !b.collided && s.player.nearphaseTestBullet(b) {
				s.player.health -= BulletDamage
			}
		}
		// draw the player
		s.player.draw(canvas, s.imd)
		// update position for AI to target
		s.playerTarget = pixel.Lerp(pixel.Vec(s.playerTarget), s.player.pos.Add(s.player.vel.Scaled(2)), 2.0/300)

		// update enemies and handle ship collisions
		for i := len(s.enemies) - 1; i >= 0; i-- {
			enemy := s.enemies[i]
			if enemy.broadphaseTestShip(s.player) && enemy.nearphaseTestShip(s.player) {
				s.player.health -= enemy.health
				enemy.health = 0
			}
			for _, bullet := range s.bullets {
				if bullet.player && !bullet.collided && enemy.nearphaseTestBullet(bullet) {
					enemy.health -= BulletDamage
					bullet.collided = true
				}
			}
			if enemy.health <= 0 {
				// enemy is dead
				copy(s.enemies[i:], s.enemies[i+1:])
				s.enemies[len(s.enemies)-1] = nil
				s.enemies = s.enemies[:len(s.enemies)-1]
				continue
			}
		}
	}

	//render
	canvas.Clear(colornames.Black)
	s.imd.Clear()
	for _, b := range s.bullets {
		b.draw(s.imd)
		b.alive--
	}
	for _, e := range s.enemies {
		e.draw(canvas, s.imd)
	}
	s.player.draw(canvas, s.imd)
	s.imd.Color = colornames.Red
	s.imd.Push(pixel.Vec(s.playerTarget))
	s.imd.Circle(5, 0)
	s.imd.Draw(canvas)
}

func Level(index int) Scene {
	imd := imdraw.New(nil)

	scene := &LevelScene{
		levelIndex: index,
		imd:        imd,
		player:     newShip("", 100),
	}

	level := assets.Levels[index]

	for _, t := range level.Teachers {
		teacher := assets.Teachers[t]
		ship := newShip(teacher.Name, level.Difficulty)
		ship.pos = pixel.Vec{X: 1920, Y: 1080}.Scaled(rand.Float64()).Sub(pixel.Vec{X: 1920 / 2, Y: 1080 / 2})
		scene.enemies = append(scene.enemies, ship)
	}

	return scene
}
