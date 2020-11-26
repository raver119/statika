import {ApiResponse} from "./ApiResponse";

export class AuthenticationResponse extends ApiResponse {
    authToken: string
    expires: number
}

export const isAuthenticationResponse = (obj: any) :obj is AuthenticationResponse => {
    return obj.authToken !== undefined && obj.expires !== undefined
}