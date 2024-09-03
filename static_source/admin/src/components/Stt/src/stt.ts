/*
 * This file is part of the Smart Home
 * Program complex distribution https://github.com/e154/smart-home
 * Copyright (C) 2024, Filippov Alex
 *
 * This library is free software: you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public
 * License as published by the Free Software Foundation; either
 * version 3 of the License, or (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
 * Library General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public
 * License along with this library.  If not, see
 * <https://www.gnu.org/licenses/>.
 */

import {ref} from "vue";
import stream from "@/api/stream";
import {UUID} from "uuid-generator-ts";

const mediaRecorder = ref()

export const startRecording = () => {
  navigator.mediaDevices
    .getUserMedia({audio: true})
    .then((stream) => {
      mediaRecorder.value = new MediaRecorder(stream);
      const recordedChunks: Blob[] = [];

      mediaRecorder.value.addEventListener('dataavailable', (e) => {
        if (e.data.size > 0) {
          recordedChunks.push(e.data);
        }
      });

      mediaRecorder.value.addEventListener('stop', async () => {
        console.log('>> Stop recording...');
        await convertAndSend(recordedChunks)
      });

      console.log('>> Start recording...')
      mediaRecorder.value.start();
    })
    .catch((error) => {
      console.error('Error accessing microphone:', error);
    });
}

export const stopRecording = () => {
  mediaRecorder.value?.stop();
};

// const downloadRecording = (data) => {
//   const downloadLink = document.createElement('a');
//   downloadLink.href = URL.createObjectURL(new Blob([data], {type: 'audio/raw'}));
//   downloadLink.download = 'recorded_audio.wav';
//   downloadLink.click();
// };

const sendToServer = (data) => {

  console.log('>> Sending to Server...');

  const reader = new FileReader();
  reader.readAsDataURL(new Blob([data], {type: 'audio/raw'})); // convert Blob to base64

  reader.onload = function () {
    stream.send({
      id: UUID.createUUID(),
      query: 'event_stt',
      body: btoa(JSON.stringify({payload: reader.result}))
    });

    console.log('>> Audio has been sent');
  }
}

// https://github.com/alphacep/vosk-server/pull/171
const AudioContext = window.AudioContext || window.webkitAudioContext;
const audioCtx = new AudioContext();
const TARGET_SAMPLE_RATE = 16000;

const floatTo16BitPCM = (output, offset, input) => {
  for (let i = 0; i < input.length; i++ , offset += 2) {
    const s = Math.max(-1, Math.min(1, input[i]));
    output.setInt16(offset, s < 0 ? s * 0x8000 : s * 0x7FFF, true);
  }
};

const encodeRAW = (samples) => {
  const buffer = new ArrayBuffer(samples.length * 2);
  const view = new DataView(buffer);
  floatTo16BitPCM(view, 0, samples);
  return view;
};

// https://stackoverflow.com/questions/27598270/resample-audio-buffer-from-44100-to-16000
const convertAndSend = async (recordedChunks) => {
  console.log('>> Converting audio...');

  // directly received by the audioprocess event from the microphone in the browser
  const arrayBuffer = await recordedChunks[0].arrayBuffer();
  const sourceAudioBuffer = await audioCtx.decodeAudioData(arrayBuffer);
  const offlineCtx = new OfflineAudioContext(sourceAudioBuffer.numberOfChannels,
    sourceAudioBuffer.duration * sourceAudioBuffer.numberOfChannels * TARGET_SAMPLE_RATE, TARGET_SAMPLE_RATE);
  const buffer = offlineCtx.createBuffer(sourceAudioBuffer.numberOfChannels, sourceAudioBuffer.length,
    sourceAudioBuffer.sampleRate);
  // Copy the source data into the offline AudioBuffer
  for (let channel = 0; channel < sourceAudioBuffer.numberOfChannels; channel++) {
    buffer.copyToChannel(sourceAudioBuffer.getChannelData(channel), channel);
  }

  // Play it from the beginning.
  const source = offlineCtx.createBufferSource();
  source.buffer = sourceAudioBuffer;
  source.connect(offlineCtx.destination);
  source.start(0);

  offlineCtx.oncomplete = (e) => {
    // use resampled.getChannelData(x) to get an Float32Array for channel x.
    const resampled = e.renderedBuffer;
    // use this float32array to send the samples to the server or whatever
    const leftFloat32Array = resampled.getChannelData(0);
    const data = encodeRAW(leftFloat32Array);
    sendToServer(data);
    // downloadRecording(data);
  }
  offlineCtx.startRendering();
};
