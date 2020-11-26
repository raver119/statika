import {MetaType} from "../../Statika";
import btoa from "btoa";

export interface UploadRequest {
    filename: string
    bucket: string
    payload: string
    meta: MetaType
}


export const bufferUploadRequest = (bucket: string, fileName: string, buffer: ArrayBuffer, meta: MetaType = undefined) :UploadRequest => {
    return {
        filename: fileName,
        bucket: bucket,
        meta: meta,
        payload: btoa(String.fromCharCode(...new Uint8Array(buffer)))
    }
}