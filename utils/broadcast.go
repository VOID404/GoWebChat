package utils

type Broadcaster interface {
	Register(chan Message)
	Unregister(chan Message)
	Send(Message)
	SendExcept(Message, ...chan Message)
}

func chMsgIn(val chan Message, list []chan Message) bool {
	for _, x := range list {
		if x == val {
			return true
		}
	}
	return false
}

type sendExcept struct {
	msg       Message
	blacklist []chan Message
}

type broadcaster struct {
	reg        chan chan Message
	unreg      chan chan Message
	input      chan Message
	inputBlack chan sendExcept
	cli        map[chan Message]bool
}

func newBroadcaster(buff int) Broadcaster {
	b := broadcaster{
		reg:        make(chan chan Message),
		unreg:      make(chan chan Message),
		input:      make(chan Message, buff),
		inputBlack: make(chan sendExcept, buff),
		cli:        make(map[chan Message]bool),
	}

	go b.run()

	return &b
}

func (b *broadcaster) run() {
	for {
		select {
		case msg := <-b.input:
			for ch := range b.cli {
				ch <- msg
			}
		case data := <-b.inputBlack:
			for ch := range b.cli {
				if !chMsgIn(ch, data.blacklist) {
					ch <- data.msg
				}
			}
		case cli := <-b.unreg:
			delete(b.cli, cli)
		case cli, ok := <-b.reg:
			if !ok {
				return
			}

			b.cli[cli] = true
		}
	}
}

func (b *broadcaster) Register(cli chan Message) {
	if b != nil {
		b.reg <- cli
	}
}
func (b *broadcaster) Unregister(cli chan Message) {
	if b != nil {
		b.unreg <- cli
	}
}
func (b *broadcaster) Send(val Message) {
	if b != nil {
		b.input <- val
	}
}

func (b *broadcaster) SendExcept(msg Message, blacklist ...chan Message) {
	if b != nil {
		b.inputBlack <- sendExcept{msg, blacklist}
	}
}
