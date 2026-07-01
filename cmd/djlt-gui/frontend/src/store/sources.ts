import { Atom } from "@effect-atom/atom";
import { Effect } from "effect";
import { AppService } from "@/services";
import type { Source } from "@/types";

const asSources = (x: unknown) => (x ?? []) as Source[];

// --- State ---

export const sourcesAtom = Atom.make<Source[]>([]);
export const sourcesLoadingAtom = Atom.make(false);
export const sourcesErrorAtom = Atom.make<string | null>(null);

// --- Operations ---

/**
 * Normalizes any error into a string for the error atom.
 */
const handleError = (e: unknown) =>
	Atom.set(sourcesErrorAtom, e instanceof Error ? e.message : String(e));

export const loadSources = Effect.gen(function* () {
	yield* Atom.set(sourcesLoadingAtom, true);
	yield* Atom.set(sourcesErrorAtom, null);

	const app = yield* AppService;
	return yield* app.listSources().pipe(
		Effect.flatMap((data) => Atom.set(sourcesAtom, asSources(data))),
		Effect.catchAll((e) => handleError(e)),
		Effect.andThen(() => Atom.set(sourcesLoadingAtom, false)),
	);
});

export const addSource = (
	name: string,
	provider: string,
	config: Record<string, string>,
) =>
	Effect.gen(function* () {
		yield* Atom.set(sourcesLoadingAtom, true);
		const app = yield* AppService;
		return yield* app.createSource(name, provider, config).pipe(
			Effect.flatMap(() => loadSources),
			Effect.catchAll((e) => handleError(e)),
			Effect.andThen(() => Atom.set(sourcesLoadingAtom, false)),
		);
	});

export const removeSource = (id: string) =>
	Effect.gen(function* () {
		yield* Atom.set(sourcesLoadingAtom, true);
		const app = yield* AppService;
		return yield* app.deleteSource(id).pipe(
			Effect.flatMap(() => loadSources),
			Effect.catchAll((e) => handleError(e)),
			Effect.andThen(() => Atom.set(sourcesLoadingAtom, false)),
		);
	});

export const updateSource = (src: Source) =>
	Effect.gen(function* () {
		yield* Atom.set(sourcesLoadingAtom, true);
		const app = yield* AppService;
		return yield* app.updateSource(src).pipe(
			Effect.flatMap(() => loadSources),
			Effect.catchAll((e) => handleError(e)),
			Effect.andThen(() => Atom.set(sourcesLoadingAtom, false)),
		);
	});
