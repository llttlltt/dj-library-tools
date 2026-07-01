# Changelog

## [1.2.1](https://github.com/llttlltt/dj-library-tools/compare/v1.2.0...v1.2.1) (2026-07-01)


### Bug Fixes

* **gui:** fix pnpm workspace config and allow build scripts ([3c6d840](https://github.com/llttlltt/dj-library-tools/commit/3c6d84060937792755c97618dd5da27607c774cf))

## [1.2.0](https://github.com/llttlltt/dj-library-tools/compare/v1.1.0...v1.2.0) (2026-07-01)


### Features

* finalize robust Effect-TS boundary and refine update UI ([84fd1db](https://github.com/llttlltt/dj-library-tools/commit/84fd1db53cb67f87f09594fcb80b31ca23e1b83f))
* **gui:** show full file path for M3U resources in query tester ([f567ad5](https://github.com/llttlltt/dj-library-tools/commit/f567ad5ebcd4a785257ecbf7ca6678b3ea7c807f))
* implement centralized reactive state management with Effect Atoms ([c7b8262](https://github.com/llttlltt/dj-library-tools/commit/c7b82628f9accdfb0b7ac35d2cdfeca2b52347e4))
* implement granular multi-group sync and consolidate UI features ([0e60a4c](https://github.com/llttlltt/dj-library-tools/commit/0e60a4c714c8fc640aef1f7d0d8dc8cee75b6a2b))
* implement granular resource capabilities and hardened GUI discovery ([cac9b6e](https://github.com/llttlltt/dj-library-tools/commit/cac9b6e504695132f882593a4f2a1739212815d2))
* implement robust Effect-TS boundary for Wails IPC and domain schemas ([e5e7b94](https://github.com/llttlltt/dj-library-tools/commit/e5e7b94fa33691e740de7c2cd361495a371b1a55))
* implement static provider capabilities and capability-aware GUI/CLI ([4f784c8](https://github.com/llttlltt/dj-library-tools/commit/4f784c80528e7e271e388c5f8a23cf55c61a6985))

## [1.1.0](https://github.com/llttlltt/dj-library-tools/compare/v1.0.0...v1.1.0) (2026-07-01)


### Features

* **gui:** overhaul UI architecture and implement auto-update engine ([448eb40](https://github.com/llttlltt/dj-library-tools/commit/448eb404448fd010095741e668f7e0b9fd861a2a))

## 1.0.0 (2026-07-01)


### Features

* add Source, Workflow, PathMap config layer with UUID storage and roundtrip tests ([6f440b3](https://github.com/llttlltt/dj-library-tools/commit/6f440b3596717003bf39b61f17d090daefd8412d))
* add Svelte frontend — Sources, Workflow editor, and Runner screens ([37e6f92](https://github.com/llttlltt/dj-library-tools/commit/37e6f9233df5d0152cfc224800518d469bfb133a))
* add workflow execution engine with dependency graph, cycle detection, and parallel Steps ([eda718a](https://github.com/llttlltt/dj-library-tools/commit/eda718af1a2c0453ff549d8fa07a5b2d15d1b4b9))
* **cli:** add --sort flag to list command ([6a0bf3f](https://github.com/llttlltt/dj-library-tools/commit/6a0bf3f1f71f4a78fddcee706c70ccabfe795a32))
* **cli:** add ascending/descending sorting and unit tests ([7114799](https://github.com/llttlltt/dj-library-tools/commit/71147995479e455087d1c223d6d3411693e6c432))
* **cli:** add column width capping and truncation to table rendering ([181f835](https://github.com/llttlltt/dj-library-tools/commit/181f835aa38efeccfce085cb986591d91624fb01))
* **cli:** add command aliases (ls, mv, rm, del, stats) ([b40ca0e](https://github.com/llttlltt/dj-library-tools/commit/b40ca0e33de719304d14d4f64e01573b2934f582))
* **cli:** add detailed verbose output for add and remove commands ([2ba647b](https://github.com/llttlltt/dj-library-tools/commit/2ba647b9af67ecf4b8bb1f4ac6485d18e4cb0b4c))
* **cli:** add folder god command (--new, --rename, --move, --remove) ([257ba79](https://github.com/llttlltt/dj-library-tools/commit/257ba79e75f06b4f73cabb20c8f09a084de9794c))
* **cli:** add global --json flag for list and stat commands ([9a8f6f0](https://github.com/llttlltt/dj-library-tools/commit/9a8f6f016eb13721740a41f9c355c5f58696bb37))
* **cli:** add global --missing and --exists flags for physical health filtering ([976d1c2](https://github.com/llttlltt/dj-library-tools/commit/976d1c2cf0f5674b425f6a84ae4a30b31b33cf38))
* **cli:** add Key column and refine ls table alignment ([bc68ce7](https://github.com/llttlltt/dj-library-tools/commit/bc68ce7382b365e43e6f21e8e65906be8ab9fe1b))
* **cli:** add ls command and global --xml flag ([7f5bf90](https://github.com/llttlltt/dj-library-tools/commit/7f5bf90d306d5e1ea726ec23b073c337b4e75350))
* **cli:** add playlist god command (--new, --add, --rename, --move, --remove) ([6a44502](https://github.com/llttlltt/dj-library-tools/commit/6a4450203fb037c65297d5b059df90eb1d58e7aa))
* **cli:** add stat command for library aggregation and analysis ([cfc0c52](https://github.com/llttlltt/dj-library-tools/commit/cfc0c524c49513d813d3b95e30c97c6474d2e2ad))
* **cli:** add visual progress bars for bulk add and remove operations ([6ac5fc3](https://github.com/llttlltt/dj-library-tools/commit/6ac5fc3a6dd48852025a868bcbfb75a3ba22b26a))
* **cli:** colorize ls output and refine table formatting ([b1864ee](https://github.com/llttlltt/dj-library-tools/commit/b1864ee15bbfa3025af25ed07f750934b5fa74d5))
* **cli:** decouple fix command and implement self-healing M3U provider ([820a08b](https://github.com/llttlltt/dj-library-tools/commit/820a08b4c76a09fb05e3e211eca3d65637620e2d))
* **cli:** enhance verbosity and move metadata logic to providers ([104e430](https://github.com/llttlltt/dj-library-tools/commit/104e430e505202116754600992e22508081cafc0))
* **cli:** expand verbose output to all core commands ([1343af6](https://github.com/llttlltt/dj-library-tools/commit/1343af617f9069e7c9f99b798d94809e27d1cabd))
* **cli:** finalize 100% architectural agnosticism by removing last provider traces ([1e93482](https://github.com/llttlltt/dj-library-tools/commit/1e934825644bf0583f1069a85bc0e89263793acd))
* **cli:** finalize agnosticism and standardize file flags ([3b3c820](https://github.com/llttlltt/dj-library-tools/commit/3b3c820842094314a2a567585ab9f0447f179662))
* **cli:** finalize global dry-run support for all modifying commands ([1a95009](https://github.com/llttlltt/dj-library-tools/commit/1a9500940219d5834536b26ec4906eca6562728c))
* **cli:** finalize verb-centric architecture and cleanup legacy commands ([574852b](https://github.com/llttlltt/dj-library-tools/commit/574852bebf1dbf998d5cc5e9f7d1399cf79c2e84))
* **cli:** harden all verbs with interface-based capability checks ([fb33d65](https://github.com/llttlltt/dj-library-tools/commit/fb33d65dabce610f4d57cb7fc4879ed22d7bf662))
* **cli:** implement --dry-run for sync and add location tests ([9da74a5](https://github.com/llttlltt/dj-library-tools/commit/9da74a597f36c38031ac3d1ad19ce44bc22fdcf5))
* **cli:** implement add and remove verbs ([675e68b](https://github.com/llttlltt/dj-library-tools/commit/675e68bde4200972601e8aa7af8cd06911f906c8))
* **cli:** implement compact table rendering for fix feedback ([e57f881](https://github.com/llttlltt/dj-library-tools/commit/e57f8814e20b34957fd4489be90681e2e169b329))
* **cli:** implement create and move verbs ([71d484a](https://github.com/llttlltt/dj-library-tools/commit/71d484abd8463656e8ed28eae9415e0612e12326))
* **cli:** implement dynamic terminal-aware table scaling ([95b3845](https://github.com/llttlltt/dj-library-tools/commit/95b384532e98797c9d026d8f88b73b00caf35336))
* **cli:** implement location-based query syntax provider/resource:query ([6941354](https://github.com/llttlltt/dj-library-tools/commit/6941354457683fb52bd316d37e391b6a4035949e))
* **cli:** implement metadata reconciliation flags and dynamic matching in sync ([35fde34](https://github.com/llttlltt/dj-library-tools/commit/35fde340888464201601d5ea6c4e1d35c414fc1f))
* **cli:** implement modify command with bulk metadata updates and file relocation ([7265ba6](https://github.com/llttlltt/dj-library-tools/commit/7265ba63ff543a58e6f96db03a3bcc501ea5ec1f))
* **cli:** implement multi-bar visual progress for parallel syncing ([c030746](https://github.com/llttlltt/dj-library-tools/commit/c03074675931bde2195b9ac75c4e1a6503bdcc58))
* **cli:** implement rename and delete verbs ([0fe04bb](https://github.com/llttlltt/dj-library-tools/commit/0fe04bb35400b45b6a74e4f5664eb829d4f1c4ea))
* **cli:** implement unified table rendering for list results ([cca6ec7](https://github.com/llttlltt/dj-library-tools/commit/cca6ec7f9092161033ca49540a26fe4c098ce9eb))
* **cli:** implement user-friendly error handling for provider sentinels ([89a0718](https://github.com/llttlltt/dj-library-tools/commit/89a07188499540008636febb035c47cc0ddeb39c))
* **cli:** move dry-run to global persistent flag ([c859eeb](https://github.com/llttlltt/dj-library-tools/commit/c859eeb0d9607dddc10f1dbde3a89e36dadf68b6))
* **cli:** move verbose to global persistent flag ([f8de36a](https://github.com/llttlltt/dj-library-tools/commit/f8de36a6d05c5520ef4738b62b27307f31bbba6c))
* **cli:** phase 2 — wire --missing/--exists flags, remove --to-file ghost ([9275f7d](https://github.com/llttlltt/dj-library-tools/commit/9275f7d15de3409ffd8455125dc4e9fdb4314fcf))
* **cli:** phase 3 — mk --populate, count returns, dry-run previews ([9c1be0b](https://github.com/llttlltt/dj-library-tools/commit/9c1be0b7bcf1d5b86075fce6589c162426cfac31))
* **cli:** refactor metadata command into update verb ([4a1610a](https://github.com/llttlltt/dj-library-tools/commit/4a1610a135bcca5e3066798dae72dbab46fc5392))
* **cli:** refactor sync command to verb-centric architecture ([ed3e4a2](https://github.com/llttlltt/dj-library-tools/commit/ed3e4a2f4051beedfbc791a36dee3906b97f15e7))
* **cli:** remove update command from root command and delete orphan file ([686a338](https://github.com/llttlltt/dj-library-tools/commit/686a338492a29962702e3c5844cb05cd7071f637))
* **cli:** rename create verb to make with mk alias ([d5e6f92](https://github.com/llttlltt/dj-library-tools/commit/d5e6f92126677fe6a3b49bd4892347039703acc4))
* **cli:** replace --xml flag with standardized --file flag ([4f9c2a4](https://github.com/llttlltt/dj-library-tools/commit/4f9c2a4b9e73d7b81c8dc354a5b900d37b4e18ea))
* **cli:** restore legacy metadata move command and integrate into djlt ([357e59a](https://github.com/llttlltt/dj-library-tools/commit/357e59a0786190e7dfd35571bd9a7b7ff58a72d1))
* **cli:** use containment policies for agnostic command validation ([44dcc60](https://github.com/llttlltt/dj-library-tools/commit/44dcc6084822212d53b7f94a580b96f7d5916f06))
* **config:** add --remove-map flag to plex config command ([9c03d0e](https://github.com/llttlltt/dj-library-tools/commit/9c03d0e15097414949dcbd1401469e93015ab3dc))
* **config:** implement persistent app config and automatic token saving ([d689edf](https://github.com/llttlltt/dj-library-tools/commit/d689edfbf1fcd4b996df9ba97c3d0871d886aab3))
* **core:** harden provider interfaces and introduce neutral metadata models ([33925b5](https://github.com/llttlltt/dj-library-tools/commit/33925b5b9413e2622cec1a98ef2b8c6389fd894e))
* **engine:** add LsPlaylists and LsFolders with NodeResult ([21210bc](https://github.com/llttlltt/dj-library-tools/commit/21210bc92c77c45115562fecf10b0df003243f27))
* **engine:** add mock engine test demonstrating Library boundary ([8bf1395](https://github.com/llttlltt/dj-library-tools/commit/8bf1395c035bf4958986ae9316deabfee9622157))
* **engine:** implement core primitives ls, stat, and modify ([15d0aae](https://github.com/llttlltt/dj-library-tools/commit/15d0aae4e91ba0355344488479ef00f147182baa))
* **engine:** introduce Library interface and decouple from Rekordbox XML ([d8d740c](https://github.com/llttlltt/dj-library-tools/commit/d8d740c93c3fba0a09fce4995e6dfb838257861e))
* **gui:** add OpenFileDialog, Plex PIN auth, and UpdateSource bindings ([3ef7032](https://github.com/llttlltt/dj-library-tools/commit/3ef70325837ff87f5efa7964675ac3d383a5ebb6))
* **gui:** migrate frontend to React + TypeScript + Tailwind + shadcn/ui; add StepDiff.Current field ([05747dd](https://github.com/llttlltt/dj-library-tools/commit/05747dd59198b16a7ba2359fec3686e486491127))
* **gui:** Query Tester — validate queries against any Source before using them in a Step ([8015bd0](https://github.com/llttlltt/dj-library-tools/commit/8015bd0af7af328672c72249589bdd7edd1b5074))
* **gui:** source editing, file picker, and Plex PIN auth in Sources screen ([1d6af6f](https://github.com/llttlltt/dj-library-tools/commit/1d6af6f8f135fa09bb9378f459afd5489209d4eb))
* **gui:** unify Workflow editor and Runner into single detail screen with Edit/Preview/Apply modes ([01c68f4](https://github.com/llttlltt/dj-library-tools/commit/01c68f4e81a7f0e411952df7644d3f61fc40ac03))
* **gui:** workflow list actions, fix track display, remove auto-preview, add Biome ([20bf54b](https://github.com/llttlltt/dj-library-tools/commit/20bf54b745bb8a20f91b119e8d068b493b892ea4))
* implement := operator for strict exact equality in queries ([bd310fa](https://github.com/llttlltt/dj-library-tools/commit/bd310fa26c439dac7ad51821c9a871dfc5ba116b))
* implement active capability gating for providers ([070c2fa](https://github.com/llttlltt/dj-library-tools/commit/070c2faccb3bdd005f4af3e816dacbab0cdb95d1))
* implement context cancellation in long-running operations ([06a62e0](https://github.com/llttlltt/dj-library-tools/commit/06a62e054d8f0f547e38cb2ab7447fa5a5a2aa72))
* implement inference-driven sort field validation ([255e405](https://github.com/llttlltt/dj-library-tools/commit/255e4055a1ac19ececdd947abdba13aaf3e85b34))
* implement multi-purpose 'fix' command and provider-agnostic repair system ([f90601d](https://github.com/llttlltt/dj-library-tools/commit/f90601db9bb9e0b4d72296aca60a0e15c0973229))
* implement rich progress feedback and CLI interrupt handling ([36e2c76](https://github.com/llttlltt/dj-library-tools/commit/36e2c76b09d923e015493ff40e9bf74780cf5454))
* implement universal dynamic table rendering for list command ([2fe6e85](https://github.com/llttlltt/dj-library-tools/commit/2fe6e85a2f374cd2787867f8111b88ce6f00c0b6))
* init Wails scaffold, add App bindings for Sources, Workflows, Preview, and Run ([4ac976e](https://github.com/llttlltt/dj-library-tools/commit/4ac976e4bc855244215ab31e37b16a239a6bd1cc))
* **m3u:** implement path normalization in fix command ([49177f3](https://github.com/llttlltt/dj-library-tools/commit/49177f384e5a234e075bfdcecc86d7dae53cc6bf))
* **media:** implement smart skip and robust error handling in transcoder ([f7450f4](https://github.com/llttlltt/dj-library-tools/commit/f7450f4dd3522ec112772cbb12324cb60b4e6c60))
* **media:** implement transcode config, ffmpeg wrapper, and path formatting ([4822468](https://github.com/llttlltt/dj-library-tools/commit/4822468d0bcd14d5926d96509c86fd10377bbf5e))
* **models:** define queryable collections and property metadata ([ac838dd](https://github.com/llttlltt/dj-library-tools/commit/ac838dd6fe9d62f4886d99ed042f672ab3217ae9))
* **models:** register 'missing' as a synthesised TrackField ([8a064b0](https://github.com/llttlltt/dj-library-tools/commit/8a064b0fb7d2d3f5775367afcda4f328b27c6d5a))
* **models:** unify track and node models under universal Resource interface ([eb11ab4](https://github.com/llttlltt/dj-library-tools/commit/eb11ab4e1dfaa1626b2c61a8ae645dc2f7fc3c22))
* **playlist:** add --output flag to fix command ([e68abd7](https://github.com/llttlltt/dj-library-tools/commit/e68abd75b07037061dc64530624a65e7aaf0a37e))
* **playlist:** add --remove and --sync track operations, rename --remove to --delete ([cb13027](https://github.com/llttlltt/dj-library-tools/commit/cb13027a18b2ebe4bd35e86992ce71d3f60b49a8))
* **playlist:** add batch processing and dry-run support ([bfb6b15](https://github.com/llttlltt/dj-library-tools/commit/bfb6b156078263ab955c2fd10b120efb6c4df3a6))
* **playlist:** add interactive removal prompt and clean up summary output ([e56f83a](https://github.com/llttlltt/dj-library-tools/commit/e56f83ab9e2bc8313e16f1fc51e829db0660cf37))
* **playlist:** add missing file tracking and reporting to fix command ([17e6671](https://github.com/llttlltt/dj-library-tools/commit/17e66716ae738e464d6ece5db943da68b6dff529))
* **playlist:** add parity test and refine fix logic for legacy script compatibility ([9c626cf](https://github.com/llttlltt/dj-library-tools/commit/9c626cfefb88f3fe6ba0065f2d1d9715393d4aca))
* **playlist:** add playlist command group and fix subcommand ([fd63506](https://github.com/llttlltt/dj-library-tools/commit/fd63506aa5afbdf3864755ac08954a0f36189d19))
* **playlist:** add priority extension resolution and automatic track pruning ([e15a155](https://github.com/llttlltt/dj-library-tools/commit/e15a1555f72db48f91366e0ab604181e0d91678c))
* **playlist:** add progress heartbeats and verbose logging to fix command ([f205186](https://github.com/llttlltt/dj-library-tools/commit/f2051861c985007e3d4a5bd9afb3cae3acaa2cca))
* **playlist:** finalize feature-playlist-hygiene and update documentation ([3f5b383](https://github.com/llttlltt/dj-library-tools/commit/3f5b383cb2d98c20253eb5f5c313ee5ede1dcf13))
* **playlist:** implement native playlist and metadata logic ([2c9708f](https://github.com/llttlltt/dj-library-tools/commit/2c9708f7a8b2ae83ea734529cb0d93f1c6026e08))
* **playlist:** implement smart metadata fallback and M3U8 parity improvements ([30bb01e](https://github.com/llttlltt/dj-library-tools/commit/30bb01ea767c3b15d2444d624791b533bfa45e54))
* **plex:** implement oauth pin flow and playlist retrieval ([cd0ebc3](https://github.com/llttlltt/dj-library-tools/commit/cd0ebc3d644d193cc5a797d9dce5e8d1c516b936))
* **provider/m3u:** implement strict native fields and update cli rendering ([b0161cf](https://github.com/llttlltt/dj-library-tools/commit/b0161cfa578eb6be408577b64fa8ab09ed2ccfc6))
* **provider/plex:** add metadata extraction and track-level filtering ([afa2b69](https://github.com/llttlltt/dj-library-tools/commit/afa2b6912cd9339add3cda0675f6bca8ed72ed25))
* **provider/plex:** aggregate tracks from all matching playlists ([745dd1b](https://github.com/llttlltt/dj-library-tools/commit/745dd1b6d06be951712d9d61a51cfdefdfcd2233))
* **provider/plex:** enforce field-based selection and add name resolution ([004efd2](https://github.com/llttlltt/dj-library-tools/commit/004efd24c176b9c7d1088a3300b1e682a4479fd9))
* **provider/plex:** implement operator-aware playlist resolution and regex filtering ([11b0d4b](https://github.com/llttlltt/dj-library-tools/commit/11b0d4b5b33da54c2afab11a6e40cfcf0e664aa7))
* **provider/plex:** support global track search when no playlist is specified ([a83d005](https://github.com/llttlltt/dj-library-tools/commit/a83d0057264a7f08e3be17883fe0f84c634b835d))
* **provider/rb:** implement full metadata compliance and enhanced verbose reporting ([6438b87](https://github.com/llttlltt/dj-library-tools/commit/6438b87ec35cd72b3c78fd57d40edf3dddcd07a7))
* **provider:** add M3U and M3U8 playlist support ([79bf633](https://github.com/llttlltt/dj-library-tools/commit/79bf633a65ae8057a2c21602ee4d2d6430e1c1d9))
* **provider:** implement 0-255 unified rating scale for plex and rekordbox ([ea1ae45](https://github.com/llttlltt/dj-library-tools/commit/ea1ae45c27d78bdb207ab9e376976cdf4e9de5f9))
* **provider:** implement capability-based verbs and universal orchestrator ([1b7db92](https://github.com/llttlltt/dj-library-tools/commit/1b7db92fb8c5520e4cce10ba241e68b051ab8fec))
* **provider:** implement GetResources and unify agnostic selection logic ([e8ce5b6](https://github.com/llttlltt/dj-library-tools/commit/e8ce5b643d0bea7ab103e14150a978f1d6bb9ca0))
* **provider:** implement SupportedResources and IdentifyGroup for full agnosticism ([212ed20](https://github.com/llttlltt/dj-library-tools/commit/212ed2080fe380feb7fea7894e23fb767964085d))
* **provider:** introduce Provider interface and refactor list command for RB ([d87cfc4](https://github.com/llttlltt/dj-library-tools/commit/d87cfc4009e13cbefbecc48c6cd4037b22b78f79))
* **provider:** move sorting business logic into providers ([2c67e5e](https://github.com/llttlltt/dj-library-tools/commit/2c67e5e725f20863fc3177962cf68b65dfaf7c6d))
* **provider:** split Provider interface into Readable, Searchable, and Writable ([c214b9f](https://github.com/llttlltt/dj-library-tools/commit/c214b9f92a57db130a168e0baac8c6aaec8ed2b9))
* **provider:** unify discovery in GetResources and harden writable interface ([f088e5c](https://github.com/llttlltt/dj-library-tools/commit/f088e5c3ac3bf037628d30654754de6bd82de3ba))
* **query:** add ability to query tracks by playlist name ([2fd7397](https://github.com/llttlltt/dj-library-tools/commit/2fd739734e19fbdeed1a5ccf42ef50dea7371b4a))
* **query:** add MatchesNode for playlist and folder node evaluation ([0fbd33c](https://github.com/llttlltt/dj-library-tools/commit/0fbd33c96eb8d8f8a16367be819b1f0df6abf4da))
* **query:** add playlistcount and support for multiple playlist filters ([9f87941](https://github.com/llttlltt/dj-library-tools/commit/9f87941eb2d63448a860942f6e962465de1f50e6))
* **query:** add query validation and helpful error for bare values ([d8ceafc](https://github.com/llttlltt/dj-library-tools/commit/d8ceafcbb0dbda773024c2ca47a4643267e6bc11))
* **query:** add shell-friendly numeric aliases (gt, lt, ge, le) and NOT keyword ([36e1f92](https://github.com/llttlltt/dj-library-tools/commit/36e1f922cd0bd2485fab54f2f390d6eb5665c290))
* **query:** add stats and path logic files ([8d4c920](https://github.com/llttlltt/dj-library-tools/commit/8d4c920474b06d01ee1ddbd3bfa362455da41748))
* **query:** add support for '-' negation, 'eq', 'ne', 'neq' aliases, and '!=' operator ([33742c4](https://github.com/llttlltt/dj-library-tools/commit/33742c47bc44b668e48489537d4b2486ea0df925))
* **query:** add support for rating, playcount, dateadded, and grouping fields ([9ef3869](https://github.com/llttlltt/dj-library-tools/commit/9ef386955726716b45aa0009d85d0c55ffeb6323))
* **query:** align 'name' and 'title' fields and improve playlist error handling ([42e6702](https://github.com/llttlltt/dj-library-tools/commit/42e67020d29029cc814ab9acd1b4b78acda4c002))
* **query:** delegate provider-specific matching to CustomMatcher interface ([e54b79d](https://github.com/llttlltt/dj-library-tools/commit/e54b79d9f256c1ae8f8c092788aaf18bf49debc6))
* **query:** exhaustive field mapping and numeric comparison operators ([4f80cd3](https://github.com/llttlltt/dj-library-tools/commit/4f80cd33a50c42ff0eed9a544bd56e00f318a708))
* **query:** implement full suite of stats including min, max, avg, jitter, and stability ([2eb9892](https://github.com/llttlltt/dj-library-tools/commit/2eb989277c1e4aecaf96efc97896569558780023))
* **query:** implement path-based resolution and stats engine ([6a80a34](https://github.com/llttlltt/dj-library-tools/commit/6a80a34b43ee1d5327a7cd1ea6c80d35195ae41f))
* **query:** implement standardized query engine with lexer and evaluator ([fe8a07f](https://github.com/llttlltt/dj-library-tools/commit/fe8a07f83bd969f1da8e7636286238ad75fc1dda))
* **query:** implement strict field validation and typo detection ([188908f](https://github.com/llttlltt/dj-library-tools/commit/188908ffbe60d026bf2b26e5c6ba8cd6a96d0c20))
* **query:** improve bare value error message with reconstructed query ([6446556](https://github.com/llttlltt/dj-library-tools/commit/6446556fa24c21b7d42803d1fc6ae28d3d2b7109))
* **query:** improve name matching and quoted string support ([3137a4e](https://github.com/llttlltt/dj-library-tools/commit/3137a4e63ed9355a3537fb617e93c7935dfc4eef))
* **query:** normalize redundancy stat to 0-100 percentage scale ([a20a814](https://github.com/llttlltt/dj-library-tools/commit/a20a814302689eac34d3c3ae366de345a1fad7e7))
* **query:** treat playlists as a collection with folder properties ([682f979](https://github.com/llttlltt/dj-library-tools/commit/682f9795df82a1058aaba8e2cc72233c2fe8c681))
* **rb:** consolidate deduplication position display ([85be908](https://github.com/llttlltt/dj-library-tools/commit/85be908ea8b47591fb84cbd25b8e64ba565b3197))
* **rb:** enhance deduplication logs with track positions ([5856575](https://github.com/llttlltt/dj-library-tools/commit/58565757c5a7d1e6e1da33a2307c1409e88aba6d))
* **rb:** finalize deduplication audit feedback ([364377c](https://github.com/llttlltt/dj-library-tools/commit/364377c43101bbdd2f4f7dd31edb84c8424436f4))
* **rb:** implement comprehensive fix logic and consolidate maintenance ([e140c16](https://github.com/llttlltt/dj-library-tools/commit/e140c1662e6f34acfb6fd002d6325743513838c7))
* **rb:** improve feedback for playlist deduplication ([e410c82](https://github.com/llttlltt/dj-library-tools/commit/e410c82ec7693316e1b3c21af34603d641e6022b))
* **rb:** refine deduplication audit table layout ([357c2ea](https://github.com/llttlltt/dj-library-tools/commit/357c2ea2090637344187a2fffb90877ebdb2ce63))
* **rb:** sync with rekordbox xml spec and standardize query engine ([0e50595](https://github.com/llttlltt/dj-library-tools/commit/0e50595022aac794e93e496a4d1efc1d8edcdc41))
* **rb:** sync with rekordbox xml spec and standardize query engine ([72af1a7](https://github.com/llttlltt/dj-library-tools/commit/72af1a7bc41f9ea677c16cfceec28ae3d05d4361))
* **rekordbox:** align struct fields with XML attribute order for cleaner diffs ([3616df5](https://github.com/llttlltt/dj-library-tools/commit/3616df5f218634f973afb370f2f60cd931b2ff3f))
* **rekordbox:** detect and preserve XML formatting ([c5292b2](https://github.com/llttlltt/dj-library-tools/commit/c5292b27a82d93dcb2db8d2a59025ae872e45586))
* **rekordbox:** extract color utilities to separate file ([025ff6a](https://github.com/llttlltt/dj-library-tools/commit/025ff6a2177486ad184dbcf23f70c82e39901c2f))
* **rekordbox:** implement high-fidelity XML emission and comprehensive testing ([03dd2ba](https://github.com/llttlltt/dj-library-tools/commit/03dd2ba1ec08dcb9a731601354142bd8201b39a7))
* **rekordbox:** implement surgical saving and playlist positioning ([a3b710e](https://github.com/llttlltt/dj-library-tools/commit/a3b710e0e28807e8974b71e477cbb074a9c7ab1a))
* **rekordbox:** match idiosyncratic attribute order for root and sub-nodes ([69a86d4](https://github.com/llttlltt/dj-library-tools/commit/69a86d4b15d6b15175e7cdb8219a03258c001f62))
* **rekordbox:** support attribute wrapping and fix omitempty issues ([700f0c6](https://github.com/llttlltt/dj-library-tools/commit/700f0c6dd2f3bcd19b3d4df54ebc836613c144f3))
* **resolver:** add path capability validation and numeric mapping ([82e97d4](https://github.com/llttlltt/dj-library-tools/commit/82e97d4117b1d999d6fe2123d7f78a6e3c51af32))
* restore hierarchical cue matching in Rekordbox core ([31a61b7](https://github.com/llttlltt/dj-library-tools/commit/31a61b70b2a498a7be6eb421d21bd6f2fea9cfce))
* **scaffold:** initialize standardized monorepo structure ([0d96544](https://github.com/llttlltt/dj-library-tools/commit/0d96544740f9268c8b20b23fe40cf259bfdddf5c))
* support track-level queries on playlists/folders (e.g. tracks/title:Oceans) ([9f3cbe3](https://github.com/llttlltt/dj-library-tools/commit/9f3cbe3c4c97b966de935d6a8731c8556414fd41))
* **sync:** add folder creation and standardize node types ([1753df2](https://github.com/llttlltt/dj-library-tools/commit/1753df200d11446ef3967fbd2fd61b9db1abc0cd))
* **sync:** add UpsertPlaylist, AddTracksToPlaylist, RenameNode, MoveNode, RemoveNode ([dfd7e1c](https://github.com/llttlltt/dj-library-tools/commit/dfd7e1c046a927315ab2c1db6b9d62a3842d4099))
* **sync:** implement agnostic Join logic for metadata reconciliation ([dd16aa4](https://github.com/llttlltt/dj-library-tools/commit/dd16aa4abc06dc1ede13c0d9d71255560fe4cf6c))
* **sync:** implement plex-to-m3u8 sync target ([5f50ffd](https://github.com/llttlltt/dj-library-tools/commit/5f50ffd9348d4874c856e74c86d9eeaca2b1322a))
* **sync:** implement plex-to-rekordbox sync engine and cli ([378be00](https://github.com/llttlltt/dj-library-tools/commit/378be0070f0eade9c86259fc9e46315fcf8e5b0c))
* **sync:** robust XML injection via sync.Engine ([4bbdcd9](https://github.com/llttlltt/dj-library-tools/commit/4bbdcd90e9fd41b1fb7bbf67cd5fcbdd571c9e58))
* **sys:** introduce system abstractions for filesystem and command execution ([3b1ebf1](https://github.com/llttlltt/dj-library-tools/commit/3b1ebf162d57f5d1e941d0b6610c75200d9ec040))
* thread context.Context through provider and orchestrator APIs ([0c2a0c8](https://github.com/llttlltt/dj-library-tools/commit/0c2a0c8778889400f611adc8138c24ce538dd4e7))
* **utils:** make location parser purely syntactic and provider-agnostic ([0939300](https://github.com/llttlltt/dj-library-tools/commit/0939300c297154d988c9f7a23501da08016b4dce))


### Bug Fixes

* **ci:** build frontend before running go tests to satisfy embed pattern ([c4db0f2](https://github.com/llttlltt/dj-library-tools/commit/c4db0f2c3ecba2866a953d3971c034049bc1e644))
* **cli:** align headers and columns in shared table utility ([f2eca7c](https://github.com/llttlltt/dj-library-tools/commit/f2eca7c939fb3a3efc6287cc43e6a71103e98a32))
* **cli:** correct header alignment in ls output ([8034c8b](https://github.com/llttlltt/dj-library-tools/commit/8034c8b95d16ede7be5aacf61f919d8ddf393ca8))
* **cli:** enforce rekordbox hierarchy constraints for tracks and nodes ([e3b07fc](https://github.com/llttlltt/dj-library-tools/commit/e3b07fca77073a750a4cb2785973d734b33cbfb6))
* **cli:** finalize precise alignment for ls table output ([0fdecde](https://github.com/llttlltt/dj-library-tools/commit/0fdecdee1986916b8d57df30b166a7f906823594))
* **cli:** left-align headers in ls output ([12c8a82](https://github.com/llttlltt/dj-library-tools/commit/12c8a8232d208631b6f3b141d6252696c2e008f8))
* **cli:** phase 1 — rm query arg, sync success msg, fix --help example ([24b80b3](https://github.com/llttlltt/dj-library-tools/commit/24b80b31f8123917f322b5f7df5f237a2946c8a4))
* **cli:** preserve heading casing in shared table utility ([fbc7f5d](https://github.com/llttlltt/dj-library-tools/commit/fbc7f5dec2124acdadbc2b8a8e33e2a2bef7ed54))
* **cli:** prevent alignment spaces from being underlined in ls headers ([29db32b](https://github.com/llttlltt/dj-library-tools/commit/29db32bd108b013cdbca3730801191e132700ddd))
* **cli:** remove duplicate playlist command registration ([cd8ccb5](https://github.com/llttlltt/dj-library-tools/commit/cd8ccb5db9340c9b2296a084a877d23590461e8f))
* **cli:** resolve header misalignment caused by ANSI color codes ([f9461ba](https://github.com/llttlltt/dj-library-tools/commit/f9461ba0b3ccebc73b6c9e105b13d742cb5089b3))
* **cli:** resolve linting warnings and align CreateGroup signature ([eac5cbb](https://github.com/llttlltt/dj-library-tools/commit/eac5cbbb9a724bccc3b2e15952c4b6181db4c9fd))
* **cli:** silence redundant error printing in root command ([6d707b4](https://github.com/llttlltt/dj-library-tools/commit/6d707b43c8d33e203c4935f40723d296ac69ad21))
* **cli:** support file-based providers in mv and rename ([f767fcb](https://github.com/llttlltt/dj-library-tools/commit/f767fcb31ee1d0323c036ed54ddae05772ca3377))
* **engine,provider:** rb/folders now correctly queries folder nodes ([5d4efcb](https://github.com/llttlltt/dj-library-tools/commit/5d4efcb94198f72a9cfdea07bcd351c0cb1f7a46))
* ensure providers are registered and config paths are respected ([e94fc3f](https://github.com/llttlltt/dj-library-tools/commit/e94fc3fdef9f5c50ff81dbec6b45e9f342c323c6))
* **gui:** clicking sidebar tab resets view to list mode; fix unused prop ([0317520](https://github.com/llttlltt/dj-library-tools/commit/031752098dfe534e3d21b7be64266d40c7b40159))
* **gui:** close QueryTester sheet when "Use this query" is clicked ([863f411](https://github.com/llttlltt/dj-library-tools/commit/863f41173d6b8d648766f2b4984101f262e7bebb))
* **gui:** fix Plex auth by resolving UUID in location strings; filter resources in QueryTester ([3b18e46](https://github.com/llttlltt/dj-library-tools/commit/3b18e4669b122373030e63c6c3ac464929bb3679))
* **gui:** make activeResult reactive so step cards render after Preview ([9bdf601](https://github.com/llttlltt/dj-library-tools/commit/9bdf601bf1f1d377bc10dd830f6bd625081bbe1c))
* **gui:** PreviewQuery supports groups, Query Tester as dedicated view, target query full width ([1684743](https://github.com/llttlltt/dj-library-tools/commit/16847433c5f3a98832aac1fadd3781db99edd904))
* **gui:** remove dev build tag split — single main.go covers both dev and production ([0ddf433](https://github.com/llttlltt/dj-library-tools/commit/0ddf4334c720b1615e268a3633e3dc034bcca62b))
* **gui:** replace window.confirm with Dialog, rename Apply to Run, add Preview and Delete to detail toolbar ([3410ed0](https://github.com/llttlltt/dj-library-tools/commit/3410ed0acc4765b2600645e5222e5046fd2a7563))
* **gui:** Run button, save clears results, min window size, nav height, full-width cards ([11dead6](https://github.com/llttlltt/dj-library-tools/commit/11dead65b7479a60129b24cee03adfaa17ade499))
* **gui:** sort workflows alphabetically in the list view ([9a45409](https://github.com/llttlltt/dj-library-tools/commit/9a4540989144f483137c4efbdbbcb950ff3e5d02))
* **gui:** source form state reset, file picker dir, Plex URL, SourceCard UX ([9340dc9](https://github.com/llttlltt/dj-library-tools/commit/9340dc9c81dec64135140fe70922ee3bf07fdc1e))
* **gui:** sticky table header and disable global app overscroll ([1d4e71d](https://github.com/llttlltt/dj-library-tools/commit/1d4e71df9db61adc3e65d000e9e781aa0934b144))
* **lint:** use tagged switch for kind in engine mock tests ([f7e42d6](https://github.com/llttlltt/dj-library-tools/commit/f7e42d63ec16aa0cfc7a6979d6c9e7d9e632a508))
* **m3u:** ensure 'fix' command respects dry-run mode for path normalization ([ca81c6e](https://github.com/llttlltt/dj-library-tools/commit/ca81c6edf01d989e194f95bbb124b0c23faf52b9))
* **main:** remove unused fmt import ([300f448](https://github.com/llttlltt/dj-library-tools/commit/300f4483fb524a11308ee15978920d2bc57dc827))
* **playlist:** avoid cross-device link error by using output dir for temp files ([2b30e3a](https://github.com/llttlltt/dj-library-tools/commit/2b30e3af1a1e2fc4dd8d6213c709d1e2ece6d35f))
* **plex:** implement playlist memberships and remove placeholder token ([54e1084](https://github.com/llttlltt/dj-library-tools/commit/54e1084b068f69f6ab0ad9bce7212e63795a4c24))
* **plex:** restore missing CheckPin and rename redeclared lsCmd ([7fa797c](https://github.com/llttlltt/dj-library-tools/commit/7fa797c5ebc44b98b87d2fc1dc1f09aaefb1fbcc))
* **provider/plex:** allow fuzzy name matching for playlist resolution ([d45fa14](https://github.com/llttlltt/dj-library-tools/commit/d45fa140b62cd32e0e53f50f4576f222d7a211e4))
* **provider/plex:** enable field validation for track and playlist queries ([c013d07](https://github.com/llttlltt/dj-library-tools/commit/c013d0791d33e6acf675f9a0d530f17a5d33edb2))
* **provider/rb:** correctly create playlists vs folders in CreateNode ([0b6f847](https://github.com/llttlltt/dj-library-tools/commit/0b6f8477aaf6ff2759a12d1dfecf6b77d3d289b6))
* **query,docs:** wire count alias; correct sync/append docs ([6462101](https://github.com/llttlltt/dj-library-tools/commit/64621013dd4e02e64e6b33a64ee3ee26487f9e71))
* **query:** correct BPM extraction and improve multi-word field matching ([762a8be](https://github.com/llttlltt/dj-library-tools/commit/762a8be712b4641b37b54d823e4c8045c4c5d677))
* **query:** exact numeric equality for playlistcount and numeric fields ([86b9cca](https://github.com/llttlltt/dj-library-tools/commit/86b9ccace1ac7169c317e533fafb14835f4b9ed8))
* **query:** fix parser bug with multi-word queries and update tests ([e172a58](https://github.com/llttlltt/dj-library-tools/commit/e172a5891b427a1cc64091364376e962319fe2c7))
* **query:** implement missing remixer and mix fields in evaluator ([47f9c3e](https://github.com/llttlltt/dj-library-tools/commit/47f9c3e4043232ba20f87d610655659b6d4464fa))
* **query:** implement robust operator parsing and multi-word token joining ([b5afd5c](https://github.com/llttlltt/dj-library-tools/commit/b5afd5c6032066f1d47aba610269ea96c5d908ab))
* **query:** remove redundant static count fields and implement path-aware validation ([a461f2d](https://github.com/llttlltt/dj-library-tools/commit/a461f2d10f681b4918bb5711eeee9901d4cdd3e8))
* **query:** support implicit AND for field comparisons ([9e71cf2](https://github.com/llttlltt/dj-library-tools/commit/9e71cf2bb91f1816fb5cac152b1479b3c66f281c))
* **rb:** improve query robustness and folder-aware maintenance ([beca78b](https://github.com/llttlltt/dj-library-tools/commit/beca78b51a1ac7198195498992586a43fc9f0f5d))
* **rekordbox:** folder nodes now map Count to models.Node.Entries ([7bd3efe](https://github.com/llttlltt/dj-library-tools/commit/7bd3efe9bc18291444f72ef323c716239017d320))
* **rekordbox:** repair surgical XML write and optional node attributes ([e6129f0](https://github.com/llttlltt/dj-library-tools/commit/e6129f0c47a9f1f1e02fe21460af436b50553ab5))
* resolve architecture migration blockers and documentation issues ([e45649c](https://github.com/llttlltt/dj-library-tools/commit/e45649cbdc0c8028d8e83f4d2a4a38b6112cb3d4))
* resolve doc link generation and finalize architecture migration ([9d07019](https://github.com/llttlltt/dj-library-tools/commit/9d07019d166e80869e6bf7cb1f298a227a800bb9))
* resolve table alignment issues caused by ANSI color codes ([ee3c680](https://github.com/llttlltt/dj-library-tools/commit/ee3c68074827314e960b72be1e1aee0cc1e01e70))
* **resolver:** inject kind filter so rb/playlists and rb/folders return distinct resource types ([cf95767](https://github.com/llttlltt/dj-library-tools/commit/cf957677d0212eed61d296f4690eacffa37d9a0f))
* restore golden baseline and clean up provider output routing ([4db79fe](https://github.com/llttlltt/dj-library-tools/commit/4db79fecbd6d9454470f859a206efb49723cb622))
* route all output through Feedback and add strict architecture guards ([8fc8956](https://github.com/llttlltt/dj-library-tools/commit/8fc89566697cf24f85f4b8f112f119576dd87bd9))
* **stat,media:** config XML path fallback and filename sanitization ([3f7f9a7](https://github.com/llttlltt/dj-library-tools/commit/3f7f9a7712e86695697d5d5d26421d852011c83d))
* **tests:** update test data to match new Rekordbox string-based numeric fields ([27cfbd9](https://github.com/llttlltt/dj-library-tools/commit/27cfbd919652f05618a34b07e7c8d684fd511988))
* **utils:** support colon separator in m3u locations ([aab2a07](https://github.com/llttlltt/dj-library-tools/commit/aab2a074dc56852587f43508fd9869d0a20cc7fd))
* **workflow:** emit sync diff as preview messages instead of calling orch.Sync in dry-run mode ([62e5faa](https://github.com/llttlltt/dj-library-tools/commit/62e5faa1a54bb573b67324c45e50b087396d0094))


### Performance Improvements

* **gui:** use virtualization in Query Tester for lag-free large result sets ([e7e4fb9](https://github.com/llttlltt/dj-library-tools/commit/e7e4fb90f3d135465cef43e97d3c5eede05a681f))

## [1.11.0](https://github.com/llttlltt/dj-library-tools/compare/v1.10.0...v1.11.0) (2026-06-29)


### Features

* **rb:** implement comprehensive fix logic and consolidate maintenance ([292fa40](https://github.com/llttlltt/dj-library-tools/commit/292fa4093c4932ff946ef64a55cd1a01143c1520))

## [1.10.0](https://github.com/llttlltt/dj-library-tools/compare/v1.9.0...v1.10.0) (2026-06-29)


### Features

* **cli:** implement compact table rendering for fix feedback ([63b9375](https://github.com/llttlltt/dj-library-tools/commit/63b937552a4f125b1df0d2effb4a8df1cad80b34))
* implement multi-purpose 'fix' command and provider-agnostic repair system ([a865af4](https://github.com/llttlltt/dj-library-tools/commit/a865af4a62b6f4be7124b80bcc6a758ef7c4a3ae))
* **m3u:** implement path normalization in fix command ([61ad3a8](https://github.com/llttlltt/dj-library-tools/commit/61ad3a813c5562e9528972be1996a803ea62a18b))
* **models:** define queryable collections and property metadata ([84afb6a](https://github.com/llttlltt/dj-library-tools/commit/84afb6a03a5ac29f5a06c3a09a93e52a55c3891d))
* **query:** add stats and path logic files ([249a4c3](https://github.com/llttlltt/dj-library-tools/commit/249a4c3c15baa67762014c7061b8a9bd0ee78d76))
* **query:** implement full suite of stats including min, max, avg, jitter, and stability ([c71c7d3](https://github.com/llttlltt/dj-library-tools/commit/c71c7d35de510ddf70eb7bb8f35efa182b6287bf))
* **query:** implement path-based resolution and stats engine ([d6cadcc](https://github.com/llttlltt/dj-library-tools/commit/d6cadccc58a5c2f126f047881d0ee35c332510e3))
* **query:** normalize redundancy stat to 0-100 percentage scale ([a9772fc](https://github.com/llttlltt/dj-library-tools/commit/a9772fc44bb035289d032bac7f4f3e25c39f3f3a))
* **query:** treat playlists as a collection with folder properties ([53f7554](https://github.com/llttlltt/dj-library-tools/commit/53f75541ff028769254b60b51aaec2312298119b))
* **rb:** consolidate deduplication position display ([6dfea05](https://github.com/llttlltt/dj-library-tools/commit/6dfea0522c30962835b692f293db6ba7acb160d1))
* **rb:** enhance deduplication logs with track positions ([cebc492](https://github.com/llttlltt/dj-library-tools/commit/cebc492e6300da0dbed3229516fcba802289f435))
* **rb:** finalize deduplication audit feedback ([c9a7ecc](https://github.com/llttlltt/dj-library-tools/commit/c9a7eccc5d8526c75bf33eb5cfdcbcb654aa1322))
* **rb:** improve feedback for playlist deduplication ([f7d0679](https://github.com/llttlltt/dj-library-tools/commit/f7d0679972723f62d4cee8d6900396a8d1a6ddb0))
* **rb:** refine deduplication audit table layout ([156975a](https://github.com/llttlltt/dj-library-tools/commit/156975a078daa61c6ac7ff0039d21ccd2b39e2b0))
* **rekordbox:** extract color utilities to separate file ([8186451](https://github.com/llttlltt/dj-library-tools/commit/8186451df051443300edf1a58d89f5b4c444cb19))
* **resolver:** add path capability validation and numeric mapping ([ebe7651](https://github.com/llttlltt/dj-library-tools/commit/ebe7651cc0acb0b23b0cdbc910ed65746dad4a6a))


### Bug Fixes

* **cli:** align headers and columns in shared table utility ([73a9e69](https://github.com/llttlltt/dj-library-tools/commit/73a9e69329af7bbeeea3c10784c6ea3543408ae4))
* **cli:** preserve heading casing in shared table utility ([d7ded35](https://github.com/llttlltt/dj-library-tools/commit/d7ded35bdba826e997b68ceddcc18d640978be51))
* **m3u:** ensure 'fix' command respects dry-run mode for path normalization ([bc04fc4](https://github.com/llttlltt/dj-library-tools/commit/bc04fc4e3fe696a9630811749783bae166281184))
* **query:** remove redundant static count fields and implement path-aware validation ([789ba4c](https://github.com/llttlltt/dj-library-tools/commit/789ba4cdc7b243ccc284748ba60f72e03a4bdf59))
* **rb:** improve query robustness and folder-aware maintenance ([35ce354](https://github.com/llttlltt/dj-library-tools/commit/35ce3545ca946c7a96f309cde661285e3fb98bca))

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
