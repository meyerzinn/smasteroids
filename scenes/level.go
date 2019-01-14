package scenes

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/jakecoffman/cp"
	"github.com/meyerzinn/smastroids/assets"
	"golang.org/x/image/colornames"
	"math/rand"
	"time"
)

const ShootingCooldown = time.Second / 2

const (
	CollisionTypePlayer cp.CollisionType = 1 << (iota + 1)
	CollisionTypeEnemy
	CollisionTypeWall
	CollisionTypeBullet
	CollisionTypeShip
)

var shipPoints = []pixel.Vec{{-12, -12}, {0, 24}, {12, -12}}
var cpShipPoints = []cp.Vector{{-12, -12}, {0, 24}, {12, -12}}

type ship struct {
	body     *cp.Body
	health   int
	label    *text.Text
	lastShot time.Time
	name     string // for debug really
}

const ShipMaxSpeed = 15.0
const ShipMaxForce = 4.0

func (s *ship) seek(target cp.Vector) cp.Vector {
	desired := target.Sub(s.body.Position())
	steer := desired.Normalize().Mult(ShipMaxSpeed).Sub(s.body.Velocity())
	if steer.LengthSq() > ShipMaxForce*ShipMaxForce {
		steer = steer.Normalize().Mult(ShipMaxForce)
	}
	return steer
}

func (s *ship) draw(canvas *pixelgl.Canvas, imd *imdraw.IMDraw) {
	imd.Color = colornames.White
	centroid := pixel.Vec(s.body.Position())
	for _, p := range shipPoints {
		imd.Push(centroid.Add(p.Rotated(s.body.Angle())))
	}
	imd.Polygon(2)
	s.label.Draw(canvas, pixel.IM.Moved(pixel.Vec(s.body.Position()).Sub(pixel.V(0, 50)).Sub(s.label.Bounds().Center())))
}

func newShip(space *cp.Space, name string, health int) *ship {
	var ship ship
	body := space.AddBody(cp.NewBody(1, cp.MomentForPoly(1, len(cpShipPoints), cpShipPoints, cp.Vector{}, 1)))
	shape := space.AddShape(cp.NewPolyShape(body, len(cpShipPoints), cpShipPoints, cp.NewTransformIdentity(), 1))
	if name == "" {
		shape.SetCollisionType(CollisionTypePlayer | CollisionTypeShip)
		shape.SetFilter(cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypePlayer|CollisionTypeShip), uint(CollisionTypeEnemy|CollisionTypeWall)))
	} else {
		shape.SetCollisionType(CollisionTypeEnemy | CollisionTypeShip)
		shape.SetFilter(cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypeEnemy|CollisionTypeShip), uint(CollisionTypePlayer|CollisionTypeWall)))
	}
	body.UserData = &ship
	ship.body = body
	ship.health = health
	ship.label = text.New(pixel.ZV, assets.FontLabel)
	_, _ = ship.label.WriteString(name)
	ship.name = name
	return &ship
}

type bullet struct {
	body     *cp.Body
	collided bool
	spawned  time.Time
}

func (b *bullet) draw(imd *imdraw.IMDraw) {
	imd.Color = colornames.White
	imd.Push(pixel.Vec(b.body.Position()))
	imd.Circle(4, 2)
}

func newBullet(space *cp.Space, player bool) *bullet {
	var bullet bullet
	body := space.AddBody(cp.NewBody(1, cp.MomentForCircle(1, 0, 4, cp.Vector{})))
	shape := space.AddShape(cp.NewCircle(body, 4, cp.Vector{}))
	if player {
		shape.SetCollisionType(CollisionTypePlayer | CollisionTypeBullet)
		shape.SetFilter(cp.NewShapeFilter(0, uint(CollisionTypePlayer|CollisionTypeBullet), uint(CollisionTypeEnemy|CollisionTypeWall)))
	} else {
		shape.SetCollisionType(CollisionTypeEnemy | CollisionTypeBullet)
		shape.SetFilter(cp.NewShapeFilter(0, uint(CollisionTypeEnemy|CollisionTypeBullet), uint(CollisionTypePlayer|CollisionTypeWall)))
	}
	body.UserData = &bullet
	bullet.body = body
	bullet.spawned = time.Now()
	return &bullet
}

type LevelScene struct {
	levelIndex   int
	space        *cp.Space
	last         time.Time
	imd          *imdraw.IMDraw
	player       *ship
	enemies      []*ship
	bullets      []*bullet
	playerTarget cp.Vector
	walls        []*cp.Shape
}

func (s *LevelScene) Render(win *pixelgl.Window, canvas *pixelgl.Canvas) {
	now := time.Now()
	dt := now.Sub(s.last).Seconds()
	if dt <= 0 {
		return
	}

	if win.Pressed(pixelgl.KeyW) {
		s.player.body.ApplyForceAtLocalPoint(cp.Vector{Y: 100}, cp.Vector{})
	}

	if win.Pressed(pixelgl.KeyA) {
		if !win.Pressed(pixelgl.KeyD) {
			s.player.body.SetAngularVelocity(ShipMaxForce / 2)
		}
	} else if win.Pressed(pixelgl.KeyD) {
		s.player.body.SetAngularVelocity(-ShipMaxForce / 2)
	}

	s.space.Step(dt)
	s.last = now
	if len(s.enemies) == 0 {
		panic("next level")
	}
	if s.player.health <= 0 {
		panic("death message")
	}
	canvas.Clear(colornames.Black)
	s.imd.Clear()
	// update bullets
	for i := len(s.bullets) - 1; i >= 0; i-- {
		// clear collided bullets
		if s.bullets[i].collided || now.Sub(s.bullets[i].spawned) > 5*time.Second {
			s.space.RemoveBody(s.bullets[i].body)
			// delete without memory leak
			copy(s.bullets[i:], s.bullets[i+1:])
			s.bullets[len(s.bullets)-1] = nil // or the zero value of T
			s.bullets = s.bullets[:len(s.bullets)-1]
			continue
		}
		s.bullets[i].draw(s.imd)
	}
	// draw the player
	s.player.draw(canvas, s.imd)
	// update position for AI to target
	s.playerTarget = cp.Vector(pixel.Lerp(pixel.Vec(s.playerTarget), pixel.Vec(s.player.body.Position()).Add(pixel.Vec(s.player.body.Velocity()).Scaled(2)), dt*1.5))

	// update enemies
	for i := len(s.enemies) - 1; i >= 0; i-- {
		if s.enemies[i].health <= 0 {
			// enemy is dead
			s.space.RemoveBody(s.enemies[i].body)
			copy(s.enemies[i:], s.enemies[i+1:])
			s.enemies[len(s.enemies)-1] = nil
			s.enemies = s.enemies[:len(s.enemies)-1]
			continue
		}
		force := s.enemies[i].seek(s.playerTarget)
		s.enemies[i].body.ApplyForceAtLocalPoint(cp.Vector{Y: 10}, cp.Vector{})
		s.enemies[i].body.SetAngle(cp.LerpConst(s.enemies[i].body.Angle(), force.ToAngle(), .1*dt))

		if now.Sub(s.enemies[i].lastShot) > ShootingCooldown {
			b := newBullet(s.space, false)
			b.body.SetPosition(s.enemies[i].body.Position())
			b.body.SetAngle(s.enemies[i].body.Angle())
			b.body.ApplyForceAtLocalPoint(cp.Vector{Y: 500}, cp.Vector{})
			s.bullets = append(s.bullets, b)
			s.enemies[i].lastShot = now
		}

		s.enemies[i].draw(canvas, s.imd)
	}
	s.imd.Color = colornames.Red
	s.imd.Push(pixel.Vec(s.playerTarget))
	s.imd.Circle(5, 0)
	s.imd.Draw(canvas)
}

func (s *LevelScene) handleShipBulletCollision(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
	shipBody, bulletBody := arb.Bodies()
	ship := shipBody.UserData.(*ship)
	bullet := bulletBody.UserData.(*bullet)
	if bullet.collided {
		return false
	}
	bullet.collided = true
	ship.health -= 20
	return true
}

func Level(index int) Scene {
	space := cp.NewSpace()
	space.SetGravity(cp.Vector{})
	space.SetDamping(.8)
	sides := []cp.Vector{
		{X: -1920 / 2, Y: -1080 / 2}, {X: -1920 / 2, Y: 1080 / 2},
		{X: 1920 / 2, Y: -1080 / 2}, {X: 1920 / 2, Y: 1080 / 2},
		{X: -1920 / 2, Y: -1080 / 2}, {X: 1920 / 2, Y: -1080 / 2},
		{X: -1920 / 2, Y: 1080 / 2}, {X: 1920 / 2, Y: 1080 / 2},
	}
	wallFilter := cp.NewShapeFilter(cp.NO_GROUP, uint(CollisionTypeWall), uint(CollisionTypeShip|CollisionTypeBullet))
	for i := 0; i < len(sides); i += 2 {
		seg := space.AddShape(cp.NewSegment(space.StaticBody, sides[i], sides[i+1], 1))
		seg.SetElasticity(1)
		seg.SetFriction(0)
		seg.SetFilter(wallFilter)
	}

	imd := imdraw.New(nil)

	scene := &LevelScene{
		levelIndex: index,
		space:      space,
		last:       time.Now(),
		imd:        imd,
		player:     newShip(space, "", 100),
	}

	handleCollision := space.NewWildcardCollisionHandler(cp.WILDCARD_COLLISION_TYPE)
	handleCollision.BeginFunc = func(arb *cp.Arbiter, space *cp.Space, userData interface{}) bool {
		return true
	}
	//shipBulletHandler := space.NewCollisionHandler(CollisionTypeShip, CollisionTypeBullet)
	//shipBulletHandler.BeginFunc = scene.handleShipBulletCollision

	level := assets.Levels[index]

	for _, t := range level.Teachers {
		teacher := assets.Teachers[t]
		ship := newShip(space, teacher.Name, level.Difficulty)
		ship.body.SetPosition(cp.Vector{X: 1920, Y: 1080}.Mult(rand.Float64()).Sub(cp.Vector{X: 1920 / 2, Y: 1080 / 2}))
		scene.enemies = append(scene.enemies, ship)
	}

	return scene
}
