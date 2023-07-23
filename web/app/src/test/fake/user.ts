import {User} from "../../api";
import {faker} from "@faker-js/faker";

export function user(): User {
  return {
    name: faker.name.fullName(),
    email: faker.internet.email(),
    id: faker.datatype.uuid(),
    role: "user",
    created_at: faker.date.recent(5).toISOString(),
    updated_at: faker.date.recent(1).toISOString()
  }
}
