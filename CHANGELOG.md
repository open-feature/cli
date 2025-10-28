# Changelog

## [0.3.8](https://github.com/open-feature/cli/compare/v0.3.7...v0.3.8) (2025-10-28)


### 🐛 Bug Fixes

* fix invalid gorelease configuration preventing new releases ([#175](https://github.com/open-feature/cli/issues/175)) ([e20442d](https://github.com/open-feature/cli/commit/e20442de78510d92980fb46d4fd28779e80c3b70))

## [0.3.7](https://github.com/open-feature/cli/compare/v0.3.6...v0.3.7) (2025-10-28)


### ✨ New Features

* add interactive prompting to manifest add command ([#174](https://github.com/open-feature/cli/issues/174)) ([9d8b2ce](https://github.com/open-feature/cli/commit/9d8b2cea4f930b0091e064fd176366a63d65e3aa))
* adds ability to access original flag keys post-generation ([#167](https://github.com/open-feature/cli/issues/167)) ([fe326f6](https://github.com/open-feature/cli/commit/fe326f6b8838f897ba3309fe09e6284758d2d8b9))
* manifest command group for adding and listing flags in manifest ([#168](https://github.com/open-feature/cli/issues/168)) ([3d3f7a6](https://github.com/open-feature/cli/commit/3d3f7a667d4dc05a0ac9d8cf9108e19cda240465))


### 📚 Documentation

* fix typo in the readme ([96d37f1](https://github.com/open-feature/cli/commit/96d37f1219deb36c4265566be6c108344da9d9eb))
* update README.md formatting for openfeature.dev website docs ([#172](https://github.com/open-feature/cli/issues/172)) ([cd4b29f](https://github.com/open-feature/cli/commit/cd4b29f82870c8cd3e81288affbd1a6f347f9849))

## [0.3.6](https://github.com/open-feature/cli/compare/v0.3.5...v0.3.6) (2025-08-29)


### 🐛 Bug Fixes

* **deps:** update module dagger.io/dagger to v0.18.10 ([#136](https://github.com/open-feature/cli/issues/136)) ([8b70612](https://github.com/open-feature/cli/commit/8b706124721dfd2a904d102235baeb445e67cce0))
* **deps:** update module dagger.io/dagger to v0.18.11 ([#142](https://github.com/open-feature/cli/issues/142)) ([2835e3c](https://github.com/open-feature/cli/commit/2835e3cf1066e8446c472bd65b01cec20c864457))
* **deps:** update module dagger.io/dagger to v0.18.12 ([#143](https://github.com/open-feature/cli/issues/143)) ([cf962d9](https://github.com/open-feature/cli/commit/cf962d967c8ada93cb0b4959d4b0bd7ae9508e08))
* **deps:** update module github.com/pterm/pterm to v0.12.81 ([#129](https://github.com/open-feature/cli/issues/129)) ([a25f90a](https://github.com/open-feature/cli/commit/a25f90a65b50d4218b905351f9ed0504d1e54fba))
* **security:** update module github.com/go-viper/mapstructure/v2 to v2.3.0 [security] ([#149](https://github.com/open-feature/cli/issues/149)) ([616b446](https://github.com/open-feature/cli/commit/616b446ca18a816c5fea89811555c30188734c11))
* **security:** update module github.com/go-viper/mapstructure/v2 to v2.4.0 [security] ([#151](https://github.com/open-feature/cli/issues/151)) ([9d635ac](https://github.com/open-feature/cli/commit/9d635ac4520b0100970b7f6f64f2d1b5b0532bc4))


### ✨ New Features

* `openfeature pull` command ([#147](https://github.com/open-feature/cli/issues/147)) ([c517e87](https://github.com/open-feature/cli/commit/c517e8722e749e296687cc9917b8e02cc7a60f8a))
* basic object flags ([#141](https://github.com/open-feature/cli/issues/141)) ([288023c](https://github.com/open-feature/cli/commit/288023c5ddd03095e6d545bf4062374758b33c82))
* **flagset:** improve validation error formatting in Load function ([#119](https://github.com/open-feature/cli/issues/119)) ([8eec779](https://github.com/open-feature/cli/commit/8eec77965ab1b14121f7492a1b08bdaadd765bd9))
* **flagset:** improve validation error formatting in Load function [#110](https://github.com/open-feature/cli/issues/110) ([8eec779](https://github.com/open-feature/cli/commit/8eec77965ab1b14121f7492a1b08bdaadd765bd9))
* integration tests for nodejs generator ([#140](https://github.com/open-feature/cli/issues/140)) ([b867485](https://github.com/open-feature/cli/commit/b8674851013d9bd6d374c78ff781e0b210aa6686))


### 🧹 Chore

* add merge group trigger to pr lint ([32f912a](https://github.com/open-feature/cli/commit/32f912a08973a1e3a5127739691ddd314f094b96))
* add renovate.json file [#122](https://github.com/open-feature/cli/issues/122) ([#124](https://github.com/open-feature/cli/issues/124)) ([83dac70](https://github.com/open-feature/cli/commit/83dac705b9bfe00725e23c07cd7154d4e0877f22))
* **deps:** update alpine docker tag to v3.22 ([#130](https://github.com/open-feature/cli/issues/130)) ([11b9c34](https://github.com/open-feature/cli/commit/11b9c34d86ab72e25387f9efb8f5453dc19f9dd1))
* **deps:** update dagger/dagger-for-github action to v7 ([#132](https://github.com/open-feature/cli/issues/132)) ([1b4d062](https://github.com/open-feature/cli/commit/1b4d0620450b94a628445a9f753b543b1159d8a6))
* **deps:** update dependency microsoft.extensions.dependencyinjection to 8.0.1 ([#127](https://github.com/open-feature/cli/issues/127)) ([3ec8443](https://github.com/open-feature/cli/commit/3ec8443823455b546245f999904923a197833e4f))
* **deps:** update dependency microsoft.extensions.dependencyinjection to 9.0.6 ([#138](https://github.com/open-feature/cli/issues/138)) ([5d7a607](https://github.com/open-feature/cli/commit/5d7a60754a612054d65ef9094fa01f40156b39b2))
* **deps:** update dependency microsoft.extensions.dependencyinjection to v9 ([#133](https://github.com/open-feature/cli/issues/133)) ([51232fe](https://github.com/open-feature/cli/commit/51232fea1cd7a9ad22cd6533d227f55778de37eb))
* **deps:** update dependency openfeature to 2.6.0 ([#131](https://github.com/open-feature/cli/issues/131)) ([35533dc](https://github.com/open-feature/cli/commit/35533dceb01f52e59f258ef5583455a938bf658f))
* **deps:** update dependency openfeature to 2.7.0 ([#144](https://github.com/open-feature/cli/issues/144)) ([2003666](https://github.com/open-feature/cli/commit/2003666eaa0637ab57337d2d657701558dc1e5f6))
* **deps:** update golangci/golangci-lint-action action to v8 ([#134](https://github.com/open-feature/cli/issues/134)) ([4f909c1](https://github.com/open-feature/cli/commit/4f909c10686183057d01c304e1495a2987c1edf5))
* fix the directory structure ([#121](https://github.com/open-feature/cli/issues/121)) ([bcd11ea](https://github.com/open-feature/cli/commit/bcd11ea9c8115c1c7fd925f13613a332359de700))
* revert golang ci lint to v6 ([594cf53](https://github.com/open-feature/cli/commit/594cf538be3eab815ab40473c20dbf551adf87f5))
* upgrade dependencies ([#123](https://github.com/open-feature/cli/issues/123)) ([79d3dce](https://github.com/open-feature/cli/commit/79d3dceb3ad306b6c04c9e3c64285b5ffec3b05a))

## [0.3.5](https://github.com/open-feature/cli/compare/v0.3.4...v0.3.5) (2025-05-20)


### 🐛 Bug Fixes

* Naming of generated java class ([#111](https://github.com/open-feature/cli/issues/111)) ([49e65c8](https://github.com/open-feature/cli/commit/49e65c828330abb732eb3b9cf85850bb5ac36531))

## [0.3.4](https://github.com/open-feature/cli/compare/v0.3.3...v0.3.4) (2025-05-14)


### ✨ New Features

* add java generator ([#107](https://github.com/open-feature/cli/issues/107)) ([9a9f11f](https://github.com/open-feature/cli/commit/9a9f11fc6c6a8ffa38870e62ac26d9f8f679825b))
* adds compare command ([#93](https://github.com/open-feature/cli/issues/93)) ([063cfca](https://github.com/open-feature/cli/commit/063cfca2d79c9f75e181422ec375e300e020e57f))
* introduce dagger for integration testing and ci ([#100](https://github.com/open-feature/cli/issues/100)) ([96f4cde](https://github.com/open-feature/cli/commit/96f4cde0f87b8daf70e02c1d4ca3bcec018fee02))

## [0.3.3](https://github.com/open-feature/cli/compare/v0.3.2...v0.3.3) (2025-04-18)


### 🐛 Bug Fixes

* use the correct json schema url in init command ([#96](https://github.com/open-feature/cli/issues/96)) ([412a117](https://github.com/open-feature/cli/commit/412a1174b5dfe9ba77e18ec57d5a761711067386))


### ✨ New Features

* add codegen for NestJS ([#99](https://github.com/open-feature/cli/issues/99)) ([5210429](https://github.com/open-feature/cli/commit/5210429e39c10c91482cb0a0a8b2f4431a0aa182))
* **csharp:** added generator and integration tests ([#97](https://github.com/open-feature/cli/issues/97)) ([ae64581](https://github.com/open-feature/cli/commit/ae645813c48b5ef10d8557406e7ab5c96ce3df69))
* Python generator ([#95](https://github.com/open-feature/cli/issues/95)) ([1f8f43a](https://github.com/open-feature/cli/commit/1f8f43ae049fcf7c4feba3edaa697329688f7343))


### 🧹 Chore

* automate project standards before push ([#94](https://github.com/open-feature/cli/issues/94)) ([e32547f](https://github.com/open-feature/cli/commit/e32547f73495a525ed4ef5e2cadd45642d6fb172))

## [0.3.2](https://github.com/open-feature/cli/compare/v0.3.1...v0.3.2) (2025-04-02)


### ✨ New Features

* add contributing guide and generator readme ([#80](https://github.com/open-feature/cli/issues/80)) ([05e094d](https://github.com/open-feature/cli/commit/05e094db68c210271205f6a043fc885d1a3c23b8)), closes [#69](https://github.com/open-feature/cli/issues/69)
* add nodejs generator ([#91](https://github.com/open-feature/cli/issues/91)) ([a40b6a4](https://github.com/open-feature/cli/commit/a40b6a4d31d6f290ccdd9475bedbbe947aad510e))
* add script to install the latest binary ([#85](https://github.com/open-feature/cli/issues/85)) ([afa46d0](https://github.com/open-feature/cli/commit/afa46d00b303de8bf34197369fe34fd6022c34b9))
* **cli:** add stability annotations to generated Markdown documentation ([#88](https://github.com/open-feature/cli/issues/88)) ([9102d13](https://github.com/open-feature/cli/commit/9102d1390ace7e3b285ae4ce38208b229de59cbf))
* consolidate logging and support debug flag ([#92](https://github.com/open-feature/cli/issues/92)) ([3fbd947](https://github.com/open-feature/cli/commit/3fbd94726d581f1911ad8e539b004dd843503ef4))


### 📚 Documentation

* add install, quick start, commands, and more to readme ([#90](https://github.com/open-feature/cli/issues/90)) ([9244276](https://github.com/open-feature/cli/commit/9244276fc47128a7a304ef22732ad5dcde38c3e8))

## [0.3.1](https://github.com/open-feature/cli/compare/v0.3.0...v0.3.1) (2025-03-18)


### 🐛 Bug Fixes

* binary name referenced in the dockerfile ([0e28e8e](https://github.com/open-feature/cli/commit/0e28e8ec3b4108eee6ae43f587201ff7cbf18020))

## [0.3.0](https://github.com/open-feature/cli/compare/v0.2.0...v0.3.0) (2025-03-18)


### ⚠ BREAKING CHANGES

* change binary name ([#82](https://github.com/open-feature/cli/issues/82))

### 🧹 Chore

* change binary name ([#82](https://github.com/open-feature/cli/issues/82)) ([fdfe561](https://github.com/open-feature/cli/commit/fdfe561d49e17017af5165dfe0eec359387935e4))

## [0.2.0](https://github.com/open-feature/cli/compare/v0.1.10...v0.2.0) (2025-03-18)


### ⚠ BREAKING CHANGES

* add init command, update cli flags, support a config file ([#71](https://github.com/open-feature/cli/issues/71))

### 🧹 Chore

* rename the checksum file ([34afca6](https://github.com/open-feature/cli/commit/34afca62ab6cf229f38b0cc81d6f6443cf1ac8ea))
* upgrade viper to 1.20 ([#78](https://github.com/open-feature/cli/issues/78)) ([6c36ee9](https://github.com/open-feature/cli/commit/6c36ee90f796cdabe318ef59aec9de3d93c3ffd5))


### 🔄 Refactoring

* add init command, update cli flags, support a config file ([#71](https://github.com/open-feature/cli/issues/71)) ([106bf9d](https://github.com/open-feature/cli/commit/106bf9ddfe93673d956487bcf84667d550543aa0))

## [0.1.10](https://github.com/open-feature/cli/compare/v0.1.9...v0.1.10) (2025-01-27)


### ✨ New Features

* add doc gen, move schema path, add tests, fix react gen ([#68](https://github.com/open-feature/cli/issues/68)) ([68a72ee](https://github.com/open-feature/cli/commit/68a72ee929b134fb787396019102ade3fae3f697))
* update golang output ([#63](https://github.com/open-feature/cli/issues/63)) ([0e7db02](https://github.com/open-feature/cli/commit/0e7db0209e13b672329fc2f4578cdb700db7b826))


### 🧹 Chore

* go mod tidy, gitignore dist folder ([1530d38](https://github.com/open-feature/cli/commit/1530d38dd3b127d80457512e8a0da87f4f38f293))
* switch base image from distroless to alpine ([#67](https://github.com/open-feature/cli/issues/67)) ([60955af](https://github.com/open-feature/cli/commit/60955af1a9fe89b62f8508ecd97284b899b50786))

## [0.1.9](https://github.com/open-feature/cli/compare/v0.1.8...v0.1.9) (2024-11-27)


### 🧹 Chore

* remove empty testutils package ([#55](https://github.com/open-feature/cli/issues/55)) ([9dc1d9f](https://github.com/open-feature/cli/commit/9dc1d9fbc3751b53956e4c61cd43df63edca9f19))
* update back to previous mkdir permissions ([#61](https://github.com/open-feature/cli/issues/61)) ([515b534](https://github.com/open-feature/cli/commit/515b5340b5d61879bf2fdb786ea38cbbe0a24247))

## [0.1.8](https://github.com/open-feature/cli/compare/v0.1.7...v0.1.8) (2024-10-31)


### 🐛 Bug Fixes

* docker publishing ([c663816](https://github.com/open-feature/cli/commit/c663816e33d0a020c1bd4db110ac0e4f451ff7b1))

## [0.1.7](https://github.com/open-feature/cli/compare/v0.1.6...v0.1.7) (2024-10-31)


### 🐛 Bug Fixes

* docker publishing ([2d24d51](https://github.com/open-feature/cli/commit/2d24d5141c0822edb7254f38efdabaa6e9b5b351))

## [0.1.6](https://github.com/open-feature/cli/compare/v0.1.5...v0.1.6) (2024-10-31)


### 📚 Documentation

* switch from code gen to cli ([#47](https://github.com/open-feature/cli/issues/47)) ([7a1f9f3](https://github.com/open-feature/cli/commit/7a1f9f304cc9c512b407b19986fbd82e3b80fe53))


### 🔄 Refactoring

* change name of go module ([#46](https://github.com/open-feature/cli/issues/46)) ([e3058db](https://github.com/open-feature/cli/commit/e3058db6d7f4feef4780df6a5f1772e05b82571a))

## [0.1.5](https://github.com/open-feature/codegen/compare/v0.1.4...v0.1.5) (2024-10-22)


### 🐛 Bug Fixes

* container copy command ([#40](https://github.com/open-feature/codegen/issues/40)) ([8448543](https://github.com/open-feature/codegen/commit/8448543fda56a3d68851cf44a4735c1902bf5b98))

## [0.1.4](https://github.com/open-feature/codegen/compare/v0.1.3...v0.1.4) (2024-10-22)


### ✨ New Features

* add version command ([#38](https://github.com/open-feature/codegen/issues/38)) ([c13a448](https://github.com/open-feature/codegen/commit/c13a4486b9b42f3e4a6f34abd43a87aecf91355e))

## [0.1.3](https://github.com/open-feature/codegen/compare/v0.1.2...v0.1.3) (2024-10-22)


### 🐛 Bug Fixes

* set github token for release process ([a2fe4aa](https://github.com/open-feature/codegen/commit/a2fe4aa33e380e86925480e7233eeed4bfb9ed90))

## [0.1.2](https://github.com/open-feature/codegen/compare/v0.1.1...v0.1.2) (2024-10-22)


### ✨ New Features

* add basic react support ([#31](https://github.com/open-feature/codegen/issues/31)) ([757ab66](https://github.com/open-feature/codegen/commit/757ab66b7fde7103ca6f5cb7f10c0632073b58d8))


### 🔄 Refactoring

* change folder, package structure; integrate with cobra ([#27](https://github.com/open-feature/codegen/issues/27)) ([850c694](https://github.com/open-feature/codegen/commit/850c694c84fad1a71722a1b1e620f1473bc2d2ab))

## [0.1.1](https://github.com/open-feature/codegen/compare/v0.1.0...v0.1.1) (2024-10-04)


### 🐛 Bug Fixes

* release permissions ([#25](https://github.com/open-feature/codegen/issues/25)) ([dc07cdf](https://github.com/open-feature/codegen/commit/dc07cdfe5487c0a22209c54d0ee195bbdcf1b5ed))

## [0.1.0](https://github.com/open-feature/codegen/compare/v0.0.1...v0.1.0) (2024-10-04)


### ⚠ BREAKING CHANGES

* lower json schema version, rename number to float ([#12](https://github.com/open-feature/codegen/issues/12))

### ✨ New Features

* Fixing problems with generated code for golang and adding sample manifest for testing. ([#22](https://github.com/open-feature/codegen/issues/22)) ([558e964](https://github.com/open-feature/codegen/commit/558e9640b8756e9cacccfdb23f136d95bd81629b))
* initial CLI for codegen with support for golang strongly typed accessors ([#13](https://github.com/open-feature/codegen/issues/13)) ([e8f3d3e](https://github.com/open-feature/codegen/commit/e8f3d3ea2815b7d5473746e71f1bedc856e723c8))
* lower json schema version, rename number to float ([#12](https://github.com/open-feature/codegen/issues/12)) ([ed844b4](https://github.com/open-feature/codegen/commit/ed844b43a3d05113b49b39a1e368d0ee3c308dc9))


### 📚 Documentation

* Add initial flag manifest schema ([#9](https://github.com/open-feature/codegen/issues/9)) ([fac35ca](https://github.com/open-feature/codegen/commit/fac35caff88e1ef9a9c5ff1e8624040d91db9307))


### 🔄 Refactoring

* change the case of the flag manifest to camel case. ([#19](https://github.com/open-feature/codegen/issues/19)) ([fbac8ce](https://github.com/open-feature/codegen/commit/fbac8ce70dda766aff437b59286beb0579aa8472))
* embed flag manifest schema into code and moves files around ([#18](https://github.com/open-feature/codegen/issues/18)) ([aa9d3b0](https://github.com/open-feature/codegen/commit/aa9d3b03f0ece5295f6ce7be1f9093ed8ee9200f))
