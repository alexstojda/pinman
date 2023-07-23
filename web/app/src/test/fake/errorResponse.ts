import {ErrorResponse} from "../../api";
import {faker} from "@faker-js/faker";

export function errorResponse(): ErrorResponse {
  return {
    status: faker.internet.httpStatusCode({
      types: ['serverError', 'clientError']
    }),
    detail: faker.hacker.phrase(),
    title: faker.random.words(3),
  }
}
