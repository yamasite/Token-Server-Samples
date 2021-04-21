var rtc = {
    // For the local audio and video tracks.
    localAudioTrack: null,
    localVideoTrack: null,
};

var options = {
    // Pass your app ID here.
    appId: "Your_App_ID",
    // Set the channel name.
    channel: "ChannelA",
    // Set the user role in the channel.
    role: "host"
};

function fetchToken(uid, channelName, tokenRole) {

    return new Promise(function (resolve) {
        axios.post('http://10.53.3.234:8082/fetch_rtc_token', {
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
    uid = 123456;
    let token = await fetchToken(uid, options.channel, 1);
    await client.join(options.appId, options.channel, token, uid);
    // Create an audio track from the audio sampled by a microphone.
    localAudioTrack = await AgoraRTC.createMicrophoneAudioTrack();
    // Create a video track from the video captured by a camera.
    rtc.localVideoTrack = await AgoraRTC.createCameraVideoTrack();
    // Publish the local audio and video tracks to the channel.
    await client.publish([localAudioTrack, rtc.localVideoTrack]);
    // Dynamically create a container in the form of a DIV element for playing the remote video track.
    const localplayerContainer = document.createElement("div");
    // Specify the ID of the DIV container. You can use the `uid` of the remote user.
    localplayerContainer.id = uid;
    localplayerContainer.textContent = "Local user " + uid;
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
            const remotePlayerContainer = document.createElement("div");
            // Specify the ID of the DIV container. You can use the `uid` of the remote user.
            remotePlayerContainer.id = user.uid.toString();
            remotePlayerContainer.textContent = "Remote user " + user.uid.toString();
            remotePlayerContainer.style.width = "640px";
            remotePlayerContainer.style.height = "480px";
            document.body.append(remotePlayerContainer);


            // Play the remote video track.
            // Pass the DIV container and the SDK dynamically creates a player in the container for playing the remote video track.
            remoteVideoTrack.play(remotePlayerContainer);

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
            const remotePlayerContainer = document.getElementById(user.uid);
            // Destroy the container.
            remotePlayerContainer.remove();
        });

    });

    client.on("token-privilege-will-expire", async function () {
        let token = await fetchToken(uid, options.channel, 1);
        await client.renewToken(token);
    });

    client.on("token-privilege-did-expire", async function () {
        let token = await fetchToken(uid, options.channel, 1);
        await rtc.client.join(options.appId, options.channel, token, 123456);
    });

}

startBasicCall()




