# EPG (Electronic Program Guide) Integration

This document describes the EPG functionality that has been integrated into the m3uproxy player.

## Features

The player now supports displaying current program information for channels that have EPG data available.

### EPG Display

When a channel is selected, the player automatically fetches and displays:

- **Program Title**: The name of the currently playing program
- **Program Time**: Start and end times in local format (HH:MM - HH:MM) 
- **Program Duration**: Duration in minutes (e.g., "60min")
- **Progress Bar**: Visual indicator showing how much of the program has elapsed
- **Program Description**: Detailed description of the current program (when available)

### Controls

- **Channel Selection**: When you change to a channel with EPG data, program information is automatically fetched and displayed
- **Info Toggle**: Press `i` key to toggle the program information overlay on/off
- **Auto-hide**: Program information automatically hides after 3 seconds of video playback
- **Auto-update**: EPG data is automatically refreshed every 30 seconds to keep information current

### Technical Details

#### EPG Endpoint
The backend provides an EPG endpoint: `/epg/{channel_id}?timestamp={timestamp}`

- `channel_id`: The channel identifier (from `tvg-id` field in M3U playlist)
- `timestamp`: Current timestamp in format `YYYYMMDDHHMMSS +ZZZZ`

#### Example Request
```bash
curl "http://admin:admin@localhost:8080/epg/RTP1.pt?timestamp=$(date "+%Y%m%d%H%M%S %z" | sed 's/ /%20/g; s/+/%2B/g; s/-/%2D/g')"
```

#### Response Format
The endpoint returns XMLTV format containing:
- Channel information
- Current program details (title, description, start/end times)

#### Files Modified
- `src/components/Overlay.js`: Extended to display program information
- `src/components/Player.js`: Added event listeners for better integration
- `src/App.js`: Added EPG fetching and integration logic
- `src/utils/EpgService.js`: New service for EPG data management
- `src/main.css`: Added styles for program information display

### Usage

1. Start the player and select a channel
2. If EPG data is available for the channel, program information will automatically appear
3. Press `i` to manually toggle program information display
4. The information updates automatically every 30 seconds
5. Program information will show/hide along with channel information when changing channels

### Responsive Design

The EPG display is responsive and adapts to mobile devices:
- Smaller font sizes on mobile
- Adjusted positioning and padding
- Shorter program descriptions on small screens

### Error Handling

- If EPG data is not available for a channel, no program information is displayed
- Network errors are logged but don't affect playback
- Invalid or malformed EPG data is handled gracefully

## Future Enhancements

Potential improvements could include:
- Full program schedule view
- Program reminders/notifications
- Program recording capabilities
- Multi-day EPG browsing
- Category-based program filtering
