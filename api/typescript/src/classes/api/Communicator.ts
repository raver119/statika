import {AuthenticationResponse} from "../entities/AuthenticationResponse";
import {AuthenticationException} from "../exceptions/AuthenticationException";
import {HttpException} from "../exceptions/HttpException";
import {AuthenticationBean} from "./AuthenticationBean";
import {EndpointsCoordinates} from "../system/EndpointsCoordinates";

export type AuthType = AuthenticationBean | AuthenticationResponse

export const communicator = (storage: EndpointsCoordinates) => {

    return {
        storage() :EndpointsCoordinates {
            return storage
        },

        post(authToken: AuthType, url: string, obj: any) :Promise<any> {
            const addr = storage.toString()

            if (!url.startsWith("/"))
                url = `/${url}`

            return fetch(`${addr}/rest/v1${url}`, {
                    method: 'POST',
                    credentials: "same-origin",
                    headers: {
                        'Accept': 'application/json',
                        'Content-Type': 'application/json',
                        'Authorization': authToken.token
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
        },

        delete(authToken: AuthType, url: string) :Promise<any> {
            const addr = storage.toString()

            // delete is special case, and might be called for an absolute address
            let input: string
            if (url.toLowerCase().startsWith("http://") || url.toLowerCase().startsWith("https://")) {
                input = url
            } else {
                if (!url.startsWith("/"))
                    url = `/${url}`

                input = `${addr}/rest/v1${url}`
            }

            return fetch(input, {
                method: 'DELETE',
                credentials: "same-origin",
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                    'Authorization': authToken.token
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
        },

        get(authToken: AuthType, url: string) :Promise<any> {
            const addr = storage.toString()

            if (!url.startsWith("/"))
                url = `/${url}`

            return fetch(`${addr}/rest/v1${url}`, {
                method: 'GET',
                credentials: "same-origin",
                headers: {
                    'Accept': 'application/json',
                    'Content-Type': 'application/json',
                    'Authorization': authToken.token
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
    }
}

export type Communicator = ReturnType<typeof communicator>