import {Communicator} from "./Communicator";
import {AuthenticationBean} from "./AuthenticationBean";
import {isUploadResponse, UploadResponse} from "../entities/UploadResponse";
import {bufferUploadRequest} from "../entities/UploadRequest";
import {DatatypeException} from "../exceptions/DatatypeException";
import {ApiResponse, isApiResponse} from "../entities/ApiResponse";
import {MetaType} from "../../Statika";


export const storageApi = (communicator: Communicator) => {
    return {
        /**
         *
         * @param bean
         * @param fileName - name of the file to be uploaded
         * @param f - ArrayBuffer with actual file content
         * @param metaInfo - optional string/string dictionary to be stored together with file
         */
        uploadFile(bean: AuthenticationBean, fileName: string, f: ArrayBuffer, metaInfo: MetaType = undefined) :Promise<UploadResponse> {
            const req = bufferUploadRequest(bean.bucket, fileName, f, metaInfo)
            return communicator.post(bean, "/file", req).then(data => {
                if (isUploadResponse(data))
                    return data as UploadResponse

                throw new DatatypeException("UploadResponse", data)
            })
        },

        deleteFile(bean: AuthenticationBean, fileName: string) :Promise<ApiResponse> {
            const addr = communicator.storage().toString()

            if (!fileName.startsWith("/"))
                fileName = `/${fileName}`

            return communicator.delete(bean, `${addr}${fileName}`).then(data => {
                if (isApiResponse(data))
                    return data as ApiResponse

                throw new DatatypeException("ApiResponse", data)
            })
        },

        listFiles() :Promise<ApiResponse> {
            return undefined
        },
    }
}

export type StorageApi = ReturnType<typeof storageApi>