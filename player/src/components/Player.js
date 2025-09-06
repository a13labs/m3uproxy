import React, { Component, createRef } from 'react';
class Player extends Component {

    constructor(props) {
        super(props);
        this.video = createRef();
        this.videoContainerRef = createRef();
        this.player = null;

        this.state = {
            EMESupported: typeof window.MediaKeys !== "undefined" && typeof window.navigator.requestMediaKeySystemAccess === "function",
            isMobile: navigator.userAgent.toLowerCase().indexOf('mobile') !== -1,
            licensingUrl: '/drm/licensing',
        };

        if (__DEV__) {
            this.state.licensingUrl = `http://${window.location.hostname}:8080/drm/licensing`;
        }
    }

    async componentDidMount() {

        let ShakaPlayer, ShakaPolyfill;

        if (__DEV__) {
            // Dynamically import the debug version in development
            const shaka = await import("shaka-player/dist/shaka-player.ui.debug");
            ShakaPlayer = shaka.Player;
            ShakaPolyfill = shaka.polyfill;
        } else {
            // Dynamically import the production version in production
            const shaka = await import("shaka-player/dist/shaka-player.ui");
            ShakaPlayer = shaka.Player;
            ShakaPolyfill = shaka.polyfill;
        }

        if (this.state.EMESupported) {
            console.log('EME Supported, skipping encrypted content');
            ShakaPolyfill.installAll();
            const keys = {};

            const username = localStorage.getItem('username');
            const password = localStorage.getItem('password');
            const headers = { Authorization: 'Basic ' + btoa(`${username}:${password}`) };

            console.log('Loading licensing from server: ' + this.state.licensingUrl);
            fetch(this.state.licensingUrl, { headers })
                .then(response => response.text())
                .then(data => {
                    const response = JSON.parse(data);
                    for (const key of response.keys) {
                        if (key.kid.length !== 32 || key.k.length !== 32) {
                            console.error('Invalid key:', key);
                            continue;
                        }
                        keys[key.kid] = key.k;
                    }
                    this.player.configure({
                        drm: {
                            clearKeys: keys
                        }
                    });
                })
                .catch(() => console.error('Error fetching keys from server'));
        } else {
            console.log('EME Not Supported');
        }

        this.player = new ShakaPlayer();
        this.player.addEventListener('error', (event) => {
            console.error('Shaka Player Error:', event.detail);
        });

        // Attach the player to the video element
        this.player.attach(this.video.current).then(() => {
            const { onReady } = this.props;
            console.log('Player attached to video element');
            
            // Add event listeners for video events
            this.video.current.addEventListener('loadstart', this.handleLoadStart);
            this.video.current.addEventListener('canplay', this.handleCanPlay);
            this.video.current.addEventListener('timeupdate', this.handleTimeUpdate);
            
            if (onReady) {
                onReady();
            }
        }).catch((err) => {
            console.error('Error attaching player:', err);
        });
    }


    componentWillUnmount() {
        if (this.video.current) {
            this.video.current.removeEventListener('loadstart', this.handleLoadStart);
            this.video.current.removeEventListener('canplay', this.handleCanPlay);
            this.video.current.removeEventListener('timeupdate', this.handleTimeUpdate);
        }
        if (this.player) {
            this.player.destroy();
        }
    }

    handlePlayerError = (err) => {
        const { onError } = this.props;

        if (err.code === 1001) {
            console.error('DRM error:', err);
        } else if (err.code === 1002) {
            console.error('Media error:', err);
        } else if (err.code === 1004) {
            console.error('Playback error:', err);
        } else if (err.code === 1003) {
            console.error('Manifest error:', err);
        } else if (err.code === 1006) {
            console.error('Key error:', err);
        } else if (err.code === 1007) {
            console.error('License error:', err);
        } else if (err.code === 1008) {
            console.error('Network error:', err);
        } else {
            console.error('Generic error:', err);
        }
        if (onError) {
            onError(err);
        }
    };

    handleLoadStart = () => {
        const { onLoadStart } = this.props;
        console.log('Video load started');
        if (onLoadStart) {
            onLoadStart();
        }
    };

    handleCanPlay = () => {
        const { onCanPlay } = this.props;
        console.log('Video can play');
        if (onCanPlay) {
            onCanPlay();
        }
    };

    handleTimeUpdate = () => {
        const { onTimeUpdate } = this.props;
        if (this.video.current && onTimeUpdate) {
            const currentTime = this.video.current.currentTime;
            onTimeUpdate(currentTime);
        }
    };

    load = (source) => {
        const { onLoad, onPlay } = this.props;

        if (source === undefined || source === null) {
            return;
        }
        if (typeof source !== 'string') {
            return;
        }
        if (source.length === 0) {
            return;
        }

        if (this.player) {
            this.player.load(source).then(() => {
                console.log('Video loaded:', source);
                if (onLoad) {
                    onLoad();
                }
                this.video.current.play().then(() => {
                    if (onPlay) {
                        onPlay();
                    }
                }).catch((err) => {
                    this.handlePlayerError(err);
                });
            }).catch((err) => {
                this.handlePlayerError(err);
            });
        }
    }

    render() {
        return (
            <div ref={this.videoContainerRef} className="player-container">
                <video
                    ref={this.video}
                    className="player"
                />
            </div>
        );
    }
}

export default Player;