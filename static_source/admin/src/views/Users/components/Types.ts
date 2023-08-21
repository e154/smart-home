import {ApiImage, ApiRole} from "@/api/stub";

export interface User {
    nickname?: string;
    firstName?: string;
    lastName?: string;
    password?: string;
    passwordRepeat?: string;
    email?: string;
    status?: string;
    lang?: string;
    image?: ApiImage;
    imageId?: number;
    role?: ApiRole;
    roleName?: string;
}
