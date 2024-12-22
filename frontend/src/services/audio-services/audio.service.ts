import {Injectable} from '@angular/core';
import RecordRTC from "recordrtc";
import * as lame from '@breezystack/lamejs';

@Injectable({
  providedIn: 'root'
})
export class AudioService {
  private recorder!: RecordRTC.StereoAudioRecorder | null;
  private stream!: MediaStream;
  private isRecording = false;
  private mp3Encoder: lame.Mp3Encoder;

  constructor() {
    this.mp3Encoder = new lame.Mp3Encoder(1, 44100, 128);
  }

  convertToMp3(blob: Blob): Promise<Blob> {
    return blob.arrayBuffer()
      .then(buffer => {
        const data = new Int16Array(buffer);
        const uint8Array = this.mp3Encoder.encodeBuffer(data);
        return new Blob([uint8Array]);
      });
  }

  startRecording() {
    if (this.isRecording) return;

    const options = {
      video: false,
      audio: true,
    }

    navigator.mediaDevices.getUserMedia(options).then(stream => {
      this.isRecording = true;
      this.stream = stream;
      this.record();
    }).catch(error => {
      console.error(error);
    })
  }

  private record() {
    this.recorder = new RecordRTC.StereoAudioRecorder(this.stream,
      {
        type: "audio",
        mimeType: "audio/wav",
        numberOfAudioChannels: 1,
        sampleRate: 44100,
        bufferSize: 4096
      }
    );
    this.recorder.record();
  }

  private stopMedia() {
    if (this.recorder) {
      this.recorder = null;
      if (this.stream) {
        this.stream.getAudioTracks().forEach(track => {
          track.stop();
        });
      }
    }
  }

  stopRecording(callback: (blob: Blob) => void) {
    this.isRecording = false;
    this.recorder?.stop((blob: Blob) => {
      this.stopMedia();
      this.convertToMp3(blob)
        .then((blob: Blob) => {
          callback(blob);
        })
        .catch(error => {
          console.log(error);
        })
    });
  }

}
