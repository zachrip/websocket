package websocket

import (
	"reflect"
	"testing"
)

func TestNewHub(t *testing.T) {
	t.Run("New Hub", func(t *testing.T) {
		hub := NewHub()

		if reflect.TypeOf(hub).String() != "*websocket.Hub" {
			t.Errorf("Wanted *websocket.Hub, got %q", reflect.TypeOf(hub))
		}
	})
}

func fn(data interface{}) {}
func TestHub_On(t *testing.T) {
	type args struct {
		event string
		fn    EventHandler
	}

	hub := NewHub()

	tests := []struct {
		name string
		h    *Hub
		args args
		want *Hub
	}{
		{
			"Add event",
			hub.On("test", fn),
			args{
				"test",
				func(data interface{}) {},
			},
			hub,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.On(tt.args.event, tt.args.fn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Hub.On() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHub_Run(t *testing.T) {
	//h := NewHub()

	t.Run("Should emit events when messages come in", func(t *testing.T) {
		//tt.h.Run()
	})
}
