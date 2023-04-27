# Changelog

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.

## [1.1.0](/git.hyperchain.cn/hyperchain/hyperchain/compare/v1.0.9...v1.1.0) (2023-04-27)


### Bug Fixes

* **go.mod:** set google.golang.org/grpc version for fix [#49](null/git.hyperchain.cn/hyperchain/hyperchain/issues/49) ([ebb3968](/git.hyperchain.cn/hyperchain/hyperchain/commit/ebb3968c14331723c04d30f8438aae56185c8099))

### [1.0.9](/git.hyperchain.cn/hyperchain/hyperchain/compare/v1.0.8...v1.0.9) (2023-03-09)

<a name="1.0.8"></a>
## [1.0.8](http://git.hyperchain.cn/hyperchain/hyperchain/compare/v1.0.7...v1.0.8) (2023-03-02)



<a name="1.0.7"></a>
## [1.0.7](http://git.hyperchain.cn/hyperchain/hyperchain/compare/v1.0.6...v1.0.7) (2023-02-16)

### [1.0.6](/git.hyperchain.cn/hyperchain/hyperchain/compare/v1.0.5...v1.0.6) (2023-02-20)

### Bug Fixes

* **controller.go:** [#42](null/git.hyperchain.cn/hyperchain/hyperchain/issues/42) Fix the error that the master cannot be stopped ([d1b4293d](/git.hyperchain.cn/hyperchain/hyperchain/commit/d1b4293d42d8d6c6d9ad373775c49b4ea99de0af))


### [1.0.5](/git.hyperchain.cn/hyperchain/hyperchain/compare/v1.0.4...v1.0.5) (2023-01-11)

### [1.0.4](/git.hyperchain.cn/hyperchain/hyperchain/compare/v1.0.3...v1.0.4) (2022-09-30)


### Features

* change VMPool ([4cb5feb](/git.hyperchain.cn/hyperchain/hyperchain/commit/4cb5feb26433a075a65693f04f7ea06a04ca29ad))
* influxdb ([dc9b3f6](/git.hyperchain.cn/hyperchain/hyperchain/commit/dc9b3f669bece21b61df979f91d78f1d513b3c96))


### Bug Fixes

* **glua:** fix repeated decode parameter and change result to Lua.UserData ([567d224](/git.hyperchain.cn/hyperchain/hyperchain/commit/567d224632abf7c8a01940a3584773fe8d5a105d))
* **go.mod:** fix [#32](null/git.hyperchain.cn/hyperchain/hyperchain/issues/32) update go.mod and go.sum ([e236c54](/git.hyperchain.cn/hyperchain/hyperchain/commit/e236c54ff78d435b375f98709d7d41356cd3671a))

### [1.0.3](/git.hyperchain.cn/hyperchain/hyperchain/compare/v1.0.2...v1.0.3) (2022-08-09)


### Features

* **verify:** add transaction sampling verification [#17](null/git.hyperchain.cn/hyperchain/hyperchain/issues/17) ([bbcbd90](/git.hyperchain.cn/hyperchain/hyperchain/commit/bbcbd90442464d97ba8810fa36d0d49d77a4a31b))


### Bug Fixes

* **gomod&manual:** update version of golang.org/x/sys and golang.org/x/net, update manual ([2bfea78](/git.hyperchain.cn/hyperchain/hyperchain/commit/2bfea78c9daf703f69128fb4f9e543eeda60ee6e))
* **network:** [#28](null/git.hyperchain.cn/hyperchain/hyperchain/issues/28) fix upload error since go 1.17 ([2919551](/git.hyperchain.cn/hyperchain/hyperchain/commit/2919551e37072bb4d5a8cc7b7641b39f4ad1126d))

### [1.0.2](/git.hyperchain.cn/hyperchain/hyperchain/compare/v1.0.1...v1.0.2) (2022-07-01)


### Features

* **doc:** close [#20](null/git.hyperchain.cn/hyperchain/hyperchain/issues/20) add development specification manual ([e4f7e28](/git.hyperchain.cn/hyperchain/hyperchain/commit/e4f7e28298262217be4300126cbb6111d0316729))


### Bug Fixes

* **benchmark&gomod:** fix hyperbench benchmark bug [#16](null/git.hyperchain.cn/hyperchain/hyperchain/issues/16) ([5627a24](/git.hyperchain.cn/hyperchain/hyperchain/commit/5627a24f4f0fe25cfc97c1518b5400337d2200c8))

### [1.0.1](/github.com/hyperbench/hyperbench/compare/v1.0.0...v1.0.1) (2022-05-20)


### Features

* **ci:** fix [#7](null/github.com/hyperbench/hyperbench/issues/7) update workflows to sync code to gitee ([8806e3e](/github.com/hyperbench/hyperbench/commit/8806e3e73bc924b3b099e4754cd677e36304e647))


### Bug Fixes

* **ci:** sync code to gitee when push to master or develop ([5457b80](/github.com/hyperbench/hyperbench/commit/5457b809b0e62d59431b5605e205510224333e60))
* **network:** fix distributed bug ([6bdf123](/github.com/hyperbench/hyperbench/commit/6bdf123d891322020990435693f59ce0e1491c6d)), closes [#12](/github.com/hyperbench/hyperbench/issues/12) [#10](/github.com/hyperbench/hyperbench/issues/10) [#11](/github.com/hyperbench/hyperbench/issues/11)


<a name="1.0.0"></a>
# 1.0.0 (2022-03-31)


### Bug Fixes

* **go.mod:** replace golang.org/x/sys ([68bb277](http://git.hyperchain.cn/hyperchain/hyperchain/commits/68bb277))
* **LICENSE:** change license protocol ([1c201c4](http://git.hyperchain.cn/hyperchain/hyperchain/commits/1c201c4))
* code modify based on code review ([410b23a](http://git.hyperchain.cn/hyperchain/hyperchain/commits/410b23a))
* missing eth transfer.From ([035db48](http://git.hyperchain.cn/hyperchain/hyperchain/commits/035db48))


### Features

* **xuperchain&plugins&tps:** add xuperchain & decouple code to plugins & modify tps ([7b6a54c](http://git.hyperchain.cn/hyperchain/hyperchain/commits/7b6a54c))
* add eth ([4769c6e](http://git.hyperchain.cn/hyperchain/hyperchain/commits/4769c6e))
* add eth ([83c3a44](http://git.hyperchain.cn/hyperchain/hyperchain/commits/83c3a44))
* add eth ([1a7934b](http://git.hyperchain.cn/hyperchain/hyperchain/commits/1a7934b))
* add eth ([ad73d70](http://git.hyperchain.cn/hyperchain/hyperchain/commits/ad73d70))
* add unittest & modify code ([ec14f2f](http://git.hyperchain.cn/hyperchain/hyperchain/commits/ec14f2f))
* **mod:** fabric ([163ba9c](http://git.hyperchain.cn/hyperchain/hyperchain/commits/163ba9c))



# Change Log

All notable changes to this project will be documented in this file. See [standard-version](https://github.com/conventional-changelog/standard-version) for commit guidelines.
