package scenes

import (
	"github.com/20zinnm/smasteroids/assets"
	"github.com/20zinnm/smasteroids/smasteroids"
	"github.com/20zinnm/smasteroids/sprites"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/jakecoffman/cp"
	"golang.org/x/image/colornames"
	"image/color"
	"math"
	"math/rand"
	"time"
)

const (
	CollisionTypePlayer = 1 << (iota + 1)
	CollisionTypeEnemy
	CollisionTypeWall
)

var Players []ControlScheme

var (
	GameBounds   = pixel.R(-1920/2, -1080/2, 1920/2, 1080/2)
	PlayerColors = []color.Color{colornames.Gold, colornames.Limegreen}

	shipCollisionPolygon = []cp.Vector{{-12, -12}, {0, 24}, {12, -12}}
)

const BoostDelay = 7 * time.Second

// Convert pixel.Vec to cp.Vector
func p2cp(v pixel.Vec) cp.Vector {
	return cp.Vector(v)
}

// Convert cp.Vector to pixel.Vec
func cp2p(v cp.Vector) pixel.Vec {
	return pixel.Vec(v)
}

var (
	playerShipFilter   = cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypePlayer), uint(CollisionTypeWall|CollisionTypeEnemy))
	playerBulletFilter = cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypePlayer), uint(CollisionTypeWall|CollisionTypeEnemy))
	enemyShipFilter    = cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypeEnemy), uint(CollisionTypeWall|CollisionTypePlayer))
	enemyBulletFilter  = cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypeEnemy), uint(CollisionTypeWall|CollisionTypePlayer))
)

type levelShip struct {
	body     *cp.Body
	sprite   *pixel.Sprite
	data     smasteroids.Ship
	lastFire time.Time
	health   float64
}

func (s *levelShip) drawHealthBar(imd *imdraw.IMDraw) {
	if s.health > .5*s.data.Health {
		imd.Color = colornames.Green
	} else if s.health > .3*s.data.Health {
		imd.Color = colornames.Yellow
	} else {
		imd.Color = colornames.Red
	}
	min := cp2p(s.body.Position()).Sub(pixel.V(32, 60))
	max := min.Add(pixel.V(s.health/s.data.Health*64, 8))
	imd.Push(min, max)
	imd.Rectangle(0)
}

func (s *levelScene) newShip(data smasteroids.Ship, enemy bool) *levelShip {
	body := s.space.AddBody(cp.NewBody(1, cp.MomentForPoly(1, len(shipCollisionPolygon), shipCollisionPolygon, cp.Vector{}, 1)))
	shipShape := s.space.AddShape(cp.NewPolyShape(body, len(shipCollisionPolygon), shipCollisionPolygon, cp.NewTransformIdentity(), 1))
	var ship = &levelShip{
		body:   body,
		sprite: pixel.NewSprite(sprites.TextureShip, sprites.TextureShip.Bounds()),
		data:   data,
	}
	if enemy {
		shipShape.SetCollisionType(CollisionTypeEnemy)
		shipShape.SetFilter(enemyShipFilter)
		body.SetPosition(p2cp(pixel.V(GameBounds.W(), GameBounds.H()).ScaledXY(pixel.V(rand.Float64(), rand.Float64())).Sub(GameBounds.Max)))
		ship.data.Health = ship.data.Health * math.Sqrt(float64(len(Players)))
		ship.health = ship.data.Health
	} else {
		shipShape.SetCollisionType(CollisionTypePlayer)
		shipShape.SetFilter(playerShipFilter)
		ship.health = data.Health
	}
	body.UserData = &ship.health
	return ship
}

type levelBullet struct {
	body    *cp.Body
	sprite  *pixel.Sprite
	health  float64
	despawn time.Time
}

func (s *levelScene) newBullet(parent *levelShip, ttl time.Duration, enemy bool) *levelBullet {
	body := s.space.AddBody(cp.NewBody(1, cp.MomentForCircle(1, 0, 4, cp.Vector{})))
	bulletShape := s.space.AddShape(cp.NewCircle(body, 4, cp.Vector{}))
	var bullet = &levelBullet{
		body:    body,
		sprite:  pixel.NewSprite(sprites.TextureBullet, sprites.TextureBullet.Bounds()),
		despawn: time.Now().Add(ttl),
	}
	if enemy {
		bulletShape.SetCollisionType(CollisionTypeEnemy)
		bulletShape.SetFilter(enemyBulletFilter)
		bullet.health = parent.data.BulletDamage * math.Sqrt(float64(len(Players)))
	} else {
		bulletShape.SetCollisionType(CollisionTypePlayer)
		bulletShape.SetFilter(playerBulletFilter)
		bullet.health = parent.data.BulletDamage
	}
	body.SetPosition(parent.body.Position())
	body.SetVelocityVector(parent.body.Velocity())
	body.SetAngle(parent.body.Angle())
	body.SetAngularVelocity(parent.body.AngularVelocity())
	body.ApplyImpulseAtLocalPoint(cp.Vector{Y: 300}, cp.Vector{})
	parent.body.ApplyImpulseAtLocalPoint(cp.Vector{Y: -5}, cp.Vector{})
	body.UserData = &bullet.health
	return bullet
}

type Controls struct {
	Left, Right, Thrust, Boost, Shoot bool
}

type ControllerInput interface {
	Controls(win *pixelgl.Window) Controls
}

type Player struct {
	Ship           *levelShip
	TargetPosition cp.Vector
	LastBoost      time.Time
	Color          color.Color
	ControlScheme
}

type levelScene struct {
	level             smasteroids.Level
	levelIndex        int
	space             *cp.Space
	players           []*Player
	enemies           []*levelShip
	labels            map[*levelShip]string
	bullets           []*levelBullet
	lastTick          time.Time
	canvas            *pixelgl.Canvas
	imd               *imdraw.IMDraw
	labelText         *text.Text
	shipSpriteBatch   *pixel.Batch
	bulletSpriteBatch *pixel.Batch
}

func (s *levelScene) deleteEnemy(i int) {
	ship := s.enemies[i]
	ship.body.EachShape(func(shape *cp.Shape) {
		s.space.RemoveShape(shape)
	})
	s.space.RemoveBody(ship.body)

	copy(s.enemies[i:], s.enemies[i+1:])
	s.enemies[len(s.enemies)-1] = nil
	s.enemies = s.enemies[:len(s.enemies)-1]
}

func (s *levelScene) deleteBullet(i int) {
	bullet := s.bullets[i]
	bullet.body.EachShape(func(shape *cp.Shape) {
		s.space.RemoveShape(shape)
	})
	s.space.RemoveBody(bullet.body)

	copy(s.bullets[i:], s.bullets[i+1:])
	s.bullets[len(s.bullets)-1] = nil
	s.bullets = s.bullets[:len(s.bullets)-1]
}

func (s *levelScene) deletePlayer(i int) {
	player := s.players[i]
	ship := player.Ship
	ship.body.EachShape(func(shape *cp.Shape) {
		s.space.RemoveShape(shape)
	})
	s.space.RemoveBody(ship.body)

	copy(s.players[i:], s.players[i+1:])
	s.players[len(s.players)-1] = nil // or the zero value of T
	s.players = s.players[:len(s.players)-1]
}

func (s *levelScene) getPlayer(ship *levelShip) (*Player, bool) {
	p, ok := s.getPlayerFromBody(ship.body)
	return p, ok
}

func (s *levelScene) getPlayerFromBody(body *cp.Body) (*Player, bool) {
	for _, p := range s.players {
		if p.Ship.body == body {
			return p, true
		}
	}
	return nil, false
}

func (s *levelScene) Render(win *pixelgl.Window) {
	// Go to death screen if all players are dead.
	if len(s.players) == 0 {
		TransitionTo(Death(s.levelIndex))
		return
	}

	// Next level if all enemies are dead or cheatcode is pressed.
	if len(s.enemies) == 0 || (win.Pressed(pixelgl.KeyJ) && win.Pressed(pixelgl.KeyA) && win.Pressed(pixelgl.KeyN) && win.Pressed(pixelgl.KeyK)) { // all enemies dead or cheatcode active
		if s.levelIndex == len(smasteroids.Levels)-1 {
			TransitionTo(NewWin())
			return
		} else {
			TransitionTo(NewTitleScene(s.levelIndex + 1))
			return
		}
	}

	// Advance physics simulation.
	now := time.Now()
	dt := now.Sub(s.lastTick).Seconds()
	s.lastTick = now
	if dt > 1.0/10 { // game is lagging, only update by 100ms
		dt = 1.0 / 10
	}
	s.space.Step(dt)

	// Render the scene.
	s.canvas.Clear(colornames.Black)
	s.imd.Clear()
	s.shipSpriteBatch.Clear()
	s.bulletSpriteBatch.Clear()

	// Update and draw players.
	for i := len(s.players) - 1; i >= 0; i-- {
		player := s.players[i]
		// Delete dead players.
		if player.Ship.health <= 0 {
			s.deletePlayer(i)
			continue
		}
		ship := player.Ship
		// Update target position for AI.
		player.TargetPosition = player.TargetPosition.Lerp(ship.body.Position().Add(ship.body.Velocity()), dt)
		// Fetch player controls.
		controls := player.Controls(win)
		// Shoot bullet.
		if controls.Shoot && ship.lastFire.Add(ship.data.Fire).Before(now) {
			s.bullets = append(s.bullets, s.newBullet(ship, 3*time.Second, false))
			ship.lastFire = now
		}
		// Apply thrust.
		if controls.Thrust {
			ship.body.ApplyForceAtLocalPoint(cp.Vector{Y: 1.5 * ship.data.Thrust}, cp.Vector{})
		}
		// Apply rotation.
		if controls.Left && !controls.Right {
			//ship.body.ApplyForceAtLocalPoint(cp.Vector{X: ship.data.Turn}, cp.Vector{Y: 1})
			ship.body.SetAngularVelocity(cp.Lerp(ship.body.AngularVelocity(), ship.data.Turn, 3*dt))
		} else if controls.Right {
			//ship.body.ApplyForceAtLocalPoint(cp.Vector{X: -ship.data.Turn}, cp.Vector{Y: 1})
			ship.body.SetAngularVelocity(cp.Lerp(ship.body.AngularVelocity(), -ship.data.Turn, 3*dt))
			//} else {
			//ship.body.SetAngularVelocity(cp.Lerp(ship.body.AngularVelocity(), 0, 3*dt))
		}
		// Apply boost.
		if controls.Boost && player.LastBoost.Add(BoostDelay).Before(now) {
			player.LastBoost = now
			ship.body.ApplyImpulseAtLocalPoint(cp.Vector{Y: 200}, cp.Vector{})
		}
		// DrawCanvas sprite--scaled down for better resolution.
		ship.sprite.DrawColorMask(s.shipSpriteBatch, pixel.IM.Scaled(pixel.ZV, 1.0/4.0).Rotated(pixel.ZV, ship.body.Angle()).Moved(cp2p(ship.body.Position())), player.Color)
		// DrawCanvas health bar.
		ship.drawHealthBar(s.imd)
	}

	// Update and draw enemies.
	for i := len(s.enemies) - 1; i >= 0; i-- {
		ship := s.enemies[i]
		if ship.health <= 0 {
			s.deleteEnemy(i)
			continue
		}
		// find closest player target to shoot at
		var target cp.Vector
		var lengthSq = math.MaxFloat64
		for _, p := range s.players {
			potentialTarget := p.TargetPosition.Sub(ship.body.Position())
			if potentialTarget.LengthSq() < lengthSq {
				target = potentialTarget
				lengthSq = potentialTarget.LengthSq()
			}
		}
		angle := target.Normalize().ToAngle() - math.Pi/2
		//ship.body.ApplyForceAtLocalPoint(cp.V)
		ship.body.SetAngle(cp.Lerp(ship.body.Angle(), angle, ship.data.Turn*dt))
		ship.body.ApplyForceAtLocalPoint(cp.Vector{Y: ship.data.Thrust}, cp.Vector{})
		if ship.lastFire.Add(ship.data.Fire).Before(now) {
			s.bullets = append(s.bullets, s.newBullet(ship, 3*time.Second, true))
			ship.lastFire = now
		}

		// DrawCanvas sprite--scaled down for better resolution.
		ship.sprite.Draw(s.shipSpriteBatch, pixel.IM.Scaled(pixel.ZV, 1.0/4.0).Rotated(pixel.ZV, ship.body.Angle()).Moved(cp2p(ship.body.Position())))

		s.labelText.Clear()
		_, _ = s.labelText.WriteString(s.labels[ship])
		s.labelText.Draw(s.canvas, pixel.IM.Moved(cp2p(ship.body.Position()).Sub(pixel.V(0, 30)).Sub(s.labelText.Bounds().Center())))

		ship.drawHealthBar(s.imd)
	}

	s.imd.Draw(s.canvas)

	// DrawCanvas bullets.
	for i := len(s.bullets) - 1; i >= 0; i-- {
		bullet := s.bullets[i]
		if bullet.health <= 0 || bullet.despawn.Before(now) {
			s.deleteBullet(i)
		} else {
			bullet.sprite.Draw(s.bulletSpriteBatch, pixel.IM.Scaled(pixel.ZV, 1.0/4.0).Rotated(pixel.ZV, bullet.body.Angle()).Moved(cp2p(bullet.body.Position())))
		}
	}
	s.shipSpriteBatch.Draw(s.canvas)
	s.bulletSpriteBatch.Draw(s.canvas)
	// Render scene to the window.
	DrawCanvas(win, s.canvas)
}

func NewLevelScene(index int) Scene {
	var scene levelScene
	scene.level = smasteroids.Levels[index]
	scene.levelIndex = index

	// initialize graphics
	scene.imd = imdraw.New(nil)
	//scene.healthCanvas = pixelgl.NewCanvas(pixel.R(0, 0, 64, 8))

	scene.canvas = pixelgl.NewCanvas(GameBounds)

	// initialize physics
	space := cp.NewSpace()
	scene.space = space
	space.SetGravity(cp.Vector{})
	space.SetDamping(.75)

	hw := GameBounds.W() / 2
	hh := GameBounds.H() / 2
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
		seg.SetFilter(cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypeWall), uint(cp.WILDCARD_COLLISION_TYPE)))
		seg.SetSensor(true)
	}

	dmg := space.NewCollisionHandler(CollisionTypePlayer, CollisionTypeEnemy)
	dmg.SeparateFunc = func(arb *cp.Arbiter, _ *cp.Space, _ interface{}) {
		a, b := arb.Bodies()
		healthA := a.UserData.(*float64)
		healthB := b.UserData.(*float64)
		temp := *healthA - math.Max(*healthB*.5, 0)
		*healthB = *healthB - math.Max(*healthA*.5, 0)
		*healthA = temp
	}

	wrapHandler := space.NewWildcardCollisionHandler(CollisionTypeWall)
	wrapHandler.BeginFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		return arb.Ignore()
	}
	wrapHandler.SeparateFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) {
		_, body := arb.Bodies()
		body.SetPosition(wrap(body.Position()))
		player, ok := scene.getPlayerFromBody(body)
		if ok {
			player.TargetPosition = player.TargetPosition.Add(player.Ship.body.Velocity())
		}
	}

	for i, p := range Players {
		scene.players = append(scene.players, &Player{
			Ship:          scene.newShip(scene.level.Player, false),
			Color:         PlayerColors[i],
			ControlScheme: p,
		})
	}

	scene.labels = make(map[*levelShip]string)

	for _, enemy := range scene.level.Enemies {
		ship := scene.newShip(enemy.Ship, true)
		scene.enemies = append(scene.enemies, ship)
		scene.labels[ship] = enemy.Name
	}
	scene.labelText = text.New(pixel.ZV, assets.FontLabel)

	scene.shipSpriteBatch = pixel.NewBatch(&pixel.TrianglesData{}, sprites.TextureShip)
	scene.bulletSpriteBatch = pixel.NewBatch(&pixel.TrianglesData{}, sprites.TextureBullet)

	return &scene
}

func wrap(pos cp.Vector) cp.Vector {
	if pos.X < -GameBounds.W()/2 {
		pos = pos.Add(cp.Vector{X: GameBounds.W()})
	} else if pos.X > GameBounds.W()/2 {
		pos = pos.Add(cp.Vector{X: -GameBounds.W()})
	}
	if pos.Y < -GameBounds.H()/2 {
		pos = pos.Add(cp.Vector{Y: GameBounds.H()})
	} else if pos.Y > GameBounds.H()/2 {
		pos = pos.Add(cp.Vector{Y: -GameBounds.H()})
	}
	return pos
}
