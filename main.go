package main

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
	log "github.com/sirupsen/logrus"
	evdev "github.com/gvalkov/golang-evdev"
)

func get_mpc_status() string {
	cmd := exec.Command("mpc", "status")
	out, _ := cmd.CombinedOutput()
	output := string(out)
	if output != "" {
		lines := strings.Split(output, "\n")
		if len(lines) > 1 {
			if strings.Contains(lines[1], "[playing]") {
				return "playing"
			} else if strings.Contains(lines[1], "[paused]") {
				return "paused"
			}
		}
	}
	return "stopped"
}

func run_command_async(cmd *exec.Cmd) {
	go func() { _ = cmd.Run() }()
}


type ClickType int

const (
	SingleClick ClickType = iota
	DoubleClick
	TripleClick
)

type ButtonHandler struct {
	PressCount     int
	LastPressTime  time.Time
	Timer          *time.Timer
	SingleClickDelay time.Duration
	OnClick         func(ClickType)
}

func NewButtonHandler(singleClickDelay time.Duration, onClick func(ClickType)) *ButtonHandler {
	return &ButtonHandler{
		SingleClickDelay: singleClickDelay,
		OnClick:          onClick,
	}
}

func (b *ButtonHandler) Press() {
	now := time.Now()

	if now.Sub(b.LastPressTime) > b.SingleClickDelay {
		b.PressCount = 0
	}
	b.PressCount++
	b.LastPressTime = now

	if b.Timer != nil {
		b.Timer.Stop()
	}

	b.Timer = time.AfterFunc(b.SingleClickDelay, func() {
		switch b.PressCount {
		case 1:
			b.OnClick(SingleClick)
		case 2:
			b.OnClick(DoubleClick)
		case 3:
			b.OnClick(TripleClick)
		}
		b.PressCount = 0
	})
}

func main() {

	// log.SetLevel( log.DebugLevel )
	log.SetLevel( log.InfoLevel )

	run_command_async(exec.Command("mpc pause"))

	devices , _ := evdev.ListInputDevices( "/dev/input/event0" )
	if len(devices) == 0 {
		fmt.Println( "No devices found!" )
		return
	}
	device := devices[0]
	fmt.Printf( "Listening on device %s...\n" , device.Name )


	buttonHandlers := make(map[int]*ButtonHandler)

	var key_states = map[int]string{
		0: "released",
		1: "pressed",
		2: "held",
	}

	const key_shift = 42
	const key_e = 18
	const key_x = 45
	const key_released = 0
	const key_pressed = 1
	const key_held = 2

	const held_count = 30

	// var e_held_fresh = true
	// var x_held_fresh = true
	var e_held_count = 0
	var x_held_count = 0

	// var startTimeMap = make(map[int]time.Time)

	// shift_pressed := false
	// e_pressed := false
	// x_pressed := false
	// shift_x_pressed := false
	// Initialize buttonHandlers for the keys you're interested in
	interestedKeys := []int{42, 18, 45} // Corresponding to key_shift, key_e, key_x in your original code
	var latest_event evdev.InputEvent
	for _, key := range interestedKeys {
		handler := NewButtonHandler(500*time.Millisecond, func(click ClickType)  {
			key_event := evdev.NewKeyEvent( &latest_event )
			key_state := key_states[ int( key_event.State ) ]
			fmt.Println( "key_event" , key_event )
			fmt.Println( "key_state" , key_state )
			switch click {
			case SingleClick:
				fmt.Println("Single Click!")
				switch latest_event.Code {
					case key_e:
						fmt.Println( "e single clicked" )
						status := get_mpc_status()
						fmt.Println(status)
						if status == "playing" {
							fmt.Println("was already playing, pressing next")
							run_command_async(exec.Command("mpc", "next"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "next"))
						} else if status == "stopped" {
							run_command_async(exec.Command("/usr/local/bin/playAll"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "restarted playing"))
						} else if status == "paused" {
							run_command_async(exec.Command("mpc", "play"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "playing"))
						} else {
							fmt.Println("was not playing, restarting playlist")
							run_command_async(exec.Command("/usr/local/bin/playAll"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "restarted playing"))
						}
					case key_x:
						fmt.Println( "x single clicked" )
				}
			case DoubleClick:
				fmt.Println("Double Click!")
				switch latest_event.Code {
					case key_e:
						fmt.Println( "e double clicked" )
						status := get_mpc_status()
						fmt.Println(status)
						if status == "playing" {
							fmt.Println("was already playing, pressing next")
							run_command_async(exec.Command("mpc", "next"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "double click - next"))
						} else if status == "stopped" {
							run_command_async(exec.Command("/usr/local/bin/playAll"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "double click - restarted playing"))
						} else if status == "paused" {
							run_command_async(exec.Command("mpc", "play"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "double click - playing"))
						} else {
							fmt.Println("was not playing, restarting playlist")
							run_command_async(exec.Command("/usr/local/bin/playAll"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "double click - restarted playing"))
						}
					case key_x:
						fmt.Println( "x double clicked" )
				}
			case TripleClick:
				fmt.Println("Triple Click!")
				switch latest_event.Code {
					case key_e:
						fmt.Println( "e triple clicked" )
						status := get_mpc_status()
						fmt.Println(status)
						if status == "playing" {
							fmt.Println("was already playing, pressing next")
							run_command_async(exec.Command("mpc", "next"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "triple click - next"))
						} else if status == "stopped" {
							run_command_async(exec.Command("/usr/local/bin/playAll"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "triple click - restarted playing"))
						} else if status == "paused" {
							run_command_async(exec.Command("mpc", "play"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "triple click - playing"))
						} else {
							fmt.Println("was not playing, restarting playlist")
							run_command_async(exec.Command("/usr/local/bin/playAll"))
							run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "triple click -restarted playing"))
						}
					case key_x:
						fmt.Println( "x triple clicked" )
				}
			}
		})
		buttonHandlers[key] = handler
	}

	for {
		events, _ := device.Read()
		for _, ev := range events {
			if ev.Type != evdev.EV_KEY { continue }
			key_event_int := int( ev.Code )
			key_event := evdev.NewKeyEvent( &ev )
			key_state := key_states[ int( key_event.State ) ]
			switch key_state {
				case "pressed":
					fmt.Println( "pressed" )
					if handler, exists := buttonHandlers[key_event_int]; exists && ev.Value == 1 { // Checking for EV_KEY press event
						latest_event = ev
						handler.Press()
					}
				case "held":
					if ev.Code == key_e {
						e_held_count += 1
						if e_held_count != held_count { continue }
						fmt.Println( "held e" , e_held_count )
					} else if ev.Code == key_x {
						x_held_count += 1
						if x_held_count != held_count { continue }
						fmt.Println( "held x" , x_held_count )
						fmt.Println( "pausing" )
						run_command_async(exec.Command("mpc" , "pause"))
						run_command_async(exec.Command("/usr/local/bin/sendPushNotification" , "paused"))
					}
				case "released":
					e_held_count = 0
					x_held_count = 0

			}

		}
	}


}
