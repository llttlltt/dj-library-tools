import { Atom } from "@effect-atom/atom";
import { Effect } from "effect";
import { AppService } from "@/services";
import type { Connection } from "@/types";

const asConnections = (x: unknown) => (x ?? []) as Connection[];

// --- State ---

export const connectionsAtom = Atom.make<Connection[]>([]);
export const connectionsLoadingAtom = Atom.make(false);
export const connectionsErrorAtom = Atom.make<string | null>(null);

// --- Operations ---

/**
 * Normalizes any error into a string for the error atom.
 */
const handleError = (e: unknown) =>
	Atom.set(connectionsErrorAtom, e instanceof Error ? e.message : String(e));

export const loadConnections = Effect.gen(function* () {
	yield* Atom.set(connectionsLoadingAtom, true);
	yield* Atom.set(connectionsErrorAtom, null);

	const app = yield* AppService;
	return yield* app.listConnections().pipe(
		Effect.flatMap((data) => Atom.set(connectionsAtom, asConnections(data))),
		Effect.catchAll((e) => handleError(e)),
		Effect.andThen(() => Atom.set(connectionsLoadingAtom, false)),
	);
});

export const addConnection = (
	name: string,
	provider: string,
	config: Record<string, string>,
) =>
	Effect.gen(function* () {
		yield* Atom.set(connectionsLoadingAtom, true);
		const app = yield* AppService;
		return yield* app.createConnection(name, provider, config).pipe(
			Effect.flatMap(() => loadConnections),
			Effect.catchAll((e) => handleError(e)),
			Effect.andThen(() => Atom.set(connectionsLoadingAtom, false)),
		);
	});

export const removeConnection = (id: string) =>
	Effect.gen(function* () {
		yield* Atom.set(connectionsLoadingAtom, true);
		const app = yield* AppService;
		return yield* app.deleteConnection(id).pipe(
			Effect.flatMap(() => loadConnections),
			Effect.catchAll((e) => handleError(e)),
			Effect.andThen(() => Atom.set(connectionsLoadingAtom, false)),
		);
	});

export const updateConnection = (conn: Connection) =>
	Effect.gen(function* () {
		yield* Atom.set(connectionsLoadingAtom, true);
		const app = yield* AppService;
		return yield* app.updateConnection(conn).pipe(
			Effect.flatMap(() => loadConnections),
			Effect.catchAll((e) => handleError(e)),
			Effect.andThen(() => Atom.set(connectionsLoadingAtom, false)),
		);
	});
