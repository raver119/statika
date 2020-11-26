import {HttpStatusCode} from "./HttpStatusCode";

export class ApiResponse {
    statusCode: HttpStatusCode

    constructor(statusCode: HttpStatusCode = HttpStatusCode.OK) {
        this.statusCode = statusCode
    }
}


export const responseOk = () => new ApiResponse(HttpStatusCode.OK)
export const responseUnauthorized = () => new ApiResponse(HttpStatusCode.UNAUTHORIZED)
export const responseError = () => new ApiResponse(HttpStatusCode.INTERNAL_SERVER_ERROR)