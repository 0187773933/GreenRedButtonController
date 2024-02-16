import sys
import time
import subprocess
from evdev import InputDevice, categorize, ecodes, KeyEvent

def get_mpc_status():
	# Run the 'mpc status' command and get its output
	result = subprocess.run(['mpc', 'status'], stdout=subprocess.PIPE)
	output = result.stdout.decode('utf-8').split('\n')

	# Check the output for playing, paused, or stopped
	if "[playing]" in output[1]:
		return "playing"
	elif "[paused]" in output[1]:
		return "paused"
	else:
		return "stopped"

dev = InputDevice('/dev/input/event0')  # replace X with the correct number

print(f"Listening on device {dev.name}...")

shift_pressed = False
shift_x_start_time = None
shift_x_pressed = False

for event in dev.read_loop():
	if event.type == ecodes.EV_KEY:
		key_event = categorize(event)
		if key_event.keycode in ['KEY_LEFTSHIFT', 'KEY_RIGHTSHIFT']:
			if key_event.keystate == KeyEvent.key_down:
				shift_pressed = True
			elif key_event.keystate == KeyEvent.key_up:
				shift_pressed = False

		if shift_pressed:
			if key_event.keycode == 'KEY_X' and key_event.keystate == KeyEvent.key_down:
				shift_x_pressed = True
				shift_x_start_time = time.time()
				print("Shift + X was pressed!" , file=sys.stderr)
			elif key_event.keycode == 'KEY_X' and key_event.keystate == KeyEvent.key_up:
				shift_x_pressed = False

			if key_event.keycode == 'KEY_E' and key_event.keystate == KeyEvent.key_down:
				status = get_mpc_status()
				print(status)
				if status == "playing":
					print("was already playing, pressing next",file=sys.stderr)
					subprocess.Popen(['mpc', 'next'])
				elif status == "stopped":
					subprocess.Popen(['/usr/local/bin/playAll'])
				elif status == "paused":
					subprocess.Popen(['mpc', 'play'])
				else:
					print("was not playing, restarting playlist",file=sys.stderr)
					subprocess.Popen(['/usr/local/bin/playAll'])

		# Check if Shift + X has been pressed for more than 6 seconds
		if shift_x_pressed and (time.time() - shift_x_start_time >= 3):
			print("Shift + X was held for 3 seconds or more. Stopping playback.",file=sys.stderr)
			subprocess.run(['mpc', 'pause'])
			shift_x_pressed = False  # Reset so we don't repeatedly execute the command
