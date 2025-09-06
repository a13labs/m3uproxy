import React, { Component, createRef } from 'react';

class Overlay extends Component {
    constructor(props) {
        super(props);
        this.channelNameRef = createRef();
        this.channelNumberRef = createRef();
        this.programInfoRef = createRef();

        this.state = {
            channelNumberVisible: true,
            channelNameVisible: true,
            programInfoVisible: false,
            channelName: "",
            channelNumber: 0,
            programTitle: "",
            programDescription: "",
            programStartTime: "",
            programEndTime: "",
            programDuration: "",
            programProgress: 0,
        };
    }

    showChannelNumber = (visible = true) => {
        this.setState({ channelNumberVisible: visible });
    }

    hideChannelNumber = () => {
        this.setState({ channelNumberVisible: false });
    }

    showChannelName = (visible = true) => {
        this.setState({ channelNameVisible: visible });
    }

    hideChannelName = () => {
        this.setState({ channelNameVisible: false });
    }

    showProgramInfo = (visible = true) => {
        this.setState({ programInfoVisible: visible });
    }

    hideProgramInfo = () => {
        this.setState({ programInfoVisible: false });
    }

    setChannelName = (name) => {
        this.setState({ channelName: name });
    }

    setChannelNumber = (number) => {
        this.setState({ channelNumber: number });
    }

    setProgramInfo = (programData) => {
        if (!programData) {
            this.setState({
                programTitle: "",
                programDescription: "",
                programStartTime: "",
                programEndTime: "",
                programDuration: "",
                programProgress: 0,
            });
            return;
        }

        const { title, description, startTime, endTime } = programData;
        const start = new Date(startTime);
        const end = new Date(endTime);
        const duration = Math.round((end - start) / 1000 / 60); // duration in minutes
        const now = new Date();
        const progress = Math.min(100, Math.max(0, (now - start) / (end - start) * 100));

        this.setState({
            programTitle: title || "",
            programDescription: description || "",
            programStartTime: start.toLocaleTimeString('pt-PT', { hour: '2-digit', minute: '2-digit' }),
            programEndTime: end.toLocaleTimeString('pt-PT', { hour: '2-digit', minute: '2-digit' }),
            programDuration: `${duration}min`,
            programProgress: progress,
        });
    }

    formatTime = (timeString) => {
        // Parse XMLTV time format: "20250906185900 +0000"
        if (!timeString) return "";
        
        const year = parseInt(timeString.substr(0, 4));
        const month = parseInt(timeString.substr(4, 2)) - 1; // months are 0-indexed
        const day = parseInt(timeString.substr(6, 2));
        const hour = parseInt(timeString.substr(8, 2));
        const minute = parseInt(timeString.substr(10, 2));
        const second = parseInt(timeString.substr(12, 2));
        
        const date = new Date(year, month, day, hour, minute, second);
        return date.toLocaleTimeString('pt-PT', { hour: '2-digit', minute: '2-digit' });
    }

    render() {
        const {
            channelNumberVisible,
            channelNameVisible,
            programInfoVisible,
            channelName,
            channelNumber,
            programTitle,
            programDescription,
            programStartTime,
            programEndTime,
            programDuration,
            programProgress
        } = this.state;

        return (
            <div className="overlay">
                <div ref={this.channelNameRef} className="channel-name" style={{
                    opacity: channelNameVisible ? 1 : 0,
                }} >{channelName}</div>
                <div ref={this.channelNumberRef} className="channel-number" style={{
                    opacity: channelNumberVisible ? 1 : 0,
                }} >{channelNumber}</div>
                
                {programInfoVisible && (
                    <div ref={this.programInfoRef} className="program-info">
                        <div className="program-title">{programTitle}</div>
                        <div className="program-time">
                            {programStartTime} - {programEndTime} ({programDuration})
                        </div>
                        <div className="program-progress">
                            <div className="progress-bar">
                                <div 
                                    className="progress-fill" 
                                    style={{ width: `${programProgress}%` }}
                                ></div>
                            </div>
                        </div>
                        {programDescription && (
                            <div className="program-description">{programDescription}</div>
                        )}
                    </div>
                )}
            </div>
        );
    }
}

export default Overlay;
