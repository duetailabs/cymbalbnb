// Copyright 2025 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

const delayShowingVideo = 6000; // milliseconds

function onGenerateVideo() {
    // hide button
    $('#generate-video-btn').addClass("d-none");
    // show spinner
    $('#video-generation-spinner').removeClass("d-none");
    // set timer to reface the component
    setTimeout(function () {

        $('#card-video-generator').addClass("d-none");
        $('#card-video-player').removeClass("d-none");

    }, delayShowingVideo);
}

function onPlayVideoClip() {
    const videoWidget = $('#video-widget');
    const videoOverlay = $('#play-video-btn');
    const videoPlayer = document.getElementById('video-player');

    console.log("playing video clip");

    // hide overlay
    videoOverlay.addClass("d-none");

    // Add controls to the video player
    videoPlayer.controls = true;

    // Attempt to play the video
    videoPlayer.play().catch(error => {
        // Autoplay might be blocked by the browser, controls are still visible
        console.warn("Video autoplay was prevented:", error);
        // The user can still click the native play button
    });

    // Remove the pointer cursor as it's no longer needed for triggering play
    videoWidget.style.cursor = 'default';

    // Remove this event listener after the first click to prevent issues
    // if the user clicks again while the video is playing/paused.
    document.getElementById('play-video-btn').removeEventListener('click', playVideo);
}

function onVideoPlayEnded() {
    const videoWidget = $('#video-widget');
    const videoOverlay = $('#play-video-btn');
    const videoPlayer = $('#video-player');

    console.log("video ended");

    videoPlayer.controls = false;
    videoOverlay.removeClass("d-none");
    videoWidget.style.cursor = 'pointer';
    document.getElementById('play-video-btn').addEventListener('click', onPlayVideoClip);
    videoPlayer.load();
};

$(document).ready(function () {
    if (window.location.href.indexOf("listing") > 0) {
        document.getElementById('generate-video-btn').addEventListener('click', onGenerateVideo);
        document.getElementById('play-video-btn').addEventListener('click', onPlayVideoClip, { once: false });
        document.getElementById('video-player').addEventListener('ended', onVideoPlayEnded);
    }
});
