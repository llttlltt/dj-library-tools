import { Atom } from "@effect-atom/atom";
import { Effect } from "effect";
import { AppService } from "@/services";
import type { ProviderInfo } from "@/types";

const asProviders = (x: unknown) => (x ?? []) as ProviderInfo[];

// --- State ---

export const providersAtom = Atom.make<ProviderInfo[]>([]);
export const providersLoadingAtom = Atom.make(false);

// --- Operations ---

export const loadProviders = Effect.gen(function* () {
	yield* Atom.set(providersLoadingAtom, true);
	const app = yield* AppService;
	return yield* app.listProviders().pipe(
		Effect.flatMap((data) => Atom.set(providersAtom, asProviders(data))),
		Effect.catchAll((e) => Effect.logError("Failed to load providers", e)),
		Effect.andThen(() => Atom.set(providersLoadingAtom, false)),
	);
});
