import {Configuration, UserApi} from "./generated";
import {Cookies} from "react-cookie";
import axios, {Axios, AxiosInstance} from "axios";

const ACCESS_TOKEN_COOKIE = 'access_token';
const REFRESH_TOKEN_COOKIE = 'refresh_token';

class MyApi {
  private cookies: Cookies;
  private axiosInstance: AxiosInstance;

  constructor() {
    this.axiosInstance = axios.create()
    this.cookies = new Cookies();
  }

  private accessToken = (name?: string, scopes?: string[]) => {
    return this.cookies.get(ACCESS_TOKEN_COOKIE);
  };

  private configuration = () => {
    const openapiConfig = new Configuration();
    openapiConfig.accessToken = this.accessToken;
    return openapiConfig;
  };

  public userApi = () => {
    return new UserApi(this.configuration());
  };
}