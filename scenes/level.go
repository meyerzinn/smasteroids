package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/jakecoffman/cp"
	"github.com/meyerzinn/smasteroids/assets"
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

const (
	ShipThrustForce = 100
	ShipTurnSpeed   = 3
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
	name           string
	health         Health
	maxHealth      Health
	ticksSinceFire int
}

func (s *LevelScene) newShip(maxHealth Health, name string, enemy bool) *Ship {
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
		body:           body,
		sprite:         pixel.NewSprite(s.shipSpriteCanvas, s.shipSpriteCanvas.Bounds()),
		health:         maxHealth,
		maxHealth:      maxHealth,
		name:           name,
		ticksSinceFire: 0,
	}
	body.UserData = &ship.health
	return ship
}

type Bullet struct {
	body   *cp.Body
	sprite *pixel.Sprite
	health Health
	alive  int
}

func (s *LevelScene) newBullet(parent *Ship, health float32, ttl int, enemy bool) *Bullet {
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
	body.ApplyForceAtLocalPoint(cp.Vector{0, 20000}, cp.Vector{})
	bullet := &Bullet{
		body:   body,
		sprite: pixel.NewSprite(s.bulletSpriteCanvas, s.bulletSpriteCanvas.Bounds()),
		health: Health(health),
		alive:  ttl,
	}
	body.UserData = &bullet.health
	return bullet
}

type LevelScene struct {
	levelIndex         int
	space              *cp.Space
	player             *Ship
	playerTarget       cp.Vector
	ships              []*Ship
	bullets            []*Bullet
	lastTick           time.Time
	canvas             *pixelgl.Canvas
	healthCanvas       *pixelgl.Canvas
	shipSpriteCanvas   *pixelgl.Canvas
	bulletSpriteCanvas *pixelgl.Canvas
	imd                *imdraw.IMDraw
	label              *text.Text
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
		Current = Death()
		return
	}
	if len(s.ships) == 1 { // all enemies dead
		if s.levelIndex == len(assets.Levels)-1 {
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
		if ship.health <= 0 {
			s.deleteShip(i)
		} else {
			ship.ticksSinceFire++
			if ship != s.player {
				angle := s.playerTarget.Sub(ship.body.Position()).Normalize().ToAngle() - math.Pi/2
				ship.body.SetAngle(cp.Lerp(ship.body.Angle(), angle, ShipTurnSpeed*dt))
				ship.body.ApplyForceAtLocalPoint(cp.Vector{Y: ShipThrustForce}, cp.Vector{})
				if ship.ticksSinceFire > 30 {
					s.bullets = append(s.bullets, s.newBullet(ship, 20, 100, true))
					ship.ticksSinceFire = 0
				}
			} else {
				if (win.Pressed(pixelgl.KeySpace)) && ship.ticksSinceFire > 20 {
					s.bullets = append(s.bullets, s.newBullet(ship, 20, 100, false))
					ship.ticksSinceFire = 0
				}
				if win.Pressed(pixelgl.KeyW) || win.Pressed(pixelgl.KeyUp) {
					ship.body.ApplyForceAtLocalPoint(cp.Vector{0, 1.5 * ShipThrustForce}, cp.Vector{})
				}
				if win.Pressed(pixelgl.KeyA) || win.Pressed(pixelgl.KeyLeft) {
					if !win.Pressed(pixelgl.KeyD) || win.Pressed(pixelgl.KeyRight) {
						ship.body.SetAngularVelocity(ShipTurnSpeed*3/4)
					}
				} else if win.Pressed(pixelgl.KeyD) || win.Pressed(pixelgl.KeyRight) {
					ship.body.SetAngularVelocity(-ShipTurnSpeed*3/4)
				}
			}
			ship.sprite.Draw(s.canvas, pixel.IM.Scaled(pixel.ZV, 1.0/4.0).Rotated(pixel.ZV, ship.body.Angle()).Moved(cp2p(ship.body.Position())))
			s.label.Clear()
			s.label.WriteString(ship.name)
			s.label.Draw(s.canvas, pixel.IM.Moved(cp2p(ship.body.Position()).Sub(pixel.V(0, 30)).Sub(s.label.Bounds().Center())))

			s.healthCanvas.Clear(colornames.Black)
			s.imd.Clear()
			if ship.health > .5*ship.maxHealth {
				s.imd.Color = colornames.Green
			} else if ship.health > .3*ship.maxHealth {
				s.imd.Color = colornames.Yellow
			} else {
				s.imd.Color = colornames.Red
			}
			s.imd.Push(pixel.V(0, 0), pixel.V(float64(ship.health)/float64(ship.maxHealth)*64, 8))
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

type Health float32

func PlayLevel(index int) *LevelScene {
	var scene LevelScene
	scene.levelIndex = index
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
		player := playerB.UserData.(*Health)
		enemy := enemyB.UserData.(*Health)
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

	// initialize ships
	leveldata := assets.Levels[index]

	scene.player = scene.newShip(Health(float64(leveldata.Difficulty)*float64(4)), "Student", false)
	scene.ships = append(scene.ships, scene.player)

	for _, t := range leveldata.Teachers {
		teacher := assets.Teachers[t]
		scene.ships = append(scene.ships, scene.newShip(Health(leveldata.Difficulty), teacher, true))
	}

	scene.label = text.New(pixel.ZV, assets.FontLabel)
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

//func newShip(space *cp.Space, ) *Ship {
//
//}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
