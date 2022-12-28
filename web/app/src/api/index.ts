import {AuthApi, Configuration, UserLogin, UsersApi, ErrorResponse} from "./generated";
import axios, {AxiosError, AxiosInstance} from "axios";

const ACCESS_TOKEN_KEY = 'access_token';
const REFRESH_DEADLINE = 5 * 1000 * 60 // 5 minutes

type Token = {
  token: string
  expires: Date
}

export class Api {
  private axiosInstance: AxiosInstance;

  constructor() {
    this.axiosInstance = axios.create()
  }

  // Retrieve the Token from local storage
  private getJwtToken(): Token | undefined {
    const tokenStr = localStorage.getItem(ACCESS_TOKEN_KEY)
    if (tokenStr) {
      return JSON.parse(tokenStr, (k, v) => {
        if (k === "expires") {
          return new Date(Date.parse(v))
        }
        return v;
      })
    }
    return undefined
  }

  // Save the Token to local storage
  public setJwtToken(token: Token) {
    localStorage.setItem(ACCESS_TOKEN_KEY, JSON.stringify(token))
  }

  public clearJwtToken() {
    localStorage.removeItem(ACCESS_TOKEN_KEY)
  }

  // Tries to refresh the token. Returns true if successful or token is still valid, false otherwise.
  public async tryTokenRefresh(): Promise<boolean> {
    const token = this.getJwtToken()
    if (token) {

      if (token.expires.getTime() < new Date(Date.now()).getTime())
        return false

      // If the token will expire in longer than REFRESH_TOKEN, we should refresh
      if (token.expires.getTime() - new Date(Date.now()).getTime() > REFRESH_DEADLINE) {
        return true
      }

      return await this.authApi().authRefreshGet().then(r => {
        if (r.status === 200) {
          this.setJwtToken({
            token: r.data.access_token,
            //TODO: UTC Date?
            expires: new Date(Date.parse(r.data.expire))
          })
          return true
        } else {
          console.error("refresh failed", r.config.data)
          return false
        }
      })
    }

    return false
  }

  public login(credentials: UserLogin): Promise<boolean> {
    return this.authApi().authLoginPost(credentials).then<boolean>(r => {
      if (r.status === 200) {
        this.setJwtToken({
          token: r.data.access_token,
          //TODO: UTC Date?
          expires: new Date(Date.parse(r.data.expire))
        })
        return true
      }
      console.error("Login failed", r.config.data)
      throw new Error((r.config.data as ErrorResponse).detail)
    })
  }

  // Method used by the generated API client to get the access_token.
  private accessToken(): string {
    const token = this.getJwtToken()
    if (token && token.expires > new Date()) {
      return token.token
    }
    return "";
  };

  private configuration(): Configuration {
    const openapiConfig = new Configuration();
    openapiConfig.accessToken = this.accessToken;
    return openapiConfig;
  };

  public userApi(): UsersApi {
    return new UsersApi(this.configuration());
  };

  public authApi(): AuthApi {
    return new AuthApi(this.configuration());
  }

  public parseError = (e: AxiosError): ErrorResponse => {
    if (e.isAxiosError && e.response)
      return e.response.data as ErrorResponse
    throw new Error("parseError: could not parse error")
}
}

export * from './generated'