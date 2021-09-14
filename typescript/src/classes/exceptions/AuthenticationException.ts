import { HttpException } from "./HttpException";

export class AuthenticationException extends HttpException {
  constructor(message: string) {
    super(message, 401);
  }
}
