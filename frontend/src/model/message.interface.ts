export interface Message {
  id?: number,
  user_id: number,
  chat_id: number,
  message_content: string,
  audio_bytes: string,
  is_audio: boolean,
  created_at: Date,
  is_user_message: boolean,
  time: string,
  sender_name: string,
}
