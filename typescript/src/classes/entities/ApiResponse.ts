import { HttpStatusCode } from "./HttpStatusCode";

export interface ApiResponse {
  statusCode: number;
  message: string;
}

export const response = (code: number, message = ""): ApiResponse => {
  return {
    statusCode: code,
    message: message,
  };
};

export const isApiResponse = (obj: any): obj is ApiResponse => {
  // TODO: add message here into validation as well
  return obj !== undefined && obj.statusCode !== undefined;
};

export const responseOk = () => response(HttpStatusCode.OK);
export const responseUnauthorized = () => response(HttpStatusCode.UNAUTHORIZED);
export const responseError = () => response(HttpStatusCode.INTERNAL_SERVER_ERROR);
