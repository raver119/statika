import {ApiResponse} from "./ApiResponse";


export class UploadResponse extends ApiResponse{
    fileName: string
}

export const isUploadResponse = (obj: any) :obj is UploadResponse => {
    return obj.fileName !== undefined
}