import { Atom } from "@effect-atom/atom";
import { Effect } from "effect";
import type { ProviderInfo } from "@/types";
import { ListProviders } from "../../wailsjs/go/gui/App";

const asProviders = (x: unknown) => (x ?? []) as ProviderInfo[];

// --- State ---

export const providersAtom = Atom.make<ProviderInfo[]>([]);
export const providersLoadingAtom = Atom.make(false);

// --- Operations ---

export const loadProviders = Effect.gen(function* () {
	yield* Atom.set(providersLoadingAtom, true);
	try {
		const data = yield* Effect.promise(() => ListProviders());
		yield* Atom.set(providersAtom, asProviders(data));
	} catch (e) {
		console.error("Failed to load providers:", e);
	} finally {
		yield* Atom.set(providersLoadingAtom, false);
	}
});
