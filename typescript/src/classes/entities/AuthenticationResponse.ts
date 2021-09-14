import { ApiResponse } from "./ApiResponse";

export interface AuthenticationResponse extends ApiResponse {
  token: string;
  expires: number;
}

export const isAuthenticationResponse = (obj: any): obj is AuthenticationResponse => {
  return obj.token !== undefined && obj.expires !== undefined;
};
