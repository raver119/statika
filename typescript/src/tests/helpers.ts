import "whatwg-fetch";
import { AuthenticationException, HttpException, isAuthenticationResponse } from "../index";
import { DatatypeException } from "../classes/exceptions/DatatypeException";

export const authorizeUpload = async (uploadKey: string, bucket: string): Promise<string> => {
  const host = process.env.API_NODE ?? "127.0.0.1";
  const port = process.env.API_PORT ?? "9191";

  return fetch(`http://${host}:${port}/rest/v1/auth/upload`, {
    method: "POST",
    headers: {
      Accept: "application/json",
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ token: uploadKey, bucket: bucket }),
  })
    .then((res) => {
      if (res.status === 401)
        return res.text().then((data) => {
          throw new AuthenticationException(data);
        });
      else if (res.status !== 200)
        return res.text().then((data) => {
          throw new HttpException(data, res.status);
        });

      return res.json();
    })
    .then((data) => {
      if (isAuthenticationResponse(data)) return data.token;

      throw new DatatypeException("AuthenticationResponse", data);
    });
};

export const httpGet = async (url: string): Promise<string> => {
  return fetch(url, {
    method: "GET",
  }).then((res) => {
    if (res.status === 401)
      return res.text().then((data) => {
        throw new AuthenticationException(data);
      });
    else if (res.status !== 200)
      return res.text().then((data) => {
        throw new HttpException(data, res.status);
      });

    return res.text();
  });
};
