import {ApiResponse} from "./ApiResponse";


export class UploadResponse extends ApiResponse{
    filename: string
}

export const isUploadResponse = (obj: any) :obj is UploadResponse => {
    return obj.filename !== undefined
}