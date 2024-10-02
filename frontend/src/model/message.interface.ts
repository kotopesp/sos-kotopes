export interface Message {
    id?: number,
    user_id: number,
    chat_id: number,
    message_content: string,
    created_at: Date,
    // UpdatedAt: Date,
    is_user_message: boolean,
    time: string,
    sender_name: string,
}