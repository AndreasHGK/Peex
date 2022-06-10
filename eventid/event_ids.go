package eventid

type EventId uint

const (
	EventItemDrop EventId = iota
	EventItemConsume
	EventMove
	EventJump
	EventTeleport
	EventToggleSneak
	EventToggleSprint
	EventChangeWorld
	EventCommandExecution
	EventTransfer
	EventChat
	EventSkinChange
	EventStartBreak
	EventBlockBreak
	EventBlockPlace
	EventBlockPick
	EventSignEdit
	EventItemPickup
	EventItemUse
	EventItemUseOnBlock
	EventItemUseOnEntity
	EventItemDamage
	EventAttackEntity
	EventPunchAir
	EventHurt
	EventHeal
	EventFoodLoss
	EventDeath
	EventRespawn
	EventQuit
	EventExperienceGain
)

var eventNames = map[string]EventId{
	"HandleItemDrop":         EventItemDrop,
	"HandleItemConsume":      EventItemConsume,
	"HandleMove":             EventMove,
	"HandleJump":             EventJump,
	"HandleTeleport":         EventTeleport,
	"HandleToggleSneak":      EventToggleSneak,
	"HandleToggleSprint":     EventToggleSprint,
	"HandleChangeWorld":      EventChangeWorld,
	"HandleCommandExecution": EventCommandExecution,
	"HandleTransfer":         EventTransfer,
	"HandleChat":             EventChat,
	"HandleSkinChange":       EventSkinChange,
	"HandleStartBreak":       EventStartBreak,
	"HandleBlockBreak":       EventBlockBreak,
	"HandleBlockPlace":       EventBlockPlace,
	"HandleBlockPick":        EventBlockPick,
	"HandleSignEdit":         EventSignEdit,
	"HandleItemPickup":       EventItemPickup,
	"HandleItemUse":          EventItemUse,
	"HandleItemUseOnBlock":   EventItemUseOnBlock,
	"HandleItemUseOnEntity":  EventItemUseOnEntity,
	"HandleItemDamage":       EventItemDamage,
	"HandleAttackEntity":     EventAttackEntity,
	"HandlePunchAir":         EventPunchAir,
	"HandleHurt":             EventHurt,
	"HandleHeal":             EventHeal,
	"HandleFoodLoss":         EventFoodLoss,
	"HandleDeath":            EventDeath,
	"HandleRespawn":          EventRespawn,
	"HandleQuit":             EventQuit,
	"HandleExperienceGain":   EventExperienceGain,
}
