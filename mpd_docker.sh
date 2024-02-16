#!/bin/bash
sudo docker rm mpd-container -f || echo "" && \
sudo docker run -dit --restart="always" \
--name mpd-container \
-p 6600:6600 \
-p 8800:8800 \
-v $(pwd)/data/config:/root/.config \
-v $(pwd)/data/music:/var/lib/mpd/music \
-v $(pwd)/data/playlists:/var/lib/mpd/playlists \
--device /dev/snd \
vimagick/mpd