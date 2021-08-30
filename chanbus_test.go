package chanbus

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	type args struct {
		size    uint
		timeout time.Duration
	}
	tests := []struct {
		name string
		args args
		send interface{}
	}{
		{
			name: "Should receive the same value",
			args: args{
				size:    1,
				timeout: 0,
			},
			send: "say, hi",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ch := New(tt.args.size, tt.args.timeout)
			nch, _ := ch.Derive(1)
			ch.Send(tt.send)
			v := <-nch
			if v != tt.send {
				t.Errorf("Receive derived chan got = %#v, want = %#v", v, tt.send)
			}
		})
	}
}
