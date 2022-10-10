package registration

import "github.com/looplab/fsm"

type RegisterFSM struct {
	FSM *fsm.FSM
}

const (
	EventRegister      = "register"
	EventGotRightToken = "got_right_token"
	EventGotWrongToken = "got_wrong_token"
	EventRepeat        = "repeat"
	EventCancel        = "cancel"

	StateIdle          = "idle"
	StateTokenAwaiting = "token_awaiting"
	StateSuccess       = "success"
	StateFailure       = "failure"
)

func NewRegisterFSM() *RegisterFSM {
	r := &RegisterFSM{}
	r.FSM = fsm.NewFSM(
		"idle",
		fsm.Events{
			{Name: EventRegister, Src: []string{StateIdle}, Dst: StateTokenAwaiting},
			{Name: EventGotRightToken, Src: []string{StateTokenAwaiting}, Dst: StateSuccess},
			{Name: EventGotWrongToken, Src: []string{StateTokenAwaiting}, Dst: StateFailure},
			{Name: EventRepeat, Src: []string{StateFailure}, Dst: StateTokenAwaiting},
			{Name: EventCancel, Src: []string{StateTokenAwaiting}, Dst: StateFailure},
		},
		fsm.Callbacks{},
	)

	return r
}
