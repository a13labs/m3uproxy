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
            const programData = this.parseEpgXml(xmlText);
            
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
        
        // Get timezone offset
        const timezoneOffset = -date.getTimezoneOffset();
        const offsetHours = Math.floor(Math.abs(timezoneOffset) / 60);
        const offsetMinutes = Math.abs(timezoneOffset) % 60;
        const offsetSign = timezoneOffset >= 0 ? '+' : '-';
        const timezone = `${offsetSign}${String(offsetHours).padStart(2, '0')}${String(offsetMinutes).padStart(2, '0')}`;
        
        return `${year}${month}${day}${hours}${minutes}${seconds} ${timezone}`;
    }

    parseEpgXml(xmlText) {
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

            // Get the first (current) programme
            const programme = programmes[0];
            
            const titleElement = programme.querySelector('title');
            const descElement = programme.querySelector('desc');
            const startTime = programme.getAttribute('start');
            const stopTime = programme.getAttribute('stop');

            if (!titleElement || !startTime || !stopTime) {
                Logger.warn('Incomplete programme data in EPG');
                return null;
            }

            return {
                title: titleElement.textContent,
                description: descElement ? descElement.textContent : '',
                startTime: this.parseXmltvTime(startTime),
                endTime: this.parseXmltvTime(stopTime)
            };
        } catch (error) {
            Logger.error('Error parsing EPG XML:', error);
            return null;
        }
    }

    parseXmltvTime(xmltvTime) {
        // Parse XMLTV time format: "20250906185900 +0000"
        if (!xmltvTime) return null;
        
        const timeString = xmltvTime.trim();
        const year = parseInt(timeString.substr(0, 4));
        const month = parseInt(timeString.substr(4, 2)) - 1; // months are 0-indexed
        const day = parseInt(timeString.substr(6, 2));
        const hour = parseInt(timeString.substr(8, 2));
        const minute = parseInt(timeString.substr(10, 2));
        const second = parseInt(timeString.substr(12, 2));
        
        // Handle timezone offset if present
        let date = new Date(year, month, day, hour, minute, second);
        
        if (timeString.length > 14) {
            const tzPart = timeString.substr(15); // Skip the space
            if (tzPart.length >= 5) {
                const sign = tzPart.charAt(0);
                const tzHours = parseInt(tzPart.substr(1, 2));
                const tzMinutes = parseInt(tzPart.substr(3, 2));
                
                let offsetMinutes = tzHours * 60 + tzMinutes;
                if (sign === '-') offsetMinutes = -offsetMinutes;
                
                // Adjust for timezone
                date = new Date(date.getTime() - offsetMinutes * 60 * 1000);
            }
        }
        
        return date;
    }
}

export default new EpgService();
