import { MetaType } from "../../Statika";
import { fromByteArray } from "base64-js";

export const bufferUploadRequest = (
  bucket: string,
  fileName: string,
  buffer: ArrayBuffer,
  meta: MetaType = undefined
) => {
  return {
    filename: fileName,
    bucket: bucket,
    meta: meta,
    payload: fromByteArray(new Uint8Array(buffer)),
  };
};

export type UploadRequest = ReturnType<typeof bufferUploadRequest>;
