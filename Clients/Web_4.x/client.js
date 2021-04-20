var rtc = {
    // For the local client.
    client: null,
    // For the local audio and video tracks.
    localAudioTrack: null,
    localVideoTrack: null,
};

var options = {
    // Pass your app ID here.
    appId: "App ID",
    // Set the channel name.
    channel: "ChannelA",
    // Pass a token if your project enables the App Certificate.
    token: null,
    // Set the user role in the channel.
    role: "audience"
};

function fetchToken(uid, channelName, tokenRole) {

    return new Promise(function (resolve) {
        axios.post('http://<URL>/fetch_rtc_token', {
            uid: uid,
            channelName: channelName,
            role: tokenRole
        }, {
            headers: {
                'Content-Type': 'application/json; charset=UTF-8'
            }
        })
            .then(function (response) {
                const token = response.data.token;
                resolve(token);
            })
            .catch(function (error) {
                console.log(error);
            });
    })
}


async function startBasicCall() {

    const client = AgoraRTC.createClient({ mode: "live", codec: "vp8" });
    client.setClientRole(options.role);
    let token = await fetchToken(123456, "ChannelA", 1);
    await client.join(options.appId, options.channel, token, 123456);
    // Create an audio track from the audio sampled by a microphone.
    localAudioTrack = await AgoraRTC.createMicrophoneAudioTrack();
    // Create a video track from the video captured by a camera.
    rtc.localVideoTrack = await AgoraRTC.createCameraVideoTrack();
    // Publish the local audio and video tracks to the channel.
    await client.publish([localAudioTrack, rtc.localVideoTrack]);
     // Dynamically create a container in the form of a DIV element for playing the remote video track.
     const localplayerContainer = document.createElement("div");
     // Specify the ID of the DIV container. You can use the `uid` of the remote user.
     localplayerContainer.id = "123456";
     localplayerContainer.style.width = "640px";
     localplayerContainer.style.height = "480px";
     document.body.append(localplayerContainer);

     rtc.localVideoTrack.play(localplayerContainer);

    console.log("publish success!");

    client.on("user-published", async (user, mediaType) => {
        // Subscribe to a remote user.
        await client.subscribe(user, mediaType);
        console.log("subscribe success");


        // If the subscribed track is video.
        if (mediaType === "video") {
            // Get `RemoteVideoTrack` in the `user` object.
            const remoteVideoTrack = user.videoTrack;
            // Dynamically create a container in the form of a DIV element for playing the remote video track.
            const playerContainer = document.createElement("div");
            // Specify the ID of the DIV container. You can use the `uid` of the remote user.
            playerContainer.id = user.uid.toString();
            playerContainer.style.width = "640px";
            playerContainer.style.height = "480px";
            document.body.append(playerContainer);


            // Play the remote video track.
            // Pass the DIV container and the SDK dynamically creates a player in the container for playing the remote video track.
            remoteVideoTrack.play(playerContainer);

            // Or just pass the ID of the DIV container.
            // remoteVideoTrack.play(playerContainer.id);
        }

        // If the subscribed track is audio.
        if (mediaType === "audio") {
            // Get `RemoteAudioTrack` in the `user` object.
            const remoteAudioTrack = user.audioTrack;
            // Play the audio track. No need to pass any DOM element.
            remoteAudioTrack.play();
        }

        client.on("user-unpublished", user => {
            // Get the dynamically created DIV container.
            const playerContainer = document.getElementById(user.uid);
            // Destroy the container.
            playerContainer.remove();
        });

    });

}

startBasicCall()




