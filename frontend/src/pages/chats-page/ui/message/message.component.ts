import {Component, Input} from '@angular/core';
import {NgIf} from '@angular/common';

@Component({
  selector: 'app-message',
  standalone: true,
  imports: [NgIf],
  templateUrl: './message.component.html',
  styleUrl: './message.component.scss'
})
export class MessageComponent {
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
}
