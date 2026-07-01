import { Atom } from "@effect-atom/atom";
import { Duration, Effect } from "effect";
import { AppService } from "@/services";
import type { config, update } from "../../wailsjs/go/models";

// --- State ---

export const versionAtom = Atom.make("v0.0.0").pipe(Atom.keepAlive);
export const updateConfigAtom = Atom.make<config.UpdateConfig | null>(
	null,
).pipe(Atom.keepAlive);

export const lastCheckedAtAtom = Atom.make((get) => {
	const cfg = get(updateConfigAtom);
	return cfg?.last_check_at
		? new Date(cfg.last_check_at).toLocaleString()
		: null;
});

export const permissionStatusAtom = Atom.make("Checking...").pipe(
	Atom.keepAlive,
);

/**
 * updateInfoAtom stores the persistent result of the last update check.
 * This is NOT reset after the button cooldown.
 */
export const updateInfoAtom = Atom.make<update.UpdateInfo | null>(null).pipe(
	Atom.keepAlive,
);

/**
 * isCheckingUpdatesAtom tracks the active Wails call progress.
 */
export const isCheckingUpdatesAtom = Atom.make(false).pipe(Atom.keepAlive);

/**
 * showCheckSuccessAtom is a transient flag for the "Latest" button state.
 * It resets after a timeout, but the updateInfo remains.
 */
export const showCheckSuccessAtom = Atom.make(false).pipe(Atom.keepAlive);

// --- Operations ---

export const loadSystemInfo = Effect.gen(function* () {
	const app = yield* AppService;
	const [version, config, status] = yield* Effect.all([
		app.getVersion(),
		app.getUpdateConfig(),
		app.getPermissionStatus(),
	]);

	yield* Atom.set(versionAtom, version as string);
	yield* Atom.set(updateConfigAtom, config as config.UpdateConfig);
	yield* Atom.set(permissionStatusAtom, status as string);
});

export const checkForUpdates = Effect.gen(function* () {
	yield* Atom.set(isCheckingUpdatesAtom, true);
	yield* Atom.set(showCheckSuccessAtom, false);
	const app = yield* AppService;
	try {
		const info = yield* app.checkForUpdate(true);
		yield* Atom.set(updateInfoAtom, info as update.UpdateInfo);

		// Always refresh config to get the new last_checked timestamp from Go
		const newConfig = yield* app.getUpdateConfig();
		yield* Atom.set(updateConfigAtom, newConfig as config.UpdateConfig);

		// End the checking state before starting the success state
		yield* Atom.set(isCheckingUpdatesAtom, false);

		// If no update is available, we show the "Latest" checkmark transiently
		if (!info.available) {
			yield* Atom.set(showCheckSuccessAtom, true);
			yield* Effect.delay(Duration.seconds(3))(Effect.void);
			yield* Atom.set(showCheckSuccessAtom, false);
		}
	} finally {
		// Safety cleanup
		yield* Atom.set(isCheckingUpdatesAtom, false);
	}
});

export const setUpdateInterval = (hours: number) =>
	Effect.gen(function* () {
		const app = yield* AppService;
		yield* app.setUpdateInterval(hours);
		yield* Atom.update(updateConfigAtom, (prev) =>
			prev ? { ...prev, check_interval_hour: hours } : null,
		);
	});

export const fixPermissions = Effect.gen(function* () {
	const app = yield* AppService;
	yield* app.fixPermissions();
	const status = yield* app.getPermissionStatus();
	yield* Atom.set(permissionStatusAtom, status as string);
});
