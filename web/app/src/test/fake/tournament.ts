import {faker} from "@faker-js/faker";
import {Tournament} from "../../api";
import * as helpers from "../../helpers";
import {location} from "./location";
import {league} from "./league";

export function tournament(): Tournament {
  const name = faker.hacker.phrase();

  return {
    name: name,
    slug: helpers.slugify(name),
    id: faker.datatype.uuid(),
    type: "multi_round_tournament",
    settings: {
      rounds: 8,
      games_per_round: 4,
      lowest_scores_dropped: 3
    },
    location: location(),
    league: league(),
    created_at: faker.date.recent(5).toISOString(),
    updated_at: faker.date.recent(1).toISOString()
  }
}
