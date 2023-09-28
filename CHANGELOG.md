# [0.8.0](https://github.com/alexstojda/pinman/compare/v0.7.0...v0.8.0) (2023-09-28)


### Bug Fixes

* **api:** Make nested objects optional, to avoid having to deep-load ownership trees ([ea2efe8](https://github.com/alexstojda/pinman/commit/ea2efe8cb697023690c21cb615993ce4fe30b0f9))


### Features

* **tournaments-frontend:** List tournaments page ([2ca7624](https://github.com/alexstojda/pinman/commit/2ca7624144efa18d8d498a537097758e46d05067))



# [0.7.0](https://github.com/alexstojda/pinman/compare/v0.6.0...v0.7.0) (2023-08-06)


### Bug Fixes

* **api:** fields should be snake case in API ([e279a3a](https://github.com/alexstojda/pinman/commit/e279a3a528066d5243d8f5375632d15dfc4dac2c))
* **api:** Fix/standardize db error handling ([1101e0f](https://github.com/alexstojda/pinman/commit/1101e0f00e74cf1f5e1c0c28984614471dc64ad2))


### Features

* **tournaments-api:** Add endpoint to create tournament ([fdce652](https://github.com/alexstojda/pinman/commit/fdce652b468dcbd9651b7d5ef823f09ddde06a55))
* **tournaments-api:** Add endpoint to list tournaments ([8d786c8](https://github.com/alexstojda/pinman/commit/8d786c8ea0c29b46c269b8ed00863a0487ac6e72))



# [0.6.0](https://github.com/alexstojda/pinman/compare/v0.5.0...v0.6.0) (2023-08-01)


### Bug Fixes

* **api:** Panic if controller is missing from API ([8921b01](https://github.com/alexstojda/pinman/commit/8921b01e40101862dde3835687c3c33f654812a0))
* **leagues-api:** Fix UUIDs in list response ([5b1fa79](https://github.com/alexstojda/pinman/commit/5b1fa79980f1ebf7cdeddb71b605a8b549a9b738))
* **locations-api:** Fix error checking logic in location api controller ([86aa067](https://github.com/alexstojda/pinman/commit/86aa067ebc058a7b90861a1ae8ffcc76650e659b))


### Features

* **leagues-api:** Include location object in league response instead of just location ID ([a378cec](https://github.com/alexstojda/pinman/commit/a378cec6268082226b4ea90c95f82bc5f7ebbc95))
* **leagues-frontend:** Show Location name & address on list page ([cac4801](https://github.com/alexstojda/pinman/commit/cac4801471ba34cdfe0e85ad0a03873c0c926995))
* **leagues-frontend:** Use locations API when creating a league ([1eedc21](https://github.com/alexstojda/pinman/commit/1eedc21401ed712777ff33918357f2f4d6e163a7))
* **leagues:** Update league API to use Locations ([aa68da1](https://github.com/alexstojda/pinman/commit/aa68da166bad5857ee62318b003d553611aa6b7a))



# [0.5.0](https://github.com/alexstojda/pinman/compare/v0.4.0...v0.5.0) (2023-07-27)


### Features

* **api-clients:** Create api client for pinball map API to fetch locations ([05fc2a6](https://github.com/alexstojda/pinman/commit/05fc2a63097853bde69343e2a181cdaf041dcd55))
* **locations:** Add API for Create/Get/List locations ([610c284](https://github.com/alexstojda/pinman/commit/610c284761426dac48672092c0651dd6bfa3e3a0))



# [0.4.0](https://github.com/alexstojda/pinman/compare/v0.3.0...v0.4.0) (2023-07-24)


### Features

* **frontend:** Add create league form ([1f655df](https://github.com/alexstojda/pinman/commit/1f655df387870b43d8311f415e318b766c36c9f8))
* **frontend:** List leagues on homepage ([87a81cf](https://github.com/alexstojda/pinman/commit/87a81cfc1e72f790826849e064351a22017b22f6))



