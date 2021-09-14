export const authenticationBean = (token: string, ...buckets: string[]) => ({
  token: token,
  bucket: buckets[0],
  buckets: buckets,
});

export type AuthenticationBean = ReturnType<typeof authenticationBean>;

export function isAuthenticationBean(obj: any): obj is AuthenticationBean {
  return obj !== undefined && typeof obj.token === "string" && obj.buckets !== undefined && Array.isArray(obj.buckets);
}
