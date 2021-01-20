import {MetaType} from "../../Statika";
import btoa from "btoa";


export const bufferUploadRequest = (bucket: string, fileName: string, buffer: ArrayBuffer, meta: MetaType = undefined) => {
    return {
        filename: fileName,
        bucket: bucket,
        meta: meta,
        payload: btoa(String.fromCharCode(...new Uint8Array(buffer)))
    }
}

export type UploadRequest = ReturnType<typeof bufferUploadRequest>