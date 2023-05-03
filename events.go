// Package peex
// This file was generated using the event generator. Do not edit.
package peex

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player/skin"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"net"
	"time"
)

type eventId uint

const (
	eventMove eventId = iota
	eventJump
	eventTeleport
	eventChangeWorld
	eventToggleSprint
	eventToggleSneak
	eventChat
	eventFoodLoss
	eventHeal
	eventHurt
	eventDeath
	eventRespawn
	eventSkinChange
	eventStartBreak
	eventBlockBreak
	eventBlockPlace
	eventBlockPick
	eventItemUse
	eventItemUseOnBlock
	eventItemUseOnEntity
	eventItemConsume
	eventAttackEntity
	eventExperienceGain
	eventPunchAir
	eventSignEdit
	eventItemDamage
	eventItemPickup
	eventItemDrop
	eventTransfer
	eventCommandExecution
	eventQuit
)

// getHandlerEvents returns which events a handler implements. Since it is impossible to distinguish actually imlemented
// methods from ones embedded using player.NopHandler, it is recommended to not embed it at all. Most Peex handlers
// won't implement player.Handler!
func getHandlerEvents(h Handler) map[eventId]struct{} {
	m := make(map[eventId]struct{})

	if _, ok := h.(eventMoveHandler); ok {
		m[eventMove] = struct{}{}
	}
	if _, ok := h.(eventJumpHandler); ok {
		m[eventJump] = struct{}{}
	}
	if _, ok := h.(eventTeleportHandler); ok {
		m[eventTeleport] = struct{}{}
	}
	if _, ok := h.(eventChangeWorldHandler); ok {
		m[eventChangeWorld] = struct{}{}
	}
	if _, ok := h.(eventToggleSprintHandler); ok {
		m[eventToggleSprint] = struct{}{}
	}
	if _, ok := h.(eventToggleSneakHandler); ok {
		m[eventToggleSneak] = struct{}{}
	}
	if _, ok := h.(eventChatHandler); ok {
		m[eventChat] = struct{}{}
	}
	if _, ok := h.(eventFoodLossHandler); ok {
		m[eventFoodLoss] = struct{}{}
	}
	if _, ok := h.(eventHealHandler); ok {
		m[eventHeal] = struct{}{}
	}
	if _, ok := h.(eventHurtHandler); ok {
		m[eventHurt] = struct{}{}
	}
	if _, ok := h.(eventDeathHandler); ok {
		m[eventDeath] = struct{}{}
	}
	if _, ok := h.(eventRespawnHandler); ok {
		m[eventRespawn] = struct{}{}
	}
	if _, ok := h.(eventSkinChangeHandler); ok {
		m[eventSkinChange] = struct{}{}
	}
	if _, ok := h.(eventStartBreakHandler); ok {
		m[eventStartBreak] = struct{}{}
	}
	if _, ok := h.(eventBlockBreakHandler); ok {
		m[eventBlockBreak] = struct{}{}
	}
	if _, ok := h.(eventBlockPlaceHandler); ok {
		m[eventBlockPlace] = struct{}{}
	}
	if _, ok := h.(eventBlockPickHandler); ok {
		m[eventBlockPick] = struct{}{}
	}
	if _, ok := h.(eventItemUseHandler); ok {
		m[eventItemUse] = struct{}{}
	}
	if _, ok := h.(eventItemUseOnBlockHandler); ok {
		m[eventItemUseOnBlock] = struct{}{}
	}
	if _, ok := h.(eventItemUseOnEntityHandler); ok {
		m[eventItemUseOnEntity] = struct{}{}
	}
	if _, ok := h.(eventItemConsumeHandler); ok {
		m[eventItemConsume] = struct{}{}
	}
	if _, ok := h.(eventAttackEntityHandler); ok {
		m[eventAttackEntity] = struct{}{}
	}
	if _, ok := h.(eventExperienceGainHandler); ok {
		m[eventExperienceGain] = struct{}{}
	}
	if _, ok := h.(eventPunchAirHandler); ok {
		m[eventPunchAir] = struct{}{}
	}
	if _, ok := h.(eventSignEditHandler); ok {
		m[eventSignEdit] = struct{}{}
	}
	if _, ok := h.(eventItemDamageHandler); ok {
		m[eventItemDamage] = struct{}{}
	}
	if _, ok := h.(eventItemPickupHandler); ok {
		m[eventItemPickup] = struct{}{}
	}
	if _, ok := h.(eventItemDropHandler); ok {
		m[eventItemDrop] = struct{}{}
	}
	if _, ok := h.(eventTransferHandler); ok {
		m[eventTransfer] = struct{}{}
	}
	if _, ok := h.(eventCommandExecutionHandler); ok {
		m[eventCommandExecution] = struct{}{}
	}
	if _, ok := h.(eventQuitHandler); ok {
		m[eventQuit] = struct{}{}
	}
	return m
}

var allEvents = map[string]eventId{
	"eventMove":             eventMove,
	"eventJump":             eventJump,
	"eventTeleport":         eventTeleport,
	"eventChangeWorld":      eventChangeWorld,
	"eventToggleSprint":     eventToggleSprint,
	"eventToggleSneak":      eventToggleSneak,
	"eventChat":             eventChat,
	"eventFoodLoss":         eventFoodLoss,
	"eventHeal":             eventHeal,
	"eventHurt":             eventHurt,
	"eventDeath":            eventDeath,
	"eventRespawn":          eventRespawn,
	"eventSkinChange":       eventSkinChange,
	"eventStartBreak":       eventStartBreak,
	"eventBlockBreak":       eventBlockBreak,
	"eventBlockPlace":       eventBlockPlace,
	"eventBlockPick":        eventBlockPick,
	"eventItemUse":          eventItemUse,
	"eventItemUseOnBlock":   eventItemUseOnBlock,
	"eventItemUseOnEntity":  eventItemUseOnEntity,
	"eventItemConsume":      eventItemConsume,
	"eventAttackEntity":     eventAttackEntity,
	"eventExperienceGain":   eventExperienceGain,
	"eventPunchAir":         eventPunchAir,
	"eventSignEdit":         eventSignEdit,
	"eventItemDamage":       eventItemDamage,
	"eventItemPickup":       eventItemPickup,
	"eventItemDrop":         eventItemDrop,
	"eventTransfer":         eventTransfer,
	"eventCommandExecution": eventCommandExecution,
	"eventQuit":             eventQuit,
}

type eventMoveHandler interface {
	HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64)
}

type eventJumpHandler interface {
	HandleJump()
}

type eventTeleportHandler interface {
	HandleTeleport(ctx *event.Context, pos mgl64.Vec3)
}

type eventChangeWorldHandler interface {
	HandleChangeWorld(before, after *world.World)
}

type eventToggleSprintHandler interface {
	HandleToggleSprint(ctx *event.Context, after bool)
}

type eventToggleSneakHandler interface {
	HandleToggleSneak(ctx *event.Context, after bool)
}

type eventChatHandler interface {
	HandleChat(ctx *event.Context, message *string)
}

type eventFoodLossHandler interface {
	HandleFoodLoss(ctx *event.Context, from int, to *int)
}

type eventHealHandler interface {
	HandleHeal(ctx *event.Context, health *float64, src world.HealingSource)
}

type eventHurtHandler interface {
	HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src world.DamageSource)
}

type eventDeathHandler interface {
	HandleDeath(src world.DamageSource, keepInv *bool)
}

type eventRespawnHandler interface {
	HandleRespawn(pos *mgl64.Vec3, w **world.World)
}

type eventSkinChangeHandler interface {
	HandleSkinChange(ctx *event.Context, skin *skin.Skin)
}

type eventStartBreakHandler interface {
	HandleStartBreak(ctx *event.Context, pos cube.Pos)
}

type eventBlockBreakHandler interface {
	HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int)
}

type eventBlockPlaceHandler interface {
	HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block)
}

type eventBlockPickHandler interface {
	HandleBlockPick(ctx *event.Context, pos cube.Pos, b world.Block)
}

type eventItemUseHandler interface {
	HandleItemUse(ctx *event.Context)
}

type eventItemUseOnBlockHandler interface {
	HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3)
}

type eventItemUseOnEntityHandler interface {
	HandleItemUseOnEntity(ctx *event.Context, e world.Entity)
}

type eventItemConsumeHandler interface {
	HandleItemConsume(ctx *event.Context, item item.Stack)
}

type eventAttackEntityHandler interface {
	HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool)
}

type eventExperienceGainHandler interface {
	HandleExperienceGain(ctx *event.Context, amount *int)
}

type eventPunchAirHandler interface {
	HandlePunchAir(ctx *event.Context)
}

type eventSignEditHandler interface {
	HandleSignEdit(ctx *event.Context, oldText, newText string)
}

type eventItemDamageHandler interface {
	HandleItemDamage(ctx *event.Context, i item.Stack, damage int)
}

type eventItemPickupHandler interface {
	HandleItemPickup(ctx *event.Context, i item.Stack)
}

type eventItemDropHandler interface {
	HandleItemDrop(ctx *event.Context, e world.Entity)
}

type eventTransferHandler interface {
	HandleTransfer(ctx *event.Context, addr *net.UDPAddr)
}

type eventCommandExecutionHandler interface {
	HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string)
}

type eventQuitHandler interface {
	HandleQuit()
}

func (s *Session) HandleMove(ctx *event.Context, newPos mgl64.Vec3, newYaw, newPitch float64) {
	s.handleEvent(eventMove, func(h Handler) {
		h.(eventMoveHandler).HandleMove(ctx, newPos, newYaw, newPitch)
	})
}

func (s *Session) HandleJump() {
	s.handleEvent(eventJump, func(h Handler) {
		h.(eventJumpHandler).HandleJump()
	})
}

func (s *Session) HandleTeleport(ctx *event.Context, pos mgl64.Vec3) {
	s.handleEvent(eventTeleport, func(h Handler) {
		h.(eventTeleportHandler).HandleTeleport(ctx, pos)
	})
}

func (s *Session) HandleChangeWorld(before, after *world.World) {
	s.handleEvent(eventChangeWorld, func(h Handler) {
		h.(eventChangeWorldHandler).HandleChangeWorld(before, after)
	})
}

func (s *Session) HandleToggleSprint(ctx *event.Context, after bool) {
	s.handleEvent(eventToggleSprint, func(h Handler) {
		h.(eventToggleSprintHandler).HandleToggleSprint(ctx, after)
	})
}

func (s *Session) HandleToggleSneak(ctx *event.Context, after bool) {
	s.handleEvent(eventToggleSneak, func(h Handler) {
		h.(eventToggleSneakHandler).HandleToggleSneak(ctx, after)
	})
}

func (s *Session) HandleChat(ctx *event.Context, message *string) {
	s.handleEvent(eventChat, func(h Handler) {
		h.(eventChatHandler).HandleChat(ctx, message)
	})
}

func (s *Session) HandleFoodLoss(ctx *event.Context, from int, to *int) {
	s.handleEvent(eventFoodLoss, func(h Handler) {
		h.(eventFoodLossHandler).HandleFoodLoss(ctx, from, to)
	})
}

func (s *Session) HandleHeal(ctx *event.Context, health *float64, src world.HealingSource) {
	s.handleEvent(eventHeal, func(h Handler) {
		h.(eventHealHandler).HandleHeal(ctx, health, src)
	})
}

func (s *Session) HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src world.DamageSource) {
	s.handleEvent(eventHurt, func(h Handler) {
		h.(eventHurtHandler).HandleHurt(ctx, damage, attackImmunity, src)
	})
}

func (s *Session) HandleDeath(src world.DamageSource, keepInv *bool) {
	s.handleEvent(eventDeath, func(h Handler) {
		h.(eventDeathHandler).HandleDeath(src, keepInv)
	})
}

func (s *Session) HandleRespawn(pos *mgl64.Vec3, w **world.World) {
	s.handleEvent(eventRespawn, func(h Handler) {
		h.(eventRespawnHandler).HandleRespawn(pos, w)
	})
}

func (s *Session) HandleSkinChange(ctx *event.Context, skin *skin.Skin) {
	s.handleEvent(eventSkinChange, func(h Handler) {
		h.(eventSkinChangeHandler).HandleSkinChange(ctx, skin)
	})
}

func (s *Session) HandleStartBreak(ctx *event.Context, pos cube.Pos) {
	s.handleEvent(eventStartBreak, func(h Handler) {
		h.(eventStartBreakHandler).HandleStartBreak(ctx, pos)
	})
}

func (s *Session) HandleBlockBreak(ctx *event.Context, pos cube.Pos, drops *[]item.Stack, xp *int) {
	s.handleEvent(eventBlockBreak, func(h Handler) {
		h.(eventBlockBreakHandler).HandleBlockBreak(ctx, pos, drops, xp)
	})
}

func (s *Session) HandleBlockPlace(ctx *event.Context, pos cube.Pos, b world.Block) {
	s.handleEvent(eventBlockPlace, func(h Handler) {
		h.(eventBlockPlaceHandler).HandleBlockPlace(ctx, pos, b)
	})
}

func (s *Session) HandleBlockPick(ctx *event.Context, pos cube.Pos, b world.Block) {
	s.handleEvent(eventBlockPick, func(h Handler) {
		h.(eventBlockPickHandler).HandleBlockPick(ctx, pos, b)
	})
}

func (s *Session) HandleItemUse(ctx *event.Context) {
	s.handleEvent(eventItemUse, func(h Handler) {
		h.(eventItemUseHandler).HandleItemUse(ctx)
	})
}

func (s *Session) HandleItemUseOnBlock(ctx *event.Context, pos cube.Pos, face cube.Face, clickPos mgl64.Vec3) {
	s.handleEvent(eventItemUseOnBlock, func(h Handler) {
		h.(eventItemUseOnBlockHandler).HandleItemUseOnBlock(ctx, pos, face, clickPos)
	})
}

func (s *Session) HandleItemUseOnEntity(ctx *event.Context, e world.Entity) {
	s.handleEvent(eventItemUseOnEntity, func(h Handler) {
		h.(eventItemUseOnEntityHandler).HandleItemUseOnEntity(ctx, e)
	})
}

func (s *Session) HandleItemConsume(ctx *event.Context, item item.Stack) {
	s.handleEvent(eventItemConsume, func(h Handler) {
		h.(eventItemConsumeHandler).HandleItemConsume(ctx, item)
	})
}

func (s *Session) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	s.handleEvent(eventAttackEntity, func(h Handler) {
		h.(eventAttackEntityHandler).HandleAttackEntity(ctx, e, force, height, critical)
	})
}

func (s *Session) HandleExperienceGain(ctx *event.Context, amount *int) {
	s.handleEvent(eventExperienceGain, func(h Handler) {
		h.(eventExperienceGainHandler).HandleExperienceGain(ctx, amount)
	})
}

func (s *Session) HandlePunchAir(ctx *event.Context) {
	s.handleEvent(eventPunchAir, func(h Handler) {
		h.(eventPunchAirHandler).HandlePunchAir(ctx)
	})
}

func (s *Session) HandleSignEdit(ctx *event.Context, oldText, newText string) {
	s.handleEvent(eventSignEdit, func(h Handler) {
		h.(eventSignEditHandler).HandleSignEdit(ctx, oldText, newText)
	})
}

func (s *Session) HandleItemDamage(ctx *event.Context, i item.Stack, damage int) {
	s.handleEvent(eventItemDamage, func(h Handler) {
		h.(eventItemDamageHandler).HandleItemDamage(ctx, i, damage)
	})
}

func (s *Session) HandleItemPickup(ctx *event.Context, i item.Stack) {
	s.handleEvent(eventItemPickup, func(h Handler) {
		h.(eventItemPickupHandler).HandleItemPickup(ctx, i)
	})
}

func (s *Session) HandleItemDrop(ctx *event.Context, e world.Entity) {
	s.handleEvent(eventItemDrop, func(h Handler) {
		h.(eventItemDropHandler).HandleItemDrop(ctx, e)
	})
}

func (s *Session) HandleTransfer(ctx *event.Context, addr *net.UDPAddr) {
	s.handleEvent(eventTransfer, func(h Handler) {
		h.(eventTransferHandler).HandleTransfer(ctx, addr)
	})
}

func (s *Session) HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string) {
	s.handleEvent(eventCommandExecution, func(h Handler) {
		h.(eventCommandExecutionHandler).HandleCommandExecution(ctx, command, args)
	})
}

func (s *Session) HandleQuit() {
	s.handleEvent(eventQuit, func(h Handler) {
		h.(eventQuitHandler).HandleQuit()
	})
	s.doQuit()
}
