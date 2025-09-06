import React, { Component, createRef } from 'react';
import { Logger } from '../utils/Logger';

class Playlist extends Component {
    constructor(props) {
        super(props);
        this.channelsRef = createRef();
        this.topTriggerRef = createRef();
        this.bottomTriggerRef = createRef();


        this.state = {
            channelsUrl: '/channels.m3u',
            playlistItems: [],
            images: {},
        };

        if (__DEV__) {
            this.state.channelsUrl = `http://${window.location.hostname}:8080/channels.m3u`;
        }

        this.scrollDir = 0;
        this.scrollSpeed = 0;
        this.animationFrameId = null;
    }

    componentDidMount() {
        this.fetchPlaylist();
        this.setupScrolling();
    }

    componentWillUnmount() {
        this.cleanupScrolling();
    }

    fetchPlaylist = () => {
        const username = localStorage.getItem('username');
        const password = localStorage.getItem('password');
        const headers = { Authorization: 'Basic ' + btoa(`${username}:${password}`) };

        Logger.info('Fetching playlist from URL:' + this.state.channelsUrl);
        fetch(this.state.channelsUrl, { headers })
            .then(response => response.text())
            .then(data => {
                const lines = data.trim().split('\n');
                const items = [];
                let item = {};

                let number = 0;
                lines.forEach((line) => {
                    if (line.startsWith("#EXTINF:")) {
                        if (item.source) items.push(item);
                        item = { id: '', name: '', logo: '', source: '', number: number++ };
                        const name = line.split(',')[1];
                        item.name = name;
                        const logoMatch = line.match(/tvg-logo="([^"]+)"/);
                        if (logoMatch) item.logo = logoMatch[1];
                        const idMatch = line.match(/tvg-id="([^"]+)"/);
                        if (idMatch) item.id = idMatch[1];
                    } else if (line && !line.startsWith("#")) {
                        item.source = line;
                    }
                });
                if (item.source) items.push(item);

                if (items.length === 0) {
                    Logger.error('No channels available');
                    this.setState({ playlistItems: [] });
                    return;
                }

                this.setState({ playlistItems: items });

                let previousChannelIndex = parseInt(localStorage.getItem('current_channel_index')) || 0;
                if (previousChannelIndex >= items.length) {
                    previousChannelIndex = 0;
                }
                if (previousChannelIndex < 0) {
                    previousChannelIndex = items.length - 1;
                }
                this.props.setCurrentChannel(items[previousChannelIndex]);
            })
            .catch((error) => {
                Logger.error('Error fetching playlist:' + error);
                this.setState({ playlistItems: [] });
            });
    };

    channelDown = () => {
        const { playlistItems } = this.state;
        if (playlistItems.length === 0) return;
        let previousChannelIndex = parseInt(localStorage.getItem('current_channel_index')) || 0;
        previousChannelIndex = (previousChannelIndex - 1 + playlistItems.length) % playlistItems.length;
        localStorage.setItem('current_channel_index', previousChannelIndex);
        this.props.setCurrentChannel(playlistItems[previousChannelIndex]);
    };

    channelUp = () => {
        const { playlistItems } = this.state;
        if (playlistItems.length === 0) return;
        let previousChannelIndex = parseInt(localStorage.getItem('current_channel_index')) || 0;
        previousChannelIndex = (previousChannelIndex + 1) % playlistItems.length;
        localStorage.setItem('current_channel_index', previousChannelIndex);
        this.props.setCurrentChannel(playlistItems[previousChannelIndex]);
    };

    changeChannel = (index) => {
        const { playlistItems } = this.state;
        if (index < 0 || index >= playlistItems.length) {
            Logger.error('Channel index out of range:', index);
            return;
        }
        localStorage.setItem('current_channel_index', index);
        this.props.setCurrentChannel(playlistItems[index]);
        Logger.info('Channel changed to:', playlistItems[index].name);
    };

    onChannelClick = (item) => {
        const { playlistItems } = this.state;
        const channelIndex = playlistItems.findIndex(channel => channel.source === item.source);
        if (channelIndex !== -1) {
            localStorage.setItem('current_channel_index', channelIndex);
            this.props.setCurrentChannel(item);
            Logger.info('Channel clicked:', item.name);
        } else {
            Logger.error('Channel not found in playlist:', item.name);
        }
    };

    setupScrolling = () => {
        const scrollContent = () => {
            if (!this.channelsRef.current) {
                this.animationFrameId = requestAnimationFrame(scrollContent);
                return;
            }
            if ((this.scrollDir === -1 || this.scrollDir === 1)) {
                this.scrollSpeed += this.scrollDir * 0.5;
                if (Math.abs(this.scrollSpeed) > 20) {
                    this.scrollSpeed = 20 * this.scrollDir;
                }
            } else {
                this.scrollSpeed *= 0.9; // Decelerate scrolling
            }
            if (this.scrollSpeed !== 0) {
                this.channelsRef.current.scrollBy({ top: this.scrollSpeed, behavior: 'auto' });
            }
            this.animationFrameId = requestAnimationFrame(scrollContent);
        };

        scrollContent();

        const handleMouseEnterTop = () => {
            this.scrollDir = -1;
        };

        const handleMouseEnterBottom = () => {
            this.scrollDir = 1;
        };

        const handleMouseLeave = () => {
            this.scrollDir = 0;
        };

        if (this.topTriggerRef.current && this.bottomTriggerRef.current) {
            this.topTriggerRef.current.addEventListener('mouseenter', handleMouseEnterTop);
            this.bottomTriggerRef.current.addEventListener('mouseenter', handleMouseEnterBottom);
            this.topTriggerRef.current.addEventListener('mouseleave', handleMouseLeave);
            this.bottomTriggerRef.current.addEventListener('mouseleave', handleMouseLeave);
        }

        if (this.channelsRef.current) {
            if (navigator.userAgent.toLowerCase().indexOf('firefox') > -1) {
                this.channelsRef.current.addEventListener('wheel', (event) => {
                    this.scrollSpeed = event.deltaY;
                });
            } else {
                this.channelsRef.current.addEventListener('mousewheel', (event) => {
                    this.scrollSpeed = event.deltaY;
                });
            }
        }

    };

    cleanupScrolling = () => {
        if (this.animationFrameId) {
            cancelAnimationFrame(this.animationFrameId);
        }
    };

    render() {
        const { playlistItems, images } = this.state;

        return (
            <div className="playlist">
                <div ref={this.topTriggerRef} className="scroll scroll_down">
                    <span>
                        <i className="bi bi-arrow-up-square-fill"></i>
                    </span>
                </div>
                <div className="channels" ref={this.channelsRef}>
                    {playlistItems.map((item, index) => (
                        <div
                            key={index}
                            className="channel"
                            onClick={() => this.onChannelClick(item)}
                        >
                            {item.logo && images[item.logo] ? (
                                <img src={images[item.logo]} alt={item.name} className="logo" />
                            ) : (
                                <span className="title">{item.name}</span>
                            )}
                        </div>
                    ))}
                </div>
                <div ref={this.bottomTriggerRef} className="scroll scroll_up">
                    <span>
                        <i className="bi bi-arrow-down-square-fill"></i>
                    </span>
                </div>
            </div>
        );
    }
}

export default Playlist;
