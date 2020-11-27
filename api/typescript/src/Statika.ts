import {isUploadResponse, UploadResponse} from "./classes/comms/UploadResponse";
import {AuthenticationException} from "./classes/exceptions/AuthenticationException";
import {HttpException} from "./classes/exceptions/HttpException";
import {isAuthenticationResponse} from "./classes/comms/AuthenticationResponse";
import {DatatypeException} from "./classes/exceptions/DatatypeException";
import {EndpointsCoordinates} from "./classes/system/EndpointsCoordinates";
import {ApiResponse, isApiResponse} from "./classes/comms/ApiResponse";
import 'whatwg-fetch'
import {bufferUploadRequest} from "./classes/comms/UploadRequest";


export type MetaType = Map<string, string>|undefined
type AuthType = string

export class Statika {
    protected storage: EndpointsCoordinates

    protected uploadToken: string | undefined
    protected assignedBucket: string | undefined

    /**
     *
     * @param storage - address of the
     * @param token - Upload token, usually generated in backend code on the fly, and fused into frontend app
     * @param bucket - Optional folder for splitting end users
     */
    constructor(storage: EndpointsCoordinates, token: string = undefined, bucket: string = undefined) {
        this.uploadToken = token
        this.assignedBucket = bucket
        this.storage = storage
    }

    protected post(authToken: AuthType, url: string, obj: any) :Promise<any> {
        let addr = this.storage.toString()

        if (!url.startsWith("/"))
            url = `/${url}`

        return fetch(`${addr}/rest/v1${url}`, {
            method: 'POST',
            credentials: "same-origin",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
                'Authorization': isAuthenticationResponse(authToken) ? authToken.token : authToken
            },
            body: JSON.stringify(obj)
        }).then(res => {
            if (res.status === 401)
                return res.text().then(data => {
                    throw new AuthenticationException(data)
                })
            else if (res.status !== 200)
                return res.text().then(data => {
                    throw new HttpException(data, res.status)
                })

            return res.json()
        })
    }

    protected delete(authToken: AuthType, url: string) :Promise<any> {
        let addr = this.storage.toString()

        if (!url.startsWith("/"))
            url = `/${url}`

        return fetch(`${addr}/rest/v1${url}`, {
            method: 'DELETE',
            credentials: "same-origin",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
                'Authorization': isAuthenticationResponse(authToken) ? authToken.token : authToken
            },
        }).then(res => {
            if (res.status === 401)
                return res.text().then(data => {
                    throw new AuthenticationException(data)
                })
            else if (res.status !== 200)
                return res.text().then(data => {
                    throw new HttpException(data, res.status)
                })

            return res.json()
        })
    }

    protected get(authToken: AuthType, url: string) :Promise<any> {
        let addr = this.storage.toString()

        if (!url.startsWith("/"))
            url = `/${url}`

        return fetch(`${addr}/rest/v1${url}`, {
            method: 'GET',
            credentials: "same-origin",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
                'Authorization': isAuthenticationResponse(authToken) ? authToken.token : authToken
            },
        }).then(res => {
            if (res.status === 401)
                return res.text().then(data => {
                    throw new AuthenticationException(data)
                })
            else if (res.status !== 200)
                return res.text().then(data => {
                    throw new HttpException(data, res.status)
                })

            return res.json()
        })
    }



    /**
     *
     * @param fileName - name of the file to be uploaded
     * @param f - ArrayBuffer with actual file content
     * @param metaInfo - optional string/string dictionary to be stored together with file
     */
    uploadFile(fileName: string, f: ArrayBuffer, metaInfo: MetaType = undefined) :Promise<UploadResponse> {
        if (this.uploadToken === undefined || this.assignedBucket === undefined)
            throw new AuthenticationException("Please authenticate first")

        const req = bufferUploadRequest(this.assignedBucket, fileName, f, metaInfo)
        return this.post(this.uploadToken, "/file", req).then(data => {
            if (isUploadResponse(data))
                return data as UploadResponse

            throw new DatatypeException("UploadResponse", data)
        })
    }

    deleteFile(fileName: string) :Promise<ApiResponse> {
        if (this.uploadToken === undefined || this.assignedBucket === undefined)
            throw new AuthenticationException("Please authenticate first")

        let addr = this.storage.toString()

        if (!fileName.startsWith("/"))
            fileName = `/${fileName}`

        return fetch(`${addr}${fileName}`, {
            method: 'DELETE',
            headers: {
                'Authorization': this.uploadToken,
            },
        }).then(res => {
            if (res.status === 401)
                return res.text().then(data => {
                    throw new AuthenticationException(data)
                })
            else if (res.status !== 200)
                return res.text().then(data => {
                    throw new HttpException(data, res.status)
                })

            return res.json()
        }).then(data => {
            if (isApiResponse(data))
                return data as ApiResponse

            throw new DatatypeException("ApiResponse", data)
        })
    }

    updateMetaInfo(fileName: string, metaInfo: MetaType = undefined) :Promise<ApiResponse> {
        return undefined
    }

    getMetaInfo(fileName: string) :Promise<MetaType> {
        return undefined
    }

    deleteMetaInfo(fileName: string) :Promise<ApiResponse> {
        return undefined
    }

    listFiles() :Promise<ApiResponse> {
        return undefined
    }

    ping() :Promise<ApiResponse> {
        return this.get(this.uploadToken, "/ping").then(data => {
            if (isApiResponse(data))
                return data as ApiResponse

            throw new DatatypeException("ApiResponse", data)
        })
    }

    setCredentials(uploadToken: AuthType, bucket: string) {
        this.uploadToken = uploadToken
        this.assignedBucket = bucket
    }
}
