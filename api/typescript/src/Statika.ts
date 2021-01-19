import {isUploadResponse, UploadResponse} from "./classes/comms/UploadResponse";
import {DatatypeException} from "./classes/exceptions/DatatypeException";
import {EndpointsCoordinates} from "./classes/system/EndpointsCoordinates";
import {ApiResponse, isApiResponse} from "./classes/comms/ApiResponse";
import 'whatwg-fetch'
import {bufferUploadRequest} from "./classes/comms/UploadRequest";
import {communicator} from "./classes/Communicator";
import {AuthenticationBean} from "./classes/api/AuthenticationBean";


export type MetaType = Map<string, string>|undefined

export const Statika = (coords: EndpointsCoordinates) => {
    const comm = communicator(coords)

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
            return comm.post(bean, "/file", req).then(data => {
                if (isUploadResponse(data))
                    return data as UploadResponse

                throw new DatatypeException("UploadResponse", data)
            })
        },

        ping(bean: AuthenticationBean) :Promise<ApiResponse> {
            return comm.get(bean, "/ping").then(data => {
                if (isApiResponse(data))
                    return data as ApiResponse

                throw new DatatypeException("ApiResponse", data)
            })
        },

        deleteFile(bean: AuthenticationBean, fileName: string) :Promise<ApiResponse> {
            const addr = comm.storage().toString()

            if (!fileName.startsWith("/"))
                fileName = `/${fileName}`

            return comm.delete(bean, `${addr}${fileName}`).then(data => {
                if (isApiResponse(data))
                    return data as ApiResponse

                throw new DatatypeException("ApiResponse", data)
            })
        },

        updateMetaInfo(fileName: string, metaInfo: MetaType = undefined) :Promise<ApiResponse> {
            return undefined
        },

        getMetaInfo(fileName: string) :Promise<MetaType> {
            return undefined
        },

        deleteMetaInfo(fileName: string) :Promise<ApiResponse> {
            return undefined
        },

        listFiles() :Promise<ApiResponse> {
            return undefined
        },
    }
}
