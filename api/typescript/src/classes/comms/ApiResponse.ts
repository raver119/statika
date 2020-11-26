import {HttpStatusCode} from "./HttpStatusCode";

export class ApiResponse {
    statusCode: HttpStatusCode
    message: string = ""

    constructor(statusCode: HttpStatusCode = HttpStatusCode.OK) {
        this.statusCode = statusCode
    }
}

export const isApiResponse = (obj: any) :obj is ApiResponse => {
    // TODO: add message here into validation as well
    return obj.statusCode !== undefined
}


export const responseOk = () => new ApiResponse(HttpStatusCode.OK)
export const responseUnauthorized = () => new ApiResponse(HttpStatusCode.UNAUTHORIZED)
export const responseError = () => new ApiResponse(HttpStatusCode.INTERNAL_SERVER_ERROR)