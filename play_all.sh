#!/bin/bash
mpc clear              # Clear the current playlist
mpc ls NowPlaying1 | mpc add  # Add songs from the folder to the playlist
mpc ls NowPlaying2 | mpc add  # Add songs from the folder to the playlist
mpc ls NowPlaying3 | mpc add  # Add songs from the folder to the playlist
mpc ls Peace-1 | mpc add  # Add songs from the folder to the playlist
mpc random on          # Turn on shuffle mode
mpc play               # Play the songs
# sudo amixer set 'Headphone' 96% # 96 = max , it bleeds out at 100%
# sudo amixer set 'Headphone' 50% # 50 = min , still playing , but you can barely hear it
# sudo amixer set 'Headphone' 70% # good true minimum for average listening

# sudo amixer set 'Headphone' 73% # good for 6 inch listening
