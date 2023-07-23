import {AuthApi, Configuration, ErrorResponse, LeaguesApi, UserLogin, UsersApi} from "./generated";
import {AxiosError} from "axios";

export const TOKEN_KEY = 'token';
export const REFRESH_DEADLINE = 5 * 1000 * 60 // 5 minutes

export type Token = {
  token: string
  expires: Date
}

export class Api {
  // Retrieve the Token from local storage
  private getJwtToken(): Token | undefined {
    const tokenStr = localStorage.getItem(TOKEN_KEY)
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
    localStorage.setItem(TOKEN_KEY, JSON.stringify(token))
  }

  public clearJwtToken() {
    localStorage.removeItem(TOKEN_KEY)
  }

  // Tries to refresh the token. Returns true if successful or token is still valid, false otherwise.
  public async tryTokenRefresh(): Promise<boolean> {
    const token = this.getJwtToken()
    if (token) {

      // If the token is already expired, return false
      if (token.expires.getTime() < new Date(Date.now()).getTime())
        return false

      // If the token won't expire within thw REFRESH_DEADLINE, token is still valid
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
    const token = Api.prototype.getJwtToken()
    if (token) {
      return token.token
    }

    return ""
  };

  private configuration(): Configuration {
    const openapiConfig = new Configuration();
    openapiConfig.basePath = process.env.REACT_APP_API_HOST
    openapiConfig.accessToken = this.accessToken;
    return openapiConfig;
  };

  public userApi(): UsersApi {
    return new UsersApi(this.configuration());
  };

  public authApi(): AuthApi {
    return new AuthApi(this.configuration());
  }

  public leaguesApi(): LeaguesApi {
    return new LeaguesApi(this.configuration());
  }

  public parseError(e: AxiosError): ErrorResponse {
    if (e.isAxiosError && e.response && e.response.data)
      return e.response.data as ErrorResponse
    throw new Error("parseError: could not parse error")
  }
}
