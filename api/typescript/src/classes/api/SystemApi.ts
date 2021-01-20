import {Communicator} from "./Communicator";
import {AuthenticationBean} from "./AuthenticationBean";
import {ApiResponse, isApiResponse} from "../entities/ApiResponse";
import {DatatypeException} from "../exceptions/DatatypeException";

export const systemApi = (communicator: Communicator) => {
    return {
        ping(bean: AuthenticationBean) :Promise<ApiResponse> {
            return communicator.get(bean, "/ping").then(data => {
                if (isApiResponse(data))
                    return data as ApiResponse

                throw new DatatypeException("ApiResponse", data)
            })
        },
    }
}

export type SystemApi = ReturnType<typeof systemApi>