import { Timestamp } from "rxjs";

export interface Message {
    ID?: number,
    UserID: number,
    ChatID: number,
    Content: string,
    CreatedAt: Date,
    UpdatedAt: Date,
    isUserMessage: boolean,
    time: string,
}
