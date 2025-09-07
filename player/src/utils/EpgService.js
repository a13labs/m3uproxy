import { Logger } from './Logger';

class EpgService {
    constructor() {
        this.baseUrl = '/epg';
        if (__DEV__) {
            this.baseUrl = `http://${window.location.hostname}:8080/epg`;
        }
    }

    async fetchCurrentProgram(channelId) {
        if (!channelId) {
            Logger.warn('No channel ID provided for EPG fetch');
            return null;
        }

        const username = localStorage.getItem('username');
        const password = localStorage.getItem('password');
        const headers = {
            Authorization: 'Basic ' + btoa(`${username}:${password}`)
        };

        // Format current timestamp for the EPG request
        const now = new Date();
        const timestamp = this.formatTimestamp(now);

        const url = `${this.baseUrl}/${encodeURIComponent(channelId)}?timestamp=${encodeURIComponent(timestamp)}`;

        try {
            Logger.info(`Fetching EPG data from: ${url}`);
            const response = await fetch(url, { headers });

            if (!response.ok) {
                Logger.error(`EPG fetch failed: ${response.status} ${response.statusText}`);
                return null;
            }

            const xmlText = await response.text();
            // Pass the request time so the parser can select the correct programme
            const programData = this.parseEpgXml(xmlText, now);

            Logger.info('EPG data fetched successfully:', programData);
            return programData;
        } catch (error) {
            Logger.error('Error fetching EPG data:', error);
            return null;
        }
    }

    formatTimestamp(date) {
        const year = date.getFullYear();
        const month = String(date.getMonth() + 1).padStart(2, '0');
        const day = String(date.getDate()).padStart(2, '0');
        const hours = String(date.getHours()).padStart(2, '0');
        const minutes = String(date.getMinutes()).padStart(2, '0');
        const seconds = String(date.getSeconds()).padStart(2, '0');

        // Use the actual current timezone offset (includes DST) and don't mutate the input date
        const offsetMinutesTotal = -date.getTimezoneOffset(); // positive for UTC+ zones
        const offsetSign = offsetMinutesTotal >= 0 ? '+' : '-';
        const absOffset = Math.abs(offsetMinutesTotal);
        const offsetHours = Math.floor(absOffset / 60);
        const offsetMinutes = absOffset % 60;
        const timezone = `${offsetSign}${String(offsetHours).padStart(2, '0')}${String(offsetMinutes).padStart(2, '0')}`;

        return `${year}${month}${day}${hours}${minutes}${seconds} ${timezone}`;
    }

    parseEpgXml(xmlText, targetTime) {
        try {
            const parser = new DOMParser();
            const xmlDoc = parser.parseFromString(xmlText, 'text/xml');

            // Check for parsing errors
            const parserError = xmlDoc.querySelector('parsererror');
            if (parserError) {
                Logger.error('XML parsing error:', parserError.textContent);
                return null;
            }

            const programmes = xmlDoc.querySelectorAll('programme');

            if (programmes.length === 0) {
                Logger.warn('No programme found in EPG data');
                return null;
            }

            // Prefer the programme that contains targetTime (or now). If none match, fall back to the first complete programme.
            const now = targetTime instanceof Date ? targetTime : new Date();
            let fallback = null;
            for (let i = 0; i < programmes.length; i++) {
                const p = programmes[i];
                const titleElement = p.querySelector('title');
                const descElement = p.querySelector('desc');
                const startAttr = p.getAttribute('start');
                const stopAttr = p.getAttribute('stop');

                if (!titleElement || !startAttr || !stopAttr) {
                    // keep scanning; incomplete entries are ignored
                    continue;
                }

                const startDate = this.parseXmlTvTime(startAttr);
                const stopDate = this.parseXmlTvTime(stopAttr);

                if (!startDate || !stopDate) {
                    continue;
                }

                if (!fallback) {
                    fallback = { titleElement, descElement, startDate, stopDate };
                }

                if (now >= startDate && now < stopDate) {
                    return {
                        title: titleElement.textContent,
                        description: descElement ? descElement.textContent : '',
                        startTime: startDate,
                        endTime: stopDate
                    };
                }
            }

            if (fallback) {
                return {
                    title: fallback.titleElement.textContent,
                    description: fallback.descElement ? fallback.descElement.textContent : '',
                    startTime: fallback.startDate,
                    endTime: fallback.stopDate
                };
            }

            Logger.warn('No usable programme found in EPG data');
            return null;
        } catch (error) {
            Logger.error('Error parsing EPG XML:', error);
            return null;
        }
    }

    parseXmlTvTime(epgTimeString) {
        // Parse EPG time format: "20250906185900 +0000" or "20250906185900"
        if (!epgTimeString) return null;

        const timeString = epgTimeString.trim();
        const year = parseInt(timeString.substr(0, 4), 10);
        const month = parseInt(timeString.substr(4, 2), 10) - 1; // months are 0-indexed
        const day = parseInt(timeString.substr(6, 2), 10);
        const hour = parseInt(timeString.substr(8, 2), 10);
        const minute = parseInt(timeString.substr(10, 2), 10);
        const second = parseInt(timeString.substr(12, 2), 10);

        // If there's a timezone part, construct an ISO-8601 string and let Date parse it
        if (timeString.length > 14) {
            // Expect a space then +HHMM or -HHMM
            const tzPart = timeString.substr(15).trim();
            if (tzPart.length >= 5) {
                const sign = tzPart.charAt(0);
                const tzHours = tzPart.substr(1, 2);
                const tzMinutes = tzPart.substr(3, 2);
                const iso = `${year}-${String(month + 1).padStart(2, '0')}-${String(day).padStart(2, '0')}T${String(hour).padStart(2, '0')}:${String(minute).padStart(2, '0')}:${String(second).padStart(2, '0')}${sign}${tzHours}:${tzMinutes}`;
                const d = new Date(iso);
                if (isNaN(d.getTime())) return null;
                return d;
            }
        }

        // No timezone provided: treat as UTC
        const d = new Date(Date.UTC(year, month, day, hour, minute, second));
        if (isNaN(d.getTime())) return null;
        return d;
    }
}

export default new EpgService();
