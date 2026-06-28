# Changelog

## [1.9.0](https://github.com/llttlltt/dj-library-tools/compare/v1.8.0...v1.9.0) (2026-06-28)


### Features

* implement := operator for strict exact equality in queries ([898267a](https://github.com/llttlltt/dj-library-tools/commit/898267a77f0ab343d4d0a877d15485842b934db3))
* implement active capability gating for providers ([fde00cb](https://github.com/llttlltt/dj-library-tools/commit/fde00cbb0d116983492323c3fec9e7e92e3e81b0))
* implement inference-driven sort field validation ([d02daf6](https://github.com/llttlltt/dj-library-tools/commit/d02daf69ae5092fe2aaba38445506754879cf7b7))
* implement universal dynamic table rendering for list command ([e96c2a8](https://github.com/llttlltt/dj-library-tools/commit/e96c2a8d229a37e8ec3ed749d803807a0e53fbca))
* **provider/m3u:** implement strict native fields and update cli rendering ([5925a2a](https://github.com/llttlltt/dj-library-tools/commit/5925a2ad47411c8a8f1abddd502d67c114018b54))
* restore hierarchical cue matching in Rekordbox core ([4a84c7d](https://github.com/llttlltt/dj-library-tools/commit/4a84c7dba98d3a02e06a26eca5ed7da317b0ad1e))


### Bug Fixes

* ensure providers are registered and config paths are respected ([fe2f22f](https://github.com/llttlltt/dj-library-tools/commit/fe2f22f366ab3353aff6235b55cca758492c5f42))
* resolve table alignment issues caused by ANSI color codes ([2b68d52](https://github.com/llttlltt/dj-library-tools/commit/2b68d52baa179e6fd6331f9405ef5dc2625dd658))

## [1.8.0](https://github.com/llttlltt/dj-library-tools/compare/v1.7.2...v1.8.0) (2026-06-28)


### Features

* **cli:** add global --missing and --exists flags for physical health filtering ([71d553b](https://github.com/llttlltt/dj-library-tools/commit/71d553b1307962c57cb7538ce007b7c77c685911))
* **cli:** decouple fix command and implement self-healing M3U provider ([a14f378](https://github.com/llttlltt/dj-library-tools/commit/a14f378584f1c88d73a7f2caa9dbf42848b8756b))
* **cli:** enhance verbosity and move metadata logic to providers ([3652b02](https://github.com/llttlltt/dj-library-tools/commit/3652b028e70eacb95ce1e6c31c8cf037bb8fc56b))
* **cli:** finalize 100% architectural agnosticism by removing last provider traces ([9fff55a](https://github.com/llttlltt/dj-library-tools/commit/9fff55a1c249e82d36049a8d83aae9513121086a))
* **cli:** finalize agnosticism and standardize file flags ([1431476](https://github.com/llttlltt/dj-library-tools/commit/14314764a92770934fc5048ceb5e79a1472c4341))
* **cli:** implement metadata reconciliation flags and dynamic matching in sync ([a9a53f1](https://github.com/llttlltt/dj-library-tools/commit/a9a53f19e226771916bd7f0ec574909eca358f5b))
* **cli:** implement modify command with bulk metadata updates and file relocation ([3efe008](https://github.com/llttlltt/dj-library-tools/commit/3efe008f6907a4ff22b94cbfda2f9dd642a74aab))
* **cli:** implement user-friendly error handling for provider sentinels ([d34600b](https://github.com/llttlltt/dj-library-tools/commit/d34600b083c04e8950f91828be4821fa40eeff7a))
* **cli:** remove update command from root command and delete orphan file ([c02f26c](https://github.com/llttlltt/dj-library-tools/commit/c02f26c647ad4bee6bfe85fa31e4462df5e6425a))
* **cli:** replace --xml flag with standardized --file flag ([133314c](https://github.com/llttlltt/dj-library-tools/commit/133314cb415b3e4567aefbcdd3343382281ae9f5))
* **cli:** use containment policies for agnostic command validation ([684cbfe](https://github.com/llttlltt/dj-library-tools/commit/684cbfe10d163264ed5d3ee1c1bb85d0758fa2b4))
* **core:** harden provider interfaces and introduce neutral metadata models ([6c5a9b3](https://github.com/llttlltt/dj-library-tools/commit/6c5a9b352500c3065706dfe915a8e58e923b8c15))
* **provider/rb:** implement full metadata compliance and enhanced verbose reporting ([46ee68a](https://github.com/llttlltt/dj-library-tools/commit/46ee68afcac287b988b048d5aede4cb6ff809576))
* **provider:** implement 0-255 unified rating scale for plex and rekordbox ([4d50bd3](https://github.com/llttlltt/dj-library-tools/commit/4d50bd3d8581c0231e2000a96b9df02bf6163ed4))
* **provider:** implement GetResources and unify agnostic selection logic ([4070b86](https://github.com/llttlltt/dj-library-tools/commit/4070b86310680116233758dc867bc6199c0c43b9))
* **provider:** implement SupportedResources and IdentifyGroup for full agnosticism ([796717e](https://github.com/llttlltt/dj-library-tools/commit/796717eb74f4b9f316dd203b40b1751620808e60))
* **provider:** move sorting business logic into providers ([9351418](https://github.com/llttlltt/dj-library-tools/commit/9351418b678dcbc860ebd5231717352f45e349ea))
* **provider:** split Provider interface into Readable, Searchable, and Writable ([24cf966](https://github.com/llttlltt/dj-library-tools/commit/24cf966df0abb074baca0ed7ace4000d50758b1b))
* **provider:** unify discovery in GetResources and harden writable interface ([9e63e8b](https://github.com/llttlltt/dj-library-tools/commit/9e63e8b61f0b5f673491a104ac3880b781f1cc78))
* **query:** delegate provider-specific matching to CustomMatcher interface ([98ee29c](https://github.com/llttlltt/dj-library-tools/commit/98ee29cc37d26d0086ccfebd0d411aaf56f2b085))
* **sync:** implement agnostic Join logic for metadata reconciliation ([e20f609](https://github.com/llttlltt/dj-library-tools/commit/e20f6094947c08287a978e158a7344f8574c9485))
* **utils:** make location parser purely syntactic and provider-agnostic ([27c7bdc](https://github.com/llttlltt/dj-library-tools/commit/27c7bdcc2372bacccc0c8a54f2ac933cd88f5b61))


### Bug Fixes

* **cli:** enforce rekordbox hierarchy constraints for tracks and nodes ([af83429](https://github.com/llttlltt/dj-library-tools/commit/af834299fd6a8d3ed2cba8f62d41d2c3b035700b))
* **cli:** resolve linting warnings and align CreateGroup signature ([32072fb](https://github.com/llttlltt/dj-library-tools/commit/32072fbeea4f41ceb505d5ef2cfce9a21415f453))
* **lint:** use tagged switch for kind in engine mock tests ([141be8b](https://github.com/llttlltt/dj-library-tools/commit/141be8b19af7e8ae75df5f2e16bd2a9b281b1e11))

## [1.7.2](https://github.com/llttlltt/dj-library-tools/compare/v1.7.1...v1.7.2) (2026-06-27)


### Bug Fixes

* **provider/rb:** correctly create playlists vs folders in CreateNode ([f6161b2](https://github.com/llttlltt/dj-library-tools/commit/f6161b282ff8f4a3cb3f67c0c5a18171fc622fb4))

## [1.7.1](https://github.com/llttlltt/dj-library-tools/compare/v1.7.0...v1.7.1) (2026-06-26)


### Bug Fixes

* **query:** implement missing remixer and mix fields in evaluator ([7aef494](https://github.com/llttlltt/dj-library-tools/commit/7aef494eca1cef369058e6fb19800284543bb0b7))
* **utils:** support colon separator in m3u locations ([5a946d2](https://github.com/llttlltt/dj-library-tools/commit/5a946d20264a9f95b6595e4cdb2f61ae1e07f94a))

## [1.7.0](https://github.com/llttlltt/dj-library-tools/compare/v1.6.0...v1.7.0) (2026-06-26)


### Features

* **cli:** rename create verb to make with mk alias ([cf0a6d9](https://github.com/llttlltt/dj-library-tools/commit/cf0a6d936bc385589cfaea1e0d3808cd2f0b6250))
* **provider:** add M3U and M3U8 playlist support ([8ef1536](https://github.com/llttlltt/dj-library-tools/commit/8ef153694895ca9533ef9f140a57a186b22ac560))


### Bug Fixes

* **cli:** support file-based providers in mv and rename ([8fd0bd0](https://github.com/llttlltt/dj-library-tools/commit/8fd0bd048f3ef442ab93d4c309fd94a162580440))
* **engine,provider:** rb/folders now correctly queries folder nodes ([c74e640](https://github.com/llttlltt/dj-library-tools/commit/c74e640197fda8a83b91fa6d5e37f0857608bd9b))
* **query,docs:** wire count alias; correct sync/append docs ([b225040](https://github.com/llttlltt/dj-library-tools/commit/b225040015e231b7a40f1fed255bef09e3547571))
* **rekordbox:** folder nodes now map Count to models.Node.Entries ([dda2626](https://github.com/llttlltt/dj-library-tools/commit/dda26265ac968905efd20e9981554cfd2c7ead72))

## [1.6.0](https://github.com/llttlltt/dj-library-tools/compare/v1.5.0...v1.6.0) (2026-06-26)


### Features

* **cli:** add --sort flag to list command ([6f29335](https://github.com/llttlltt/dj-library-tools/commit/6f293354077b52e9e3f3b906245fe4b6b3a9f6ee))
* **cli:** add ascending/descending sorting and unit tests ([cd0d2f7](https://github.com/llttlltt/dj-library-tools/commit/cd0d2f7f1a3daffde5bacb69bb9aef0c4f300efe))
* **cli:** add column width capping and truncation to table rendering ([a92f0b3](https://github.com/llttlltt/dj-library-tools/commit/a92f0b313b8d3a3cd9ba7ca998d0576c94709cd2))
* **cli:** harden all verbs with interface-based capability checks ([0c23d0c](https://github.com/llttlltt/dj-library-tools/commit/0c23d0cca3b94c3dab1b62434b6c82c4d7b1e810))
* **cli:** implement dynamic terminal-aware table scaling ([87bf290](https://github.com/llttlltt/dj-library-tools/commit/87bf2906981b4af00ff963ab43afcb546642cc4a))
* **cli:** implement unified table rendering for list results ([2b3ef50](https://github.com/llttlltt/dj-library-tools/commit/2b3ef50aaffd59fccb186b81f14e4603803dd719))
* **models:** unify track and node models under universal Resource interface ([22d3476](https://github.com/llttlltt/dj-library-tools/commit/22d347696a2c8a0aefa894184e24b7f134c89d5c))
* **provider/plex:** add metadata extraction and track-level filtering ([8657fe9](https://github.com/llttlltt/dj-library-tools/commit/8657fe90588e85574c7319d55815a4c29d6cafa7))
* **provider/plex:** aggregate tracks from all matching playlists ([dea840b](https://github.com/llttlltt/dj-library-tools/commit/dea840bcd97e0038a4e959b704441879e5139f9d))
* **provider/plex:** enforce field-based selection and add name resolution ([e82a85e](https://github.com/llttlltt/dj-library-tools/commit/e82a85ebb01c334453cf4f0bdf9d8685cc937dd8))
* **provider/plex:** implement operator-aware playlist resolution and regex filtering ([144d41f](https://github.com/llttlltt/dj-library-tools/commit/144d41f0c0c428bfdd243219d4c5570538315dce))
* **provider/plex:** support global track search when no playlist is specified ([1a9178d](https://github.com/llttlltt/dj-library-tools/commit/1a9178df4b7fe9d4f464e17eb4870cf84647efd6))
* **provider:** implement capability-based verbs and universal orchestrator ([df508cb](https://github.com/llttlltt/dj-library-tools/commit/df508cb87f251907104116c6030cb0c8ba5568be))
* **query:** align 'name' and 'title' fields and improve playlist error handling ([6b3185b](https://github.com/llttlltt/dj-library-tools/commit/6b3185b45f7e6fcb339fe369d50eaefb4871b83d))


### Bug Fixes

* **cli:** silence redundant error printing in root command ([f8f3148](https://github.com/llttlltt/dj-library-tools/commit/f8f31481586cdd87df5b0d7946f9623c0cd24052))
* **main:** remove unused fmt import ([915be4e](https://github.com/llttlltt/dj-library-tools/commit/915be4eb28d3dfbb8dc37058754536ddea5f9b1f))
* **provider/plex:** allow fuzzy name matching for playlist resolution ([369e50d](https://github.com/llttlltt/dj-library-tools/commit/369e50d0c9837344e81ba69f827aceb53a68fe68))
* **provider/plex:** enable field validation for track and playlist queries ([46843e3](https://github.com/llttlltt/dj-library-tools/commit/46843e359ea58fd3d23c1e210215a80d15e2a70f))
* **query:** support implicit AND for field comparisons ([990c012](https://github.com/llttlltt/dj-library-tools/commit/990c012394db1c660d1d8a7c02f8704d3b58d488))

## [1.5.0](https://github.com/llttlltt/dj-library-tools/compare/v1.4.0...v1.5.0) (2026-06-26)


### Features

* **cli:** add command aliases (ls, mv, rm, del, stats) ([c3619b7](https://github.com/llttlltt/dj-library-tools/commit/c3619b706fac2e4924c5141a1ec065c27375dd4b))
* **cli:** add detailed verbose output for add and remove commands ([a82c275](https://github.com/llttlltt/dj-library-tools/commit/a82c2756f0e9fd5cd1aa2c474163217b4aa60b26))
* **cli:** add global --json flag for list and stat commands ([3f0825f](https://github.com/llttlltt/dj-library-tools/commit/3f0825f1748697beaae24142552b168916a880ad))
* **cli:** add visual progress bars for bulk add and remove operations ([f09f010](https://github.com/llttlltt/dj-library-tools/commit/f09f01016aaa55eb14ce94f89a6d0155694e803a))
* **cli:** expand verbose output to all core commands ([5ed106a](https://github.com/llttlltt/dj-library-tools/commit/5ed106a3f1dca47f9348fc4bf711d839ce3b3997))
* **cli:** finalize global dry-run support for all modifying commands ([aa0a2a0](https://github.com/llttlltt/dj-library-tools/commit/aa0a2a0b153334fc59fc6d6cbf0115484e79bc82))
* **cli:** move dry-run to global persistent flag ([85adf39](https://github.com/llttlltt/dj-library-tools/commit/85adf396e0286ab34d5d80e2ef7912caf859ee9a))
* **cli:** move verbose to global persistent flag ([80c54c0](https://github.com/llttlltt/dj-library-tools/commit/80c54c0059d8dd177c487d8397d49c5d9ab8049e))
* **engine:** add mock engine test demonstrating Library boundary ([3fad172](https://github.com/llttlltt/dj-library-tools/commit/3fad172e84198da8a58719eaa7d07acca540a56b))
* **engine:** introduce Library interface and decouple from Rekordbox XML ([ff69614](https://github.com/llttlltt/dj-library-tools/commit/ff696142ce5f6fc39a2e4040494bf25ba8e66e0f))
* **provider:** introduce Provider interface and refactor list command for RB ([d31786f](https://github.com/llttlltt/dj-library-tools/commit/d31786fca3c13bf020b313464fefd99cad83c3e4))
* **query:** add query validation and helpful error for bare values ([001093a](https://github.com/llttlltt/dj-library-tools/commit/001093a76f980d846b64d6ecc65bf70db506f3f9))
* **query:** implement strict field validation and typo detection ([536c6e2](https://github.com/llttlltt/dj-library-tools/commit/536c6e22b93253139db5914bb2d19b28629dc3f8))
* **query:** improve bare value error message with reconstructed query ([3f12fb9](https://github.com/llttlltt/dj-library-tools/commit/3f12fb9f0223ef133f1ae590397321592c0bf5eb))
* **sys:** introduce system abstractions for filesystem and command execution ([c4b0795](https://github.com/llttlltt/dj-library-tools/commit/c4b07950c60f1859aa1f9b454b59b8b8a1f3fdaf))

## [1.4.0](https://github.com/llttlltt/dj-library-tools/compare/v1.3.0...v1.4.0) (2026-06-26)


### Features

* **cli:** finalize verb-centric architecture and cleanup legacy commands ([2df73b3](https://github.com/llttlltt/dj-library-tools/commit/2df73b39ae78477f60298c977a6e59c6c0c80896))
* **cli:** implement add and remove verbs ([b969bb6](https://github.com/llttlltt/dj-library-tools/commit/b969bb6756265813951b681afc1a60d8e0a998a5))
* **cli:** implement create and move verbs ([328b06d](https://github.com/llttlltt/dj-library-tools/commit/328b06dee4733dd1ce1eaa8a6f39ff529f0f68cf))
* **cli:** implement rename and delete verbs ([c3507a9](https://github.com/llttlltt/dj-library-tools/commit/c3507a994b164ea30aab0e28d7839427e505aa50))
* **cli:** refactor metadata command into update verb ([64a012a](https://github.com/llttlltt/dj-library-tools/commit/64a012ad83e0a58d1e3b2f7df74354c1c78b18bd))
* **cli:** refactor sync command to verb-centric architecture ([0eb787f](https://github.com/llttlltt/dj-library-tools/commit/0eb787f40eb691cdfce9d5a5d69cad0f8d5a1187))
* **query:** improve name matching and quoted string support ([2fcbcb9](https://github.com/llttlltt/dj-library-tools/commit/2fcbcb920ffcc01893decccdf1c6922c648e161d))
* **sync:** add folder creation and standardize node types ([529d1c5](https://github.com/llttlltt/dj-library-tools/commit/529d1c5ba845f4b4ca3c3264f02c91f97a6deb4d))

## [1.3.0](https://github.com/llttlltt/dj-library-tools/compare/v1.2.0...v1.3.0) (2026-06-26)


### Features

* **rb:** sync with rekordbox xml spec and standardize query engine ([db25650](https://github.com/llttlltt/dj-library-tools/commit/db25650478c1355895dd7131ed0e68701c89f72c))
* **rb:** sync with rekordbox xml spec and standardize query engine ([5883dc1](https://github.com/llttlltt/dj-library-tools/commit/5883dc1e50df36fca8ad6208e210f78461349f8b))
* **rekordbox:** implement high-fidelity XML emission and comprehensive testing ([cfaeebe](https://github.com/llttlltt/dj-library-tools/commit/cfaeebeb1fa11e9ee6c4eca3e40841d667701ff2))


### Bug Fixes

* **rekordbox:** repair surgical XML write and optional node attributes ([d3b5def](https://github.com/llttlltt/dj-library-tools/commit/d3b5def84bf41499bcce2066b3f335e4234052f1))

## [1.2.0](https://github.com/llttlltt/dj-library-tools/compare/v1.1.0...v1.2.0) (2026-06-24)


### Features

* **rekordbox:** match idiosyncratic attribute order for root and sub-nodes ([488932e](https://github.com/llttlltt/dj-library-tools/commit/488932ecd4ab64bd1aef917491b7600f7b9f4cd7))

## [1.1.0](https://github.com/llttlltt/dj-library-tools/compare/v1.0.0...v1.1.0) (2026-06-24)


### Features

* **cli:** add folder god command (--new, --rename, --move, --remove) ([b5620e5](https://github.com/llttlltt/dj-library-tools/commit/b5620e50c2814a6bee294658e4fefa9630e59696))
* **cli:** add Key column and refine ls table alignment ([3c5d202](https://github.com/llttlltt/dj-library-tools/commit/3c5d202ef1d0dbf1dcd9279e2bc37e907aec01d4))
* **cli:** add ls command and global --xml flag ([928428b](https://github.com/llttlltt/dj-library-tools/commit/928428b5521440e6d9bd99ac08bd9b5a8e0659be))
* **cli:** add playlist god command (--new, --add, --rename, --move, --remove) ([b736a86](https://github.com/llttlltt/dj-library-tools/commit/b736a860cbf7847233521a6f2178d9ac6a4de649))
* **cli:** add stat command for library aggregation and analysis ([b2f42f8](https://github.com/llttlltt/dj-library-tools/commit/b2f42f8a4851b033ada63bdf97049c6fbf9ee6f0))
* **cli:** colorize ls output and refine table formatting ([40129b5](https://github.com/llttlltt/dj-library-tools/commit/40129b57ecf730453dd78fb3da7050d710b3e5a8))
* **cli:** implement --dry-run for sync and add location tests ([027facc](https://github.com/llttlltt/dj-library-tools/commit/027facc9c235c05786a45e2e4d897a988508e4cc))
* **cli:** implement location-based query syntax provider/resource:query ([7b556cc](https://github.com/llttlltt/dj-library-tools/commit/7b556cc585ae27b8e36c7fa4875a093663a429cb))
* **cli:** implement multi-bar visual progress for parallel syncing ([ca0e548](https://github.com/llttlltt/dj-library-tools/commit/ca0e548458bc87378d0e8820271456e98cea5491))
* **cli:** restore legacy metadata move command and integrate into djlt ([f0cdec0](https://github.com/llttlltt/dj-library-tools/commit/f0cdec09e1c8f8bf61d9823d0a6b86a31a8e532f))
* **config:** add --remove-map flag to plex config command ([8f87edc](https://github.com/llttlltt/dj-library-tools/commit/8f87edcdcfade16fedb3ac4f667aafc1fd879714))
* **config:** implement persistent app config and automatic token saving ([67cbcb1](https://github.com/llttlltt/dj-library-tools/commit/67cbcb183862cca944eb3d41cb75912af8750c5f))
* **engine:** add LsPlaylists and LsFolders with NodeResult ([a6d4cea](https://github.com/llttlltt/dj-library-tools/commit/a6d4cea5a598121270ffabb211bcb0dfa3d8a04c))
* **engine:** implement core primitives ls, stat, and modify ([2826339](https://github.com/llttlltt/dj-library-tools/commit/28263399f061f4c8bff9f78d227454917746a644))
* **media:** implement smart skip and robust error handling in transcoder ([8d723d3](https://github.com/llttlltt/dj-library-tools/commit/8d723d30dedde71cdc88fcf9d1f992abe0945e31))
* **media:** implement transcode config, ffmpeg wrapper, and path formatting ([683b225](https://github.com/llttlltt/dj-library-tools/commit/683b2256240dfc94f39242aeb0ee2df8b3f9876c))
* **playlist:** add --output flag to fix command ([2b7d82e](https://github.com/llttlltt/dj-library-tools/commit/2b7d82efd2da3fe870669257a6d73689ebcac3b0))
* **playlist:** add --remove and --sync track operations, rename --remove to --delete ([bf11984](https://github.com/llttlltt/dj-library-tools/commit/bf119844f33b624e1d9d4da73047a01e85fd26fb))
* **playlist:** add batch processing and dry-run support ([2130546](https://github.com/llttlltt/dj-library-tools/commit/213054634c317b73ca15f6c5a97cba146593b010))
* **playlist:** add interactive removal prompt and clean up summary output ([6fb1f74](https://github.com/llttlltt/dj-library-tools/commit/6fb1f74646a07623b02e79a7d4deb0416158ee12))
* **playlist:** add missing file tracking and reporting to fix command ([b16e9cc](https://github.com/llttlltt/dj-library-tools/commit/b16e9cccf53e94931711658b1e0fcff6bbc2cbe2))
* **playlist:** add parity test and refine fix logic for legacy script compatibility ([5d5c79d](https://github.com/llttlltt/dj-library-tools/commit/5d5c79dc553a2eca2753ee397489c7c70237dc16))
* **playlist:** add playlist command group and fix subcommand ([af6dd16](https://github.com/llttlltt/dj-library-tools/commit/af6dd1631ded08c9deb4af7b094caa35e072d14c))
* **playlist:** add priority extension resolution and automatic track pruning ([a9bd3c4](https://github.com/llttlltt/dj-library-tools/commit/a9bd3c41a0fdf3c2218acd868273e0e8b1bf0352))
* **playlist:** add progress heartbeats and verbose logging to fix command ([d34acb9](https://github.com/llttlltt/dj-library-tools/commit/d34acb9d8a02647a67f4fd7f1f058b3cf9d4d80c))
* **playlist:** finalize feature-playlist-hygiene and update documentation ([853ee65](https://github.com/llttlltt/dj-library-tools/commit/853ee657b4e512651aea0a4a8f63306a0609eb69))
* **playlist:** implement native playlist and metadata logic ([d95902e](https://github.com/llttlltt/dj-library-tools/commit/d95902e6d00ad80aec91d3c97712da272b6e29f2))
* **playlist:** implement smart metadata fallback and M3U8 parity improvements ([6a5c9a4](https://github.com/llttlltt/dj-library-tools/commit/6a5c9a4c38142583861b68544f8ddffe26face01))
* **plex:** implement oauth pin flow and playlist retrieval ([eb382d7](https://github.com/llttlltt/dj-library-tools/commit/eb382d7da388b6b2c5a6f94b58dcac837f404cc7))
* **query:** add ability to query tracks by playlist name ([0940f5f](https://github.com/llttlltt/dj-library-tools/commit/0940f5f1160cae91edcb1e62e8a436c0aca85b84))
* **query:** add MatchesNode for playlist and folder node evaluation ([ee67c6f](https://github.com/llttlltt/dj-library-tools/commit/ee67c6f1d4af007c1de635b6749e41587c7b5fe7))
* **query:** add playlistcount and support for multiple playlist filters ([8d436d0](https://github.com/llttlltt/dj-library-tools/commit/8d436d02f4b286f4075e81771952ff8f03daaeb0))
* **query:** add shell-friendly numeric aliases (gt, lt, ge, le) and NOT keyword ([e94ca86](https://github.com/llttlltt/dj-library-tools/commit/e94ca862940e0afe791d20ab35e6da349af9a732))
* **query:** add support for '-' negation, 'eq', 'ne', 'neq' aliases, and '!=' operator ([fd645a5](https://github.com/llttlltt/dj-library-tools/commit/fd645a5f388c59a1f002fbb8990b5863eefe54fb))
* **query:** add support for rating, playcount, dateadded, and grouping fields ([28fcd78](https://github.com/llttlltt/dj-library-tools/commit/28fcd788bea95864b1da8f27ede0b3adde4609f9))
* **query:** exhaustive field mapping and numeric comparison operators ([f9b68f8](https://github.com/llttlltt/dj-library-tools/commit/f9b68f88b9d7ec8184993f8e5b3a4b2c4d8685d7))
* **query:** implement standardized query engine with lexer and evaluator ([ef169f0](https://github.com/llttlltt/dj-library-tools/commit/ef169f014508bc964f27d9a5e7419edf442e1cdf))
* **rekordbox:** align struct fields with XML attribute order for cleaner diffs ([fc3f0d0](https://github.com/llttlltt/dj-library-tools/commit/fc3f0d02cefa3418aa05bae1ffcc3a843cf012bf))
* **rekordbox:** detect and preserve XML formatting ([1483855](https://github.com/llttlltt/dj-library-tools/commit/1483855af0458224a90626e3689957409dde268e))
* **rekordbox:** implement surgical saving and playlist positioning ([85f7076](https://github.com/llttlltt/dj-library-tools/commit/85f70764878f797778d469765ed710cd415c4189))
* **rekordbox:** support attribute wrapping and fix omitempty issues ([f8bb0af](https://github.com/llttlltt/dj-library-tools/commit/f8bb0af7b7f7c0700d17772a65a228d24bb32cd7))
* **scaffold:** initialize standardized monorepo structure ([98c64a2](https://github.com/llttlltt/dj-library-tools/commit/98c64a2d21e3a02d73eac85b099d71f4a8e9561d))
* **sync:** add UpsertPlaylist, AddTracksToPlaylist, RenameNode, MoveNode, RemoveNode ([2c0bc70](https://github.com/llttlltt/dj-library-tools/commit/2c0bc709174f669426c12b73dc6711c4814bd3c1))
* **sync:** implement plex-to-m3u8 sync target ([9d4054f](https://github.com/llttlltt/dj-library-tools/commit/9d4054f422cdddcb2750a695db15d94dc7ea13de))
* **sync:** implement plex-to-rekordbox sync engine and cli ([e04b98a](https://github.com/llttlltt/dj-library-tools/commit/e04b98acc21f322133c6a03bc9ee11662c3f5c0e))
* **sync:** robust XML injection via sync.Engine ([4ef6532](https://github.com/llttlltt/dj-library-tools/commit/4ef653260712dfcaebc40529a04967a15c32364c))


### Bug Fixes

* **cli:** correct header alignment in ls output ([8634162](https://github.com/llttlltt/dj-library-tools/commit/86341621195d1b2dab094de0dffc483e181afc80))
* **cli:** finalize precise alignment for ls table output ([54af0fb](https://github.com/llttlltt/dj-library-tools/commit/54af0fbf486617b9362464f8469fb472f5c46498))
* **cli:** left-align headers in ls output ([6ef414c](https://github.com/llttlltt/dj-library-tools/commit/6ef414cf691a47d59f9193a318a05b0071443b8a))
* **cli:** prevent alignment spaces from being underlined in ls headers ([90d1836](https://github.com/llttlltt/dj-library-tools/commit/90d1836250220b326c9aa002ae9d584478c44672))
* **cli:** remove duplicate playlist command registration ([f1d2a56](https://github.com/llttlltt/dj-library-tools/commit/f1d2a562478fb03e38897b41e9f91ddf83f36560))
* **cli:** resolve header misalignment caused by ANSI color codes ([eb4c6fb](https://github.com/llttlltt/dj-library-tools/commit/eb4c6fb0050557b269a04625a3d22c23652ebfcd))
* **playlist:** avoid cross-device link error by using output dir for temp files ([8247e4e](https://github.com/llttlltt/dj-library-tools/commit/8247e4ef4888fe1655cb7849121750a76dd6fabc))
* **plex:** restore missing CheckPin and rename redeclared lsCmd ([aa26cdd](https://github.com/llttlltt/dj-library-tools/commit/aa26cddb72ed01903a64e44a0ca334d9c5fee1bf))
* **query:** correct BPM extraction and improve multi-word field matching ([161dfcc](https://github.com/llttlltt/dj-library-tools/commit/161dfccf26016a0c2d885189f0893b830e30f22a))
* **query:** exact numeric equality for playlistcount and numeric fields ([4547176](https://github.com/llttlltt/dj-library-tools/commit/454717626b37552ac5f07de91a49ed83ec8cf8e4))
* **query:** fix parser bug with multi-word queries and update tests ([cc0edb6](https://github.com/llttlltt/dj-library-tools/commit/cc0edb65945550a4519399439ed777461e364979))
* **query:** implement robust operator parsing and multi-word token joining ([b14468d](https://github.com/llttlltt/dj-library-tools/commit/b14468dddc6f4456150f45a5aba5c92ba49fc5ea))
* **stat,media:** config XML path fallback and filename sanitization ([73afa4d](https://github.com/llttlltt/dj-library-tools/commit/73afa4d8b03c7a8790f0d627ed7c9894b5ea27ba))
* **tests:** update test data to match new Rekordbox string-based numeric fields ([2400cc0](https://github.com/llttlltt/dj-library-tools/commit/2400cc0357576cce42efcc7557bc160e4dc9dbb6))

## 1.0.0 (2026-06-24)


### Features

* **cli:** add folder god command (--new, --rename, --move, --remove) ([b5620e5](https://github.com/llttlltt/dj-library-tools/commit/b5620e50c2814a6bee294658e4fefa9630e59696))
* **cli:** add Key column and refine ls table alignment ([3c5d202](https://github.com/llttlltt/dj-library-tools/commit/3c5d202ef1d0dbf1dcd9279e2bc37e907aec01d4))
* **cli:** add ls command and global --xml flag ([928428b](https://github.com/llttlltt/dj-library-tools/commit/928428b5521440e6d9bd99ac08bd9b5a8e0659be))
* **cli:** add playlist god command (--new, --add, --rename, --move, --remove) ([b736a86](https://github.com/llttlltt/dj-library-tools/commit/b736a860cbf7847233521a6f2178d9ac6a4de649))
* **cli:** add stat command for library aggregation and analysis ([b2f42f8](https://github.com/llttlltt/dj-library-tools/commit/b2f42f8a4851b033ada63bdf97049c6fbf9ee6f0))
* **cli:** colorize ls output and refine table formatting ([40129b5](https://github.com/llttlltt/dj-library-tools/commit/40129b57ecf730453dd78fb3da7050d710b3e5a8))
* **cli:** implement --dry-run for sync and add location tests ([027facc](https://github.com/llttlltt/dj-library-tools/commit/027facc9c235c05786a45e2e4d897a988508e4cc))
* **cli:** implement location-based query syntax provider/resource:query ([7b556cc](https://github.com/llttlltt/dj-library-tools/commit/7b556cc585ae27b8e36c7fa4875a093663a429cb))
* **cli:** implement multi-bar visual progress for parallel syncing ([ca0e548](https://github.com/llttlltt/dj-library-tools/commit/ca0e548458bc87378d0e8820271456e98cea5491))
* **cli:** restore legacy metadata move command and integrate into djlt ([f0cdec0](https://github.com/llttlltt/dj-library-tools/commit/f0cdec09e1c8f8bf61d9823d0a6b86a31a8e532f))
* **config:** add --remove-map flag to plex config command ([8f87edc](https://github.com/llttlltt/dj-library-tools/commit/8f87edcdcfade16fedb3ac4f667aafc1fd879714))
* **config:** implement persistent app config and automatic token saving ([67cbcb1](https://github.com/llttlltt/dj-library-tools/commit/67cbcb183862cca944eb3d41cb75912af8750c5f))
* **engine:** add LsPlaylists and LsFolders with NodeResult ([a6d4cea](https://github.com/llttlltt/dj-library-tools/commit/a6d4cea5a598121270ffabb211bcb0dfa3d8a04c))
* **engine:** implement core primitives ls, stat, and modify ([2826339](https://github.com/llttlltt/dj-library-tools/commit/28263399f061f4c8bff9f78d227454917746a644))
* **media:** implement smart skip and robust error handling in transcoder ([8d723d3](https://github.com/llttlltt/dj-library-tools/commit/8d723d30dedde71cdc88fcf9d1f992abe0945e31))
* **media:** implement transcode config, ffmpeg wrapper, and path formatting ([683b225](https://github.com/llttlltt/dj-library-tools/commit/683b2256240dfc94f39242aeb0ee2df8b3f9876c))
* **playlist:** add --output flag to fix command ([2b7d82e](https://github.com/llttlltt/dj-library-tools/commit/2b7d82efd2da3fe870669257a6d73689ebcac3b0))
* **playlist:** add --remove and --sync track operations, rename --remove to --delete ([bf11984](https://github.com/llttlltt/dj-library-tools/commit/bf119844f33b624e1d9d4da73047a01e85fd26fb))
* **playlist:** add batch processing and dry-run support ([2130546](https://github.com/llttlltt/dj-library-tools/commit/213054634c317b73ca15f6c5a97cba146593b010))
* **playlist:** add interactive removal prompt and clean up summary output ([6fb1f74](https://github.com/llttlltt/dj-library-tools/commit/6fb1f74646a07623b02e79a7d4deb0416158ee12))
* **playlist:** add missing file tracking and reporting to fix command ([b16e9cc](https://github.com/llttlltt/dj-library-tools/commit/b16e9cccf53e94931711658b1e0fcff6bbc2cbe2))
* **playlist:** add parity test and refine fix logic for legacy script compatibility ([5d5c79d](https://github.com/llttlltt/dj-library-tools/commit/5d5c79dc553a2eca2753ee397489c7c70237dc16))
* **playlist:** add playlist command group and fix subcommand ([af6dd16](https://github.com/llttlltt/dj-library-tools/commit/af6dd1631ded08c9deb4af7b094caa35e072d14c))
* **playlist:** add priority extension resolution and automatic track pruning ([a9bd3c4](https://github.com/llttlltt/dj-library-tools/commit/a9bd3c41a0fdf3c2218acd868273e0e8b1bf0352))
* **playlist:** add progress heartbeats and verbose logging to fix command ([d34acb9](https://github.com/llttlltt/dj-library-tools/commit/d34acb9d8a02647a67f4fd7f1f058b3cf9d4d80c))
* **playlist:** finalize feature-playlist-hygiene and update documentation ([853ee65](https://github.com/llttlltt/dj-library-tools/commit/853ee657b4e512651aea0a4a8f63306a0609eb69))
* **playlist:** implement native playlist and metadata logic ([d95902e](https://github.com/llttlltt/dj-library-tools/commit/d95902e6d00ad80aec91d3c97712da272b6e29f2))
* **playlist:** implement smart metadata fallback and M3U8 parity improvements ([6a5c9a4](https://github.com/llttlltt/dj-library-tools/commit/6a5c9a4c38142583861b68544f8ddffe26face01))
* **plex:** implement oauth pin flow and playlist retrieval ([eb382d7](https://github.com/llttlltt/dj-library-tools/commit/eb382d7da388b6b2c5a6f94b58dcac837f404cc7))
* **query:** add ability to query tracks by playlist name ([0940f5f](https://github.com/llttlltt/dj-library-tools/commit/0940f5f1160cae91edcb1e62e8a436c0aca85b84))
* **query:** add MatchesNode for playlist and folder node evaluation ([ee67c6f](https://github.com/llttlltt/dj-library-tools/commit/ee67c6f1d4af007c1de635b6749e41587c7b5fe7))
* **query:** add playlistcount and support for multiple playlist filters ([8d436d0](https://github.com/llttlltt/dj-library-tools/commit/8d436d02f4b286f4075e81771952ff8f03daaeb0))
* **query:** add shell-friendly numeric aliases (gt, lt, ge, le) and NOT keyword ([e94ca86](https://github.com/llttlltt/dj-library-tools/commit/e94ca862940e0afe791d20ab35e6da349af9a732))
* **query:** add support for '-' negation, 'eq', 'ne', 'neq' aliases, and '!=' operator ([fd645a5](https://github.com/llttlltt/dj-library-tools/commit/fd645a5f388c59a1f002fbb8990b5863eefe54fb))
* **query:** add support for rating, playcount, dateadded, and grouping fields ([28fcd78](https://github.com/llttlltt/dj-library-tools/commit/28fcd788bea95864b1da8f27ede0b3adde4609f9))
* **query:** exhaustive field mapping and numeric comparison operators ([f9b68f8](https://github.com/llttlltt/dj-library-tools/commit/f9b68f88b9d7ec8184993f8e5b3a4b2c4d8685d7))
* **query:** implement standardized query engine with lexer and evaluator ([ef169f0](https://github.com/llttlltt/dj-library-tools/commit/ef169f014508bc964f27d9a5e7419edf442e1cdf))
* **scaffold:** initialize standardized monorepo structure ([98c64a2](https://github.com/llttlltt/dj-library-tools/commit/98c64a2d21e3a02d73eac85b099d71f4a8e9561d))
* **sync:** add UpsertPlaylist, AddTracksToPlaylist, RenameNode, MoveNode, RemoveNode ([2c0bc70](https://github.com/llttlltt/dj-library-tools/commit/2c0bc709174f669426c12b73dc6711c4814bd3c1))
* **sync:** implement plex-to-m3u8 sync target ([9d4054f](https://github.com/llttlltt/dj-library-tools/commit/9d4054f422cdddcb2750a695db15d94dc7ea13de))
* **sync:** implement plex-to-rekordbox sync engine and cli ([e04b98a](https://github.com/llttlltt/dj-library-tools/commit/e04b98acc21f322133c6a03bc9ee11662c3f5c0e))
* **sync:** robust XML injection via sync.Engine ([4ef6532](https://github.com/llttlltt/dj-library-tools/commit/4ef653260712dfcaebc40529a04967a15c32364c))


### Bug Fixes

* **cli:** correct header alignment in ls output ([8634162](https://github.com/llttlltt/dj-library-tools/commit/86341621195d1b2dab094de0dffc483e181afc80))
* **cli:** finalize precise alignment for ls table output ([54af0fb](https://github.com/llttlltt/dj-library-tools/commit/54af0fbf486617b9362464f8469fb472f5c46498))
* **cli:** left-align headers in ls output ([6ef414c](https://github.com/llttlltt/dj-library-tools/commit/6ef414cf691a47d59f9193a318a05b0071443b8a))
* **cli:** prevent alignment spaces from being underlined in ls headers ([90d1836](https://github.com/llttlltt/dj-library-tools/commit/90d1836250220b326c9aa002ae9d584478c44672))
* **cli:** remove duplicate playlist command registration ([f1d2a56](https://github.com/llttlltt/dj-library-tools/commit/f1d2a562478fb03e38897b41e9f91ddf83f36560))
* **cli:** resolve header misalignment caused by ANSI color codes ([eb4c6fb](https://github.com/llttlltt/dj-library-tools/commit/eb4c6fb0050557b269a04625a3d22c23652ebfcd))
* **playlist:** avoid cross-device link error by using output dir for temp files ([8247e4e](https://github.com/llttlltt/dj-library-tools/commit/8247e4ef4888fe1655cb7849121750a76dd6fabc))
* **plex:** restore missing CheckPin and rename redeclared lsCmd ([aa26cdd](https://github.com/llttlltt/dj-library-tools/commit/aa26cddb72ed01903a64e44a0ca334d9c5fee1bf))
* **query:** correct BPM extraction and improve multi-word field matching ([161dfcc](https://github.com/llttlltt/dj-library-tools/commit/161dfccf26016a0c2d885189f0893b830e30f22a))
* **query:** exact numeric equality for playlistcount and numeric fields ([4547176](https://github.com/llttlltt/dj-library-tools/commit/454717626b37552ac5f07de91a49ed83ec8cf8e4))
* **query:** fix parser bug with multi-word queries and update tests ([cc0edb6](https://github.com/llttlltt/dj-library-tools/commit/cc0edb65945550a4519399439ed777461e364979))
* **query:** implement robust operator parsing and multi-word token joining ([b14468d](https://github.com/llttlltt/dj-library-tools/commit/b14468dddc6f4456150f45a5aba5c92ba49fc5ea))
* **stat,media:** config XML path fallback and filename sanitization ([73afa4d](https://github.com/llttlltt/dj-library-tools/commit/73afa4d8b03c7a8790f0d627ed7c9894b5ea27ba))

## 1.0.0 (2026-06-24)


### Features

* **cli:** add folder god command (--new, --rename, --move, --remove) ([ca6a459](https://github.com/llttlltt/dj-library-tools/commit/ca6a459bbff8b9f66ac2fef8c790ce4300027d07))
* **cli:** add Key column and refine ls table alignment ([76cd463](https://github.com/llttlltt/dj-library-tools/commit/76cd463f1ac7616788b6566ccbd81c37198692dc))
* **cli:** add ls command and global --xml flag ([f8dfddc](https://github.com/llttlltt/dj-library-tools/commit/f8dfddc657822c85fe3a1186f280c9d20649dd19))
* **cli:** add playlist god command (--new, --add, --rename, --move, --remove) ([d2019e6](https://github.com/llttlltt/dj-library-tools/commit/d2019e60566fd2300bc0afe77cae47fc26d3e0ef))
* **cli:** add stat command for library aggregation and analysis ([d5e8280](https://github.com/llttlltt/dj-library-tools/commit/d5e82808eb9606ee3d257529fe931d1493de093c))
* **cli:** colorize ls output and refine table formatting ([16b0081](https://github.com/llttlltt/dj-library-tools/commit/16b008163f689246103b41551b4c3a5f92783ea7))
* **cli:** implement --dry-run for sync and add location tests ([b943416](https://github.com/llttlltt/dj-library-tools/commit/b943416cacb1d771ff610dcebfc5603bafc91e1c))
* **cli:** implement location-based query syntax provider/resource:query ([f000cca](https://github.com/llttlltt/dj-library-tools/commit/f000ccab4174bb9d9b969805911f90f48566aea0))
* **cli:** implement multi-bar visual progress for parallel syncing ([7b46f6b](https://github.com/llttlltt/dj-library-tools/commit/7b46f6b890514527c1354841c63b97c682779be0))
* **cli:** restore legacy metadata move command and integrate into djlt ([a074d9d](https://github.com/llttlltt/dj-library-tools/commit/a074d9d4201a2cb325900d1f6db8549e22e96580))
* **config:** add --remove-map flag to plex config command ([e53bee5](https://github.com/llttlltt/dj-library-tools/commit/e53bee5c131953646936bf70dc5f1c7182794ae5))
* **config:** implement persistent app config and automatic token saving ([aee3b02](https://github.com/llttlltt/dj-library-tools/commit/aee3b02c67650a10d92042da7c341499aa8c337f))
* **engine:** add LsPlaylists and LsFolders with NodeResult ([6714864](https://github.com/llttlltt/dj-library-tools/commit/6714864665df88d6d70b3604993d754bca0e3446))
* **engine:** implement core primitives ls, stat, and modify ([fdf2dfd](https://github.com/llttlltt/dj-library-tools/commit/fdf2dfdf0d448384ef3f3054190ced03fe151515))
* **media:** implement smart skip and robust error handling in transcoder ([a6d0a98](https://github.com/llttlltt/dj-library-tools/commit/a6d0a983e83d0a0f84badf5cddaf2a9761bbd83a))
* **media:** implement transcode config, ffmpeg wrapper, and path formatting ([ad50a1a](https://github.com/llttlltt/dj-library-tools/commit/ad50a1ab32fe28f193ceb96a5add1515f464b478))
* **playlist:** add --output flag to fix command ([59882c7](https://github.com/llttlltt/dj-library-tools/commit/59882c7f05070f0e68a5170578eab95fbbc150ad))
* **playlist:** add --remove and --sync track operations, rename --remove to --delete ([887b5db](https://github.com/llttlltt/dj-library-tools/commit/887b5db13e24b9e6fa9271da9570b99381f487fa))
* **playlist:** add batch processing and dry-run support ([76f9dc6](https://github.com/llttlltt/dj-library-tools/commit/76f9dc6ab3837bb844453085f0a8d016451c9ac6))
* **playlist:** add interactive removal prompt and clean up summary output ([61a76f2](https://github.com/llttlltt/dj-library-tools/commit/61a76f2455241a76c424ae150a857d853570a191))
* **playlist:** add missing file tracking and reporting to fix command ([6dd6823](https://github.com/llttlltt/dj-library-tools/commit/6dd682328a55b7d64c1226ca92c9398be548bc66))
* **playlist:** add parity test and refine fix logic for legacy script compatibility ([de9d3f5](https://github.com/llttlltt/dj-library-tools/commit/de9d3f5e3f3244f805adc3cea0f9ab59829a0d94))
* **playlist:** add playlist command group and fix subcommand ([f797640](https://github.com/llttlltt/dj-library-tools/commit/f79764097bb2a2a0bb70b5333f6602fe83003a8c))
* **playlist:** add priority extension resolution and automatic track pruning ([22d2d7b](https://github.com/llttlltt/dj-library-tools/commit/22d2d7bfdead08b4a41bc9a7ef8d5a566ea1779c))
* **playlist:** add progress heartbeats and verbose logging to fix command ([c452e2f](https://github.com/llttlltt/dj-library-tools/commit/c452e2f82fcbd6b69a1e17f94a6ca9161f726e72))
* **playlist:** finalize feature-playlist-hygiene and update documentation ([ddce649](https://github.com/llttlltt/dj-library-tools/commit/ddce64956ba1fba9e7eb83b699d25c9475bc64e9))
* **playlist:** implement native playlist and metadata logic ([0975b80](https://github.com/llttlltt/dj-library-tools/commit/0975b8048aaeb27bca4eba2b458312ac05a0c53b))
* **playlist:** implement smart metadata fallback and M3U8 parity improvements ([801eb45](https://github.com/llttlltt/dj-library-tools/commit/801eb452cd78c124d11ec6f56c524cb9df745228))
* **plex:** implement oauth pin flow and playlist retrieval ([4f80a29](https://github.com/llttlltt/dj-library-tools/commit/4f80a2998deeb95cffdbab9422bfd038c15693d0))
* **query:** add ability to query tracks by playlist name ([49e0282](https://github.com/llttlltt/dj-library-tools/commit/49e028245a0d8f757c008e6f52ffa89986b20966))
* **query:** add MatchesNode for playlist and folder node evaluation ([1d0bfbd](https://github.com/llttlltt/dj-library-tools/commit/1d0bfbd7951a0df60c89a4ad7f0a1d102d9bf398))
* **query:** add playlistcount and support for multiple playlist filters ([990ac75](https://github.com/llttlltt/dj-library-tools/commit/990ac7516da1f43fe78d74852084f59beed2b421))
* **query:** add shell-friendly numeric aliases (gt, lt, ge, le) and NOT keyword ([2dd3ffe](https://github.com/llttlltt/dj-library-tools/commit/2dd3ffe826d018196767e4c88e93c870174d656b))
* **query:** add support for '-' negation, 'eq', 'ne', 'neq' aliases, and '!=' operator ([6104a5c](https://github.com/llttlltt/dj-library-tools/commit/6104a5c16010008ba0f8693f623968c648b465b5))
* **query:** add support for rating, playcount, dateadded, and grouping fields ([ad5df27](https://github.com/llttlltt/dj-library-tools/commit/ad5df276b2ccad781ef0d6ea87eea361d746b024))
* **query:** exhaustive field mapping and numeric comparison operators ([e73f878](https://github.com/llttlltt/dj-library-tools/commit/e73f878cc706dec917126618264f3b0da0524e64))
* **query:** implement standardized query engine with lexer and evaluator ([fa01d0a](https://github.com/llttlltt/dj-library-tools/commit/fa01d0a93655e2a2b892fb66e688fc4999e379be))
* **scaffold:** initialize standardized monorepo structure ([b344dce](https://github.com/llttlltt/dj-library-tools/commit/b344dce51f130de8ecf2a1b0a0802f7e83f7c90d))
* **sync:** add UpsertPlaylist, AddTracksToPlaylist, RenameNode, MoveNode, RemoveNode ([3c8bee0](https://github.com/llttlltt/dj-library-tools/commit/3c8bee002668da91fe389f5c4d323367a41f2ff6))
* **sync:** implement plex-to-m3u8 sync target ([058b7fa](https://github.com/llttlltt/dj-library-tools/commit/058b7fa97a3e433b283818bb1e4f8bab0b8ec263))
* **sync:** implement plex-to-rekordbox sync engine and cli ([51eb069](https://github.com/llttlltt/dj-library-tools/commit/51eb06980c3a74b24b33470da9901bdd5091e9f7))
* **sync:** robust XML injection via sync.Engine ([41c815e](https://github.com/llttlltt/dj-library-tools/commit/41c815e9f87b57957b8e6581dfffcf4f585cdb14))


### Bug Fixes

* **cli:** correct header alignment in ls output ([267261d](https://github.com/llttlltt/dj-library-tools/commit/267261d7f4a98422ad1547c66fb8a8b0a881936f))
* **cli:** finalize precise alignment for ls table output ([f73ae39](https://github.com/llttlltt/dj-library-tools/commit/f73ae3904fcbd12cdd6332d9977f8c22771c0f3d))
* **cli:** left-align headers in ls output ([6f2e6ac](https://github.com/llttlltt/dj-library-tools/commit/6f2e6acb3da483cf4af3a49afcc80302de98ce4e))
* **cli:** prevent alignment spaces from being underlined in ls headers ([f35465e](https://github.com/llttlltt/dj-library-tools/commit/f35465e5468c583ad49ba86d17009763f9284149))
* **cli:** remove duplicate playlist command registration ([9b3eeb3](https://github.com/llttlltt/dj-library-tools/commit/9b3eeb308aa03d0bcff759bdb2e535a004b63a97))
* **cli:** resolve header misalignment caused by ANSI color codes ([de6a1eb](https://github.com/llttlltt/dj-library-tools/commit/de6a1eb021ad52809d462aea6fa6771bc463ef48))
* **playlist:** avoid cross-device link error by using output dir for temp files ([d6c1001](https://github.com/llttlltt/dj-library-tools/commit/d6c10016dcdd4f067dc00c4b418c682dbe6ccfd6))
* **plex:** restore missing CheckPin and rename redeclared lsCmd ([082dfb7](https://github.com/llttlltt/dj-library-tools/commit/082dfb78faa5fa4d0ee26ff745643375f40a4f1d))
* **query:** correct BPM extraction and improve multi-word field matching ([68a5a79](https://github.com/llttlltt/dj-library-tools/commit/68a5a7909e708c90981f200b8adbfbadc6686c11))
* **query:** exact numeric equality for playlistcount and numeric fields ([b7cffc6](https://github.com/llttlltt/dj-library-tools/commit/b7cffc6c988aea4cb8a39ab61669ae32a3cad317))
* **query:** fix parser bug with multi-word queries and update tests ([ef15a4a](https://github.com/llttlltt/dj-library-tools/commit/ef15a4aad94adfb69debbe84f0f533f925f2bde8))
* **query:** implement robust operator parsing and multi-word token joining ([e183ea1](https://github.com/llttlltt/dj-library-tools/commit/e183ea19bc724c63fc833aba3c9eda2d48be7e60))
* **stat,media:** config XML path fallback and filename sanitization ([7325663](https://github.com/llttlltt/dj-library-tools/commit/7325663e6f84d026a9396d23c9b8a73d2e715e23))

## 1.0.0 (2026-06-23)


### Features

* **cli:** add folder god command (--new, --rename, --move, --remove) ([633c0b1](https://github.com/llttlltt/dj-library-tools/commit/633c0b11f8db96a16da1382813af956667dc74a5))
* **cli:** add Key column and refine ls table alignment ([76cd463](https://github.com/llttlltt/dj-library-tools/commit/76cd463f1ac7616788b6566ccbd81c37198692dc))
* **cli:** add ls command and global --xml flag ([f8dfddc](https://github.com/llttlltt/dj-library-tools/commit/f8dfddc657822c85fe3a1186f280c9d20649dd19))
* **cli:** add playlist god command (--new, --add, --rename, --move, --remove) ([4ad232b](https://github.com/llttlltt/dj-library-tools/commit/4ad232bc296d1b3ac6a86d1bdefa9b2842e50390))
* **cli:** add stat command for library aggregation and analysis ([d5e8280](https://github.com/llttlltt/dj-library-tools/commit/d5e82808eb9606ee3d257529fe931d1493de093c))
* **cli:** colorize ls output and refine table formatting ([16b0081](https://github.com/llttlltt/dj-library-tools/commit/16b008163f689246103b41551b4c3a5f92783ea7))
* **cli:** implement --dry-run for sync and add location tests ([b943416](https://github.com/llttlltt/dj-library-tools/commit/b943416cacb1d771ff610dcebfc5603bafc91e1c))
* **cli:** implement location-based query syntax provider/resource:query ([f000cca](https://github.com/llttlltt/dj-library-tools/commit/f000ccab4174bb9d9b969805911f90f48566aea0))
* **cli:** implement multi-bar visual progress for parallel syncing ([7b46f6b](https://github.com/llttlltt/dj-library-tools/commit/7b46f6b890514527c1354841c63b97c682779be0))
* **cli:** restore legacy metadata move command and integrate into djlt ([a074d9d](https://github.com/llttlltt/dj-library-tools/commit/a074d9d4201a2cb325900d1f6db8549e22e96580))
* **config:** add --remove-map flag to plex config command ([d6918bb](https://github.com/llttlltt/dj-library-tools/commit/d6918bb92c0cb06df6929af7a3ddc785fd85ad78))
* **config:** implement persistent app config and automatic token saving ([aee3b02](https://github.com/llttlltt/dj-library-tools/commit/aee3b02c67650a10d92042da7c341499aa8c337f))
* **engine:** add LsPlaylists and LsFolders with NodeResult ([583366d](https://github.com/llttlltt/dj-library-tools/commit/583366d2679f9559541e62b5946a0e75fdb5c757))
* **engine:** implement core primitives ls, stat, and modify ([fdf2dfd](https://github.com/llttlltt/dj-library-tools/commit/fdf2dfdf0d448384ef3f3054190ced03fe151515))
* **media:** implement smart skip and robust error handling in transcoder ([a6d0a98](https://github.com/llttlltt/dj-library-tools/commit/a6d0a983e83d0a0f84badf5cddaf2a9761bbd83a))
* **media:** implement transcode config, ffmpeg wrapper, and path formatting ([ad50a1a](https://github.com/llttlltt/dj-library-tools/commit/ad50a1ab32fe28f193ceb96a5add1515f464b478))
* **playlist:** add --output flag to fix command ([59882c7](https://github.com/llttlltt/dj-library-tools/commit/59882c7f05070f0e68a5170578eab95fbbc150ad))
* **playlist:** add --remove and --sync track operations, rename --remove to --delete ([8d3341b](https://github.com/llttlltt/dj-library-tools/commit/8d3341b1832597cf4c4860449cdde24eed27ce91))
* **playlist:** add batch processing and dry-run support ([76f9dc6](https://github.com/llttlltt/dj-library-tools/commit/76f9dc6ab3837bb844453085f0a8d016451c9ac6))
* **playlist:** add interactive removal prompt and clean up summary output ([61a76f2](https://github.com/llttlltt/dj-library-tools/commit/61a76f2455241a76c424ae150a857d853570a191))
* **playlist:** add missing file tracking and reporting to fix command ([6dd6823](https://github.com/llttlltt/dj-library-tools/commit/6dd682328a55b7d64c1226ca92c9398be548bc66))
* **playlist:** add parity test and refine fix logic for legacy script compatibility ([de9d3f5](https://github.com/llttlltt/dj-library-tools/commit/de9d3f5e3f3244f805adc3cea0f9ab59829a0d94))
* **playlist:** add playlist command group and fix subcommand ([f797640](https://github.com/llttlltt/dj-library-tools/commit/f79764097bb2a2a0bb70b5333f6602fe83003a8c))
* **playlist:** add priority extension resolution and automatic track pruning ([22d2d7b](https://github.com/llttlltt/dj-library-tools/commit/22d2d7bfdead08b4a41bc9a7ef8d5a566ea1779c))
* **playlist:** add progress heartbeats and verbose logging to fix command ([c452e2f](https://github.com/llttlltt/dj-library-tools/commit/c452e2f82fcbd6b69a1e17f94a6ca9161f726e72))
* **playlist:** finalize feature-playlist-hygiene and update documentation ([ddce649](https://github.com/llttlltt/dj-library-tools/commit/ddce64956ba1fba9e7eb83b699d25c9475bc64e9))
* **playlist:** implement native playlist and metadata logic ([0975b80](https://github.com/llttlltt/dj-library-tools/commit/0975b8048aaeb27bca4eba2b458312ac05a0c53b))
* **playlist:** implement smart metadata fallback and M3U8 parity improvements ([801eb45](https://github.com/llttlltt/dj-library-tools/commit/801eb452cd78c124d11ec6f56c524cb9df745228))
* **plex:** implement oauth pin flow and playlist retrieval ([4f80a29](https://github.com/llttlltt/dj-library-tools/commit/4f80a2998deeb95cffdbab9422bfd038c15693d0))
* **query:** add ability to query tracks by playlist name ([49e0282](https://github.com/llttlltt/dj-library-tools/commit/49e028245a0d8f757c008e6f52ffa89986b20966))
* **query:** add MatchesNode for playlist and folder node evaluation ([bb01f23](https://github.com/llttlltt/dj-library-tools/commit/bb01f237dadeb29c444c4e28de6064c33cea8036))
* **query:** add playlistcount and support for multiple playlist filters ([990ac75](https://github.com/llttlltt/dj-library-tools/commit/990ac7516da1f43fe78d74852084f59beed2b421))
* **query:** add shell-friendly numeric aliases (gt, lt, ge, le) and NOT keyword ([2dd3ffe](https://github.com/llttlltt/dj-library-tools/commit/2dd3ffe826d018196767e4c88e93c870174d656b))
* **query:** add support for '-' negation, 'eq', 'ne', 'neq' aliases, and '!=' operator ([6104a5c](https://github.com/llttlltt/dj-library-tools/commit/6104a5c16010008ba0f8693f623968c648b465b5))
* **query:** add support for rating, playcount, dateadded, and grouping fields ([ad5df27](https://github.com/llttlltt/dj-library-tools/commit/ad5df276b2ccad781ef0d6ea87eea361d746b024))
* **query:** exhaustive field mapping and numeric comparison operators ([e73f878](https://github.com/llttlltt/dj-library-tools/commit/e73f878cc706dec917126618264f3b0da0524e64))
* **query:** implement standardized query engine with lexer and evaluator ([fa01d0a](https://github.com/llttlltt/dj-library-tools/commit/fa01d0a93655e2a2b892fb66e688fc4999e379be))
* **scaffold:** initialize standardized monorepo structure ([b344dce](https://github.com/llttlltt/dj-library-tools/commit/b344dce51f130de8ecf2a1b0a0802f7e83f7c90d))
* **sync:** add UpsertPlaylist, AddTracksToPlaylist, RenameNode, MoveNode, RemoveNode ([3f025c0](https://github.com/llttlltt/dj-library-tools/commit/3f025c023b999ced6509b2e05ee5d0e0cfdd8182))
* **sync:** implement plex-to-m3u8 sync target ([058b7fa](https://github.com/llttlltt/dj-library-tools/commit/058b7fa97a3e433b283818bb1e4f8bab0b8ec263))
* **sync:** implement plex-to-rekordbox sync engine and cli ([51eb069](https://github.com/llttlltt/dj-library-tools/commit/51eb06980c3a74b24b33470da9901bdd5091e9f7))
* **sync:** robust XML injection via sync.Engine ([f36a749](https://github.com/llttlltt/dj-library-tools/commit/f36a74935021c25c5bde8959887ef3350809f75c))


### Bug Fixes

* **cli:** correct header alignment in ls output ([267261d](https://github.com/llttlltt/dj-library-tools/commit/267261d7f4a98422ad1547c66fb8a8b0a881936f))
* **cli:** finalize precise alignment for ls table output ([f73ae39](https://github.com/llttlltt/dj-library-tools/commit/f73ae3904fcbd12cdd6332d9977f8c22771c0f3d))
* **cli:** left-align headers in ls output ([6f2e6ac](https://github.com/llttlltt/dj-library-tools/commit/6f2e6acb3da483cf4af3a49afcc80302de98ce4e))
* **cli:** prevent alignment spaces from being underlined in ls headers ([f35465e](https://github.com/llttlltt/dj-library-tools/commit/f35465e5468c583ad49ba86d17009763f9284149))
* **cli:** remove duplicate playlist command registration ([9b3eeb3](https://github.com/llttlltt/dj-library-tools/commit/9b3eeb308aa03d0bcff759bdb2e535a004b63a97))
* **cli:** resolve header misalignment caused by ANSI color codes ([de6a1eb](https://github.com/llttlltt/dj-library-tools/commit/de6a1eb021ad52809d462aea6fa6771bc463ef48))
* **playlist:** avoid cross-device link error by using output dir for temp files ([d6c1001](https://github.com/llttlltt/dj-library-tools/commit/d6c10016dcdd4f067dc00c4b418c682dbe6ccfd6))
* **plex:** restore missing CheckPin and rename redeclared lsCmd ([082dfb7](https://github.com/llttlltt/dj-library-tools/commit/082dfb78faa5fa4d0ee26ff745643375f40a4f1d))
* **query:** correct BPM extraction and improve multi-word field matching ([68a5a79](https://github.com/llttlltt/dj-library-tools/commit/68a5a7909e708c90981f200b8adbfbadc6686c11))
* **query:** exact numeric equality for playlistcount and numeric fields ([b7cffc6](https://github.com/llttlltt/dj-library-tools/commit/b7cffc6c988aea4cb8a39ab61669ae32a3cad317))
* **query:** fix parser bug with multi-word queries and update tests ([ef15a4a](https://github.com/llttlltt/dj-library-tools/commit/ef15a4aad94adfb69debbe84f0f533f925f2bde8))
* **query:** implement robust operator parsing and multi-word token joining ([e183ea1](https://github.com/llttlltt/dj-library-tools/commit/e183ea19bc724c63fc833aba3c9eda2d48be7e60))
* **stat,media:** config XML path fallback and filename sanitization ([2435d58](https://github.com/llttlltt/dj-library-tools/commit/2435d58a57c545809615f0047b5551e7d15bfc2a))
