import { User } from "./user.interface";

export interface Chat {
    id: number,
    chat_type?: string,
    users: User[],
    title?: string,
    last_message?: {
        message_content: string,
        created_at: Date,
        time?: string,
        user_id: number,
        is_read: boolean,
        sender_name: string,
      };
    unread_count: number;
    created_at: Date;
}

export interface ResponseUser {
    id: number;
    username: string;
    firstname?: string;
    lastname?: string;
    description?: string;
    photo?: string;
  }