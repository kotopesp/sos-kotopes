import { User } from "./user.interface";

export interface Chat {
    ID: number,
    ChatType?: string,
    // isdeleted?: boolean,
    Users: User[],
}

export interface ResponseUser {
    id: number;
    username: string;
    firstname?: string;
    lastname?: string;
    description?: string;
    photo?: string;
  }