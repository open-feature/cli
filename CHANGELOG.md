# Changelog

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
