package peex

import (
	"github.com/andreashgk/peex/eventid"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/entity/damage"
	"github.com/df-mc/dragonfly/server/entity/healing"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"net"
	"time"
)

var _ player.Handler = (*Session)(nil)

func (s *Session) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	s.handleEvent(eventid.EventMove, func(h Handler) {
		h.HandleMove(ctx, newPos, newYaw, newPitch)
	})
}

func (s *Session) HandleJump() {
	s.handleEvent(eventid.EventJump, func(h Handler) {
		h.HandleJump()
	})
}

func (s *Session) HandleTeleport(ctx *event.Context, pos mgl64.Vec3) {
	s.handleEvent(eventid.EventTeleport, func(h Handler) {
		h.HandleTeleport(ctx, pos)
	})
}

func (s *Session) HandleChangeWorld(before, after *world.World) {
	//TODO implement me

}

func (s *Session) HandleToggleSprint(ctx *event.Context, after bool) {
	//TODO implement me

}

func (s *Session) HandleToggleSneak(ctx *event.Context, after bool) {
	//TODO implement me

}

func (s *Session) HandleChat(ctx *event.Context, message *string) {
	//TODO implement me

}

func (s *Session) HandleFoodLoss(ctx *event.Context, from, to int) {
	//TODO implement me

}

func (s *Session) HandleHeal(ctx *event.Context, health *float64, src healing.Source) {
	//TODO implement me

}

func (s *Session) HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src damage.Source) {
	//TODO implement me

}

func (s *Session) HandleDeath(src damage.Source) {
	//TODO implement me

}

func (s *Session) HandleRespawn(pos *mgl64.Vec3, w **world.World) {
	//TODO implement me

}

func (s *Session) HandleSkinChange(ctx *event.Context, skin *skin.Skin) {
	//TODO implement me

}

func (s *Session) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	//TODO implement me

}

func (s *Session) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack) {
	//TODO implement me

}

func (s *Session) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	//TODO implement me

}

func (s *Session) HandleBlockPick(ctx *event.Context, pos cube.Pos, b world.Block) {
	//TODO implement me

}

func (s *Session) HandleItemUse(ctx *event.Context) {
	//TODO implement me

}

func (s *Session) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	//TODO implement me

}

func (s *Session) HandleItemUseOnEntity(ctx *event.Context, e world.Entity) {
	//TODO implement me

}

func (s *Session) HandleItemConsume(ctx *event.Context, item item.Stack) {
	//TODO implement me

}

func (s *Session) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	//TODO implement me

}

func (s *Session) HandleExperienceGain(ctx *event.Context, amount *int) {
	//TODO implement me

}

func (s *Session) HandlePunchAir(ctx *event.Context) {
	//TODO implement me

}

func (s *Session) HandleSignEdit(ctx *event.Context, oldText, newText string) {
	//TODO implement me

}

func (s *Session) HandleItemDamage(ctx *event.Context, i item.Stack, damage int) {
	//TODO implement me

}

func (s *Session) HandleItemPickup(ctx *event.Context, i item.Stack) {
	//TODO implement me

}

func (s *Session) HandleItemDrop(ctx *event.Context, e *entity.Item) {
	//TODO implement me

}

func (s *Session) HandleTransfer(ctx *event.Context, addr *net.UDPAddr) {
	//TODO implement me

}

func (s *Session) HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string) {
	//TODO implement me

}

func (s *Session) HandleQuit() {
	s.handleEvent(eventid.EventQuit, func(h Handler) {
		h.HandleQuit()
	})
	delete(s.m.sessions, s.p.UUID())
}
