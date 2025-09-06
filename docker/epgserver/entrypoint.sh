#!/bin/sh

# permissions
echo "Fixing /config permissions. . ."
chown -R epguser:epguser \
  /config

# chown the app directory, but not node_modules
echo "Fixing app permissions. . ."
find /app/epg -maxdepth 1 ! -name node_modules ! -name epg -exec chown -R epguser:epguser '{}' \;

cd /app/epg || exit

export PATH=$PATH:/app/epg/node_modules/.bin

npm run api:load  
exec \
  pm2-runtime pm2.config.js
