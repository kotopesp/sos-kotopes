import {AfterViewInit, Component, ElementRef, Input, ViewChild} from '@angular/core';
import {NgIf} from '@angular/common';
import WaveSurfer from "wavesurfer.js";

@Component({
  selector: 'app-message',
  standalone: true,
  imports: [NgIf],
  templateUrl: './message.component.html',
  styleUrl: './message.component.scss'
})
export class MessageComponent implements AfterViewInit {
  @Input()
  answer!: boolean;
  @Input()
  messageContent!: string;
  @Input()
  messageTime!: string;
  @Input()
  name!: string;
  @Input()
  isAudio!: boolean;
  @Input()
  audioBytes!: string | null;

  audioBlobUrl = '';

  @ViewChild("wave") container!: ElementRef;
  waveSurfer!: WaveSurfer;
  waveSurferReady = false;
  audioPlaying = false;
  currentPlaybackRate = 1;

  ngAfterViewInit(): void {
    if (this.isAudio) {
      const options = {
        "container": this.container.nativeElement,
        url: this.getAudioUrl()!,
        "cursorColor": "#333",
        "progressColor": "#ffffff",
        "waveColor": "#8c8c8c",
        "barGap": 1,
        "barHeight": 1,
        "barMinHeight": 1,
        "barRadius": 4,
        "barWidth": 5,
        "height": 50,
        "fillParent": true,
      }

      this.waveSurfer = WaveSurfer.create(options);

      this.waveSurfer.on('ready', () => {
        this.waveSurferReady = true;
      });

      this.waveSurfer.on('pause', () => {
        this.audioPlaying = false;
      });
    }
  }

  playAudio() {
    console.log("Play audio");
    this.waveSurfer.playPause();
    this.audioPlaying = this.waveSurfer.isPlaying();
  }

  getAudioUrl() {
    const blob = this.getBlob();
    if (blob !== null) {
      if (this.audioBlobUrl === '') {
        console.log("Creating new url for audio blob");
        this.audioBlobUrl = URL.createObjectURL(blob);
      }
      return this.audioBlobUrl;
    } else {
      return null;
    }
  }

  speedUp() {
    if (this.audioPlaying) {
      const playbackRate = this.waveSurfer.getPlaybackRate();
      if (playbackRate < 2) {
        this.currentPlaybackRate = playbackRate + 0.25;
        this.waveSurfer.setPlaybackRate(this.currentPlaybackRate);
      }
    }
  }

  speedDown() {
    if (this.audioPlaying) {
      const playbackRate = this.waveSurfer.getPlaybackRate();
      if (playbackRate > 0.25) {
        this.currentPlaybackRate = playbackRate - 0.25;
        this.waveSurfer.setPlaybackRate(this.currentPlaybackRate);
      }
    }
  }

  getBlob(): Blob | null {
    if (this.audioBytes === null) {
      return null;
    }
    const arrayBuffer = this.base64ToArrayBuffer(this.audioBytes);
    return new Blob([arrayBuffer])
  }

  private base64ToArrayBuffer(base64: string): ArrayBuffer {
    const binaryString = atob(base64);
    const bytes = new Uint8Array(binaryString.length);
    for (let i = 0; i < binaryString.length; i++) {
      bytes[i] = binaryString.charCodeAt(i);
    }
    return bytes.buffer;
  }

  resetPlaybackRate() {
    if (this.audioPlaying) {
      this.currentPlaybackRate = 1;
      this.waveSurfer.setPlaybackRate(1);
    }
  }
}
