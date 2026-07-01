import { Atom } from "@effect-atom/atom";
import { Effect } from "effect";
import type { Source } from "@/types";
import {
	CreateSource,
	DeleteSource,
	ListSources,
	UpdateSource,
} from "../../wailsjs/go/gui/App";

const asSources = (x: unknown) => (x ?? []) as Source[];

// --- State ---

export const sourcesAtom = Atom.make<Source[]>([]);
export const sourcesLoadingAtom = Atom.make(false);
export const sourcesErrorAtom = Atom.make<string | null>(null);

// --- Operations ---

export const loadSources = Effect.gen(function* () {
	yield* Atom.set(sourcesLoadingAtom, true);
	yield* Atom.set(sourcesErrorAtom, null);

	try {
		const data = yield* Effect.promise(() => ListSources());
		yield* Atom.set(sourcesAtom, asSources(data));
	} catch (e) {
		yield* Atom.set(sourcesErrorAtom, String(e));
	} finally {
		yield* Atom.set(sourcesLoadingAtom, false);
	}
});

export const addSource = (
	name: string,
	provider: string,
	config: Record<string, string>,
) =>
	Effect.gen(function* () {
		yield* Atom.set(sourcesLoadingAtom, true);
		try {
			yield* Effect.promise(() => CreateSource(name, provider, config));
			yield* loadSources;
		} catch (e) {
			yield* Atom.set(sourcesErrorAtom, String(e));
		} finally {
			yield* Atom.set(sourcesLoadingAtom, false);
		}
	});

export const removeSource = (id: string) =>
	Effect.gen(function* () {
		yield* Atom.set(sourcesLoadingAtom, true);
		try {
			yield* Effect.promise(() => DeleteSource(id));
			yield* loadSources;
		} catch (e) {
			yield* Atom.set(sourcesErrorAtom, String(e));
		} finally {
			yield* Atom.set(sourcesLoadingAtom, false);
		}
	});

export const updateSource = (src: Source) =>
	Effect.gen(function* () {
		yield* Atom.set(sourcesLoadingAtom, true);
		try {
			yield* Effect.promise(() => UpdateSource(src as never));
			yield* loadSources;
		} catch (e) {
			yield* Atom.set(sourcesErrorAtom, String(e));
		} finally {
			yield* Atom.set(sourcesLoadingAtom, false);
		}
	});
