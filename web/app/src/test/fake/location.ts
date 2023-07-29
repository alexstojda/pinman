import {Location} from "../../api";
import {faker} from "@faker-js/faker";
import * as helpers from "../../helpers";

export function location(): Location {
  const name = faker.hacker.phrase();

  return {
    id: faker.datatype.uuid(),
    name: name,
    slug: helpers.slugify(name),
    address: faker.address.streetAddress(),
    pinball_map_id: faker.datatype.number(),
    created_at: faker.date.recent(5).toISOString(),
    updated_at: faker.date.recent(1).toISOString()
  }
}
