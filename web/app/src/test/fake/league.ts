import {League} from "../../api";
import {faker} from "@faker-js/faker";
import * as helpers from "../../helpers";

export function league(): League {
  const name = faker.hacker.phrase();

  return {
    id: faker.datatype.uuid(),
    name: name,
    slug: helpers.slugify(name),
    locationId: faker.datatype.uuid(),
    ownerId: faker.datatype.uuid(),
    created_at: faker.date.recent(5).toISOString(),
    updated_at: faker.date.recent(1).toISOString()
  }
}
