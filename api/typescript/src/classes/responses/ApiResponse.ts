
export class ApiResponse {
    statusCode: number

    constructor(statusCode: number = 200) {
        this.statusCode = statusCode
    }
}