
export class HttpException extends Error {
    constructor(public statusText: string,
                public statusCode: number)
    {
        super(statusText);
    }
}