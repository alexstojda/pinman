import {Api, Token, TOKEN_KEY} from "./api";
import {faker} from "@faker-js/faker";
import {AuthApi, ErrorResponse, TokenResponse, UsersApi} from "./generated";
import axios from "axios";

jest.mock('./generated')

jest.mock('axios');
const mockedAxios = axios as jest.Mocked<typeof axios>;

beforeEach(() => {
  mockedAxios.get.mockReset()
})

describe('Api.tryTokenRefresh()', () => {
  it('returns true if token is in localStorage and not expired', async () => {
    const localStorageSpy = jest.spyOn(Storage.prototype, 'getItem')
    const token: Token = {
      token: faker.random.alphaNumeric(64),
      expires: new Date(Date.now() + 1000 * 60 * 10)
    }
    localStorageSpy.mockReturnValueOnce(JSON.stringify(token))

    const result = await new Api().tryTokenRefresh()
    expect(localStorageSpy).toBeCalledTimes(1)
    expect(result).toBe(true)
  })
  it('returns false if token is expired', async () => {
    const localStorageSpy = jest.spyOn(Storage.prototype, 'getItem')
    const token: Token = {
      token: faker.random.alphaNumeric(64),
      expires: new Date(Date.now() - 1000 * 60 * 10)
    }
    localStorageSpy.mockReturnValueOnce(JSON.stringify(token))

    const result = await new Api().tryTokenRefresh()
    expect(localStorageSpy).toBeCalledTimes(1)
    expect(result).toBe(false)
  })
  it('returns false if token is not in localStorage', async () => {
    const localStorageSpy = jest.spyOn(Storage.prototype, 'getItem')
    localStorageSpy.mockReturnValueOnce(null)

    const result = await new Api().tryTokenRefresh()
    expect(localStorageSpy).toBeCalledTimes(1)
    expect(result).toBe(false)
  })
  it('returns true if token is refreshed successfully', async () => {
    const localStorageGetSpy = jest.spyOn(Storage.prototype, 'getItem')
    const token: Token = {
      token: faker.random.alphaNumeric(64),
      expires: new Date(Date.now() + 1000 * 60 * 4)
    }
    localStorageGetSpy.mockReturnValueOnce(JSON.stringify(token))

    const localStorageSetSpy = jest.spyOn(Storage.prototype, 'setItem')

    const mockedAuthApi = jest.mocked(AuthApi)
    const refreshResult: TokenResponse = {
      access_token: faker.random.alphaNumeric(64),
      expire: new Date(Date.now() + 60 * 60 * 1000).toISOString()
    }
    mockedAuthApi.prototype.authRefreshGet.mockResolvedValue({
      config: {},
      data: refreshResult,
      headers: {},
      status: 200,
      statusText: "Ok",
    })
    const refreshedToken: Token = {
      token: refreshResult.access_token,
      expires: new Date(Date.parse(refreshResult.expire))
    }

    const result = await new Api().tryTokenRefresh()
    expect(localStorageGetSpy).toBeCalledTimes(1)
    expect(localStorageSetSpy).toBeCalledWith(TOKEN_KEY, JSON.stringify(refreshedToken))
    expect(result).toBe(true)
  })
  it('returns false if token refresh fails', async () => {
    const localStorageGetSpy = jest.spyOn(Storage.prototype, 'getItem')
    const token: Token = {
      token: faker.random.alphaNumeric(64),
      expires: new Date(Date.now() + 1000 * 60 * 4)
    }
    localStorageGetSpy.mockReturnValueOnce(JSON.stringify(token))

    const mockedAuthApi = jest.mocked(AuthApi)
    const refreshResult: ErrorResponse = {
      status: 401,
      detail: "token is expired",
      title: "Unauthorized"
    }
    mockedAuthApi.prototype.authRefreshGet.mockResolvedValue({
      config: {},
      data: refreshResult as any,
      headers: {},
      status: 401,
      statusText: "Unauthorized",
    })

    const consErrSpy = jest.spyOn(console, 'error').mockImplementation(jest.fn())

    const result = await new Api().tryTokenRefresh()
    expect(localStorageGetSpy).toBeCalledTimes(1)
    expect(consErrSpy).toBeCalledTimes(1)
    expect(result).toBe(false)
  })
})

describe('Api.login', () => {
  it('succeeds with valid credentials',async () => {
    const mockedAuthApi = jest.mocked(AuthApi)
    const loginResult: TokenResponse = {
      access_token: faker.random.alphaNumeric(64),
      expire: new Date(Date.now() + 60 * 60 * 1000).toISOString()
    }
    mockedAuthApi.prototype.authLoginPost.mockResolvedValue({
      config: {},
      data: loginResult,
      headers: {},
      status: 200,
      statusText: "Ok",
    })

    const setJwtSpy = jest.spyOn(Api.prototype, 'setJwtToken')

    const result = await (new Api()).login({
      username: faker.internet.email(),
      password: faker.random.alphaNumeric(16)
    })

    expect(setJwtSpy).toBeCalledWith({
      token: loginResult.access_token,
      expires: new Date(Date.parse(loginResult.expire))
    })
    expect(result).toBe(true)
  })
  it('fails with invalid credentials',async () => {
    const mockedAuthApi = jest.mocked(AuthApi)
    const loginResult: ErrorResponse = {
      status: 401,
      detail: "invalid username or password",
      title: "Unauthorized"
    }
    mockedAuthApi.prototype.authLoginPost.mockResolvedValue({
      config: {
        data: loginResult
      },
      data: loginResult as any,
      headers: {},
      status: 401,
      statusText: "Unauthorized",
    })

    const setJwtSpy = jest.spyOn(Api.prototype, 'setJwtToken')
    await expect(new Api().login({
      username: faker.internet.email(),
      password: faker.random.alphaNumeric(16)
    })).rejects.toThrowError()
    expect(setJwtSpy).toBeCalledTimes(0)
  })
})
describe('Api.clearJwtToken', () => {
  it('should call localStorage.removeItem', () => {
    const removeItemSpy = jest.spyOn(Storage.prototype, 'removeItem')
    new Api().clearJwtToken()
    expect(removeItemSpy).toBeCalledTimes(1)
  })
})
describe('Api.parseError', () => {
  it('should return ErrorResponse if valid data is provided', () => {
    const errorResponse: ErrorResponse = {
      detail: "page does not exist",
      status: 404,
      title: "Not Found"
    }

    const response = new Api().parseError({
      config: {},
      isAxiosError: true,
      message: "An error occured",
      name: "Error",
      toJSON: jest.fn(),
      response: {
        config: {},
        status: 404,
        headers: {},
        statusText: 'Not Found',
        data: errorResponse
      }
    })

    expect(response).toEqual(errorResponse)
  })
  it('should throw error if invalid data is provided ', async () => {
    expect(() => {
      new Api().parseError({
        config: {},
        isAxiosError: false,
        message: "An error occured",
        name: "Error",
        toJSON: jest.fn(),
      })
    }).toThrowError()
  })
})

describe('Api.userApi', () => {
  it('should return an instance of UsersApi', () => {
    const userApi = new Api().userApi()
    expect(userApi).toBeInstanceOf(UsersApi)
  })
})

describe('Api.authApi', () => {
  it('should return an instance of AuthApi', () => {
    const authApi = new Api().authApi()
    expect(authApi).toBeInstanceOf(AuthApi)
  })
})

describe('Api.accessToken', () => {
  it('returns the token if it was available in storage', () => {
    const localStorageSpy = jest.spyOn(Storage.prototype, 'getItem')
    const token: Token = {
      token: faker.random.alphaNumeric(64),
      expires: new Date(Date.now() + 1000 * 60 * 10)
    }
    localStorageSpy.mockReturnValueOnce(JSON.stringify(token))

    const mockUsersApi = jest.mocked(UsersApi)
    mockUsersApi.mockImplementation((c) => {
      expect(c).toBeDefined()

      const accessToken = typeof c!.accessToken === 'function'
        ? c!.accessToken()
        : undefined;

      expect(accessToken).toBe(token.token)

      return {} as UsersApi
    })

    new Api().userApi()
  })
  it('returns empty string if it was not available in storage', () => {
    const localStorageSpy = jest.spyOn(Storage.prototype, 'getItem')
    localStorageSpy.mockReturnValueOnce(null)

    const mockUsersApi = jest.mocked(UsersApi)
    mockUsersApi.mockImplementation((c) => {
      expect(c).toBeDefined()
      const accessToken = typeof c!.accessToken === 'function'
        ? c!.accessToken()
        : undefined;

      expect(accessToken).toBe("")

      return {} as UsersApi
    })

    new Api().userApi()
  })
})
