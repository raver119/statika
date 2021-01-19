

export const authenticationBean = (token: string, bucket: string) => {
    return {
        token: token,
        bucket: bucket
    }
}

export type AuthenticationBean = ReturnType<typeof authenticationBean>