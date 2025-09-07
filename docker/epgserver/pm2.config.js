const grab = process.env.SITE
  ? `npm run grab -- --site=${process.env.SITE} ${process.env.CLANG ? `--lang=${process.env.CLANG}` : ''
  } --output=${process.env.OUTPUT || 'public/epg.xml'}`
  : `npm run grab -- --channels=channels.xml --output=${process.env.OUTPUT || 'public/epg.xml'}`


const apps = [
  {
    name: 'serve',
    script: 'npx serve -- public',
    instances: 1,
    watch: false,
    autorestart: true
  },
  {
    name: 'grab',
    script: `npx chronos -e "${grab}" -p "${process.env.CRON_SCHEDULE}" -l`,
    instances: 1,
    watch: false,
    autorestart: true
  }
];

if (process.env.RUN_AT_STARTUP === 'true') {
  apps.push({
    name: 'grab-at-startup',
    script: grab,
    instances: 1,
    autorestart: false,
    watch: false,
    max_restarts: 1
  });
}

module.exports = { apps };