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



# [0.3.0](https://github.com/alexstojda/pinman/compare/v0.2.6...v0.3.0) (2023-07-24)


### Features

* **backend:** Add create league endpoint ([efbbe0a](https://github.com/alexstojda/pinman/commit/efbbe0a8961d6626c56cf6bb6d791a3e91a5dacc))
* **backend:** Add get league endpoint ([31632f0](https://github.com/alexstojda/pinman/commit/31632f08c559d51961ac0343efcf011729e70062))
* **backend:** Add list leagues endpoint ([3b7302c](https://github.com/alexstojda/pinman/commit/3b7302ccc29732f234e574632bab98aa844f6e08))



## [0.2.6](https://github.com/alexstojda/pinman/compare/v0.2.5...v0.2.6) (2023-01-06)


### Bug Fixes

* Add code to correct viper bug ([f3565ce](https://github.com/alexstojda/pinman/commit/f3565ce55e2c41d04f880e2ba3c6c75d028faced))
* **backend:** Fix serving static app ([4c61fb9](https://github.com/alexstojda/pinman/commit/4c61fb9a4424224941fb0c24d222564daaf278a7))
* **backend:** Run migrations on every server start ([19c715c](https://github.com/alexstojda/pinman/commit/19c715c04287d8c98a0df79bd599ae0b5cdbcd5f))
* **frontend:** Set base path to match static serve route ([3c793e6](https://github.com/alexstojda/pinman/commit/3c793e65a3aeb35752952c854da857d09b46e292))
* Use RAILWAY_STATIC_URL during build ([320e84e](https://github.com/alexstojda/pinman/commit/320e84eaf77868b6cce4fbf93720403fd9628e83))



