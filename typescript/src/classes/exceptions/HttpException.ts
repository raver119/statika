
export class HttpException extends Error {
    statusText: string
    statusCode: number

    constructor(statusText: string,
                statusCode: number) {

        let message = `API response code: ${statusCode}`
        if (statusText !== undefined && statusText.length > 0)
            message = `${message}; Message: ${statusText}`

        super(message);

        this.statusText = statusText
        this.statusCode = statusCode
    }
}