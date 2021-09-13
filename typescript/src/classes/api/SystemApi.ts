import {Communicator} from "./Communicator";
import {authenticationBean, AuthenticationBean} from "./AuthenticationBean";
import {ApiResponse, isApiResponse} from "../entities/ApiResponse";
import {DatatypeException} from "../exceptions/DatatypeException";
import {AuthenticationResponse} from "../entities/AuthenticationResponse";

export const systemApi = (communicator: Communicator) => {
    return {
        async ping(bean: AuthenticationBean) :Promise<ApiResponse> {
            return communicator.get(bean, "/ping").then(data => {
                if (isApiResponse(data))
                    return data as ApiResponse

                throw new DatatypeException("ApiResponse", data)
            })
        },

        async issueToken(upload_key: string, ...buckets: string[]): Promise<AuthenticationBean> {
            return communicator.post(undefined, `/auth/upload`, {token: upload_key, buckets: buckets})
                .then((response: AuthenticationResponse) => authenticationBean(response.token, ...buckets))
        }
    }
}

export type SystemApi = ReturnType<typeof systemApi>