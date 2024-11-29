package types

type AttackerCommand int

const (
	CMD_NULL AttackerCommand = iota
	CMD_CONTINUE
	CMD_RETURN
	CMD_ABORT
	CMD_SKIP
	CMD_ROLE_TO_NORMAL
	CMD_ROLE_TO_ATTACKER
	CMD_EXIT
	CMD_UPDATE_STATE
)
