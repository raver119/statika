import {isUploadResponse, UploadResponse} from "./classes/responses/UploadResponse";
import {AuthenticationException} from "./classes/exceptions/AuthenticationException";
import {HttpException} from "./classes/exceptions/HttpException";
import {isAuthenticationResponse} from "./classes/responses/AuthenticationResponse";
import {DatatypeException} from "./classes/exceptions/DatatypeException";
import {EndpointsCoordinates} from "./classes/system/EndpointsCoordinates";
import {ApiResponse} from "./classes/responses/ApiResponse";


export type MetaType = Map<string, string>|undefined
type AuthType = string

class AsynchronousApi {
    protected endpoints: EndpointsCoordinates

    protected uploadToken: string
    protected assignedBucket: string

    protected constructor(token: string, bucket: string) {
        this.uploadToken = token
        this.assignedBucket = bucket
    }

    protected post(authToken: AuthType, url: string, obj: any) :Promise<any> {
        let addr = this.endpoints.toString()

        return fetch(`${addr}/rest/v1${url}`, {
            method: 'POST',
            credentials: "same-origin",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
                'Authorization': isAuthenticationResponse(authToken) ? authToken.authToken : authToken
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
        let addr = this.endpoints.toString()

        return fetch(`${addr}/rest/v1${url}`, {
            method: 'DELETE',
            credentials: "same-origin",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
                'Authorization': isAuthenticationResponse(authToken) ? authToken.authToken : authToken
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
        let addr = this.endpoints.toString()

        return fetch(`${addr}/rest/v1${url}`, {
            method: 'GET',
            credentials: "same-origin",
            headers: {
                'Accept': 'application/json',
                'Content-Type': 'application/json',
                'Authorization': isAuthenticationResponse(authToken) ? authToken.authToken : authToken
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
     * @param f
     * @param metaInfo - optional string/string dictionary to be stored together with file
     */
    uploadFile(f: ArrayBuffer, metaInfo: MetaType = undefined) :Promise<UploadResponse> {
        return this.post("", "/file", f).then(data => {
            if (isUploadResponse(data))
                return data as UploadResponse

            throw new DatatypeException("UploadResponse", data)
        })
    }

    deleteFile(fileName: string) :Promise<ApiResponse> {
        return undefined
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
        return undefined
    }
}


export class Statika extends AsynchronousApi {

    /**
     *
     * @param token - Upload token, usually generated in backend code on the fly, and fused into frontend app
     * @param bucket - Optional folder for splitting end users
     */
    constructor(token: string, bucket: string = "") {
        super(token, bucket)
    }

}