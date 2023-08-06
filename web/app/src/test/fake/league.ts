import {League} from "../../api";
import {faker} from "@faker-js/faker";
import * as helpers from "../../helpers";
import {location} from "./location";

export function league(): League {
  const name = faker.hacker.phrase();

  return {
    id: faker.datatype.uuid(),
    name: name,
    slug: helpers.slugify(name),
    location: location(),
    owner_id: faker.datatype.uuid(),
    created_at: faker.date.recent(5).toISOString(),
    updated_at: faker.date.recent(1).toISOString()
  }
}
