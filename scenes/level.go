package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/jakecoffman/cp"
	"gitlab.com/meyerzinn/smasteroids/assets"
	"gitlab.com/meyerzinn/smasteroids/smasteroids"
	"golang.org/x/image/colornames"
	"math"
	"math/rand"
	"time"
)

const (
	CollisionTypePlayer = 1 << (iota + 1)
	CollisionTypeEnemy
	CollisionTypeWall
)

// Convert pixel Vector to chipmunk Vector
func p2cp(v pixel.Vec) cp.Vector {
	return cp.Vector(v)
}

// Convert chipmunk Vector to pixel Vector
func cp2p(v cp.Vector) pixel.Vec {
	return pixel.Vec(v)
}

var (
	playerFilter = cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypePlayer), uint(CollisionTypeWall|CollisionTypeEnemy))
	enemyFilter  = cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypeEnemy), uint(CollisionTypeWall|CollisionTypePlayer))
)

var shipVertices = []cp.Vector{{-12, -12}, {0, 24}, {12, -12}}

type Ship struct {
	body           *cp.Body
	sprite         *pixel.Sprite
	data           smasteroids.Ship
	ticksSinceFire int
	health         float64
	//name           string
	//health         Health
	//health      Health
}

func (s *LevelScene) newShip(data smasteroids.Ship, enemy bool) *Ship {
	body := s.space.AddBody(cp.NewBody(1, cp.MomentForPoly(1, len(shipVertices), shipVertices, cp.Vector{}, 1)))
	shipShape := s.space.AddShape(cp.NewPolyShape(body, len(shipVertices), shipVertices, cp.NewTransformIdentity(), 1))
	if enemy {
		shipShape.SetCollisionType(CollisionTypeEnemy)
		shipShape.SetFilter(enemyFilter)
		body.SetPosition(p2cp(pixel.V(CanvasBounds.W(), CanvasBounds.H()).Scaled(rand.Float64()).Sub(CanvasBounds.Max)))
	} else {
		shipShape.SetCollisionType(CollisionTypePlayer)
		shipShape.SetFilter(playerFilter)
	}

	ship := &Ship{
		body:   body,
		sprite: pixel.NewSprite(s.shipSpriteCanvas, s.shipSpriteCanvas.Bounds()),
		data:   data,
		health: data.Health,
		//health:      health,
		//name:           name,
		//ticksSinceFire: 0,
	}
	body.UserData = &ship.health
	return ship
}

type Bullet struct {
	body   *cp.Body
	sprite *pixel.Sprite
	health float64
	alive  int
}

func (s *LevelScene) newBullet(parent *Ship, ttl int, enemy bool) *Bullet {
	body := s.space.AddBody(cp.NewBody(1, cp.MomentForCircle(1, 0, 4, cp.Vector{})))
	bulletShape := s.space.AddShape(cp.NewCircle(body, 4, cp.Vector{}))
	if enemy {
		bulletShape.SetCollisionType(CollisionTypeEnemy)
		bulletShape.SetFilter(enemyFilter)
	} else {
		bulletShape.SetCollisionType(CollisionTypePlayer)
		bulletShape.SetFilter(playerFilter)
	}
	body.SetPosition(parent.body.Position())
	body.SetVelocityVector(parent.body.Velocity())
	body.SetAngle(parent.body.Angle())
	body.SetAngularVelocity(parent.body.AngularVelocity())
	body.ApplyForceAtLocalPoint(cp.Vector{Y: 20000}, cp.Vector{})
	bullet := &Bullet{
		body:   body,
		sprite: pixel.NewSprite(s.bulletSpriteCanvas, s.bulletSpriteCanvas.Bounds()),
		health: parent.data.BulletDamage,
		alive:  ttl,
	}
	body.UserData = &bullet.health
	return bullet
}

type LevelScene struct {
	level              smasteroids.Level
	levelIndex         int
	space              *cp.Space
	player             *Ship
	playerTarget       cp.Vector
	ships              []*Ship
	labels             map[*Ship]string
	bullets            []*Bullet
	lastTick           time.Time
	canvas             *pixelgl.Canvas
	healthCanvas       *pixelgl.Canvas
	shipSpriteCanvas   *pixelgl.Canvas
	bulletSpriteCanvas *pixelgl.Canvas
	imd                *imdraw.IMDraw
	labelText          *text.Text
}

func (s *LevelScene) deleteShip(i int) {
	ship := s.ships[i]
	ship.body.EachShape(func(shape *cp.Shape) {
		s.space.RemoveShape(shape)
	})
	s.space.RemoveBody(ship.body)

	if s.player == ship {
		s.player = nil
	}

	copy(s.ships[i:], s.ships[i+1:])
	s.ships[len(s.ships)-1] = nil // or the zero value of T
	s.ships = s.ships[:len(s.ships)-1]
}

func (s *LevelScene) deleteBullet(i int) {
	bullet := s.bullets[i]
	bullet.body.EachShape(func(shape *cp.Shape) {
		s.space.RemoveShape(shape)
	})
	s.space.RemoveBody(bullet.body)

	copy(s.bullets[i:], s.bullets[i+1:])
	s.bullets[len(s.bullets)-1] = nil // or the zero value of T
	s.bullets = s.bullets[:len(s.bullets)-1]
}

func (s *LevelScene) Render(win *pixelgl.Window) {
	if s.player == nil { // player dead
		Current = Death(s.levelIndex)
		return
	}
	if len(s.ships) == 1 { // all enemies dead
		if s.levelIndex == len(smasteroids.Levels)-1 {
			Current = Win()
		} else {
			Current = TitleScene(s.levelIndex + 1)
			return
		}
	}

	// Advance physics simulation
	now := time.Now()
	dt := now.Sub(s.lastTick).Seconds()
	s.lastTick = now
	if dt > 1 { // game was paused, only update by 1 "second"
		dt = 1
	}
	s.space.Step(dt)

	// render + handle inputs
	s.canvas.Clear(colornames.Black)

	s.playerTarget = s.playerTarget.Lerp(s.player.body.Position().Add(s.player.body.Velocity()), .1)
	for i := len(s.ships) - 1; i >= 0; i-- {
		ship := s.ships[i]
		if ship.health <= 0.0 {
			s.deleteShip(i)
		} else {
			ship.ticksSinceFire++
			if ship != s.player {
				angle := s.playerTarget.Sub(ship.body.Position()).Normalize().ToAngle() - math.Pi/2
				ship.body.SetAngle(cp.Lerp(ship.body.Angle(), angle, ship.data.Turn*dt))
				ship.body.ApplyForceAtLocalPoint(cp.Vector{Y: ship.data.Thrust}, cp.Vector{})
				if ship.ticksSinceFire > ship.data.Fire {
					s.bullets = append(s.bullets, s.newBullet(ship, 100, true))
					ship.ticksSinceFire = 0
				}
				s.labelText.Clear()
				_, _ = s.labelText.WriteString(s.labels[ship])
				s.labelText.Draw(s.canvas, pixel.IM.Moved(cp2p(ship.body.Position()).Sub(pixel.V(0, 30)).Sub(s.labelText.Bounds().Center())))

			} else {
				if (win.Pressed(pixelgl.KeySpace) || win.Pressed(pixelgl.KeyEnter)) && ship.ticksSinceFire > 20 {
					s.bullets = append(s.bullets, s.newBullet(ship, 100, false))
					ship.ticksSinceFire = 0
				}
				if win.Pressed(pixelgl.KeyW) || win.Pressed(pixelgl.KeyUp) {
					ship.body.ApplyForceAtLocalPoint(cp.Vector{0, 1.5 * ship.data.Thrust}, cp.Vector{})
				}
				if win.Pressed(pixelgl.KeyA) || win.Pressed(pixelgl.KeyLeft) {
					if !win.Pressed(pixelgl.KeyD) || win.Pressed(pixelgl.KeyRight) {
						ship.body.SetAngularVelocity(cp.Lerp(ship.body.AngularVelocity(), 3, 4*dt))
					}
				} else if win.Pressed(pixelgl.KeyD) || win.Pressed(pixelgl.KeyRight) {
					ship.body.SetAngularVelocity(cp.Lerp(ship.body.AngularVelocity(), -3, 4*dt))
				} else {
					ship.body.SetAngularVelocity(cp.Lerp(ship.body.AngularVelocity(), 0, 4*dt))
				}
			}
			// draw ship
			ship.sprite.Draw(s.canvas, pixel.IM.Scaled(pixel.ZV, 1.0/4.0).Rotated(pixel.ZV, ship.body.Angle()).Moved(cp2p(ship.body.Position())))

			// draw health bar
			s.healthCanvas.Clear(colornames.Black)
			s.imd.Clear()
			if ship.health > .5*ship.data.Health {
				s.imd.Color = colornames.Green
			} else if ship.health > .3*ship.data.Health {
				s.imd.Color = colornames.Yellow
			} else {
				s.imd.Color = colornames.Red
			}
			s.imd.Push(pixel.V(0, 0), pixel.V(ship.health/ship.data.Health*64, 8))
			s.imd.Rectangle(0)
			s.imd.Draw(s.healthCanvas)
			s.healthCanvas.Draw(s.canvas, pixel.IM.Moved(cp2p(ship.body.Position()).Sub(pixel.V(0, 40)).Sub(pixel.V(0, s.healthCanvas.Bounds().H()/2))))
		}
	}

	for i := len(s.bullets) - 1; i >= 0; i-- {
		bullet := s.bullets[i]
		bullet.alive--
		if bullet.health <= 0 || bullet.alive <= 0 {
			s.deleteBullet(i)
		} else {
			bullet.sprite.Draw(s.canvas, pixel.IM.Scaled(pixel.ZV, 1.0/4.0).Rotated(pixel.ZV, bullet.body.Angle()).Moved(cp2p(bullet.body.Position())))
		}
	}

	Draw(win, s.canvas)
}

func PlayLevel(index int) *LevelScene {
	var scene LevelScene
	scene.level = smasteroids.Levels[index]
	// initialize graphics
	scene.imd = imdraw.New(nil)
	scene.healthCanvas = pixelgl.NewCanvas(pixel.R(0, 0, 64, 8))
	scene.shipSpriteCanvas = pixelgl.NewCanvas(pixel.R(-48, -48, 48, 96))
	scene.bulletSpriteCanvas = pixelgl.NewCanvas(pixel.R(-16, -16, 16, 16))
	// player sprite
	scene.imd.Color = colornames.White
	for _, v := range shipVertices {
		scene.imd.Push(cp2p(v).Scaled(4))
	}
	scene.imd.Polygon(8)
	scene.imd.Draw(scene.shipSpriteCanvas)
	scene.imd.Reset()
	// bullet sprite
	scene.imd.Color = colornames.White
	scene.imd.Push(pixel.ZV)
	scene.imd.Circle(16, 4)
	scene.imd.Draw(scene.bulletSpriteCanvas)
	scene.imd.Reset()

	canvas := pixelgl.NewCanvas(CanvasBounds)
	scene.canvas = canvas

	// initialize physics
	space := cp.NewSpace()
	scene.space = space
	space.SetGravity(cp.Vector{})
	space.SetDamping(.8)

	hw := CanvasBounds.W() / 2
	hh := CanvasBounds.H() / 2
	sides := []cp.Vector{
		{-hw, -hh}, {-hw, hh},
		{hw, -hh}, {hw, hh},
		{-hw, -hh}, {hw, -hh},
		{-hw, hh}, {hw, hh},
	}

	for i := 0; i < len(sides); i += 2 {
		seg := space.AddShape(cp.NewSegment(space.StaticBody, sides[i], sides[i+1], 1))
		seg.SetElasticity(1)
		seg.SetFriction(0)
		seg.SetCollisionType(CollisionTypeWall)
		seg.SetFilter(cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypeWall), uint(CollisionTypePlayer|CollisionTypeEnemy)))
	}

	damageHandler := space.NewCollisionHandler(CollisionTypePlayer, CollisionTypeEnemy)
	damageHandler.BeginFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		playerB, enemyB := arb.Bodies()
		player := playerB.UserData.(*float64)
		enemy := enemyB.UserData.(*float64)
		temp := *player - *enemy*.5
		*enemy = *enemy - *player*.5
		*player = temp
		return true
	}

	wrapHandler := space.NewWildcardCollisionHandler(CollisionTypeWall)
	wrapHandler.BeginFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		return arb.Ignore()
	}
	wrapHandler.SeparateFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) {
		_, body := arb.Bodies()
		body.SetPosition(wrap(body.Position()))
		if body == scene.player.body {
			scene.playerTarget = body.Position()
		}
	}

	scene.player = scene.newShip(scene.level.Player, false)
	scene.ships = append(scene.ships, scene.player)
	scene.labels = make(map[*Ship]string)

	for _, enemy := range scene.level.Enemies {
		ship := scene.newShip(enemy.Ship, true)
		scene.ships = append(scene.ships, ship)
		scene.labels[ship] = enemy.Name
	}
	scene.labelText = text.New(pixel.ZV, assets.FontLabel)

	return &scene
}

func wrap(pos cp.Vector) cp.Vector {
	if pos.X < -CanvasBounds.W()/2 {
		pos = pos.Add(cp.Vector{X: CanvasBounds.W()})
	} else if pos.X > CanvasBounds.W()/2 {
		pos = pos.Add(cp.Vector{X: -CanvasBounds.W()})
	}
	if pos.Y < -CanvasBounds.H()/2 {
		pos = pos.Add(cp.Vector{Y: CanvasBounds.H()})
	} else if pos.Y > CanvasBounds.H()/2 {
		pos = pos.Add(cp.Vector{Y: -CanvasBounds.H()})
	}
	return pos
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
